#!/bin/bash
set -e

echo "🚀 Iniciando Ollama..."
ollama serve &

echo "⏳ Aguardando Ollama iniciar..."

# Substituímos o curl por 'ollama list'
until ollama list >/dev/null 2>&1; do
  sleep 2
done

MODEL=${LLM_MODEL:-qwen-2.5-3b}

echo "📦 Modelo configurado: $MODEL"

if ! ollama list | grep -q "$MODEL"; then
  echo "⬇️ Baixando modelo $MODEL..."
  ollama pull "$MODEL"
else
  echo "✅ Modelo $MODEL já existe"
fi

echo "✅ Ollama pronto!"
wait