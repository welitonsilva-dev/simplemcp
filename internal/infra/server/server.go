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
func New(addr, apiKey string, limitIP, limitGlobal int, window, timeout time.Duration, agentUseCase *agent.AgentUseCase) *Server {
	h := NewHandler(agentUseCase)
	limiter := newRateLimiter(limitIP, limitGlobal, window)

	mux := http.NewServeMux()

	// /health — sem autenticação, mas com rate limit e timeout para evitar abusos
	// ordem: rateLimit → timeout → handler
	mux.HandleFunc("/health",
		rateLimitMiddleware(limiter,
			timeoutMiddleware(timeout, h.Health),
		))

	// /v1/do — pipeline completo de middlewares
	// ordem: apiKey → rateLimit → timeout → handler
	mux.HandleFunc("/v1/do", apiKeyMiddleware(apiKey,
		rateLimitMiddleware(limiter,
			timeoutMiddleware(timeout, h.Do),
		),
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
