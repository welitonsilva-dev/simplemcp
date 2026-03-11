# SimpleMCP — Guia de Setup e Uso

## Sumário

- [Inicialização do Ambiente](#inicialização-do-ambiente)
- [Verificação dos Containers](#verificação-dos-containers)
- [Rotas da API](#rotas-da-api)
  - [Chat de Comando (MCP)](#chat-de-comando-mcp)
  - [Chat do Modelo LLM (Ollama)](#chat-do-modelo-llm-ollama)

---

## Inicialização do Ambiente

### 1. Limpar o ambiente Docker

```bash
cd simplemcp/docker && sh docker-clean.sh
```

### 2. Voltar para a raiz do projeto

```bash
cd ../
```

### 3. Registrar tools nativas e plugins externos

```bash
go generate ./...
```

> Esse comando escaneia o projeto e registra automaticamente todas as tools nativas e plugins externos disponíveis.

### 4. Subir os containers

```bash
docker compose up
```

Aguarde o pull completo do modelo de IA local. Quando concluído, pressione `d` para sair da visualização dos logs do container.

---

## Verificação dos Containers

Verifique se os containers subiram corretamente:

```bash
docker ps
```

Caso os containers não tenham iniciado na primeira execução, rode:

```bash
docker compose up -d
```

---

## Rotas da API

### Chat de Comando (MCP)

Envie comandos em linguagem natural para o servidor MCP.

**Rota:** `POST http://localhost:8081/v1/do`

**Body:**
```json
{
  "message": "comando em linguagem natural. Ex.: quais tools existe"
}
```

**Resposta `200`:**
```json
{
  "results": "step 1 [tool_list]: nativas:\n  - tool_list\n  - fs_cd\n  - fs_mkdir\n  - fs_rm\n  - fs_rmdir\n  - fs_touch\n  - echo\n  - fs_list\n  - fs_rmrf\n  - double_echo\n\nplugins:\n  - docker_ps\n  - hello\n\n"
}
```

---

### Chat do Modelo LLM (Ollama)

Interaja diretamente com o modelo de linguagem local via Ollama.

**Rota:** `POST http://localhost:11434/api/generate`

**Body:**
```json
{
  "model": "llama3",
  "prompt": "diga 3 palavras",
  "stream": false
}
```

**Resposta `200`:**
```json
{
  "model": "llama3",
  "created_at": "2026-03-08T19:59:54.044299424Z",
  "response": "Fácil!\n\n1. Amigo\n2. Feliz\n3. Voador",
  "done": true,
  "done_reason": "stop",
  "total_duration": 9853913643,
  "load_duration": 5667100991,
  "prompt_eval_count": 16,
  "prompt_eval_duration": 1054232592,
  "eval_count": 19,
  "eval_duration": 3122188127
}
```

> O campo `context` retornado pelo Ollama contém os tokens da conversa e pode ser reutilizado em requisições subsequentes para manter o histórico de contexto.