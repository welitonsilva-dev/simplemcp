package server

import (
	"net/http"
	"time"

	"simplemcp/internal/usecase/agent"
)

// Server encapsula o servidor HTTP.
type Server struct {
	httpServer *http.Server
}

// New cria e configura o servidor com todas as rotas.
func New(addr, apiKey string, limitIP, limitGlobal int, window time.Duration, agentUseCase *agent.AgentUseCase) *Server {
	h := NewHandler(agentUseCase)
	limiter := newRateLimiter(limitIP, limitGlobal, window)

	mux := http.NewServeMux()

	// /health — livre, sem autenticação e sem rate limit
	mux.HandleFunc("/health", rateLimitMiddleware(limiter, h.Health))

	// /v1/chat — protegido por API Key + rate limit
	mux.HandleFunc("/v1/chat", apiKeyMiddleware(apiKey,
		rateLimitMiddleware(limiter, h.Chat),
	))

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
