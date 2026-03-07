#!/bin/bash

# =============================================================
#  docker-clean.sh — Reset completo do Docker
#  Limpa containers, imagens, volumes, redes e cache
# =============================================================

RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}"
echo "╔══════════════════════════════════════════════╗"
echo "║        🐳  Docker Full Reset Script          ║"
echo "╚══════════════════════════════════════════════╝"
echo -e "${NC}"

# Confirmação antes de executar
echo -e "${RED}⚠️  ATENÇÃO: Isso vai apagar TUDO do Docker!${NC}"
echo -e "${YELLOW}   • Todos os containers (rodando ou parados)"
echo -e "   • Todas as imagens"
echo -e "   • Todos os volumes (incluindo dados persistidos)"
echo -e "   • Todas as redes customizadas"
echo -e "   • Todo o cache de build${NC}"
echo ""
read -p "Tem certeza? Digite 'sim' para continuar: " CONFIRM

if [ "$CONFIRM" != "sim" ]; then
  echo -e "${GREEN}Operação cancelada.${NC}"
  exit 0
fi

echo ""
echo -e "${CYAN}🔄 Iniciando limpeza...${NC}"
echo ""

# 1. Parar todos os containers em execução
echo -e "${YELLOW}[1/7] Parando todos os containers em execução...${NC}"
RUNNING=$(docker ps -q)
if [ -n "$RUNNING" ]; then
  docker stop $RUNNING
  echo -e "${GREEN}      ✔ Containers parados.${NC}"
else
  echo "      Nenhum container em execução."
fi

# 2. Remover todos os containers (incluindo parados)
echo -e "${YELLOW}[2/7] Removendo todos os containers...${NC}"
CONTAINERS=$(docker ps -aq)
if [ -n "$CONTAINERS" ]; then
  docker rm -f $CONTAINERS
  echo -e "${GREEN}      ✔ Containers removidos.${NC}"
else
  echo "      Nenhum container encontrado."
fi

# 3. Remover todas as imagens
echo -e "${YELLOW}[3/7] Removendo todas as imagens...${NC}"
IMAGES=$(docker images -q)
if [ -n "$IMAGES" ]; then
  docker rmi -f $IMAGES
  echo -e "${GREEN}      ✔ Imagens removidas.${NC}"
else
  echo "      Nenhuma imagem encontrada."
fi

# 4. Remover todos os volumes
echo -e "${YELLOW}[4/7] Removendo todos os volumes...${NC}"
VOLUMES=$(docker volume ls -q)
if [ -n "$VOLUMES" ]; then
  docker volume rm -f $VOLUMES
  echo -e "${GREEN}      ✔ Volumes removidos.${NC}"
else
  echo "      Nenhum volume encontrado."
fi

# 5. Remover todas as redes customizadas
echo -e "${YELLOW}[5/7] Removendo redes customizadas...${NC}"
docker network prune -f
echo -e "${GREEN}      ✔ Redes removidas.${NC}"

# 6. Limpar cache de build do BuildKit
echo -e "${YELLOW}[6/7] Limpando cache de build...${NC}"
docker builder prune -af
echo -e "${GREEN}      ✔ Cache de build limpo.${NC}"

# 7. Limpeza geral com system prune (garante que nada ficou para trás)
echo -e "${YELLOW}[7/7] Executando limpeza geral (system prune)...${NC}"
docker system prune -af --volumes
echo -e "${GREEN}      ✔ Limpeza geral concluída.${NC}"

echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════╗"
echo -e "║   ✅  Docker limpo! Estado: zero absoluto.   ║"
echo -e "╚══════════════════════════════════════════════╝${NC}"
echo ""

# Exibe o estado final
echo -e "${CYAN}📊 Estado atual do Docker:${NC}"
echo "── Containers: $(docker ps -aq | wc -l)"
echo "── Imagens:    $(docker images -q | wc -l)"
echo "── Volumes:    $(docker volume ls -q | wc -l)"
echo "── Redes:      $(docker network ls -q | wc -l) (apenas as padrão)"
echo ""echo -e "${GREEN}Tudo limpo! Você pode começar do zero agora.${NC}"