package llm

import "fmt"

// plannerPrompt decide qual tool chamar nesta iteração.
// Prompt curto e direto — modelos pequenos como qwen2.5:7b
// seguem instruções muito melhor quando o prompt é simples.
func plannerPrompt(history, tools string) string {
	return fmt.Sprintf(`Você é um agente que executa ferramentas.

Ferramentas disponíveis:
%s

Histórico:
%s

Qual ferramenta executar agora? Responda APENAS com JSON:
{"tool": "nome_da_tool", "params": {}, "confidence": 0.9}

Se nenhuma ferramenta for necessária, responda:
{"tool": "none", "params": {}, "confidence": 1.0}

Responda SOMENTE com o JSON. Sem explicações.`, tools, history)
}

// finalizerPrompt gera a resposta final em linguagem natural.
// Chamado APÓS a tool ser executada, para encerrar o loop.
func finalizerPrompt(history string) string {
	return fmt.Sprintf(`Você é um assistente. Com base no histórico abaixo, escreva uma resposta
curta e direta ao usuário resumindo o que foi feito.

Histórico:
%s

Responda em linguagem natural, sem JSON. Máximo 2 frases.`, history)
}
