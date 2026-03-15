# 🔁 Fluxo de Execução — humancli-server + Ollama

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
                    │                         │
                    └──────────────┬──────────┘
                                   ▼
                ┌──────────────────────────┐
                │ Ollama pronto            │
                │ servidor ativo           │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Inicia humancli-server   │
                │ (Handler HTTP)           │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ humancli-client envia    │
                │ prompt em linguagem      │
                │ natural via HTTP         │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Pipeline limpa e         │
                │ normaliza o input        │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ LLM (Ollama) interpreta  │
                │ e cria plano de          │
                │ execução detalhado       │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Agent recebe plano       │
                │ e inicia execução        │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Executa Tools (nativas   │
                │ ou plugins) conforme     │
                │ plano gerado             │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Agent consolida          │
                │ resultados               │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ humancli-server retorna  │
                │ resposta ao client       │
                └──────────────────────────┘
```