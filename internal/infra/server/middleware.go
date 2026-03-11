package server

import (
	"net/http"
)

// apiKeyMiddleware protege uma rota exigindo o header X-API-Key.
// Rotas sem esse middleware (ex: /health) ficam livres.
func apiKeyMiddleware(apiKey string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != apiKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
