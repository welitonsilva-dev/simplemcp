package pipeline

import (
	"strings"
)

// CleanInput remove ruídos comuns da entrada do usuário
func CleanInput(input string) string {

	input = strings.TrimSpace(input)

	// remove quebras de linha
	input = strings.ReplaceAll(input, "\n", " ")

	// remove espaços duplicados
	input = strings.Join(strings.Fields(input), " ")

	return input
}
