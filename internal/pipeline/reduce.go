package pipeline

import (
	"strings"
)

// Lista mais completa de stopwords
var stopwords = map[string]bool{
	// Português
	"a": true, "o": true, "e": true, "eu": true, "de": true, "do": true, "da": true, "olá": true,
	"dos": true, "das": true, "em": true, "no": true, "na": true, "nos": true, "nas": true,
	"por": true, "para": true, "com": true, "sem": true, "sobre": true, "entre": true,
	"um": true, "uma": true, "uns": true, "umas": true, "que": true, "quem": true,
	"como": true, "quando": true, "onde": true, "qual": true, "quais": true,
	"se": true, "não": true, "mais": true, "menos": true, "muito": true, "pouco": true,
	"ele": true, "ela": true, "eles": true, "elas": true, "isso": true, "aquilo": true,
	"este": true, "esta": true, "estes": true, "estas": true, "isto": true,
	// Inglês
	"the": true, "and": true, "of": true, "to": true, "in": true, "an": true,
	"for": true, "on": true, "with": true, "at": true, "by": true, "from": true,
	"is": true, "are": true, "was": true, "were": true, "be": true, "been": true,
	"this": true, "that": true, "these": true, "those": true,
	"it": true, "its": true, "as": true, "but": true, "or": true, "if": true,
}

// Função que reduz o input removendo stopwords
func ReduceInput(text string) string {
	words := strings.Fields(text) // quebra o texto em palavras
	result := []string{}

	for _, w := range words {
		lower := strings.ToLower(w)
		// se não for stopword, adiciona
		if !stopwords[lower] {
			result = append(result, w)
		}
	}

	return strings.Join(result, " ")
}
