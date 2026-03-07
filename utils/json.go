package utils

import (
	"encoding/json"
	"regexp"
	"strings"
)

// CleanJSON tenta extrair um JSON válido de uma string, removendo ruídos comuns como blocos de código e texto adicional
func CleanJSON(input string) string {
	input = strings.TrimSpace(input)

	input = strings.ReplaceAll(input, "```json", "")
	input = strings.ReplaceAll(input, "```", "")

	re := regexp.MustCompile(`(?s)//.*?\n|/\*.*?\*/`)
	input = re.ReplaceAllString(input, "")

	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")

	if start >= 0 && end >= 0 {
		return input[start : end+1]
	}

	input = removeTrailingCommas(input)

	if json.Valid([]byte(input)) {
		return input
	}

	return input
}

func removeTrailingCommas(input string) string {
	re := regexp.MustCompile(`,\s*([}\]])`)
	return re.ReplaceAllString(input, "$1")
}
