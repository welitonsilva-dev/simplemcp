# humancli-server

> Núcleo do ecossistema HumanCLI — um servidor de agente que interpreta linguagem natural e executa ferramentas de forma autônoma, com suporte a LLM local e provedores de IA externos.
>
> *Fale o que quer. O agente entende, age, observa o resultado e decide o próximo passo.*

---

## Sumário

- [O que é o humancli-server?](#o-que-é-o-humancli-server)
- [Ecossistema HumanCLI](#ecossistema-humancli)
- [Como funciona — Loop ReAct](#como-funciona--loop-react)
- [Rotas da API](#rotas-da-api)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Tools Nativas](#tools-nativas)
- [Sistema de Plugins](#sistema-de-plugins)
- [Provedores de IA](#provedores-de-ia)
- [Sessões e Contexto](#sessões-e-contexto)
- [Pipeline de Entrada](#pipeline-de-entrada)
- [Segurança](#segurança)
- [Configuração — Variáveis de Ambiente](#configuração--variáveis-de-ambiente)
- [Setup e Execução](#setup-e-execução)
- [Tecnologias](#tecnologias)
- [Diferencial](#diferencial)
- [Roadmap](#roadmap)
- [Contribuindo](#contribuindo)

---

## O que é o humancli-server?

O **humancli-server** é o núcleo do ecossistema HumanCLI: um servidor HTTP escrito em Go que recebe mensagens em linguagem natural, raciocina sobre elas usando um LLM (local ou externo) e executa ferramentas (tools) de forma autônoma até concluir a tarefa solicitada.

Você escreve o que quer. O agente entende, age, observa o resultado e decide o próximo passo — repetindo o ciclo até concluir a tarefa.

> *"Cria a pasta projeto e adiciona um README"* → o agente chama `fs_mkdir`, depois `fs_touch`, confirma os resultados e responde em linguagem natural.

Com o sistema de plugins, o servidor pode ser expandido para qualquer coisa: automações, integrações com serviços externos, comandos personalizados, rotinas de trabalho.

---

## Ecossistema HumanCLI

| Projeto | Papel | Status |
|---|---|---|
| **humancli-server** | Servidor de agente. Hospeda o LLM, registra as tools nativas e plugins externos, e executa o loop ReAct. | ✅ Este repositório |
| **humancli-plugins** | Repositório de plugins externos. Cada subpasta é uma tool independente que se integra ao servidor via SDK. | ✅ [Repositório separado](https://github.com/welitonsilva-dev/humancli-plugins) |
| **humancli-client** | Interface de linha de comando. Ponto de entrada onde o usuário digita comandos em linguagem natural e vê o agente agir em tempo real via SSE. | ✅ [Repositório separado](https://github.com/welitonsilva-dev/humancli-client) |

---

## Como funciona — Loop ReAct

O agente opera em um loop **ReAct (Reason + Act)**. A cada iteração, o LLM observa o histórico completo da conversa e decide autonomamente a próxima ação — ao contrário de um dispatcher simples que gera um plano fixo e o executa cegamente.

```
usuário → pipeline → [LLM raciocina → executa tool → observa resultado] → resposta final
                      └──────────────── loop até final=true ────────────────┘
```

### Detalhe interno por iteração

O loop usa **dois prompts separados** para lidar com modelos pequenos (como qwen2.5:7b):

1. **plannerPrompt** — pergunta ao LLM qual tool executar agora. Resposta esperada: `{"tool": "nome", "params": {}, "confidence": 0.9}`. Formato simples e direto para que modelos compactos sigam com consistência.
2. **finalizerPrompt** — chamado apenas ao encerrar o loop. Pede ao LLM que resuma em linguagem natural o que foi feito, com base no histórico completo.

### Proteções do loop

| Proteção | Comportamento |
|---|---|
| **Confidence guard** | Tools destrutivas (`fs_rm`, `fs_rmdir`, `fs_rmrf`) são bloqueadas quando `confidence` do LLM está abaixo de `CONFIDENCE_THRESHOLD`. O agente pede ao usuário que seja mais específico. |
| **Max iterations** | O loop encerra após `AGENT_MAX_ITERATIONS` ciclos (padrão: 10). Evita loops infinitos mesmo se o LLM não convergir. |
| **Deduplicação** | Se o LLM tentar chamar a mesma tool com os mesmos parâmetros pela segunda vez, o loop encerra automaticamente para evitar repetições. |
| **Tool inexistente** | O erro é registrado no histórico e o LLM decide o próximo passo (pode tentar uma alternativa ou encerrar). |
| **Tool com falha** | O erro também entra no histórico. O LLM pode tentar uma abordagem diferente antes de encerrar. |

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
  LLM decide: tool="none" → encerrar
  finalizerPrompt → "Pasta 'projeto' criada com README.md."
  Loop encerra.
```

---

## Rotas da API

### `POST /v1/do` — Resposta consolidada (JSON)

Executa o loop ReAct completo e retorna a resposta final quando o loop encerrar. Indicado para clientes que não suportam SSE ou que preferem aguardar a resposta completa.

**Headers obrigatórios:**
```
Content-Type: application/json
X-API-Key: <sua-api-key>
```

**Body:**
```json
{
  "session_id": "cli-a3f8b21c",
  "message": "cria a pasta projeto e adiciona um README"
}
```

**Resposta `200`:**
```json
{
  "results": [
    { "tool": "fs_mkdir", "output": "Diretório criado: /app/host/projeto", "error": "" },
    { "tool": "fs_touch", "output": "Arquivo criado/atualizado: /app/host/projeto/README.md", "error": "" }
  ],
  "final_message": "Pasta 'projeto' criada com README.md dentro."
}
```

---

### `POST /v1/stream` — Streaming em tempo real (SSE)

Executa o loop ReAct e emite um evento SSE por iteração — sem esperar o loop terminar. Indicado para interfaces interativas como o humancli-client, onde o usuário vê o agente agindo em tempo real.

**Headers obrigatórios:**
```
Content-Type: application/json
X-API-Key: <sua-api-key>
Accept: text/event-stream
```

**Body:** mesmo formato de `/v1/do`

**Formato dos eventos SSE:**
```
data: {"type":"step","tool":"fs_mkdir","output":"Diretório criado: /app/host/projeto","iteration":1}

data: {"type":"step","tool":"fs_touch","output":"Arquivo criado: /app/host/projeto/README.md","iteration":2}

data: {"type":"final","message":"Pasta 'projeto' criada com README.md.","iteration":3}
```

**Tipos de evento:**

| Tipo | Campos | Descrição |
|---|---|---|
| `step` | `tool`, `output`, `error`, `iteration` | Resultado de uma iteração do loop |
| `final` | `message`, `iteration` | Resposta final em linguagem natural |
| `error` | `error` | Erro fatal durante o processamento |

> **Nota sobre timeout:** `/v1/stream` não possui timeout fixo — o loop encerra por conta própria via `final=true` ou ao atingir `AGENT_MAX_ITERATIONS`. Já `/v1/do` respeita o `REQUEST_TIMEOUT` configurado.

---

### `GET /health` — Status do servidor

Rota pública (sem autenticação) que retorna o status do servidor. Usada pelo humancli-client no comando `humancli health`.

**Resposta `200`:**
```json
{ "status": "ok" }
```

---

## Estrutura do Projeto

```
humancli-server/
├── cmd/server/
│   └── main.go                   — entrypoint: wiring de todas as dependências
├── internal/
│   ├── adapter/
│   │   ├── llm/
│   │   │   ├── Client.go         — adaptador HTTP para o Ollama (/api/generate)
│   │   │   ├── parser.go         — parse do JSON retornado pelo LLM (plan)
│   │   │   └── prompt.go         — plannerPrompt e finalizerPrompt
│   │   ├── pipeline/
│   │   │   ├── pipeline.go       — orquestrador dos 5 passos de pré-processamento
│   │   │   ├── validator.go      — rejeita entradas inválidas (tamanho, vazio)
│   │   │   ├── sanitizer.go      — bloqueia intenções perigosas no input
│   │   │   ├── cleaner.go        — remove ruído (espaços, caracteres especiais)
│   │   │   ├── normalize.go      — padroniza formato do texto
│   │   │   └── optimizer.go      — melhora clareza do input para o LLM
│   │   └── tools/
│   │       ├── registry.go       — singleton Registry com detecção de origem (native/plugin)
│   │       └── native/
│   │           ├── state.go      — CwdState compartilhado + ResolvePath + ToHostPath
│   │           ├── tool_list.go  — tool_list: lista todas as tools registradas
│   │           ├── echo/
│   │           │   ├── echo.go          — echo: repete mensagem
│   │           │   └── double_echo.go   — double_echo: duplica mensagem
│   │           └── filesystem/
│   │               ├── cd.go     — fs_cd: muda diretório de trabalho
│   │               ├── list.go   — fs_list: lista arquivos e diretórios
│   │               ├── mkdir.go  — fs_mkdir: cria diretório
│   │               ├── touch.go  — fs_touch: cria arquivo ou atualiza timestamp
│   │               ├── mr.go     — fs_rm: remove arquivo (com confirmação)
│   │               └── rmdir.go  — fs_rmdir: remove diretório vazio (com confirmação)
│   ├── domain/
│   │   ├── message/
│   │   │   └── message.go        — tipos UserMessage, AgentResponse, StreamEvent, StepResult
│   │   ├── plan/
│   │   │   └── plan.go           — tipo Plan (resposta do plannerPrompt) com IsFinal/IsUnknown
│   │   ├── session/
│   │   │   └── session.go        — tipo Session (ID, History, UpdatedAt) e interface Store
│   │   └── tool/
│   │       ├── tool.go           — interface Tool (Name, Description, Execute)
│   │       └── registry.go       — interface ToolRegistry (Get, All)
│   ├── infra/
│   │   ├── config/
│   │   │   └── config.go         — Config struct e Load() via variáveis de ambiente
│   │   ├── logger/
│   │   │   └── logger.go         — logger com arquivo rotacionado em LOG_DIR
│   │   ├── server/
│   │   │   ├── server.go         — HTTP server com mux e rotas
│   │   │   ├── handler.go        — handlers Do (JSON) e Stream (SSE) e Health
│   │   │   ├── middleware.go     — apiKeyMiddleware (header X-API-Key)
│   │   │   ├── ratelimit.go      — rate limiter por IP e global
│   │   │   └── timeout.go        — timeoutMiddleware para /v1/do
│   │   └── session/
│   │       ├── memory_store.go   — armazenamento em memória com GC automático por TTL
│   │       └── sqlite_store.go   — armazenamento persistente SQLite (WAL mode, UPSERT)
│   └── usecase/
│       └── agent/
│           └── agent.go          — AgentUseCase: loop ReAct, confidence guard, deduplicação
├── sdk/
│   └── sdk.go                    — SDK público exposto para plugins (interface Tool + Register)
├── scripts/
│   ├── genplugins/main.go        — gerador: escaneia humancli-plugins e gera imports
│   └── gentools/main.go          — gerador: escaneia tools nativas e gera imports
├── generate.go                   — diretivas go:generate para os dois scripts
├── docker/
│   ├── ollama/
│   │   └── entrypoint.sh         — inicia Ollama, aguarda API, verifica/baixa modelo, inicia servidor
│   └── docker-clean.sh           — remove containers, volumes e imagens do projeto
├── docker-compose.yml            — orquestra containers ollama + humancli-server
├── Dockerfile                    — build multi-stage do humancli-server
├── Dockerfile.standalone         — build sem dependência do Ollama (para provedores externos)
├── exemple.env                   — template de variáveis de ambiente com documentação inline
├── cli.sh                        — wrapper shell que injeta HOST_ROOT e CONTAINER_CWD
└── cli.ps1                       — wrapper PowerShell equivalente para Windows
```

---

## Tools Nativas

Tools registradas automaticamente quando o servidor inicia. Todas residem em `internal/adapter/tools/native/` e utilizam o `CwdState` compartilhado para resolver caminhos relativos.

### `tool_list`

Lista todas as tools registradas, separadas por nativas e plugins.

| Parâmetro | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| — | — | — | Não aceita parâmetros |

**Resposta:**
```json
{
  "message": "11 ferramentas disponíveis (9 nativas, 2 plugins)",
  "native": ["tool_list", "fs_cd", "fs_list", "fs_mkdir", "fs_touch", "fs_rm", "fs_rmdir", "echo", "double_echo"],
  "plugins": ["hello", "docker_ps"]
}
```

---

### `fs_cd` — Mudar diretório

Muda o diretório de trabalho compartilhado entre todas as tools (`CwdState`). O novo CWD é mantido em memória e usado como base para todos os caminhos relativos nas chamadas seguintes.

| Parâmetro | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `path` | string | ✅ | Caminho absoluto, relativo ao CWD atual, ou com `~` para home |

Suporta `~` (home do usuário), caminhos relativos (`../pasta`) e absolutos. Compatível com Windows, Linux e macOS.

---

### `fs_list` — Listar arquivos

Lista arquivos e diretórios, incluindo arquivos ocultos. Usa `ls -a` no Linux/macOS e `dir /a` no Windows. O output inclui o caminho legível do host (sem o prefixo `/app/host`).

| Parâmetro | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `path` | string | ❌ | Diretório a listar. Se omitido, usa o CWD atual |

**Resposta:**
```json
{
  "message": "Encontrei 4 itens em /home/user/projeto",
  "items": ["README.md", "main.go", "go.mod", ".gitignore"]
}
```

---

### `fs_mkdir` — Criar diretório

Cria um diretório no sistema de arquivos. Com `parents=true` comporta-se como `mkdir -p`, criando todos os diretórios intermediários necessários.

| Parâmetro | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `path` | string | ✅ | Caminho do diretório a criar (absoluto ou relativo ao CWD) |
| `parents` | bool | ❌ | Se `true`, cria diretórios intermediários. Padrão: `false` |

---

### `fs_touch` — Criar arquivo ou atualizar timestamp

Cria um arquivo vazio se não existir, ou atualiza os timestamps (`atime` e `mtime`) de um arquivo existente — comportamento idêntico ao comando `touch` do Unix.

| Parâmetro | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `path` | string | ✅ | Caminho do arquivo (absoluto ou relativo ao CWD) |

---

### `fs_rm` — Remover arquivo

Remove um arquivo do sistema. Esta é uma operação **destrutiva** e exige confirmação dupla:

1. O LLM deve enviar `confirmed: true` somente se o usuário tiver expressado consentimento explícito.
2. O confidence guard do agente também atua: se `confidence < CONFIDENCE_THRESHOLD`, a execução é bloqueada.

| Parâmetro | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `path` | string | ✅ | Caminho do arquivo a remover |
| `confirmed` | bool | ✅ | Deve ser `true` para executar. Sem ele, retorna solicitação de confirmação |

> Só remove arquivos. Para diretórios, use `fs_rmdir` (vazio) ou `fs_rmrf` (com conteúdo).

---

### `fs_rmdir` — Remover diretório vazio

Remove um diretório vazio. Também é uma operação **destrutiva** com o mesmo mecanismo de confirmação dupla do `fs_rm`.

| Parâmetro | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `path` | string | ✅ | Caminho do diretório a remover |
| `confirmed` | bool | ✅ | Deve ser `true` para executar |

> Falha se o diretório contiver arquivos. Use `fs_rmrf` para remover com conteúdo.

---

### `echo` — Repetir mensagem

Repete exatamente o texto recebido no parâmetro `message`. Útil para testar o fluxo de execução de tools e validar a comunicação entre componentes.

| Parâmetro | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `message` | string | ✅ | Texto a ser repetido |

---

### `double_echo` — Duplicar mensagem

Duplica o texto recebido, concatenando-o consigo mesmo separado por vírgula. Útil para testar e depurar o fluxo do agente.

| Parâmetro | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `message` | string | ✅ | Texto a ser duplicado |

---

## Sistema de Plugins

O humancli-server é extensível por design. Plugins são módulos Go independentes que implementam a interface `sdk.Tool` e se registram automaticamente no servidor sem necessidade de alterar o código principal.

### SDK público

```go
// sdk/sdk.go
type Tool interface {
    Name()        string
    Description() string
    Execute(params map[string]interface{}) (any, error)
}

func Register(t Tool) {
    tools.GlobalRegistry().Register(t)
}
```

| Método | Retorno | Responsabilidade |
|---|---|---|
| `Name()` | string | Identificador único da tool. Usado pelo LLM para chamá-la. |
| `Description()` | string | Instruções para o LLM: quando acionar, parâmetros aceitos, comportamento esperado. |
| `Execute()` | (any, error) | Lógica da tool. Recebe params como `map[string]interface{}` e retorna resultado ou erro. |

### Como os plugins são carregados

O gerador `scripts/genplugins/main.go` é invocado via `go generate ./...`. Ele escaneia todas as subpastas de `humancli-plugins/` e gera automaticamente um arquivo `cmd/server/plugins.go` com os imports necessários:

```go
// Code generated by go generate. DO NOT EDIT.
package main

import (
    _ "github.com/weliton/humancli-plugins/hello"
    _ "github.com/weliton/humancli-plugins/meu_plugin"
)
```

O `init()` de cada plugin é executado na inicialização do binário, registrando a tool no `GlobalRegistry`.

### Registry — Detecção de origem

O `Registry` detecta automaticamente se uma tool é nativa ou plugin inspecionando o `PkgPath` via `reflect`:

```go
pkgPath := reflect.TypeOf(t).Elem().PkgPath()
origin := OriginPlugin
if strings.Contains(pkgPath, "tools/native") {
    origin = OriginNative
}
```

Isso permite que `tool_list` separe as tools por categoria na resposta.

### Criando um plugin

Consulte o [PLUGIN_README.md](./PLUGIN_README.md) para o guia completo passo a passo, ou o repositório [humancli-plugins](https://github.com/welitonsilva-dev/humancli-plugins) para exemplos funcionais.

---

## Provedores de IA

O HumanCLI foi projetado para ser **agnóstico ao provedor de IA**. A variável `HUMANCLI_PROVIDER` define qual LLM será usado — sem alterar uma linha de código ou qualquer plugin existente.

```env
# Provedores implementados
HUMANCLI_PROVIDER=ollama      # LLM local via Ollama ou LM Studio
HUMANCLI_PROVIDER=groq        # Llama 3.3, Mixtral (tier gratuito disponível)

# Provedores planejados (não implementados ainda)
HUMANCLI_PROVIDER=openai      # GPT-4o, GPT-4o-mini, GPT-3.5-turbo
HUMANCLI_PROVIDER=anthropic   # Claude 3.5 Haiku, Claude 3.5 Sonnet
HUMANCLI_PROVIDER=openrouter  # Acesso a dezenas de modelos via uma só API
```

### Configuração por provedor

```env
# Ollama (padrão — gratuito, privado, sem internet)
HUMANCLI_PROVIDER=ollama
HUMANCLI_MODEL=qwen2.5:7b     # Outros: llama3.2, mistral, phi3
OLLAMA_URL=http://ollama:11434

# OpenAI (pago)
HUMANCLI_PROVIDER=openai
HUMANCLI_MODEL=gpt-4o-mini
LLM_API_KEY=sk-...

# Anthropic (pago)
HUMANCLI_PROVIDER=anthropic
HUMANCLI_MODEL=claude-3-5-haiku-20241022
LLM_API_KEY=sk-ant-...

# Groq (tier gratuito)
HUMANCLI_PROVIDER=groq
HUMANCLI_MODEL=llama-3.3-70b-versatile
LLM_API_KEY=gsk_...

# OpenRouter (multi-modelo)
HUMANCLI_PROVIDER=openrouter
HUMANCLI_MODEL=meta-llama/llama-3.1-8b-instruct:free
LLM_API_KEY=sk-or-...

# URL customizada (LM Studio, Ollama com proxy, API OpenAI-compatível)
LLM_BASE_URL=http://localhost:1234/v1
```

---

## Sessões e Contexto

O servidor mantém o histórico de cada conversa por `session_id`. Isso permite que o agente lembre de ações anteriores dentro de uma mesma sessão — por exemplo, saber que já criou uma pasta e não tentar criá-la de novo.

### Stores disponíveis

#### MemoryStore (padrão)

Armazenamento em memória. Rápido, sem dependências. Sessões são perdidas quando o servidor reinicia. Um goroutine de GC roda a cada minuto e expira sessões inativas há mais de `SESSION_TTL_MINUTES`.

```env
# MemoryStore é o padrão quando SESSION_DB_PATH está vazio
SESSION_DB_PATH=
SESSION_TTL_MINUTES=30
```

#### SQLiteStore (persistente)

Armazenamento em SQLite usando `modernc.org/sqlite` (CGo-free — funciona em containers sem gcc). O histórico das conversas sobrevive a reinicializações do servidor. O banco usa WAL mode para suportar leituras e escritas simultâneas sem travamentos.

```env
SESSION_DB_PATH=data/sessions.db
SESSION_TTL_MINUTES=60
```

O schema é criado automaticamente na primeira execução:

```sql
CREATE TABLE IF NOT EXISTS sessions (
    id         TEXT PRIMARY KEY,
    history    TEXT    NOT NULL DEFAULT '[]',
    updated_at INTEGER NOT NULL
)
```

O GC do SQLiteStore roda a cada 5 minutos e remove registros com `updated_at` anterior ao TTL.

---

## Pipeline de Entrada

Antes de enviar o input do usuário ao LLM, o texto passa por um pipeline de 5 etapas em sequência:

| Etapa | Responsável | O que faz |
|---|---|---|
| 1. `validate` | `validator.go` | Rejeita entradas vazias ou acima de `INPUT_MAX_LENGTH` (padrão: 1000 caracteres) |
| 2. `sanitize` | `sanitizer.go` | Bloqueia intenções potencialmente perigosas no input bruto |
| 3. `clean` | `cleaner.go` | Remove espaços em excesso, caracteres especiais e ruído |
| 4. `normalize` | `normalize.go` | Padroniza capitalização, pontuação e formato |
| 5. `optimize` | `optimizer.go` | Melhora a clareza do texto para maximizar a qualidade da resposta do LLM |

---

## Segurança

### Autenticação por API Key

Todas as rotas de agente (`/v1/do` e `/v1/stream`) exigem o header `X-API-Key` com o valor de `API_KEY` definido no `.env`. A rota `/health` é pública. A API Key é obrigatória — o servidor rejeita iniciar sem ela.

```bash
# Requisição autenticada
curl -X POST http://localhost:8081/v1/do \
  -H "Content-Type: application/json" \
  -H "X-API-Key: supersecretkey" \
  -d '{"session_id": "cli-abc123", "message": "lista os arquivos"}'
```

### Rate Limiting

O servidor implementa rate limiting em duas camadas configuráveis via `.env`:

| Variável | Padrão | Descrição |
|---|---|---|
| `RATE_LIMIT_PER_IP` | `10` | Máximo de requisições por IP na janela de tempo |
| `RATE_LIMIT_GLOBAL` | `50` | Máximo de requisições globais na janela de tempo |
| `RATE_LIMIT_WINDOW` | `60` | Janela de tempo em segundos |

### Timeout de requisição

A rota `/v1/do` possui timeout configurável via `REQUEST_TIMEOUT` (padrão: 180 segundos). Requisições que excedem esse tempo são encerradas automaticamente. A rota `/v1/stream` não possui timeout fixo, pois o loop SSE encerra por conta própria.

### Confidence Guard

Tools destrutivas (`fs_rm`, `fs_rmdir`, `fs_rmrf`) são bloqueadas quando o LLM retorna `confidence` abaixo de `CONFIDENCE_THRESHOLD` (padrão: 0.8 = 80%). O agente responde ao usuário pedindo mais clareza antes de executar a ação.

### Confirmação dupla em tools destrutivas

Além do confidence guard, as tools destrutivas implementam uma verificação interna do parâmetro `confirmed`. Se o LLM não incluir `confirmed: true` no payload, a tool retorna uma solicitação de confirmação explícita sem executar a operação.

### Isolamento via Docker

O disco do host é mapeado como volume somente de referência dentro do container em `/app/host`. O servidor opera dentro desse caminho, sem acesso direto ao sistema de arquivos do host fora do ponto de montagem configurado.

---

## Configuração — Variáveis de Ambiente

Copie `exemple.env` para `.env` e ajuste as variáveis:

```bash
cp exemple.env .env
```

| Variável | Padrão | Descrição |
|---|---|---|
| `API_KEY` | — | **Obrigatório.** Chave de autenticação do servidor |
| `SERVER_ADDR` | `:8081` | Endereço e porta do servidor HTTP |
| `HUMANCLI_PROVIDER` | `ollama` | Provedor de LLM (`ollama`, `openai`, `anthropic`, `groq`, `openrouter`) |
| `HUMANCLI_MODEL` | `qwen2.5:7b` | Modelo a usar no provedor escolhido |
| `LLM_API_KEY` | — | API Key para provedores externos. Vazio para Ollama |
| `LLM_BASE_URL` | — | URL base customizada (LM Studio, proxy, etc.) |
| `OLLAMA_URL` | `http://ollama:11434` | URL do servidor Ollama (usado apenas com `PROVIDER=ollama`) |
| `INPUT_MAX_LENGTH` | `1000` | Máximo de caracteres aceitos por requisição |
| `RATE_LIMIT_PER_IP` | `10` | Requisições por IP por janela de tempo |
| `RATE_LIMIT_GLOBAL` | `50` | Requisições globais por janela de tempo |
| `RATE_LIMIT_WINDOW` | `60` | Janela de tempo em segundos para rate limiting |
| `REQUEST_TIMEOUT` | `180` | Timeout em segundos para requisições em `/v1/do` |
| `CONFIDENCE_THRESHOLD` | `0.8` | Limiar mínimo de confidence para tools destrutivas (0.0 a 1.0) |
| `AGENT_MAX_ITERATIONS` | `10` | Máximo de iterações do loop ReAct por requisição |
| `SESSION_TTL_MINUTES` | `30` | Tempo em minutos para expirar sessões inativas |
| `SESSION_DB_PATH` | — | Caminho para arquivo SQLite de sessões. Vazio = memória |
| `HOST_ROOT` | `/` | Raiz do disco do host mapeado no container |
| `HOST_CWD` | `/` | Diretório de trabalho inicial do host |
| `LOG_DIR` | `/app/logs` | Diretório onde os logs são gravados |

---

## Setup e Execução

### Pré-requisitos

- Docker e Docker Compose
- Go 1.21+ (apenas para desenvolvimento ou para adicionar plugins)
- Ambos os repositórios clonados lado a lado:

```
projetos/
├── humancli-server/
└── humancli-plugins/
```

### 1. Configure o ambiente

```bash
cp exemple.env .env
# Edite .env: defina API_KEY e ajuste o modelo se necessário
```

### 2. Registre tools e plugins

```bash
go generate ./...
```

Este comando executa dois geradores:
- `scripts/genplugins/main.go` — escaneia `humancli-plugins/` e gera os imports dos plugins
- `scripts/gentools/main.go` — escaneia as tools nativas e gera os imports

### 3. Limpe o ambiente Docker (se necessário)

```bash
cd docker && sh docker-clean.sh && cd ..
```

### 4. Suba os containers

```bash
docker compose up
```

Na primeira execução, o Ollama baixará o modelo configurado em `HUMANCLI_MODEL`. Aguarde o download completo. Quando concluído, o servidor estará acessível em `http://localhost:8081`.

### 5. Verifique os containers

```bash
docker ps
# Se os containers não iniciaram na primeira execução:
docker compose up -d
```

### Desenvolvimento local (sem Docker)

```bash
# Compile e execute diretamente
go build ./...
./humancli-server
```

Certifique-se de que um servidor Ollama esteja rodando localmente e que `OLLAMA_URL` aponte para ele.

### Interagir via curl

```bash
# Verificar saúde do servidor
curl http://localhost:8081/health

# Enviar comando ao agente (resposta consolidada)
curl -X POST http://localhost:8081/v1/do \
  -H "Content-Type: application/json" \
  -H "X-API-Key: supersecretkey" \
  -d '{"session_id": "teste-1", "message": "lista os arquivos aqui"}'

# Streaming em tempo real
curl -X POST http://localhost:8081/v1/stream \
  -H "Content-Type: application/json" \
  -H "X-API-Key: supersecretkey" \
  -H "Accept: text/event-stream" \
  -d '{"session_id": "teste-1", "message": "cria a pasta demo"}'

# Interagir diretamente com o Ollama
curl -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{"model": "qwen2.5:7b", "prompt": "diga olá", "stream": false}'
```

---

## Tecnologias

| Componente | Tecnologia |
|---|---|
| Linguagem | Go 1.21+ |
| LLM local | Ollama (qwen2.5:7b, llama3.2, mistral, phi3, entre outros) |
| Persistência de sessões | SQLite via `modernc.org/sqlite` (CGo-free) |
| Transporte | HTTP + SSE (Server-Sent Events) |
| Containerização | Docker + Docker Compose |
| Build | `go generate` + multi-stage Dockerfile |

---

## Diferencial

| Característica | HumanCLI | Alternativas (Aider, Shell AI, etc.) |
|---|---|---|
| Loop ReAct — age e observa resultados | ✅ | ❌ (maioria executa plano fixo) |
| Local e offline por padrão | ✅ | ❌ (maioria usa APIs externas) |
| Sistema de plugins extensível via SDK | ✅ | ❌ |
| SSE — feedback em tempo real por iteração | ✅ | ❌ |
| Sessões persistentes com contexto multi-turno | ✅ | ❌ |
| Confidence guard + confirmação em tools destrutivas | ✅ | ❌ |
| Agnóstico ao provedor de IA | ✅ | Parcial |
| Desenvolvido em Go (performance, binário único) | ✅ | ❌ (maioria em Python) |
| Arquitetura cliente/servidor desacoplada | ✅ | ❌ |

---

## Roadmap

### Novas categorias de tools (curto prazo)

- **Git** — `git_init`, `git_commit`, `git_checkout`, `git_push`, `git_status`, `git_log`
- **Docker** — `docker_build`, `docker_run`, `docker_ps`, `docker_stop`, `docker_logs`
- **Sistema/Processos** — matar processos, verificar portas, uso de memória e CPU
- **Rede/HTTP** — fazer requisições, verificar conectividade, inspecionar headers
- **fs_rmrf** — remover diretório com conteúdo (já mencionada no código, aguardando implementação)

### Suporte a múltiplos provedores de IA (médio prazo)

Expansão da camada `AIProvider` para suporte nativo a Anthropic, OpenAI, Gemini, Groq e OpenRouter. A troca será feita via `HUMANCLI_PROVIDER` no `.env`, sem alterar nenhum plugin.

### Melhorias no agente (médio prazo)

- Suporte a múltiplas tools por iteração (plano paralelo)
- Tool scheduling com dependências entre steps
- Memória longa entre sessões diferentes

---

## Contribuindo

### Criando uma nova tool nativa

1. Crie um arquivo `.go` em `internal/adapter/tools/native/<categoria>/`
2. Implemente a interface `tool.Tool` (`Name`, `Description`, `Execute`)
3. Registre no `init()` com `tools.GlobalRegistry().Register(&SuaTool{})`
4. Rode `go generate ./...` para atualizar os imports

### Criando um plugin externo

Consulte o [PLUGIN_README.md](./PLUGIN_README.md) para o guia completo.

### Interface Tool

```go
type Tool interface {
    Name()        string
    Description() string
    Execute(params map[string]interface{}) (any, error)
}
```

A `Description()` é lida diretamente pelo LLM — ela determina quando e como a tool será acionada. Inclua palavras-chave em linguagem natural, exemplos de uso, descrição dos parâmetros e comportamentos esperados.

---

*HumanCLI — porque você não deveria precisar lembrar de comando nenhum.*