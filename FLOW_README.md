# 🔁 Fluxo de Execução — humancli-server

## Inicialização (Docker)

```text
┌──────────────────────────┐
│   Container inicia       │
│   (Docker start)         │
└──────────────┬───────────┘
               │
               ▼
┌──────────────────────────┐
│ Inicia Ollama            │
│ (ollama serve)           │
└──────────────┬───────────┘
               │
               ▼
┌──────────────────────────┐
│ Aguarda API responder    │
│ loop até funcionar       │
└──────────────┬───────────┘
               │
               ▼
┌──────────────────────────┐
│ Verifica/baixa modelo    │
│ (ollama list / pull)     │
└──────────────┬───────────┘
               │
               ▼
┌──────────────────────────┐
│ Inicia humancli-server   │
│ (Handler HTTP)           │
└──────────────────────────┘
```

---

## Loop ReAct — Por Requisição

O agente opera em um loop **ReAct (Reason + Act)**. A cada iteração o LLM
observa o histórico completo e decide a próxima ação de forma autônoma —
ao contrário de um dispatcher simples que gera um plano fixo e executa cegamente.

```text
humancli-client
      │
      │  POST /do { message: "cria a pasta projeto e adiciona README" }
      │
      ▼
┌─────────────────────┐
│ Pipeline            │
│ limpa e normaliza   │
│ o input             │
└──────────┬──────────┘
           │
           ▼
┌──────────────────────────────────────────────────────┐
│                   LOOP ReAct                         │
│                                                      │
│  história = ["usuário: cria a pasta..."]             │
│                                                      │
│  ┌────────────────────────────────────────────┐      │
│  │  LLM raciocina sobre o histórico completo  │      │
│  └────────────────┬───────────────────────────┘      │
│                   │                                  │
│         ┌─────────┴──────────┐                       │
│         │                    │                       │
│         ▼                    ▼                       │
│  ┌─────────────┐    ┌─────────────────────┐          │
│  │ final: true │    │ steps: [tool, params]│          │
│  │             │    │                     │          │
│  │ encerra o   │    │ executa a tool      │          │
│  │ loop        │    │                     │          │
│  └──────┬──────┘    └──────────┬──────────┘          │
│         │                     │                      │
│         │           resultado entra no histórico     │
│         │           história = [..., "tool X → Y"]  │
│         │           próxima iteração                 │
│         │                                            │
└─────────┼────────────────────────────────────────────┘
          │
          ▼
┌─────────────────────┐
│ AgentResponse       │
│ results + message   │
└─────────────────────┘
          │
          ▼
   humancli-client
```

### Exemplo de execução em 3 iterações

```
Iteração 1:
  LLM recebe: "usuário: cria a pasta projeto e adiciona README"
  LLM decide: chamar fs_mkdir { path: "projeto" }
  Resultado: "pasta criada"

Iteração 2:
  LLM recebe: histórico + "tool fs_mkdir retornou: pasta criada"
  LLM decide: chamar fs_touch { path: "projeto/README.md" }
  Resultado: "arquivo criado"

Iteração 3:
  LLM recebe: histórico + "tool fs_touch retornou: arquivo criado"
  LLM decide: final=true, message="Pasta 'projeto' criada com README.md."
  Loop encerra.
```

---

## Proteções do Loop

| Proteção | Comportamento |
|---|---|
| **Confidence guard** | Tools destrutivas (`fs_rm`, `fs_rmdir`) são bloqueadas se confidence < threshold |
| **Max iterations** | Loop encerra após `AGENT_MAX_ITERATIONS` ciclos (padrão: 10) |
| **Tool inexistente** | Erro entra no histórico, LLM decide o próximo passo |
| **Tool com falha** | Erro entra no histórico, LLM pode tentar alternativa |
