package agent

import (
	"fmt"

	"simplemcp/internal/logger"
	"simplemcp/internal/tools"
)

// Run recebe input e o plano e executa as tools
func Run(plan Plan) (string, error) {
	// se o plano não tiver steps
	if len(plan.Steps) == 0 {
		return "não entendi", nil
	}

	results := ""
	for i, step := range plan.Steps {
		tool, exists := tools.GlobalRegistry.Get(step.Tool)
		if !exists {
			results += fmt.Sprintf("step %d: tool '%s' não existe\n", i+1, step.Tool)
			continue
		}

		if step.Params == nil {
			step.Params = map[string]interface{}{}
		}

		result, err := tool.Execute(step.Params)
		if err != nil {
			logger.Error("failed to execute tool '%s': %v", step.Tool, err.Error())
			results += fmt.Sprintf("step %d: erro ao executar tool '%s'\n", i+1, step.Tool)
			continue
		}

		// converte resultado para string
		str, ok := result.(string)
		if !ok {
			results += fmt.Sprintf("step %d: resultado da tool '%s' não é string\n", i+1, step.Tool)
			continue
		}

		results += fmt.Sprintf("step %d [%s]: %s\n", i+1, step.Tool, str)
	}

	return results, nil
}
