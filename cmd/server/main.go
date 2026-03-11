package main

import (
	"log"
	"os"

	"simplemcp/internal/adapter/llm"
	"simplemcp/internal/adapter/pipeline"
	"simplemcp/internal/adapter/tools"
	"simplemcp/internal/infra/config"
	"simplemcp/internal/infra/logger"
	"simplemcp/internal/infra/server"
	"simplemcp/internal/usecase/agent"

	// pacotes de ferramentas nativas
	_ "simplemcp/internal/adapter/tools/native/echo"
	_ "simplemcp/internal/adapter/tools/native/filesystem"

	// pacotes de ferramentas externas/plugins
	_ "github.com/weliton/simplemcpplugins/dockercmd"
	_ "github.com/weliton/simplemcpplugins/hello"
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
