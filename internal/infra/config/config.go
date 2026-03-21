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
	MaxIterations       int
	SessionTTL          time.Duration

	// SessionDBPath é o caminho para o arquivo SQLite de sessões.
	// Quando definido, as sessões sobrevivem a reinicializações do servidor.
	// Quando vazio, o servidor usa armazenamento em memória (padrão).
	//
	// Exemplo: SESSION_DB_PATH=data/sessions.db
	SessionDBPath string
}

func Load() *Config {
	inputMaxLength, _ := strconv.Atoi(getEnv("INPUT_MAX_LENGTH", "1000"))
	if inputMaxLength == 0 {
		inputMaxLength = 1000
	}

	rateLimitIP, _ := strconv.Atoi(getEnv("RATE_LIMIT_PER_IP", "10"))
	if rateLimitIP == 0 {
		rateLimitIP = 10
	}

	rateLimitGlobal, _ := strconv.Atoi(getEnv("RATE_LIMIT_GLOBAL", "50"))
	if rateLimitGlobal == 0 {
		rateLimitGlobal = 50
	}

	rateLimitWindow, _ := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW", "60"))
	if rateLimitWindow == 0 {
		rateLimitWindow = 60
	}

	requestTimeout, _ := strconv.Atoi(getEnv("REQUEST_TIMEOUT", "120"))
	if requestTimeout == 0 {
		requestTimeout = 120
	}

	confidenceThreshold, err := strconv.ParseFloat(getEnv("CONFIDENCE_THRESHOLD", "0.8"), 64)
	if err != nil {
		confidenceThreshold = 0.8
	}

	maxIterations, _ := strconv.Atoi(getEnv("AGENT_MAX_ITERATIONS", "10"))
	if maxIterations == 0 {
		maxIterations = 10
	}

	sessionTTLMin, _ := strconv.Atoi(getEnv("SESSION_TTL_MINUTES", "30"))
	if sessionTTLMin == 0 {
		sessionTTLMin = 30
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
		SessionTTL:          time.Duration(sessionTTLMin) * time.Minute,
		SessionDBPath:       getEnv("SESSION_DB_PATH", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
