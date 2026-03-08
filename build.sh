#!/bin/bash

cd "$(dirname "$0")/.."

if [ -d "simplemcpplugins" ]; then
  echo "simplemcpplugins encontrado, build completo"
  docker build -f simplemcp/Dockerfile .
else
  echo "simplemcpplugins não encontrado, build só simplemcp"
  docker build -f simplemcp/Dockerfile.standalone simplemcp/
fi 