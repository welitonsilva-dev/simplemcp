package server

import (
	"net/http"
	"sync"
	"time"
)

// rateLimiter controla requisições por IP e globalmente.
type rateLimiter struct {
	mu          sync.Mutex
	perIP       map[string]*bucket
	global      *bucket
	limitIP     int
	limitGlobal int
	window      time.Duration
}

// bucket conta requisições dentro de uma janela de tempo.
type bucket struct {
	count   int
	resetAt time.Time
}

// newRateLimiter cria um rateLimiter com os limites e janela definidos.
func newRateLimiter(limitIP, limitGlobal int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		perIP:       make(map[string]*bucket),
		global:      &bucket{resetAt: time.Now().Add(window)},
		limitIP:     limitIP,
		limitGlobal: limitGlobal,
		window:      window,
	}
}

// allow verifica se a requisição do IP pode prosseguir.
// Retorna false se o limite por IP ou global foi atingido.
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// reset global se a janela expirou
	if now.After(rl.global.resetAt) {
		rl.global.count = 0
		rl.global.resetAt = now.Add(rl.window)
	}

	// verifica limite global
	if rl.global.count >= rl.limitGlobal {
		return false
	}

	// busca ou cria bucket do IP
	b, exists := rl.perIP[ip]
	if !exists || now.After(b.resetAt) {
		rl.perIP[ip] = &bucket{count: 0, resetAt: now.Add(rl.window)}
		b = rl.perIP[ip]
	}

	// verifica limite por IP
	if b.count >= rl.limitIP {
		return false
	}

	// contabiliza requisição
	b.count++
	rl.global.count++

	return true
}

// rateLimitMiddleware aplica os limites de requisição por IP e global.
func rateLimitMiddleware(limiter *rateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if !limiter.allow(ip) {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}
