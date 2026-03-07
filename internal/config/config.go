package config

import (
	"os"
	"strconv"
)

type Config struct {
	Provider       string
	Model          string
	OllamaURL      string
	InputMaxLength int
}

func Load() *Config {
	inputMaxLength, err := strconv.Atoi(getEnv("INPUT_MAX_LENGTH", "1000"))
	if err != nil {
		inputMaxLength = 1000
	}

	return &Config{
		Provider:       getEnv("LLM_PROVIDER", "ollama"),
		Model:          getEnv("LLM_MODEL", "qwen2.5:7b"),
		OllamaURL:      getEnv("OLLAMA_URL", "http://ollama:11434"),
		InputMaxLength: inputMaxLength,
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
