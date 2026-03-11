package pipeline

import (
	"strings"
	"unicode"
)

// accentMap substitui letras acentuadas pela letra base sem dependência externa.
var accentMap = strings.NewReplacer(
	"á", "a", "à", "a", "â", "a", "ã", "a", "ä", "a",
	"é", "e", "è", "e", "ê", "e", "ë", "e",
	"í", "i", "ì", "i", "î", "i", "ï", "i",
	"ó", "o", "ò", "o", "ô", "o", "õ", "o", "ö", "o",
	"ú", "u", "ù", "u", "û", "u", "ü", "u",
	"ç", "c", "ñ", "n",
	"Á", "a", "À", "a", "Â", "a", "Ã", "a", "Ä", "a",
	"É", "e", "È", "e", "Ê", "e", "Ë", "e",
	"Í", "i", "Ì", "i", "Î", "i", "Ï", "i",
	"Ó", "o", "Ò", "o", "Ô", "o", "Õ", "o", "Ö", "o",
	"Ú", "u", "Ù", "u", "Û", "u", "Ü", "u",
	"Ç", "c", "Ñ", "n",
)

// normalize padroniza o formato sem remover palavras.
// Substitui o reducer original que removia stopwords —
// perigoso para LLM (ex: "não delete" → "delete").
func normalize(input string) (string, error) {
	result := strings.ToLower(input)
	result = accentMap.Replace(result)

	// remove combining marks residuais (codepoints separados)
	var b strings.Builder
	b.Grow(len(result))
	for _, r := range result {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}

	return strings.TrimSpace(b.String()), nil
}
