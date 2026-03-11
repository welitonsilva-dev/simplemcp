package server

import (
	"context"
	"net/http"
	"time"
)

// timeoutMiddleware cancela a requisição se exceder o tempo limite.
func timeoutMiddleware(timeout time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		done := make(chan struct{})

		go func() {
			next(w, r.WithContext(ctx))
			close(done)
		}()

		select {
		case <-done:
			// requisição concluída normalmente
		case <-ctx.Done():
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
		}
	}
}
