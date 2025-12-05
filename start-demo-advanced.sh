#!/bin/bash

# Couleurs
GREEN='\033[0;32m'
CYAN='\033[0;36m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${CYAN}╔════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║                                                ║${NC}"
echo -e "${CYAN}║     🏰  AETHER ENGINE - DEMO AVANCÉE  ⚔️      ║${NC}"
echo -e "${CYAN}║                                                ║${NC}"
echo -e "${CYAN}╚════════════════════════════════════════════════╝${NC}"
echo ""

# Vérifier si Go est installé
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go n'est pas installé!${NC}"
    echo "Téléchargez-le depuis: https://golang.org/dl/"
    exit 1
fi

echo -e "${YELLOW}🔨 Compilation de la démo avancée...${NC}"
go build -o bin/demo-advanced cmd/demo-advanced/main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Compilation réussie!${NC}"
    echo ""
    echo -e "${CYAN}🎮 Lancement du jeu...${NC}"
    echo ""
    ./bin/demo-advanced
else
    echo -e "${RED}❌ Erreur de compilation${NC}"
    exit 1
fi
