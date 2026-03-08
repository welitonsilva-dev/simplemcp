package sdk

import "simplemcp/internal/tools"

type Tool interface {
	Name() string
	Description() string
	Execute(params map[string]interface{}) (any, error)
}

func Register(t Tool) {
	tools.GlobalRegistry().Register(t)
}
