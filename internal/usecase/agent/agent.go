package agent

import (
	"fmt"
	"strings"

	"humancli-server/internal/adapter/llm"
	"humancli-server/internal/adapter/pipeline"
	"humancli-server/internal/domain/message"
	"humancli-server/internal/domain/tool"
	"humancli-server/internal/infra/logger"
)

const defaultMaxIterations = 10

// destructiveTools lista as tools que exigem confidence mínima para executar.
// Protege contra execuções destrutivas quando o LLM está inseguro sobre a intenção.
var destructiveTools = map[string]bool{
	"fs_rm":    true,
	"fs_rmdir": true,
	"fs_rmrf":  true,
}

// AgentUseCase orquestra o loop ReAct (Reason + Act) do agente:
//
//  1. Pipeline sanitiza o input
//  2. LLM raciocina sobre o histórico e decide a próxima ação
//  3. Tool é executada e o resultado entra no histórico
//  4. Loop repete até o LLM encerrar ou atingir maxIterations
type AgentUseCase struct {
	pipeline            *pipeline.Pipeline
	llm                 *llm.Client
	registry            tool.ToolRegistry
	confidenceThreshold float64
	maxIterations       int
}

// New cria um AgentUseCase com todas as dependências injetadas.
// maxIterations define o limite de ciclos do loop ReAct (0 usa o padrão de 10).
func New(p *pipeline.Pipeline, l *llm.Client, r tool.ToolRegistry, confidenceThreshold float64, maxIterations int) *AgentUseCase {
	if maxIterations <= 0 {
		maxIterations = defaultMaxIterations
	}
	return &AgentUseCase{
		pipeline:            p,
		llm:                 l,
		registry:            r,
		confidenceThreshold: confidenceThreshold,
		maxIterations:       maxIterations,
	}
}

// Execute roda o loop ReAct para a mensagem do usuário.
//
// A cada iteração, o LLM recebe o histórico completo da conversa
// (input original + resultados anteriores) e decide entre:
//   - Chamar uma tool e continuar o loop
//   - Encerrar e retornar a resposta final ao usuário
func (a *AgentUseCase) Execute(msg message.UserMessage) (*message.AgentResponse, error) {
	// 1. sanitiza o input via pipeline
	clean, err := a.pipeline.Process(msg.Content)
	if err != nil {
		logger.Error("pipeline error: %v", err)
		return nil, fmt.Errorf("entrada inválida: %w", err)
	}

	// histórico acumula: input original + "tool X retornou: Y" a cada iteração
	history := []string{fmt.Sprintf("usuário: %s", clean)}
	var results []message.StepResult
	toolsCtx := buildToolsContext(a.registry)

	for i := 0; i < a.maxIterations; i++ {
		logger.Info("agente: iteração %d/%d", i+1, a.maxIterations)

		// 2. LLM raciocina sobre o histórico e gera o próximo plano
		plan, err := a.llm.Plan(strings.Join(history, "\n"), toolsCtx)
		if err != nil {
			logger.Error("llm plan error: %v", err)
			return nil, fmt.Errorf("erro ao gerar plano: %w", err)
		}

		// 3. LLM decidiu encerrar — retorna resposta final
		if plan.IsFinal() {
			logger.Info("agente: loop encerrado pelo LLM na iteração %d", i+1)
			return &message.AgentResponse{
				Results:      results,
				FinalMessage: plan.FinalMessage,
			}, nil
		}

		// 4. plano desconhecido — nenhuma tool correspondente
		if plan.IsUnknown() {
			logger.Info("agente: intenção desconhecida na iteração %d", i+1)
			return &message.AgentResponse{
				Results: []message.StepResult{
					{Tool: "unknown", Output: "não entendi o comando"},
				},
			}, nil
		}

		// 5. executa o step da iteração atual
		step := plan.Steps[0]

		// confidence guard — bloqueia tools destrutivas se o LLM estiver inseguro
		if destructiveTools[step.Tool] && plan.Confidence < a.confidenceThreshold {
			logger.Error("confidence guard: tool '%s' bloqueada (confidence %.2f < %.2f)",
				step.Tool, plan.Confidence, a.confidenceThreshold)

			blocked := message.StepResult{
				Tool: step.Tool,
				Output: map[string]any{
					"requires_confirmation": true,
					"tool":                  step.Tool,
					"confidence":            plan.Confidence,
					"message": fmt.Sprintf(
						"comando bloqueado: não entendi com clareza suficiente (%.0f%% de certeza). "+
							"Seja mais específico ou confirme explicitamente a ação.",
						plan.Confidence*100,
					),
				},
			}
			results = append(results, blocked)
			// encerra o loop — não faz sentido continuar após bloqueio de segurança
			return &message.AgentResponse{Results: results}, nil
		}

		t, exists := a.registry.Get(step.Tool)
		if !exists {
			logger.Error("tool '%s' não encontrada (iteração %d)", step.Tool, i+1)
			errResult := message.StepResult{
				Tool:  step.Tool,
				Error: fmt.Sprintf("tool '%s' não existe", step.Tool),
			}
			results = append(results, errResult)
			history = append(history, fmt.Sprintf("tool %s falhou: %s", step.Tool, errResult.Error))
			continue
		}

		if step.Params == nil {
			step.Params = map[string]any{}
		}

		output, err := t.Execute(step.Params)
		if err != nil {
			logger.Error("tool '%s' error: %v (iteração %d)", step.Tool, err, i+1)
			errResult := message.StepResult{
				Tool:  step.Tool,
				Error: err.Error(),
			}
			results = append(results, errResult)
			// informa o LLM sobre o erro para que ele decida o próximo passo
			history = append(history, fmt.Sprintf("tool %s retornou erro: %s", step.Tool, err.Error()))
			continue
		}

		// sucesso — registra o resultado e atualiza o histórico
		results = append(results, message.StepResult{
			Tool:   step.Tool,
			Output: output,
		})
		history = append(history, fmt.Sprintf("tool %s retornou: %v", step.Tool, output))
		logger.Info("agente: tool '%s' executada com sucesso (iteração %d)", step.Tool, i+1)
	}

	// limite de iterações atingido — retorna o que foi coletado
	logger.Info("agente: limite de %d iterações atingido", a.maxIterations)
	return &message.AgentResponse{
		Results:      results,
		FinalMessage: fmt.Sprintf("limite de %d iterações atingido", a.maxIterations),
	}, nil
}

// buildToolsContext formata a lista de tools disponíveis para o prompt do LLM.
func buildToolsContext(r tool.ToolRegistry) string {
	var ctx string
	for _, t := range r.All() {
		ctx += fmt.Sprintf("%s → %s\n", t.Name(), t.Description())
	}
	return ctx
}
