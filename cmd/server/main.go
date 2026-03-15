package main

import (
	"log"
	"os"

	"humancli-server/internal/adapter/llm"
	"humancli-server/internal/adapter/pipeline"
	"humancli-server/internal/adapter/tools"
	"humancli-server/internal/infra/config"
	"humancli-server/internal/infra/logger"
	"humancli-server/internal/infra/server"
	"humancli-server/internal/usecase/agent"

	// pacotes de ferramentas nativas
	_ "humancli-server/internal/adapter/tools/native/echo"
	_ "humancli-server/internal/adapter/tools/native/filesystem"

	// pacotes de ferramentas externas/plugins
	_ "github.com/weliton/humancli-plugins/dockercmd"
	_ "github.com/weliton/humancli-plugins/hello"
)

func main() {
	logger.Info("🚀 iniciando servidor")

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
	llmClient := llm.NewClient(cfg.OllamaURL, cfg.Model)

	logger.Info("🚀 configurando pipeline")
	pipe := pipeline.New()

	logger.Info("🚀 configurando registry de tools")
	registry := tools.GlobalRegistry()

	logger.Info("🚀 configurando agente")
	agentUseCase := agent.New(pipe, llmClient, registry, cfg.ConfidenceThreshold)

	logger.Info("🚀 configurando servidor HTTP")
	srv := server.New(cfg.Addr, cfg.APIKey, cfg.RateLimitIP, cfg.RateLimitGlobal, cfg.RateLimitWindow, cfg.RequestTimeout, agentUseCase)

	logger.Info("🚀 MCP Server running on %s", cfg.Addr)
	if err := srv.Start(); err != nil {
		logger.Error("fatal: %v", err)
	}
}
