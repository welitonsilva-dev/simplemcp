package message

// UserMessage representa a entrada do usuário recebida pelo servidor.
type UserMessage struct {
	// SessionID identifica a sessão do usuário.
	// Usado para rastrear contexto entre mensagens (ex: histórico, SQLite futuramente).
	SessionID string `json:"session_id"`

	// Content é o texto enviado pelo usuário.
	Content string `json:"message"`
}

// AgentResponse é a resposta final enviada ao usuário após o loop do agente encerrar.
type AgentResponse struct {
	// Results contém os resultados de cada step executado durante o loop.
	Results []StepResult `json:"results"`

	// FinalMessage é a resposta em linguagem natural gerada pelo agente
	// ao encerrar o loop. Pode estar vazia se o agente encerrou por limite de iterações.
	FinalMessage string `json:"final_message,omitempty"`
}

// StepResult é o resultado de um único step executado durante o loop do agente.
type StepResult struct {
	// Tool é o nome da ferramenta executada.
	Tool string `json:"tool"`

	// Output é o resultado da execução.
	Output any `json:"output,omitempty"`

	// Error contém a mensagem de erro, se houver.
	Error string `json:"error,omitempty"`
}
