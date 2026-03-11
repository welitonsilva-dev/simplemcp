package pipeline

import "strings"

// verbosePhrases são frases de cortesia que não agregam ao comando.
var verbosePhrases = []string{
	"por favor,",
	"por favor ",
	"você poderia ",
	"você pode ",
	"sera que ",
	"gostaria que voce ",
	"gostaria que ",
	"preciso que voce ",
	"quero que voce ",
	"pode fazer ",
	"consegue ",
}

// optimize remove rodeios preservando o significado do comando.
// Não remove stopwords isoladas como "não", "sem", "nunca".
func optimize(input string) (string, error) {
	result := input

	for _, phrase := range verbosePhrases {
		if strings.HasPrefix(result, phrase) {
			result = strings.TrimPrefix(result, phrase)
			result = strings.TrimSpace(result)
		}
	}

	if len(result) > 0 {
		result = strings.ToUpper(result[:1]) + result[1:]
	}

	return strings.TrimSpace(result), nil
}
