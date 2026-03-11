package message

// UserMessage representa a entrada do usuário recebida pelo servidor.
type UserMessage struct {
	// SessionID identifica a sessão do usuário.
	// Usado para rastrear contexto entre mensagens (ex: histórico, SQLite futuramente).
	SessionID string `json:"session_id"`

	// Content é o texto enviado pelo usuário.
	Content string `json:"message"`
}

// AgentResponse é a resposta final enviada ao usuário
// após o plano ser executado.
type AgentResponse struct {
	// Results contém os resultados de cada step executado.
	Results []StepResult `json:"results"`
}

// StepResult é o resultado de um único step do plano
// formatado para o usuário final.
type StepResult struct {
	// Tool é o nome da ferramenta executada.
	Tool string `json:"tool"`

	// Output é o resultado da execução.
	Output any `json:"output"`

	// Error contém a mensagem de erro, se houver.
	Error string `json:"error,omitempty"`
}
