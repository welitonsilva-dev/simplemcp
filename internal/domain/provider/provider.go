package provider

import "humancli-server/internal/domain/plan"

// Provider define a interface usada pela camada de use case para chamar o LLM.
type Provider interface {
	Generate(prompt string) (string, error)
	Plan(history, tools string) (*plan.ExecutionPlan, error)
	Finalize(history string) (string, error)
}
