package tool

// Tool é a interface que toda ferramenta deve implementar.
// Nativa ou plugin externo — o contrato é o mesmo.
type Tool interface {
	// Name retorna o identificador único da tool.
	// É o nome usado pela LLM no campo "tool" do JSON.
	Name() string

	// Description retorna a descrição da tool.
	// É enviada para a LLM como contexto para ela decidir qual usar.
	Description() string

	// Execute executa a tool com os parâmetros fornecidos.
	// params vem diretamente do campo "params" do ExecutionPlan.
	Execute(params map[string]any) (any, error)
}
