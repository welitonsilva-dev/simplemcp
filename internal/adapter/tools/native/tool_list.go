package native

import (
	"fmt"

	"humancli-server/internal/adapter/tools"
)

func init() {
	tools.GlobalRegistry().Register(&ToolList{})
}

type ToolList struct{}

// ToolListResult é o output estruturado da tool tool_list.
type ToolListResult struct {
	Message string   `json:"message"`
	Native  []string `json:"native"`
	Plugins []string `json:"plugins"`
}

func (t *ToolList) Name() string {
	return "tool_list"
}

func (t *ToolList) Description() string {
	return `
Prioridade de interpretação:

Listar tools disponíveis
Palavras associadas:
- listar tools
- quais tools existem
- tools disponíveis
- listar ferramentas
- quais ferramentas tenho
- list tools

→ usar ferramenta "tool_list"

Descrição:
Lista todas as ferramentas registradas, separadas por nativas e plugins.

Parâmetros:
- nenhum
`
}

func (t *ToolList) Execute(params map[string]interface{}) (any, error) {
	native := tools.GlobalRegistry().ListByOrigin(tools.OriginNative)
	plugins := tools.GlobalRegistry().ListByOrigin(tools.OriginPlugin)

	// extrai só os nomes de cada grupo
	nativeNames := make([]string, 0, len(native))
	for _, tool := range native {
		nativeNames = append(nativeNames, tool.Name())
	}

	pluginNames := make([]string, 0, len(plugins))
	for _, tool := range plugins {
		pluginNames = append(pluginNames, tool.Name())
	}

	return ToolListResult{
		Message: fmt.Sprintf("%d ferramentas disponíveis (%d nativas, %d plugins)",
			len(native)+len(plugins), len(native), len(plugins)),
		Native:  nativeNames,
		Plugins: pluginNames,
	}, nil
}
