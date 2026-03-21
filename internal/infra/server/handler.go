package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"humancli-server/internal/domain/message"
	"humancli-server/internal/infra/logger"
	"humancli-server/internal/usecase/agent"
)

// Handler processa as requisições HTTP e delega ao AgentUseCase.
type Handler struct {
	agent *agent.AgentUseCase
}

// NewHandler cria um Handler com o usecase injetado.
func NewHandler(a *agent.AgentUseCase) *Handler {
	return &Handler{agent: a}
}

// Do executa o agente e retorna a resposta consolidada em JSON.
// Use quando o cliente não suporta SSE ou prefere esperar a resposta completa.
//
// POST /v1/do
// Body: { "session_id": "...", "message": "..." }
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

// Stream executa o agente e transmite cada evento via Server-Sent Events (SSE).
// O cliente recebe um evento por iteração do loop — sem esperar o loop terminar.
//
// POST /v1/stream
// Body: { "session_id": "...", "message": "..." }
//
// Formato dos eventos SSE:
//
//	data: {"type":"step","tool":"fs_mkdir","output":"pasta criada","iteration":1}
//	data: {"type":"step","tool":"fs_touch","output":"arquivo criado","iteration":2}
//	data: {"type":"final","message":"Pasta criada com README.md.","iteration":3}
func (h *Handler) Stream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// verifica se o cliente suporta SSE (flusher)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	var msg message.UserMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.Error("stream decode error: %v", err)
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// headers SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // desativa buffer do nginx se houver proxy

	// emit envia cada evento como uma linha SSE e faz flush imediato
	emit := func(event message.StreamEvent) {
		data, err := json.Marshal(event)
		if err != nil {
			logger.Error("stream marshal error: %v", err)
			return
		}
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	if err := h.agent.ExecuteStream(msg, emit); err != nil {
		logger.Error("agent stream error: %v", err)
		emit(message.StreamEvent{Type: "error", Error: err.Error()})
	}
}

// Health retorna o status do servidor.
//
// GET /health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
