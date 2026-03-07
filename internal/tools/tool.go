package tools

// Tool define o contrato que todas as ferramentas devem implementar
type Tool interface {
	Name() string
	Description() string
	Execute(params map[string]interface{}) (any, error)
}
