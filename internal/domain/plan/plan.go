package plan

// ExecutionPlan é o plano gerado pela LLM.
// Representa o JSON retornado pelo Ollama após interpretar o prompt do usuário.
type ExecutionPlan struct {
	// Steps é a lista de ferramentas a executar, em ordem.
	Steps []ToolCall `json:"steps"`

	// Confidence é o grau de certeza da LLM sobre o plano (0.0 a 1.0).
	Confidence float64 `json:"confidence"`
}

// ToolCall representa um único passo do plano:
// qual tool executar e com quais parâmetros.
type ToolCall struct {
	// Tool é o nome da ferramenta — deve corresponder a Tool.Name().
	Tool string `json:"tool"`

	// Params são os parâmetros a passar para Tool.Execute().
	Params map[string]any `json:"params"`
}

// PlanResult é o resultado da execução de um ToolCall.
// O usecase retorna uma lista de PlanResult para o handler.
type PlanResult struct {
	// Tool é o nome da ferramenta executada.
	Tool string `json:"tool"`

	// Output é o resultado retornado por Tool.Execute().
	Output any `json:"output"`

	// Err contém o erro caso a execução tenha falhado.
	// nil se a execução foi bem-sucedida.
	Err error `json:"error,omitempty"`
}

// IsUnknown retorna true se o plano não encontrou ferramenta correspondente.
func (p *ExecutionPlan) IsUnknown() bool {
	return len(p.Steps) == 1 && p.Steps[0].Tool == "unknown"
}
