## 🔁 Fluxo de Execução MCP + Ollama

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
                │ Inicia MCP Server        │
                │ (Handler HTTP)           │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Usuário envia prompt     │
                │ via HTTP request         │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Planner cria plano       │
                │ de execução detalhado    │
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
                │ Consulta LLM (Ollama)   │
                │ para validação ou       │
                │ expansão do plano       │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Executa Tools (nativas   │
                │ / plugins) conforme plano│
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ Agent consolida          │
                │ resultados e resposta    │
                └──────────────┬───────────┘
                               │
                               ▼
                ┌──────────────────────────┐
                │ MCP Server retorna       │
                │ resposta final ao usuário│
                └──────────────────────────┘