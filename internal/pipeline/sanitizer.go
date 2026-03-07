package pipeline

import (
	"errors"
	"strings"
)

var blockedWords = []string{
	"ignore instructions",
	"system prompt",
	"rm -rf",
	"sudo",
}

// SanitizeInput verifica se o input contém conteúdo perigoso
func SanitizeInput(input string) (string, error) {

	lower := strings.ToLower(input)

	for _, word := range blockedWords {

		if strings.Contains(lower, word) {
			return "", errors.New("entrada bloqueada por conter conteúdo potencialmente perigoso")
		}

	}

	return input, nil
}
