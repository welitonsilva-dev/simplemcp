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
	// --- logger ---
	logger.Info("🚀 iniciando servidor")

	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "logs"
	}
	if err := logger.Init(logDir); err != nil {
		log.Fatal(err)
	}

	// --- config ---
	logger.Info("🚀 carregando configurações")
	cfg := config.Load()

	// --- adapters ---
	logger.Info("🚀 configurando cliente LLM")
	llmClient := llm.NewClient(cfg.OllamaURL, cfg.Model)

	logger.Info("🚀 configurando pipeline")
	pipe := pipeline.New()

	logger.Info("🚀 configurando registry de tools")
	registry := tools.GlobalRegistry()

	// --- usecase ---
	logger.Info("🚀 configurando agente")
	agentUseCase := agent.New(pipe, llmClient, registry)

	// --- server ---
	logger.Info("🚀 configurando servidor HTTP")
	srv := server.New(cfg.Address, agentUseCase)

	logger.Info("🚀 MCP Server running on %s", cfg.Address)
	if err := srv.Start(); err != nil {
		logger.Error("fatal: %v", err)
	}
}
