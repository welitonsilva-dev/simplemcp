package tools

import (
	"fmt"
	"sort"
	"strings"
)

var GlobalRegistry = NewRegistry()

// Registry armazena e gerencia todas as ferramentas disponíveis
type Registry struct {
	tools map[string]Tool
}

// NewRegistry cria um novo registry de ferramentas
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// Register adiciona uma nova tool ao registry
func (r *Registry) Register(t Tool) error {

	name := t.Name()

	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool '%s' já registrada", name)
	}

	r.tools[name] = t
	return nil
}

// Get retorna uma tool pelo nome
func (r *Registry) Get(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

// List retorna todas as tools registradas
func (r *Registry) List() []Tool {

	list := make([]Tool, 0, len(r.tools))

	for _, t := range r.tools {
		list = append(list, t)
	}

	return list
}

// Names retorna os nomes das tools ordenados
func (r *Registry) Names() []string {

	names := make([]string, 0, len(r.tools))

	for name := range r.tools {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// AvailableTools retorna uma lista formatada para prompt de LLM
func (r *Registry) AvailableTools() string {

	var builder strings.Builder

	names := r.Names()

	for _, name := range names {

		tool := r.tools[name]

		builder.WriteString(
			fmt.Sprintf("%s → %s\n",
				tool.Name(),
				tool.Description(),
			),
		)
	}

	return builder.String()
}
