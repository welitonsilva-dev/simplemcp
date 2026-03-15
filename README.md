# humancli-server

> Um agente com LLM local que executa qualquer ferramenta em linguagem natural — extensível, privado e sem dependência de nuvem.
>
> *Fale o que quer. O agente entende, age, observa o resultado e decide o próximo passo.*

---

## Sumário

- [O que é o humancli-server?](#o-que-é-o-humancli-server)
- [Ecossistema HumanCLI](#ecossistema-humancli)
- [Como funciona](#como-funciona)
- [O que já existe](#o-que-já-existe)
- [Provedores de IA](#provedores-de-ia)
- [Diferencial](#diferencial)
- [Roadmap](#roadmap)
- [Tecnologias](#tecnologias)
- [Contribuindo](#contribuindo)

---

## O que é o humancli-server?

O **humancli-server** é o núcleo do ecossistema HumanCLI: um servidor de agente que interpreta linguagem natural e executa ferramentas (tools) de forma autônoma.

Você escreve o que quer. O agente entende, age, observa o resultado e decide o próximo passo — repetindo o ciclo até concluir a tarefa.

> *"Cria a pasta projeto e adiciona um README"*, *"Sobe o ambiente Docker"*, *"Lista os arquivos maiores que 100MB"* — sem lembrar de um único comando.

Com o sistema de plugins, o servidor pode ser expandido para qualquer coisa: automações, integrações, comandos personalizados, rotinas de trabalho.

---

## Ecossistema HumanCLI

| Projeto | Papel | Status |
|---|---|---|
| **humancli-server** | Servidor de agente. Hospeda o LLM local, registra as tools e executa o loop ReAct. | ✅ Este repositório |
| **humancli-client** | Interface do usuário. Ponto de entrada onde você digita comandos em linguagem natural. | 🔜 Em desenvolvimento |

---

## Como funciona

O agente opera em um loop **ReAct (Reason + Act)**:

```
usuário → pipeline → [LLM raciocina → executa tool → observa resultado] → resposta
                      └──────────────── loop até concluir ────────────────┘
```

A cada iteração, o LLM recebe o histórico completo da conversa — input original mais todos os resultados anteriores — e decide de forma autônoma entre chamar mais uma ferramenta ou encerrar com uma resposta ao usuário.

Isso é o que diferencia o HumanCLI de um dispatcher simples: o agente **observa o que aconteceu** e **adapta os próximos passos**, podendo encadear ações, recuperar de erros e confirmar conclusão antes de responder.

Veja o [FLOW_README.md](./FLOW_README.md) para o fluxo detalhado com exemplos.

---

## O que já existe

### Agente ReAct
O núcleo do servidor implementa o loop ReAct completo: cada iteração envia o histórico ao LLM, que decide a próxima ação ou encerra o loop. O número máximo de iterações é configurável via `AGENT_MAX_ITERATIONS`.

### Servidor HTTP
Servidor HTTP que recebe mensagens em linguagem natural e retorna a resposta consolidada do agente após o loop concluir.

### Sistema de Plugins
Plugins são módulos Go independentes. Basta criar um pacote seguindo a interface `sdk.Tool` e o servidor o reconhece automaticamente via `go generate`. Isso significa que você pode criar:

- Comandos de sistema (Git, Docker, filesystem)
- Automações personalizadas
- Integrações com serviços externos
- Qualquer ferramenta que sua imaginação permitir

### Categoria: Filesystem
Tools nativas para interação com o sistema de arquivos:

- `fs_mkdir` — criar pastas
- `fs_touch` — criar arquivos
- `fs_list` — listar diretórios
- `fs_mv` — mover/renomear arquivos
- `fs_rmdir` — remover pastas (requer alta confidence)

### Proteções de segurança
- **Confidence guard** — tools destrutivas são bloqueadas quando o LLM está inseguro sobre a intenção
- **Max iterations** — limite configurável de ciclos para evitar loops infinitos

### Projeto de Exemplo
Plugin `hello` criado com base no [PLUGIN_README.md](./PLUGIN_README.md), que serve como referência para novos contribuidores.

---

## Provedores de IA

O HumanCLI foi projetado para ser **agnóstico ao provedor de IA**. Hoje funciona com LLM local via Ollama. A arquitetura usa uma interface Go simples (`AIProvider`) para que qualquer provedor possa ser integrado sem alterar tools ou plugins.

```env
HUMANCLI_PROVIDER=ollama      # padrão — Ollama, LM Studio, etc.
```

Suporte a provedores externos está no roadmap:

```env
# planejado — não implementado ainda
HUMANCLI_PROVIDER=anthropic
HUMANCLI_PROVIDER=openai
HUMANCLI_PROVIDER=gemini
HUMANCLI_PROVIDER=groq
HUMANCLI_PROVIDER=openrouter
```

---

## Diferencial

| Característica | HumanCLI | Alternativas (Aider, Shell AI, etc.) |
|---|---|---|
| Loop ReAct — age e observa resultados | ✅ | ❌ (maioria executa plano fixo) |
| Local e offline por padrão | ✅ | ❌ (maioria usa APIs externas) |
| Sistema de plugins extensível | ✅ | ❌ |
| Vai além de comandos de sistema | ✅ | ❌ |
| Desenvolvido em Go | ✅ | ❌ (maioria em Python) |
| Sem dependência de IDE | ✅ | Parcial |
| Arquitetura cliente/servidor | ✅ | ❌ |

---

## Roadmap

### Novas categorias de tools (curto prazo)

- **Git** — `git_init`, `git_commit`, `git_checkout`, `git_push`, `git_status`
- **Docker** — `docker_build`, `docker_run`, `docker_ps`, `docker_stop`
- **Sistema/Processos** — matar processos, verificar portas, uso de memória
- **Rede/HTTP** — fazer requisições, verificar conectividade, inspecionar headers

### humancli-client (médio prazo)

```bash
humancli "faz um commit com as mudanças de hoje"
humancli "sobe o docker do projeto"
humancli "qual processo está usando a porta 3000?"
```

### Suporte a múltiplos provedores de IA (médio prazo)

Expansão da camada `AIProvider` para suporte a Anthropic, OpenAI, Gemini, Groq e OpenRouter. A troca será feita via variável de ambiente, sem alterar nenhum plugin existente.

---

## Tecnologias

- **Linguagem:** Go
- **LLM (atual):** Ollama / LM Studio (local, sem custo)
- **LLM (futuro):** Anthropic, OpenAI, Gemini, Groq, OpenRouter
- **Transporte atual:** HTTP
- **Transporte futuro:** humancli-client (CLI nativa)

---

## Contribuindo

O projeto tem um sistema de plugins documentado no [PLUGIN_README.md](./PLUGIN_README.md). Para criar uma nova tool:

```go
type Tool interface {
    Name()        string
    Description() string
    Execute(params map[string]interface{}) (any, error)
}
```

Registre com `sdk.Register(&SuaTool{})` e rode `go generate ./...` no humancli-server.

---

*HumanCLI — porque você não deveria precisar lembrar de comando nenhum.*
