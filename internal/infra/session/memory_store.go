package session

import (
	"sync"
	"time"

	domain "humancli-server/internal/domain/session"
)

// MemoryStore é a implementação em memória do session.Store.
// Simples e eficiente para uso single-node. Sessões são expiradas
// automaticamente após o TTL para evitar crescimento ilimitado de memória.
//
// Para persistência entre reinicializações do servidor, substitua por
// uma implementação SQLite ou Redis que implemente a mesma interface.
type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]*domain.Session
	ttl      time.Duration
}

// NewMemoryStore cria um MemoryStore com o TTL informado.
// Sessões não acessadas por mais de ttl são expiradas automaticamente.
func NewMemoryStore(ttl time.Duration) *MemoryStore {
	store := &MemoryStore{
		sessions: make(map[string]*domain.Session),
		ttl:      ttl,
	}
	go store.gcLoop()
	return store
}

// Get retorna a sessão pelo ID. Cria uma sessão nova se não existir.
func (s *MemoryStore) Get(id string) *domain.Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, exists := s.sessions[id]
	if !exists {
		sess = &domain.Session{
			ID:        id,
			History:   []string{},
			UpdatedAt: time.Now(),
		}
		s.sessions[id] = sess
	}
	return sess
}

// Save persiste a sessão atualizada no store.
func (s *MemoryStore) Save(sess *domain.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[sess.ID] = sess
}

// Delete remove a sessão do store.
func (s *MemoryStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
}

// gcLoop roda em background e expira sessões inativas a cada minuto.
func (s *MemoryStore) gcLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.evictExpired()
	}
}

func (s *MemoryStore) evictExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()
	cutoff := time.Now().Add(-s.ttl)
	for id, sess := range s.sessions {
		if sess.UpdatedAt.Before(cutoff) {
			delete(s.sessions, id)
		}
	}
}
