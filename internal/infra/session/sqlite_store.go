package session

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	domain "humancli-server/internal/domain/session"
	"humancli-server/internal/infra/logger"

	_ "modernc.org/sqlite"
)

// SQLiteStore é a implementação persistente do session.Store.
// O histórico sobrevive a reinicializações do servidor — ao contrário
// do MemoryStore, que perde tudo ao parar o processo.
//
// Usa modernc.org/sqlite (CGo-free) para funcionar em qualquer ambiente,
// incluindo containers sem gcc instalado.
type SQLiteStore struct {
	db  *sql.DB
	ttl time.Duration
}

// NewSQLiteStore abre (ou cria) o banco SQLite no caminho informado e
// inicializa o schema. Sessões inativas por mais de ttl são expiradas
// automaticamente por um GC em background.
//
// Exemplo:
//
//	store, err := session.NewSQLiteStore("data/sessions.db", 30*time.Minute)
func NewSQLiteStore(path string, ttl time.Duration) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir sqlite: %w", err)
	}

	// WAL mode: leituras e escritas simultâneas sem travar o banco
	if _, err := db.Exec(`PRAGMA journal_mode=WAL`); err != nil {
		return nil, fmt.Errorf("falha ao ativar WAL: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("falha na migração: %w", err)
	}

	store := &SQLiteStore{db: db, ttl: ttl}
	go store.gcLoop()
	return store, nil
}

// migrate cria as tabelas necessárias se ainda não existirem.
// Idempotente — seguro chamar a cada inicialização.
func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id         TEXT PRIMARY KEY,
			history    TEXT    NOT NULL DEFAULT '[]',
			updated_at INTEGER NOT NULL
		)
	`)
	return err
}

// Get retorna a sessão pelo ID.
// Se a sessão não existir, retorna uma nova (sem persistir ainda).
func (s *SQLiteStore) Get(id string) *domain.Session {
	row := s.db.QueryRow(
		`SELECT history, updated_at FROM sessions WHERE id = ?`, id,
	)

	var historyJSON string
	var updatedAtUnix int64

	err := row.Scan(&historyJSON, &updatedAtUnix)
	if err == sql.ErrNoRows {
		// sessão nova — não persiste ainda, Save() fará isso
		return &domain.Session{
			ID:        id,
			History:   []string{},
			UpdatedAt: time.Now(),
		}
	}
	if err != nil {
		logger.Error("sqlite get session '%s': %v", id, err)
		return &domain.Session{
			ID:        id,
			History:   []string{},
			UpdatedAt: time.Now(),
		}
	}

	var history []string
	if err := json.Unmarshal([]byte(historyJSON), &history); err != nil {
		logger.Error("sqlite unmarshal history '%s': %v", id, err)
		history = []string{}
	}

	return &domain.Session{
		ID:        id,
		History:   history,
		UpdatedAt: time.Unix(updatedAtUnix, 0),
	}
}

// Save persiste a sessão no banco. Usa UPSERT para criar ou atualizar.
func (s *SQLiteStore) Save(sess *domain.Session) {
	historyJSON, err := json.Marshal(sess.History)
	if err != nil {
		logger.Error("sqlite marshal history '%s': %v", sess.ID, err)
		return
	}

	_, err = s.db.Exec(`
		INSERT INTO sessions (id, history, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			history    = excluded.history,
			updated_at = excluded.updated_at
	`, sess.ID, string(historyJSON), sess.UpdatedAt.Unix())

	if err != nil {
		logger.Error("sqlite save session '%s': %v", sess.ID, err)
	}
}

// Delete remove a sessão do banco.
func (s *SQLiteStore) Delete(id string) {
	if _, err := s.db.Exec(`DELETE FROM sessions WHERE id = ?`, id); err != nil {
		logger.Error("sqlite delete session '%s': %v", id, err)
	}
}

// Close fecha a conexão com o banco. Chame ao encerrar o servidor.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// gcLoop roda em background e expira sessões inativas a cada 5 minutos.
func (s *SQLiteStore) gcLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.evictExpired()
	}
}

func (s *SQLiteStore) evictExpired() {
	cutoff := time.Now().Add(-s.ttl).Unix()
	res, err := s.db.Exec(`DELETE FROM sessions WHERE updated_at < ?`, cutoff)
	if err != nil {
		logger.Error("sqlite gc error: %v", err)
		return
	}
	if n, _ := res.RowsAffected(); n > 0 {
		logger.Info("sqlite gc: %d sessão(ões) expirada(s)", n)
	}
}
