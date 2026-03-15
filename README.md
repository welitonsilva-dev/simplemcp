# humancli-server

> Um agente com LLM local que executa qualquer ferramenta em linguagem natural — extensível, privado e sem dependência de nuvem.
>
> *Fale o que quer. O agente entende, executa e aprende novas habilidades com plugins.*

---

## Sumário

- [O que é o humancli-server?](#o-que-é-o-humancli-server)
- [Ecossistema HumanCLI](#ecossistema-humancli)
- [Objetivo](#objetivo)
- [O que já existe](#o-que-já-existe)
  - [Servidor HTTP](#servidor-http)
  - [Sistema de Plugins](#sistema-de-plugins)
  - [Categoria: Filesystem](#categoria-filesystem)
  - [Projeto de Exemplo](#projeto-de-exemplo)
- [Provedores de IA](#provedores-de-ia)
- [Diferencial](#diferencial)
- [Roadmap](#roadmap)
- [Tecnologias](#tecnologias)
- [Contribuindo](#contribuindo)

---

## O que é o humancli-server?

O **humancli-server** é o núcleo do ecossistema HumanCLI: um servidor de agente que interpreta linguagem natural e executa ferramentas (tools). Hoje roda com LLM 100% local, sem depender de nenhuma API externa — e foi projetado desde o início para suportar qualquer provedor de IA com uma simples troca de configuração.

Você escreve o que quer. O agente entende e faz.

> *"Faz um commit com as mudanças de hoje"*, *"Sobe o ambiente Docker"*, *"Lista os arquivos maiores que 100MB"* — tudo isso sem lembrar de um único comando.

Mas vai além de executar comandos do sistema. Com o sistema de plugins, o servidor pode ser expandido para qualquer coisa que você imaginar: automações, integrações, comandos personalizados, rotinas de trabalho — o limite é a criatividade.

### Como funciona

```
Usuário → HTTP → Pipeline → AIProvider (LLM) → executa tool → resposta
```

O servidor recebe a mensagem, normaliza o texto via pipeline de sanitização, envia ao provedor de IA configurado para gerar um plano de execução, e então aciona as ferramentas correspondentes. O provedor é intercambiável — hoje Ollama local, amanhã qualquer outra coisa.

---

## Ecossistema HumanCLI

O HumanCLI é composto por duas partes que trabalham juntas:

| Projeto | Papel | Status |
|---|---|---|
| **humancli-server** | Servidor de agente. Hospeda o LLM local, registra as tools e executa os planos. | ✅ Este repositório |
| **humancli-client** | Interface do usuário. O ponto de entrada onde você digita comandos em linguagem natural. | 🔜 Em desenvolvimento |

O **humancli-client** será a experiência final do usuário — uma interface leve que se comunica com o servidor e torna o uso ainda mais fluido, seja via terminal, CLI nativa ou outra forma de interação.

---

## Objetivo

Criar um ecossistema extensível onde qualquer pessoa possa interagir com ferramentas técnicas usando linguagem natural, com foco em:

- **Liberdade** — não apenas comandos do sistema, mas qualquer automação ou integração que você quiser criar
- **Flexibilidade** — use o LLM que preferir: local e gratuito hoje, pago e na nuvem amanhã, com uma linha de configuração
- **Privacidade** — por padrão roda 100% local, sem envio de dados para servidores externos
- **Extensibilidade** — qualquer desenvolvedor pode criar e registrar novas tools via sistema de plugins
- **Acessibilidade** — reduzir a barreira de uso de ferramentas técnicas como Git, Docker, CLI e muito mais
- **Performance** — construído em Go para leveza e velocidade nativa

---

## O que já existe

### Servidor HTTP
O humancli-server já funciona como um servidor HTTP que recebe mensagens em linguagem natural, interpreta a intenção do usuário e executa a tool correspondente.

### Sistema de Plugins
O projeto possui um sistema de plugins que permite expandir as tools sem alterar o núcleo do servidor. Basta criar um novo pacote Go seguindo a interface padrão e registrá-lo — o servidor já o reconhece automaticamente.

Isso significa que você pode criar:
- Comandos de sistema (Git, Docker, filesystem)
- Automações personalizadas
- Integrações com serviços externos
- Qualquer ferramenta que sua imaginação permitir

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

## Provedores de IA

O HumanCLI foi projetado para ser **agnóstico ao provedor de IA**. Hoje funciona com LLM local via Ollama — sem custo, sem internet, sem envio de dados. No futuro, trocar de provedor será tão simples quanto mudar uma variável de ambiente:

```env
HUMANCLI_PROVIDER=local       # padrão — Ollama, LM Studio, etc.
HUMANCLI_PROVIDER=anthropic   # Claude (Anthropic)
HUMANCLI_PROVIDER=openai      # GPT-4 e família
HUMANCLI_PROVIDER=gemini      # Google Gemini
HUMANCLI_PROVIDER=groq        # Groq (rápido e gratuito com limites)
HUMANCLI_PROVIDER=openrouter  # OpenRouter (acesso a vários modelos)
```

A arquitetura é baseada em uma interface Go simples (`AIProvider`), o que significa que qualquer provedor — pago, gratuito, local ou na nuvem — pode ser integrado sem alterar o resto do sistema.

Isso dá liberdade real ao usuário:

- Máquina potente? Rode tudo local, de graça e com privacidade total.
- Hardware limitado? Use uma API paga só quando precisar.
- Quer experimentar modelos diferentes? Troque com uma linha.

---

## Diferencial

| Característica | HumanCLI | Alternativas (Aider, Shell AI, etc.) |
|---|---|---|
| Local e offline por padrão | ✅ | ❌ (maioria usa APIs externas) |
| Troca de provedor por config | ✅ | ❌ |
| Sistema de plugins extensível | ✅ | ❌ |
| Vai além de comandos de sistema | ✅ | ❌ |
| Desenvolvido em Go | ✅ | ❌ (maioria em Python) |
| Sem dependência de IDE | ✅ | Parcial |
| Arquitetura cliente/servidor | ✅ | ❌ |

Enquanto ferramentas como **Claude Code**, **Cursor** e **Aider** dependem de APIs pagas ou estão presas em IDEs, o HumanCLI começa local e cresce junto com as suas necessidades.

---

## Roadmap

### Novas categorias de tools (curto prazo)
Expansão do ecossistema com novas categorias nativas:

- **Git** — `git_init`, `git_commit`, `git_checkout`, `git_push`, `git_status`, etc.
- **Docker** — `docker_build`, `docker_run`, `docker_ps`, `docker_stop`, etc.
- **Sistema/Processos** — matar processos, verificar portas, uso de memória, etc.
- **Rede/HTTP** — fazer requisições, verificar conectividade, inspecionar headers, etc.

### humancli-client (médio prazo)
Desenvolvimento do cliente oficial — a interface do usuário que se conecta ao servidor e oferece uma experiência fluida de uso em linguagem natural.

```bash
humancli "faz um commit com as mudanças de hoje"
humancli "sobe o docker do projeto"
humancli "qual processo está usando a porta 3000?"
```

### Suporte a múltiplos provedores de IA (médio prazo)

Expansão da camada `AIProvider` para suportar provedores externos além do Ollama local. A interface já está desenhada para isso — é uma questão de implementar os adaptadores:

- **Anthropic** — Claude API
- **OpenAI** — GPT-4 e família
- **Google** — Gemini
- **Groq** — modelos rápidos com tier gratuito
- **OpenRouter** — acesso unificado a dezenas de modelos

A troca entre provedores será feita via variável de ambiente, sem alterar nenhum plugin ou tool existente.

---

## Tecnologias

- **Linguagem:** Go
- **LLM (atual):** Ollama / LM Studio (local, sem custo)
- **LLM (futuro):** Anthropic, OpenAI, Gemini, Groq, OpenRouter e outros
- **Transporte atual:** HTTP
- **Transporte futuro:** humancli-client (CLI nativa)

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

*HumanCLI — porque você não deveria precisar lembrar de comando nenhum.*