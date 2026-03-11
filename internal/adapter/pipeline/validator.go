package pipeline

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	maxInputLength = 2000
	minInputLength = 2
)

func validate(input string) (string, error) {
	trimmed := strings.TrimSpace(input)

	if trimmed == "" {
		return "", fmt.Errorf("mensagem vazia")
	}
	if len([]rune(trimmed)) < minInputLength {
		return "", fmt.Errorf("mensagem muito curta")
	}
	if len(trimmed) > maxInputLength {
		return "", fmt.Errorf("mensagem muito longa: máximo %d caracteres", maxInputLength)
	}
	if !utf8.ValidString(trimmed) {
		return "", fmt.Errorf("encoding inválido: esperado UTF-8")
	}

	return trimmed, nil
}
