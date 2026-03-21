package session

import "time"

// Session armazena o histórico de uma conversa entre o usuário e o agente.
// Cada SessionID representa uma conversa independente — o agente lembra
// de tudo que aconteceu desde o início daquela sessão.
type Session struct {
	// ID é o identificador único da sessão, fornecido pelo cliente.
	ID string

	// History acumula o histórico completo da conversa:
	// inputs do usuário + resultados de tools de todas as iterações anteriores.
	// É passado inteiro ao LLM a cada nova mensagem.
	History []string

	// UpdatedAt registra a última vez que a sessão foi usada.
	// Usado para expirar sessões inativas.
	UpdatedAt time.Time
}

// Append adiciona uma entrada ao histórico e atualiza o timestamp.
func (s *Session) Append(entry string) {
	s.History = append(s.History, entry)
	s.UpdatedAt = time.Now()
}

// Store define o contrato de persistência de sessões.
// A implementação padrão é em memória; pode ser trocada por SQLite, Redis, etc.
type Store interface {
	// Get retorna a sessão pelo ID. Cria uma nova se não existir.
	Get(id string) *Session

	// Save persiste a sessão atualizada.
	Save(s *Session)

	// Delete remove a sessão do store.
	Delete(id string)
}
