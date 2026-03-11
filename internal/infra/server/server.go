package server

import (
	"net/http"

	"simplemcp/internal/usecase/agent"
)

// Server encapsula o servidor HTTP.
type Server struct {
	httpServer *http.Server
}

// New cria e configura o servidor com todas as rotas.
func New(addr string, apiKey string, agentUseCase *agent.AgentUseCase) *Server {
	h := NewHandler(agentUseCase)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.Health)

	mux.HandleFunc("/v1/chat", apiKeyMiddleware(apiKey, h.Chat))

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

// Start inicia o servidor HTTP.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
