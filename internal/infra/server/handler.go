package server

import (
	"encoding/json"
	"net/http"

	"humancli-server/internal/domain/message"
	"humancli-server/internal/infra/logger"
	"humancli-server/internal/usecase/agent"
)

// Handler processa as requisições HTTP e delega ao usecase.
type Handler struct {
	agent *agent.AgentUseCase
}

// NewHandler cria um Handler com o usecase injetado.
func NewHandler(a *agent.AgentUseCase) *Handler {
	return &Handler{agent: a}
}

// /Do recebe o prompt do usuário e retorna o resultado da execução.
func (h *Handler) Do(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg message.UserMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.Error("handler decode error: %v", err)
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	response, err := h.agent.Execute(msg)
	if err != nil {
		logger.Error("agent execute error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Health retorna o status do servidor.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
