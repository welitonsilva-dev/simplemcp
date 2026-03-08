# docker-clean.sh — Documentação de Processo e Fluxo

> Reset completo do ambiente Docker: remove containers, imagens, volumes, redes e cache.

---

## Visão Geral

O `docker-clean.sh` realiza uma limpeza **total e irreversível** do Docker local. Ao final da execução, o ambiente retorna ao estado zero — sem containers, imagens, volumes ou redes customizadas.

**Localização no projeto:** `docker/docker-clean.sh`

---

## Pré-requisitos

- Docker instalado e acessível no terminal
- Usuário com permissão para executar comandos Docker (grupo `docker` ou `root`)
- Bash disponível (`/bin/bash`)
- Permissão de execução no arquivo:

```bash
chmod +x docker/docker-clean.sh
```

---

## Fluxo de Execução

### Fase 1 — Confirmação

Antes de qualquer ação, o script exibe um aviso listando tudo que será apagado e solicita confirmação explícita. Apenas a entrada `sim` prossegue; qualquer outro valor cancela a operação.

```
⚠️  ATENÇÃO: Isso vai apagar TUDO do Docker!
   • Todos os containers (rodando ou parados)
   • Todas as imagens
   • Todos os volumes (incluindo dados persistidos)
   • Todas as redes customizadas
   • Todo o cache de build

Tem certeza? Digite 'sim' para continuar:
```

```bash
if [ "$CONFIRM" != "sim" ]; then
  echo "Operação cancelada."
  exit 0
fi
```

---

### Fase 2 — Limpeza (7 Etapas Sequenciais)

As etapas 1–4 verificam se existem recursos antes de agir, evitando erros em ambientes já limpos.

| Etapa | Ação | Comando |
|-------|------|---------|
| `[1/7]` | Parar containers em execução | `docker stop $(docker ps -q)` |
| `[2/7]` | Remover todos os containers | `docker rm -f $(docker ps -aq)` |
| `[3/7]` | Remover todas as imagens | `docker rmi -f $(docker images -q)` |
| `[4/7]` | Remover todos os volumes | `docker volume rm -f $(docker volume ls -q)` |
| `[5/7]` | Remover redes customizadas | `docker network prune -f` |
| `[6/7]` | Limpar cache de build (BuildKit) | `docker builder prune -af` |
| `[7/7]` | Limpeza geral (garantia final) | `docker system prune -af --volumes` |

> A etapa 7 é uma salvaguarda: garante que nada ficou para trás caso alguma etapa anterior não tenha coberto todos os recursos.

---

### Fase 3 — Verificação Final

Após a limpeza, o script exibe o estado atual do ambiente:

```
📊 Estado atual do Docker:
── Containers: 0
── Imagens:    0
── Volumes:    0
── Redes:      3 (apenas as padrão)
```

---

## Diagrama de Fluxo

```
[ Início ]
    │
    ▼
[ Exibe aviso e solicita confirmação ]
    │
    ├── input ≠ 'sim' ──► [ "Operação cancelada" ] ──► [ Fim ]
    │
    ▼  input = 'sim'
[ 1/7 ] Parar containers em execução
    │
[ 2/7 ] Remover todos os containers
    │
[ 3/7 ] Remover todas as imagens
    │
[ 4/7 ] Remover todos os volumes
    │
[ 5/7 ] Remover redes customizadas
    │
[ 6/7 ] Limpar cache de build
    │
[ 7/7 ] Limpeza geral (system prune --volumes)
    │
    ▼
[ Exibe estado final do Docker ]
    │
[ Fim ]
```

---

## Comportamento Condicional (Etapas 1–4)

Para cada recurso, o script verifica a existência antes de remover:

```bash
RUNNING=$(docker ps -q)
if [ -n "$RUNNING" ]; then
  docker stop $RUNNING
else
  echo "Nenhum container em execução."
fi
```

Isso torna o script seguro para execução mesmo em ambientes já limpos, sem gerar erros.

---

## Saída de Cores no Terminal

| Cor | Uso |
|-----|-----|
| 🔴 Vermelho | Avisos críticos e cabeçalho de atenção |
| 🟡 Amarelo | Indicadores de etapa em progresso |
| 🟢 Verde | Confirmações de sucesso |
| 🔵 Ciano | Cabeçalhos e estado final |

---

## Como Executar

A partir da raiz do projeto:

```bash
./docker/docker-clean.sh
```

Ou diretamente do diretório `docker/`:

```bash
cd docker && bash docker-clean.sh
```

---

## ⚠️ Avisos Importantes

- **Irreversível:** não há desfazer após a confirmação
- **Dados perdidos permanentemente:** volumes incluem bancos de dados e arquivos de configuração persistidos
- **Não execute em produção**
- Redes padrão do Docker (`bridge`, `host`, `none`) **não são removidas** — apenas redes customizadas
