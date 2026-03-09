# SimpleMCP

> Um servidor MCP local, extensível e offline para execução de comandos via linguagem natural.

---

## Sumário

- [O que é o SimpleMCP?](#o-que-é-o-simplemcp)
- [Objetivo](#objetivo)
- [O que já existe](#o-que-já-existe)
  - [Servidor HTTP](#servidor-http)
  - [Sistema de Plugins](#sistema-de-plugins)
  - [Categoria: Filesystem](#categoria-filesystem)
  - [Projeto de Exemplo](#projeto-de-exemplo)
- [Diferencial](#diferencial)
- [Roadmap e Expectativas de Mudança](#roadmap-e-expectativas-de-mudança)
  - [Novas categorias de tools](#novas-categorias-de-tools-curto-prazo)
  - [CLI](#cli-médio-prazo)
  - [Suporte a múltiplos provedores de IA](#suporte-a-múltiplos-provedores-de-ia-médiolongo-prazo)
- [Tecnologias](#tecnologias)
- [Contribuindo](#contribuindo)

---

## O que é o SimpleMCP?

O SimpleMCP é um servidor MCP (**Model Context Protocol**) desenvolvido em Go, que permite controlar o sistema operacional, ferramentas de desenvolvimento e serviços através de linguagem natural humana — sem depender de nenhuma API externa ou conexão com a internet.

A ideia central é simples: você fala o que quer fazer, e a IA entende e executa.

> *"Faz um commit com as mudanças de hoje"*, *"Sobe o ambiente Docker"*, *"Lista os arquivos maiores que 100MB"* — tudo isso sem lembrar de um único comando.

---

## Objetivo

Criar um ecossistema extensível de ferramentas (tools) que permitam ao usuário interagir com o computador, projetos e serviços usando linguagem natural, com foco em:

- **Privacidade** — IA rodando 100% local, sem envio de dados para servidores externos
- **Extensibilidade** — qualquer desenvolvedor pode criar e adicionar novas tools via sistema de plugins
- **Acessibilidade** — reduzir a barreira de uso de ferramentas técnicas como Git, Docker e CLI
- **Performance** — construído em Go para leveza e velocidade nativa

---

## O que já existe

### Servidor HTTP
O SimpleMCP já funciona como um servidor HTTP que recebe requisições em linguagem natural, interpreta a intenção do usuário e executa a tool correspondente.

### Sistema de Plugins
O projeto possui um sistema de plugins que permite expandir as tools sem alterar o núcleo do servidor. Basta criar um novo pacote Go seguindo a interface padrão e registrá-lo — o servidor já o reconhece automaticamente.

### Categoria: Filesystem
Já existe a categoria `filesystem` com tools para interação com o sistema de arquivos:

- Criar pastas
- Mover arquivos
- Listar diretórios
- Renomear pastas
- Entre outras

### Projeto de Exemplo
Existe um plugin de exemplo (`hello`) criado com base no [PLUGIN_README.md](./PLUGIN_README.md), que serve como referência para novos contribuidores criarem suas próprias tools.

---

## Diferencial

| Característica | SimpleMCP | Alternativas (Aider, Shell AI, etc.) |
|---|---|---|
| 100% local e offline | ✅ | ❌ (maioria usa APIs externas) |
| Sistema de plugins | ✅ | ❌ |
| Protocolo MCP nativo | ✅ | Parcial |
| Desenvolvido em Go | ✅ | ❌ (maioria em Python) |
| Sem dependência de IDE | ✅ | Parcial |

Enquanto ferramentas como **Claude Code**, **Cursor** e **Aider** dependem de APIs pagas ou estão presas em IDEs, o SimpleMCP roda inteiramente na máquina do usuário, com modelos de linguagem locais como os suportados pelo **Ollama** e similares.

Além disso, o protocolo MCP está crescendo rapidamente como padrão universal de integração entre IAs e ferramentas — e o SimpleMCP já nasce sobre esse protocolo.

---

## Roadmap e Expectativas de Mudança

### Novas categorias de tools (curto prazo)
Expansão do ecossistema de tools com novas categorias, cada uma com comandos separados:

- **Git** — `git_init`, `git_commit`, `git_checkout`, `git_push`, `git_status`, etc.
- **Docker** — `docker_build`, `docker_run`, `docker_ps`, `docker_stop`, etc.
- **Sistema/Processos** — matar processos, verificar portas, uso de memória, etc.
- **Rede/HTTP** — fazer requisições, verificar conectividade, inspecionar headers, etc.

### CLI (médio prazo)
Transformar o servidor HTTP em uma aplicação CLI nativa, permitindo uso direto no terminal sem necessidade de um servidor rodando em segundo plano.

```bash
simplemcp "faz um commit com as mudanças de hoje"
```

### Suporte a múltiplos provedores de IA (médio/longo prazo)
Integração com provedores de IA externos como alternativa à IA local, respeitando a escolha e o hardware do usuário:

```env
MCP_PROVIDER=local       # padrão atual (Ollama, LM Studio)
MCP_PROVIDER=anthropic   # Claude API
MCP_PROVIDER=openai      # GPT-4
MCP_PROVIDER=gemini      # Google AI
```

A arquitetura será baseada em uma interface Go simples (`AIProvider`), mantendo o restante do sistema agnóstico ao provedor utilizado. Isso permite:

- Usuários com hardware limitado usarem APIs pagas
- Usuários com boas máquinas rodarem tudo localmente e de graça
- Possibilidade futura de uso misto (tarefas simples local, tarefas complexas na nuvem)

---

## Tecnologias

- **Linguagem:** Go
- **Protocolo:** MCP (Model Context Protocol)
- **IA local:** Ollama / LM Studio (ou similar)
- **Transporte atual:** HTTP
- **Transporte futuro:** CLI nativo

---

## Contribuindo

O projeto possui um sistema de plugins documentado no [PLUGIN_README.md](./PLUGIN_README.md). Para criar uma nova tool, basta implementar a interface padrão:

```go
type Tool interface {
    Name() string
    Description() string
    Execute(params map[string]interface{}) (any, error)
}
```

Registre com `sdk.Register(&SuaTool{})` e o servidor já a reconhece automaticamente.

---

*SimpleMCP — porque você não deveria precisar lembrar de comando nenhum.*
