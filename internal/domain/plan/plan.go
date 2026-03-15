package plan

// ExecutionPlan é o plano gerado pela LLM em cada iteração do loop do agente.
// Pode representar uma ação a executar ou a decisão de encerrar a conversa.
type ExecutionPlan struct {
	// Steps é a lista de ferramentas a executar nesta iteração.
	// Contém exatamente um step por iteração no loop ReAct.
	Steps []ToolCall `json:"steps"`

	// Confidence é o grau de certeza da LLM sobre o plano (0.0 a 1.0).
	Confidence float64 `json:"confidence"`

	// Final indica que o agente decidiu encerrar o loop e retornar ao usuário.
	Final bool `json:"final"`

	// FinalMessage é a resposta em linguagem natural enviada ao usuário
	// quando Final=true. Só é preenchido quando Final=true.
	FinalMessage string `json:"final_message,omitempty"`
}

// ToolCall representa uma única chamada de ferramenta dentro do plano.
type ToolCall struct {
	// Tool é o nome da ferramenta — deve corresponder a Tool.Name().
	Tool string `json:"tool"`

	// Params são os parâmetros a passar para Tool.Execute().
	Params map[string]any `json:"params"`
}

// IsFinal retorna true quando o agente decidiu encerrar o loop.
// Isso ocorre quando a LLM retorna Final=true ou quando nenhum step foi gerado.
func (p *ExecutionPlan) IsFinal() bool {
	return p.Final || len(p.Steps) == 0
}

// IsUnknown retorna true se o plano não encontrou ferramenta correspondente.
func (p *ExecutionPlan) IsUnknown() bool {
	return len(p.Steps) == 1 && p.Steps[0].Tool == "unknown"
}
