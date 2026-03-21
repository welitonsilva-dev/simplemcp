package server

import (
	"net/http"
	"time"

	"humancli-server/internal/usecase/agent"
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

	// /health — sem autenticação, com rate limit e timeout
	mux.HandleFunc("/health",
		rateLimitMiddleware(limiter,
			timeoutMiddleware(timeout, h.Health),
		))

	// /v1/do — resposta completa em JSON (sem streaming)
	// ordem: apiKey → rateLimit → timeout → handler
	mux.HandleFunc("/v1/do", apiKeyMiddleware(apiKey,
		rateLimitMiddleware(limiter,
			timeoutMiddleware(timeout, h.Do),
		),
	))

	// /v1/stream — resposta em tempo real via SSE
	// sem timeout fixo — o loop encerra por conta própria ou por max_iterations
	mux.HandleFunc("/v1/stream", apiKeyMiddleware(apiKey,
		rateLimitMiddleware(limiter, h.Stream),
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
