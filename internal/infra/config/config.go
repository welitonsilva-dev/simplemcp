package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr       string
	APIKey     string
	Provider   string
	Model      string
	LLMAPIKey  string
	LLMBaseURL string
	OllamaURL  string

	InputMaxLength      int
	RateLimitIP         int
	RateLimitGlobal     int
	RateLimitWindow     time.Duration
	RequestTimeout      time.Duration
	ConfidenceThreshold float64

	// MaxIterations define o limite de ciclos do loop ReAct por requisição.
	// Evita loops infinitos em cenários onde o LLM não encerra por conta própria.
	// Padrão: 10
	MaxIterations int
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

	requestTimeout, err := strconv.Atoi(getEnv("REQUEST_TIMEOUT", "120"))
	if err != nil {
		requestTimeout = 120
	}

	confidenceThreshold, err := strconv.ParseFloat(getEnv("CONFIDENCE_THRESHOLD", "0.8"), 64)
	if err != nil {
		confidenceThreshold = 0.8
	}

	maxIterations, err := strconv.Atoi(getEnv("AGENT_MAX_ITERATIONS", "10"))
	if err != nil {
		maxIterations = 10
	}

	return &Config{
		Addr:                getEnv("SERVER_ADDR", ":8081"),
		APIKey:              getEnv("API_KEY", ""),
		Provider:            getEnv("HUMANCLI_PROVIDER", "ollama"),
		Model:               getEnv("HUMANCLI_MODEL", "qwen2.5:7b"),
		LLMAPIKey:           getEnv("LLM_API_KEY", ""),
		LLMBaseURL:          getEnv("LLM_BASE_URL", ""),
		OllamaURL:           getEnv("OLLAMA_URL", "http://ollama:11434"),
		InputMaxLength:      inputMaxLength,
		RateLimitIP:         rateLimitIP,
		RateLimitGlobal:     rateLimitGlobal,
		RateLimitWindow:     time.Duration(rateLimitWindow) * time.Second,
		RequestTimeout:      time.Duration(requestTimeout) * time.Second,
		ConfidenceThreshold: confidenceThreshold,
		MaxIterations:       maxIterations,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
