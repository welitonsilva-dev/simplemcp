package native

import (
	"fmt"
	"strings"

	"simplemcp/internal/tools"
)

func init() {
	tools.GlobalRegistry().Register(&ToolList{})
}

type ToolList struct{}

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

	var sb strings.Builder

	sb.WriteString("nativas:\n")
	if len(native) == 0 {
		sb.WriteString("  nenhuma\n")
	}
	for _, tool := range native {
		sb.WriteString(fmt.Sprintf("  - %s\n", tool.Name()))
	}

	sb.WriteString("\nplugins:\n")
	if len(plugins) == 0 {
		sb.WriteString("  nenhum\n")
	}
	for _, tool := range plugins {
		sb.WriteString(fmt.Sprintf("  - %s\n", tool.Name()))
	}

	return sb.String(), nil
}
