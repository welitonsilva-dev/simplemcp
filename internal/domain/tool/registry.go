// internal/domain/tool/registry.go
package tool

// ToolRegistry é a interface de acesso ao registro de tools.
// O usecase depende desta interface — nunca da implementação concreta.
type ToolRegistry interface {
	// Get retorna uma tool pelo nome.
	// Retorna false se a tool não existir.
	Get(name string) (Tool, bool)

	// All retorna todas as tools registradas.
	// Usado pelo Planner para montar o contexto enviado à LLM.
	All() []Tool
}
