package tools

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

type ToolOrigin string

const (
	OriginNative ToolOrigin = "native"
	OriginPlugin ToolOrigin = "plugin"
)

type registeredTool struct {
	tool   Tool
	origin ToolOrigin
}

type Registry struct {
	mu    sync.RWMutex
	tools map[string]registeredTool
}

var (
	globalRegistry *Registry
	once           sync.Once
)

func GlobalRegistry() *Registry {
	once.Do(func() {
		globalRegistry = &Registry{
			tools: make(map[string]registeredTool),
		}
	})
	return globalRegistry
}

// NewRegistry mantém compatibilidade — sempre retorna o singleton
func NewRegistry() *Registry {
	return GlobalRegistry()
}

func (r *Registry) Register(t Tool) error {
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

func (r *Registry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.tools[name]
	return entry.tool, ok
}

func (r *Registry) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]Tool, 0, len(r.tools))
	for _, entry := range r.tools {
		list = append(list, entry.tool)
	}
	return list
}

func (r *Registry) ListByOrigin(origin ToolOrigin) []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []Tool
	for _, entry := range r.tools {
		if entry.origin == origin {
			list = append(list, entry.tool)
		}
	}
	return list
}

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

func (r *Registry) AvailableTools() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var builder strings.Builder
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		entry := r.tools[name]
		builder.WriteString(fmt.Sprintf("%s → %s\n", entry.tool.Name(), entry.tool.Description()))
	}
	return builder.String()
}
