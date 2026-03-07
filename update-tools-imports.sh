#!/bin/bash

# =============================================================
#  update-tool-imports.sh
#  Detecta todos os subpacotes em internal/tools/native/
#  e atualiza o bloco de imports blank no cmd/main.go
#
#  Uso: bash update-tool-imports.sh
#  Execute na raiz do projeto.
# =============================================================

set -e

MAIN_FILE="cmd/main.go"
TOOLS_DIR="internal/tools/native"
MARKER="pacotes de ferramentas nativas"

# ── Validações ────────────────────────────────────────────────

if [ ! -f "$MAIN_FILE" ]; then
  echo "❌ $MAIN_FILE não encontrado. Execute na raiz do projeto."
  exit 1
fi

if [ ! -d "$TOOLS_DIR" ]; then
  echo "❌ Diretório $TOOLS_DIR não encontrado."
  exit 1
fi

MODULE=$(grep '^module ' go.mod | awk '{print $2}')
if [ -z "$MODULE" ]; then
  echo "❌ Não foi possível detectar o module name no go.mod."
  exit 1
fi

echo "📦 Module  : $MODULE"
echo "📁 Scanning: $TOOLS_DIR"
echo ""

# ── Detecta subpacotes e injeta imports ───────────────────────

ADDED=0
SKIPPED=0

for dir in "$TOOLS_DIR"/*/; do
  [ -d "$dir" ] || continue

  gofiles=$(find "$dir" -maxdepth 1 -name "*.go" | head -1)
  [ -z "$gofiles" ] && continue

  # Remove trailing slash e monta import path
  pkg="${MODULE}/${dir%/}"

  if grep -qF "\"$pkg\"" "$MAIN_FILE"; then
    echo "   ⏭  Já existe : $pkg"
    SKIPPED=$((SKIPPED + 1))
  else
    # Injeta linha após o marcador
    sed -i "/$MARKER/a \\\\t_ \"$pkg\"" "$MAIN_FILE"
    echo "   ➕ Adicionado: $pkg"
    ADDED=$((ADDED + 1))
  fi
done

echo ""
echo "────────────────────────────────────────"
echo "✅ Concluído: $ADDED adicionado(s), $SKIPPED já existia(m)"
echo "────────────────────────────────────────"
echo ""

# ── Mostra o bloco final no main.go ──────────────────────────

echo "📄 Bloco atual em $MAIN_FILE:"
echo ""
grep -n "pacotes de ferramentas\|_ \"$MODULE" "$MAIN_FILE"
echo ""
echo "💡 Para nova tool: crie internal/tools/native/<nome>/ e rode este script."