package agent

import (
	"encoding/json"
	"fmt"
	"strings"

	"humancli-server/internal/adapter/llm"
	"humancli-server/internal/adapter/pipeline"
	"humancli-server/internal/domain/message"
	domainSession "humancli-server/internal/domain/session"
	"humancli-server/internal/domain/tool"
	"humancli-server/internal/infra/logger"
)

const defaultMaxIterations = 10

// destructiveTools exigem confidence mínima para executar.
var destructiveTools = map[string]bool{
	"fs_rm":    true,
	"fs_rmdir": true,
	"fs_rmrf":  true,
}

// AgentUseCase orquestra o loop ReAct com dois prompts separados:
//
//  1. plannerPrompt  → LLM decide qual tool chamar (JSON simples)
//  2. finalizerPrompt → LLM gera resposta final em linguagem natural
//
// Separar os prompts resolve o problema de modelos pequenos (7b) que
// não conseguem escolher consistentemente entre dois formatos JSON diferentes.
type AgentUseCase struct {
	pipeline            *pipeline.Pipeline
	llm                 *llm.Client
	registry            tool.ToolRegistry
	sessions            domainSession.Store
	confidenceThreshold float64
	maxIterations       int
}

// New cria um AgentUseCase com todas as dependências injetadas.
func New(
	p *pipeline.Pipeline,
	l *llm.Client,
	r tool.ToolRegistry,
	sessions domainSession.Store,
	confidenceThreshold float64,
	maxIterations int,
) *AgentUseCase {
	if maxIterations <= 0 {
		maxIterations = defaultMaxIterations
	}
	return &AgentUseCase{
		pipeline:            p,
		llm:                 l,
		registry:            r,
		sessions:            sessions,
		confidenceThreshold: confidenceThreshold,
		maxIterations:       maxIterations,
	}
}

// Execute roda o loop e retorna a resposta consolidada. Usado em /v1/do.
func (a *AgentUseCase) Execute(msg message.UserMessage) (*message.AgentResponse, error) {
	var results []message.StepResult
	var finalMessage string

	emit := func(event message.StreamEvent) {
		if event.Type == "step" {
			results = append(results, message.StepResult{
				Tool:   event.Tool,
				Output: event.Output,
				Error:  event.Error,
			})
		}
		if event.Type == "final" {
			finalMessage = event.Message
		}
	}

	if err := a.run(msg, emit); err != nil {
		return nil, err
	}

	return &message.AgentResponse{
		Results:      results,
		FinalMessage: finalMessage,
	}, nil
}

// ExecuteStream roda o loop emitindo eventos SSE em tempo real. Usado em /v1/stream.
func (a *AgentUseCase) ExecuteStream(msg message.UserMessage, emit func(message.StreamEvent)) error {
	return a.run(msg, emit)
}

// run é o núcleo compartilhado entre Execute e ExecuteStream.
//
// Fluxo por iteração:
//  1. plannerPrompt → LLM retorna {"tool": "...", "params": {}}
//  2. Se tool="none" → chama finalizerPrompt e encerra
//  3. Senão → executa a tool, adiciona resultado ao histórico, repete
func (a *AgentUseCase) run(msg message.UserMessage, emit func(message.StreamEvent)) error {
	clean, err := a.pipeline.Process(msg.Content)
	if err != nil {
		return fmt.Errorf("entrada inválida: %w", err)
	}

	sess := a.sessions.Get(msg.SessionID)
	sess.Append(fmt.Sprintf("usuário: %s", clean))

	toolsCtx := buildToolsContext(a.registry)

	// deduplicação: chave = "tool:params_json"
	executed := make(map[string]bool)

	for i := 0; i < a.maxIterations; i++ {
		logger.Info("agente: iteração %d/%d (sessão: %s)", i+1, a.maxIterations, msg.SessionID)

		// 1. planner decide a próxima tool
		plan, err := a.llm.Plan(strings.Join(sess.History, "\n"), toolsCtx)
		if err != nil {
			logger.Error("llm plan error: %v", err)
			return a.finalize(sess, emit, "não consegui processar o comando")
		}

		// 2. tool="none" ou final=true → gera resposta e encerra
		if plan.IsFinal() || plan.IsUnknown() {
			return a.finalize(sess, emit, "")
		}

		step := plan.Steps[0]

		// 3. deduplicação — evita loop com a mesma chamada
		paramsJSON, _ := json.Marshal(step.Params)
		dedupKey := fmt.Sprintf("%s:%s", step.Tool, paramsJSON)
		if executed[dedupKey] {
			logger.Info("agente: '%s' repetida — encerrando", step.Tool)
			return a.finalize(sess, emit, "")
		}
		executed[dedupKey] = true

		// 4. confidence guard para tools destrutivas
		if destructiveTools[step.Tool] && plan.Confidence < a.confidenceThreshold {
			msg := fmt.Sprintf(
				"comando bloqueado: não entendi com clareza suficiente (%.0f%% de certeza). Seja mais específico.",
				plan.Confidence*100,
			)
			sess.Append(fmt.Sprintf("agente: %s", msg))
			a.sessions.Save(sess)
			emit(message.StreamEvent{Type: "final", Message: msg, Iteration: i + 1})
			return nil
		}

		// 5. executa a tool
		t, exists := a.registry.Get(step.Tool)
		if !exists {
			errMsg := fmt.Sprintf("tool '%s' não existe", step.Tool)
			sess.Append(fmt.Sprintf("tool %s falhou: %s", step.Tool, errMsg))
			emit(message.StreamEvent{Type: "step", Tool: step.Tool, Error: errMsg, Iteration: i + 1})
			continue
		}

		if step.Params == nil {
			step.Params = map[string]any{}
		}

		output, err := t.Execute(step.Params)
		if err != nil {
			logger.Error("tool '%s' error: %v", step.Tool, err)
			sess.Append(fmt.Sprintf("tool %s retornou erro: %s", step.Tool, err.Error()))
			emit(message.StreamEvent{Type: "step", Tool: step.Tool, Error: err.Error(), Iteration: i + 1})
			continue
		}

		// 6. sucesso — emite evento e atualiza histórico
		sess.Append(fmt.Sprintf("tool %s retornou: %v", step.Tool, output))
		emit(message.StreamEvent{Type: "step", Tool: step.Tool, Output: output, Iteration: i + 1})
		logger.Info("agente: tool '%s' OK (iteração %d)", step.Tool, i+1)
	}

	// limite atingido
	return a.finalize(sess, emit, "")
}

// finalize chama o finalizerPrompt para gerar a resposta final e salva a sessão.
// Se fallback não for vazio, usa ele diretamente sem chamar o LLM.
func (a *AgentUseCase) finalize(sess *domainSession.Session, emit func(message.StreamEvent), fallback string) error {
	msg := fallback
	if msg == "" {
		var err error
		msg, err = a.llm.Finalize(strings.Join(sess.History, "\n"))
		if err != nil {
			logger.Error("finalize error: %v", err)
			msg = "tarefa concluída"
		}
	}

	sess.Append(fmt.Sprintf("agente: %s", msg))
	a.sessions.Save(sess)
	emit(message.StreamEvent{Type: "final", Message: msg})
	return nil
}

// buildToolsContext formata a lista de tools disponíveis para o prompt.
func buildToolsContext(r tool.ToolRegistry) string {
	var ctx string
	for _, t := range r.All() {
		ctx += fmt.Sprintf("%s → %s\n", t.Name(), t.Description())
	}
	return ctx
}
