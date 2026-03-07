package pipeline

import (
	"fmt"

	"simplemcp/internal/config"
)

// OptimizeInput executa todo pipeline antes do LLM
func OptimizeInput(input string) (string, error) {
	cfg := config.Load()

	input = CleanInput(input)

	input = ReduceInput(input)

	input, err := SanitizeInput(input)
	if err != nil {
		return "", err
	}

	max := cfg.InputMaxLength
	if len(input) > max {
		input = input[:max]
	}

	return input, nil
}

// DebugPipeline mostra as etapas (útil para log)
func DebugPipeline(input string) {

	fmt.Println("INPUT ORIGINAL:", input)

	clean := CleanInput(input)
	fmt.Println("CLEAN:", clean)

	safe, err := SanitizeInput(clean)
	if err != nil {
		fmt.Println("SANITIZE ERROR:", err)
		return
	}

	fmt.Println("SANITIZED:", safe)
}
