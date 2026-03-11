package agent

import (
	"fmt"

	"simplemcp/internal/adapter/llm"
	"simplemcp/internal/adapter/pipeline"
	"simplemcp/internal/domain/message"
	"simplemcp/internal/domain/tool"
	"simplemcp/internal/infra/logger"
)

// AgentUseCase orquestra o fluxo completo:
// pipeline → plan → execute → respond
type AgentUseCase struct {
	pipeline *pipeline.Pipeline
	llm      *llm.Client
	registry tool.ToolRegistry
}

// New cria um AgentUseCase com todas as dependências injetadas.
func New(p *pipeline.Pipeline, l *llm.Client, r tool.ToolRegistry) *AgentUseCase {
	return &AgentUseCase{
		pipeline: p,
		llm:      l,
		registry: r,
	}
}

// Execute recebe a mensagem do usuário e retorna a resposta após executar o plano.
func (a *AgentUseCase) Execute(msg message.UserMessage) (*message.AgentResponse, error) {
	// 1. pipeline: limpa, sanitiza e normaliza o input
	clean, err := a.pipeline.Process(msg.Content)
	if err != nil {
		logger.Error("pipeline error: %v", err)
		return nil, fmt.Errorf("entrada inválida: %w", err)
	}

	// 2. monta lista de tools disponíveis e envia ao LLM
	toolsCtx := buildToolsContext(a.registry)
	plan, err := a.llm.Plan(clean, toolsCtx)
	if err != nil {
		logger.Error("llm plan error: %v", err)
		return nil, fmt.Errorf("erro ao gerar plano: %w", err)
	}

	// 3. plano desconhecido — nenhuma tool correspondente
	if plan.IsUnknown() {
		return &message.AgentResponse{
			Results: []message.StepResult{
				{Tool: "unknown", Output: "não entendi o comando"},
			},
		}, nil
	}

	// 4. executa cada step do plano
	var results []message.StepResult
	for i, step := range plan.Steps {
		t, exists := a.registry.Get(step.Tool)
		if !exists {
			logger.Error("tool '%s' não encontrada (step %d)", step.Tool, i+1)
			results = append(results, message.StepResult{
				Tool:  step.Tool,
				Error: fmt.Sprintf("tool '%s' não existe", step.Tool),
			})
			continue
		}

		if step.Params == nil {
			step.Params = map[string]any{}
		}

		output, err := t.Execute(step.Params)
		if err != nil {
			logger.Error("tool '%s' error: %v", step.Tool, err)
			results = append(results, message.StepResult{
				Tool:  step.Tool,
				Error: err.Error(),
			})
			continue
		}

		results = append(results, message.StepResult{
			Tool:   step.Tool,
			Output: output,
		})
	}

	return &message.AgentResponse{Results: results}, nil
}

// buildToolsContext formata a lista de tools para o prompt do LLM.
func buildToolsContext(r tool.ToolRegistry) string {
	var ctx string
	for _, t := range r.All() {
		ctx += fmt.Sprintf("%s → %s\n", t.Name(), t.Description())
	}
	return ctx
}
