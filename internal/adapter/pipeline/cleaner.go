package pipeline

import (
	"regexp"
	"strings"
)

var (
	// múltiplos espaços → um espaço
	multipleSpaces = regexp.MustCompile(`\s+`)

	// caracteres de controle e invisíveis (exceto espaço normal)
	controlChars = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)

	// pontuação repetida ex: "!!!", "???"
	repeatedPunctuation = regexp.MustCompile(`([!?.]){2,}`)
)

// clean remove ruídos comuns da entrada do usuário sem alterar
// o significado ou remover palavras.
func clean(input string) (string, error) {
	result := strings.ReplaceAll(input, "\n", " ")
	result = strings.ReplaceAll(result, "\r", " ")
	result = strings.ReplaceAll(result, "\t", " ")
	result = controlChars.ReplaceAllString(result, "")
	result = repeatedPunctuation.ReplaceAllString(result, "$1")
	result = multipleSpaces.ReplaceAllString(result, " ")
	return strings.TrimSpace(result), nil
}
