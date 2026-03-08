package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"simplemcp/internal/agent"
	"simplemcp/internal/config"
	"simplemcp/internal/llm"
	"simplemcp/internal/logger"
	"simplemcp/internal/tools"

	// pacotes de ferramentas nativas
	_ "simplemcp/internal/tools/native/filesystem"
	_ "simplemcp/internal/tools/native/testinterno"

	// pacotes de ferramentas externas/plugins
	_ "github.com/weliton/simplemcpplugins/dockercmd"
	_ "github.com/weliton/simplemcpplugins/hello"
)

func main() {

	logger.Info("🚀 iniciando main")
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "logs"
	}
	if err := logger.Init(logDir); err != nil {
		log.Fatal(err)
	}

	logger.Info("🚀 carregando configurações")
	cfg := config.Load()

	logger.Info("🚀 configurando modelo llm")
	model := cfg.Model
	if model == "" {
		model = "llama3"
	}

	logger.Info("🚀 configurando URL do Ollama")
	ollamaURL := cfg.OllamaURL
	if ollamaURL == "" {
		ollamaURL = "http://ollama:11434"
	}

	logger.Info("🚀 configurando cliente LLM")
	llmClient := llm.Client{
		BaseURL: ollamaURL,
		Model:   model,
	}

	logger.Info("🚀 inicializando registro de ferramentas")
	toolRegistry := tools.GlobalRegistry

	// -------------------------------
	// Endpoints
	// -------------------------------

	logger.Info("🚀 configurando endpoints /v2/chat")
	// Chat v2 (Planner multi-step && agente executa o plano)
	http.HandleFunc("/v2/chat", func(w http.ResponseWriter, r *http.Request) {

		var body struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			logger.Error("invalid body: %v", err)
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		// Cria planner
		planner := agent.NewPlanner(llmClient)

		plan, err := planner.Generate(body.Message, toolRegistry().AvailableTools())
		if err != nil {
			logger.Error("failed to generate plan: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Agente executa o plano e retorna resultados
		results, err := agent.Run(*plan)
		if err != nil {
			logger.Error("failed to execute plan: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retorna JSON completo com steps e parâmetros
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"results": results,
		})
		if err != nil {
			logger.Error("failed to encode response: %v", err)
			return
		}
	})

	logger.Info("🚀 MCP Server running on :8081")
	logger.Error("fatal: %v", http.ListenAndServe(":8081", nil))
}
