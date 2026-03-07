package main

import (
	"encoding/json"
	"log"
	"net/http"

	"simplemcp/internal/agent"
	"simplemcp/internal/config"
	"simplemcp/internal/llm"
	"simplemcp/internal/server"
	"simplemcp/internal/tools"
	// pacotes de ferramentas nativas
	_ "simplemcp/internal/tools/native/testinterno"
	_ "simplemcp/internal/tools/native/filesystem"
)

func main() {
	cfg := config.Load()

	model := cfg.Model
	if model == "" {
		model = "llama3"
	}

	ollamaURL := cfg.OllamaURL
	if ollamaURL == "" {
		ollamaURL = "http://ollama:11434"
	}

	log.Println("🧠 Model:", model)
	log.Println("🔗 Ollama:", ollamaURL)

	llmClient := llm.Client{
		BaseURL: ollamaURL,
		Model:   model,
	}

	toolRegistry := tools.GlobalRegistry

	srv := server.NewServer(toolRegistry.List())

	// -------------------------------
	// Endpoints
	// -------------------------------

	// MCP tradicional
	http.HandleFunc("/mcp", srv.TestTools)

	// Chat v1 (Planner multi-step)
	http.HandleFunc("/v1/chat", func(w http.ResponseWriter, r *http.Request) {

		var body struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		// Cria planner
		planner := agent.NewPlanner(llmClient)

		plan, err := planner.Generate(body.Message, toolRegistry.AvailableTools())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retorna JSON completo com steps e parâmetros
		json.NewEncoder(w).Encode(plan)
	})

	// Chat v2 (Planner multi-step && agente executa o plano)
	http.HandleFunc("/v2/chat", func(w http.ResponseWriter, r *http.Request) {

		var body struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		// Cria planner
		planner := agent.NewPlanner(llmClient)

		plan, err := planner.Generate(body.Message, toolRegistry.AvailableTools())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Agente executa o plano e retorna resultados
		results, err := agent.Run(*plan)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retorna JSON completo com steps e parâmetros
		json.NewEncoder(w).Encode(map[string]interface{}{
			"results": results,
		})
	})

	log.Println("🚀 MCP Server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
