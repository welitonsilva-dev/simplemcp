package message

// UserMessage representa a entrada do usuário recebida pelo servidor.
type UserMessage struct {
	// SessionID identifica a conversa. O agente carrega o histórico desta sessão
	// antes de processar — permitindo contexto entre múltiplas mensagens.
	// Se vazio, a requisição é tratada como sessão anônima sem persistência.
	SessionID string `json:"session_id"`

	// Content é o texto enviado pelo usuário.
	Content string `json:"message"`
}

// AgentResponse é a resposta final enviada ao usuário após o loop encerrar.
// Usado na rota não-streaming /v1/do.
type AgentResponse struct {
	// Results contém os resultados de cada step executado durante o loop.
	Results []StepResult `json:"results"`

	// FinalMessage é a resposta em linguagem natural gerada pelo agente.
	FinalMessage string `json:"final_message,omitempty"`
}

// StepResult é o resultado de um único step executado durante o loop do agente.
type StepResult struct {
	Tool   string `json:"tool"`
	Output any    `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

// StreamEvent é um evento enviado ao cliente via SSE durante o loop.
// Cada iteração emite um evento — o cliente vê o progresso em tempo real.
type StreamEvent struct {
	// Type classifica o evento:
	//   "step"    — uma tool foi executada (durante o loop)
	//   "final"   — o agente encerrou (última mensagem do stream)
	//   "error"   — erro fatal que encerrou o loop
	Type string `json:"type"`

	// Tool é o nome da ferramenta executada. Preenchido quando Type="step".
	Tool string `json:"tool,omitempty"`

	// Output é o resultado da tool. Preenchido quando Type="step".
	Output any `json:"output,omitempty"`

	// Error contém a mensagem de erro. Preenchido quando Type="step" com falha ou Type="error".
	Error string `json:"error,omitempty"`

	// Message é a resposta final em linguagem natural. Preenchido quando Type="final".
	Message string `json:"message,omitempty"`

	// Iteration indica em qual iteração do loop este evento foi gerado.
	Iteration int `json:"iteration,omitempty"`
}
