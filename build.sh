#!/bin/bash

cd "$(dirname "$0")/.."

if [ -d "humancli-plugins" ]; then
  echo "humancli-plugins encontrado, build completo"
  docker build -f humancli-server/Dockerfile .
else
  echo "humancli-plugins não encontrado, build só humancli-server"
  docker build -f humancli-server/Dockerfile.standalone humancli-server/
fi 