package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr            string
	APIKey          string
	Provider        string
	Model           string
	OllamaURL       string
	InputMaxLength  int
	RateLimitIP     int
	RateLimitGlobal int
	RateLimitWindow time.Duration
}

func Load() *Config {
	inputMaxLength, err := strconv.Atoi(getEnv("INPUT_MAX_LENGTH", "1000"))
	if err != nil {
		inputMaxLength = 1000
	}

	rateLimitIP, err := strconv.Atoi(getEnv("RATE_LIMIT_PER_IP", "10"))
	if err != nil {
		rateLimitIP = 10
	}

	rateLimitGlobal, err := strconv.Atoi(getEnv("RATE_LIMIT_GLOBAL", "50"))
	if err != nil {
		rateLimitGlobal = 50
	}

	rateLimitWindow, err := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW", "60"))
	if err != nil {
		rateLimitWindow = 60
	}

	return &Config{
		Addr:            getEnv("SERVER_ADDR", ":8081"),
		APIKey:          getEnv("API_KEY", ""),
		Provider:        getEnv("LLM_PROVIDER", "ollama"),
		Model:           getEnv("LLM_MODEL", "qwen2.5:7b"),
		OllamaURL:       getEnv("OLLAMA_URL", "http://ollama:11434"),
		InputMaxLength:  inputMaxLength,
		RateLimitIP:     rateLimitIP,
		RateLimitGlobal: rateLimitGlobal,
		RateLimitWindow: time.Duration(rateLimitWindow) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
