package main

import (
	"log"
	"os"

	"humancli-server/internal/adapter/llm"
	"humancli-server/internal/adapter/pipeline"
	"humancli-server/internal/adapter/tools"
	domainSession "humancli-server/internal/domain/session"
	"humancli-server/internal/infra/config"
	"humancli-server/internal/infra/logger"
	"humancli-server/internal/infra/server"
	infraSession "humancli-server/internal/infra/session"
	"humancli-server/internal/usecase/agent"

	"github.com/joho/godotenv"

	// ferramentas nativas
	_ "humancli-server/internal/adapter/tools/native/echo"
	_ "humancli-server/internal/adapter/tools/native/filesystem"

	// plugins externos
	_ "github.com/weliton/humancli-plugins/hello"
)

func main() {
	logger.Info("🚀 iniciando servidor")

	if err := godotenv.Load(); err != nil {
		logger.Info("⚠️  arquivo .env não encontrado ou erro ao carregar: %v", err)
	}

	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "logs"
	}
	if err := logger.Init(logDir); err != nil {
		log.Fatal(err)
	}

	logger.Info("🚀 carregando configurações")
	cfg := config.Load()

	if cfg.APIKey == "" {
		log.Fatal("API_KEY não definida no .env — defina uma chave para proteger o servidor")
	}

	logger.Info("🚀 configurando cliente LLM")
	llmClient := llm.NewProvider(cfg)

	logger.Info("🚀 configurando pipeline")
	pipe := pipeline.New()

	logger.Info("🚀 configurando registry de tools")
	registry := tools.GlobalRegistry()

	logger.Info("🚀 configurando session store")
	sessions := buildSessionStore(cfg)

	logger.Info("🚀 configurando agente (max_iterations=%d)", cfg.MaxIterations)
	agentUseCase := agent.New(
		pipe,
		llmClient,
		registry,
		sessions,
		cfg.ConfidenceThreshold,
		cfg.MaxIterations,
	)

	logger.Info("🚀 configurando servidor HTTP")
	srv := server.New(
		cfg.Addr,
		cfg.APIKey,
		cfg.RateLimitIP,
		cfg.RateLimitGlobal,
		cfg.RateLimitWindow,
		cfg.RequestTimeout,
		agentUseCase,
	)

	logger.Info("🚀 humancli-server rodando em %s", cfg.Addr)
	if err := srv.Start(); err != nil {
		logger.Error("fatal: %v", err)
	}
}

// buildSessionStore escolhe o store de sessões conforme a configuração.
// SQLite é o padrão quando SESSION_DB_PATH está definido.
// Fallback para memória se o caminho estiver vazio ou o banco falhar ao abrir.
func buildSessionStore(cfg *config.Config) domainSession.Store {
	if cfg.SessionDBPath != "" {
		store, err := infraSession.NewSQLiteStore(cfg.SessionDBPath, cfg.SessionTTL)
		if err != nil {
			logger.Error("falha ao abrir SQLite (%s): %v — usando memória", cfg.SessionDBPath, err)
		} else {
			logger.Info("sessões persistidas em SQLite: %s (TTL: %s)", cfg.SessionDBPath, cfg.SessionTTL)
			return store
		}
	}

	logger.Info("sessões em memória (TTL: %s)", cfg.SessionTTL)
	return infraSession.NewMemoryStore(cfg.SessionTTL)
}
