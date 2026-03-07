# 📦 Fluxo de Inicialização do Ollama no Container

## 💡 Resumo do que o script `.sh` faz

- Inicia o servidor do Ollama
- Espera até a API estar disponível (`localhost:11434`)
- Define qual modelo usar (via variável de ambiente ou padrão)
- Verifica se o modelo já está instalado
- Baixa o modelo se necessário
- Mantém o servidor rodando dentro do container

---

## 🔁 Fluxo de Execução

```text
                ┌──────────────────────────┐
                │   Container inicia       │
                │   (Docker start)         │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Inicia servidor do       │
                │ Ollama                   │
                │ (ollama serve)           │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Espera API responder     │
                │ curl localhost:11434     │
                │ loop até funcionar       │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Define modelo            │
                │ ENV LLM_MODEL            │
                │ ou modelo padrão         │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Verifica se modelo       │
                │ já está instalado        │
                │ (ollama list)            │
                └──────────────┬───────────┘
                               │
                     ┌─────────┴─────────┐
                     │                   │
                     ▼                   ▼
        ┌──────────────────────┐  ┌──────────────────────┐
        │ Modelo já existe     │  │ Modelo não existe    │
        │                      │  │                      │
        │ continua execução    │  │ ollama pull MODEL    │
        └───────────┬──────────┘  └───────────┬──────────┘
                    │                        │
                    └──────────────┬─────────┘
                                   ▼
                ┌──────────────────────────┐
                │ Ollama pronto            │
                │ servidor ativo           │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Mantém container vivo    │
                │ (wait)                   │
                └──────────────────────────┘