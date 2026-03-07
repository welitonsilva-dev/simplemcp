#!/bin/bash

# =============================================================
#  cli.sh — Wrapper Linux/Mac para o MCP CLI
#  Detecta disco e diretório atual automaticamente
# =============================================================

# Diretório raiz do disco (sempre / no Linux e Mac)
HOST_ROOT="/"

# Diretório atual do usuário no host
HOST_CWD="$(pwd)"

# Caminho relativo do CWD dentro do container
# Ex: /home/user/projetos/repo1 → /app/host/home/user/projetos/repo1
CONTAINER_CWD="/app/host${HOST_CWD}"

export HOST_ROOT
export HOST_CWD
export CONTAINER_CWD

# Sobe o container com o disco mapeado dinamicamente
docker compose run --rm \
  -e HOST_ROOT="$HOST_ROOT" \
  -e HOST_CWD="$HOST_CWD" \
  -e CONTAINER_CWD="$CONTAINER_CWD" \
  -v "${HOST_ROOT}:/app/host" \
  -w "$CONTAINER_CWD" \
  mcp-server "$@"