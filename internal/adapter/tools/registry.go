package tools

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"simplemcp/internal/domain/tool"
)

// ToolOrigin identifica a origem de uma tool registrada.
type ToolOrigin string

const (
	OriginNative ToolOrigin = "native"
	OriginPlugin ToolOrigin = "plugin"
)

type registeredTool struct {
	tool   tool.Tool
	origin ToolOrigin
}

// Registry é a implementação concreta de domain/tool.ToolRegistry.
type Registry struct {
	mu    sync.RWMutex
	tools map[string]registeredTool
}

var (
	globalRegistry *Registry
	once           sync.Once
)

// GlobalRegistry retorna o singleton do registry.
// Usado nos init() de cada pacote de tool.
func GlobalRegistry() *Registry {
	once.Do(func() {
		globalRegistry = &Registry{
			tools: make(map[string]registeredTool),
		}
	})
	return globalRegistry
}

// Register adiciona uma tool ao registry detectando a origem automaticamente.
func (r *Registry) Register(t tool.Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := t.Name()
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool '%s' já registrada", name)
	}

	pkgPath := reflect.TypeOf(t).Elem().PkgPath()
	origin := OriginPlugin
	if strings.Contains(pkgPath, "tools/native") {
		origin = OriginNative
	}

	r.tools[name] = registeredTool{tool: t, origin: origin}
	return nil
}

// Get retorna uma tool pelo nome. Satisfaz domain/tool.ToolRegistry.
func (r *Registry) Get(name string) (tool.Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.tools[name]
	return entry.tool, ok
}

// All retorna todas as tools registradas. Satisfaz domain/tool.ToolRegistry.
func (r *Registry) All() []tool.Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]tool.Tool, 0, len(r.tools))
	for _, entry := range r.tools {
		list = append(list, entry.tool)
	}
	return list
}

// ListByOrigin retorna tools filtradas por origem.
func (r *Registry) ListByOrigin(origin ToolOrigin) []tool.Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []tool.Tool
	for _, entry := range r.tools {
		if entry.origin == origin {
			list = append(list, entry.tool)
		}
	}
	return list
}

// Names retorna os nomes de todas as tools em ordem alfabética.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
