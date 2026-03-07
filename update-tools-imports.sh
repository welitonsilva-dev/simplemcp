#!/bin/bash

# =============================================================
#  update-tool-imports.sh
#  Detecta todos os subpacotes em internal/tools/native/
#  e internal/tools/plugins/ e atualiza o bloco de imports
#  blank no cmd/main.go
#
#  Uso: bash update-tool-imports.sh
#  Execute na raiz do projeto.
# =============================================================

set -e

MAIN_FILE="cmd/main.go"
NATIVE_DIR="internal/tools/native"
PLUGINS_DIR="internal/tools/plugins"
MARKER_NATIVE="pacotes de ferramentas nativas"
MARKER_PLUGINS="pacotes de ferramentas plugins"

# ── Validações ────────────────────────────────────────────────

if [ ! -f "$MAIN_FILE" ]; then
  echo "❌ $MAIN_FILE não encontrado. Execute na raiz do projeto."
  exit 1
fi

if [ ! -d "$NATIVE_DIR" ] && [ ! -d "$PLUGINS_DIR" ]; then
  echo "❌ Nenhum diretório de tools encontrado ($NATIVE_DIR ou $PLUGINS_DIR)."
  exit 1
fi

MODULE=$(grep '^module ' go.mod | awk '{print $2}')
if [ -z "$MODULE" ]; then
  echo "❌ Não foi possível detectar o module name no go.mod."
  exit 1
fi

echo "📦 Module  : $MODULE"
echo ""

# ── Função para escanear e injetar imports ───────────────────

ADDED=0
SKIPPED=0

inject_imports() {
  local dir=$1
  local label=$2
  local marker=$3

  if [ ! -d "$dir" ]; then
    echo "⚠️  Diretório não encontrado, pulando: $dir"
    return
  fi

  echo "📁 Scanning $label: $dir"

  for subdir in "$dir"/*/; do
    [ -d "$subdir" ] || continue

    gofiles=$(find "$subdir" -maxdepth 1 -name "*.go" | head -1)
    [ -z "$gofiles" ] && continue

    pkg="${MODULE}/${subdir%/}"

    if grep -qF "\"$pkg\"" "$MAIN_FILE"; then
      echo "   ⏭  Já existe : $pkg"
      SKIPPED=$((SKIPPED + 1))
    else
      sed -i "/$marker/a \\\\t_ \"$pkg\"" "$MAIN_FILE"
      echo "   ➕ Adicionado: $pkg"
      ADDED=$((ADDED + 1))
    fi
  done

  echo ""
}

# ── Executa para native e plugins ────────────────────────────

inject_imports "$NATIVE_DIR" "native" "$MARKER_NATIVE"
inject_imports "$PLUGINS_DIR" "plugins" "$MARKER_PLUGINS"

echo "────────────────────────────────────────"
echo "✅ Concluído: $ADDED adicionado(s), $SKIPPED já existia(m)"
echo "────────────────────────────────────────"
echo ""

# ── Mostra o bloco final no main.go ──────────────────────────

echo "📄 Bloco atual em $MAIN_FILE:"
echo ""
grep -n "pacotes de ferramentas\|_ \"$MODULE" "$MAIN_FILE"
echo ""
echo "💡 Para nova tool: crie $NATIVE_DIR/<nome>/ ou $PLUGINS_DIR/<nome>/ e rode este script."