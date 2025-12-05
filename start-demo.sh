#!/bin/bash

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}╔══════════════════════════════════════╗${NC}"
echo -e "${CYAN}║   AETHER ENGINE - QUICK START       ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════╝${NC}"
echo ""

# Vérifier si Go est installé
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go n'est pas installé${NC}"
    echo "Installe Go depuis: https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✅ Go détecté: $(go version)${NC}"
echo ""

# Compiler la démo
echo -e "${YELLOW}🔨 Compilation de la démo...${NC}"
if go build -o bin/demo cmd/demo/main.go; then
    echo -e "${GREEN}✅ Compilation réussie!${NC}"
else
    echo -e "${RED}❌ Erreur de compilation${NC}"
    exit 1
fi

echo ""
echo -e "${CYAN}🎮 Lancement de la démo...${NC}"
echo ""

# Lancer la démo
./bin/demo
