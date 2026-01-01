#!/bin/bash

echo "Aether Engine - Character Creation API (Sans BDD)"
echo "================================================="

# Démarrer le serveur en arrière-plan
echo "Démarrage du serveur..."
go run cmd/main.go &
SERVER_PID=$!

# Attendre que le serveur soit prêt
sleep 3

echo ""
echo "Serveur démarré sur le port 8080"
echo ""

# Fonction pour appeler l'API
call_api() {
    echo "$1"
    echo "---"
    curl -s "$2" | jq '.' 2>/dev/null || curl -s "$2"
    echo ""
    echo ""
}

# Test des endpoints fonctionnels
call_api "1. LISTER LES CLASSES DISPONIBLES" "http://localhost:8080/aether/v1/characters/classes"

call_api "2. DÉTAILS DE LA CLASSE MAGE" "http://localhost:8080/aether/v1/characters/classes/MAGE"

call_api "3. TOUS LES TEMPLATES DE CLASSES" "http://localhost:8080/aether/v1/characters/templates"

echo "CRÉATION D'UN MAGE (EN MÉMOIRE UNIQUEMENT)"
echo "---"
curl -s -X POST "http://localhost:8080/aether/v1/characters/create" \
  -H "Content-Type: application/json" \
  -d '{
    "nom": "Gandalf le Gris",
    "classe": "MAGE",
    "position": {"x": 2, "y": 3}
  }' | jq '.' 2>/dev/null || curl -s -X POST "http://localhost:8080/aether/v1/characters/create" \
  -H "Content-Type: application/json" \
  -d '{
    "nom": "Gandalf le Gris",
    "classe": "MAGE", 
    "position": {"x": 2, "y": 3}
  }'
echo ""
echo ""

echo "CRÉATION D'UN PALADIN (EN MÉMOIRE UNIQUEMENT)"
echo "---"
curl -s -X POST "http://localhost:8080/aether/v1/characters/create" \
  -H "Content-Type: application/json" \
  -d '{
    "nom": "Arthas le Justicier",
    "classe": "PALADIN",
    "position": {"x": 1, "y": 1}
  }' | jq '.' 2>/dev/null || curl -s -X POST "http://localhost:8080/aether/v1/characters/create" \
  -H "Content-Type: application/json" \
  -d '{
    "nom": "Arthas le Justicier",
    "classe": "PALADIN",
    "position": {"x": 1, "y": 1}
  }'
echo ""
echo ""

echo "CRÉATION D'UN ARCHER (EN MÉMOIRE UNIQUEMENT)"
echo "---"  
curl -s -X POST "http://localhost:8080/aether/v1/characters/create" \
  -H "Content-Type: application/json" \
  -d '{
    "nom": "Legolas le Précis",
    "classe": "ARCHER",
    "position": {"x": 3, "y": 2}
  }' | jq '.' 2>/dev/null || curl -s -X POST "http://localhost:8080/aether/v1/characters/create" \
  -H "Content-Type: application/json" \
  -d '{
    "nom": "Legolas le Précis", 
    "classe": "ARCHER",
    "position": {"x": 3, "y": 2}
  }'
echo ""
echo ""

echo "IMPORTANT : PERSONNAGES CRÉÉS EN MÉMOIRE UNIQUEMENT"
echo "   → Pas de sauvegarde en base de données"
echo "   → Angular devra les envoyer à votre API BDD externe"
echo ""

echo "ENDPOINTS DISPONIBLES (SANS BDD):"
echo "  GET  /aether/v1/characters/classes      - Liste des classes"
echo "  GET  /aether/v1/characters/classes/{id} - Détails d'une classe"  
echo "  GET  /aether/v1/characters/templates    - Tous les templates"
echo "  POST /aether/v1/characters/create       - Créer un personnage (en mémoire)"
echo ""

echo "ENDPOINTS NON IMPLÉMENTÉS (À FAIRE AVEC API BDD EXTERNE):"
echo "  GET    /aether/v1/characters/           - Tous les personnages"
echo "  GET    /aether/v1/characters/{id}       - Un personnage par ID"
echo "  PUT    /aether/v1/characters/{id}       - Mettre à jour"
echo "  DELETE /aether/v1/characters/{id}       - Supprimer"
echo "  GET    /aether/v1/characters/player/{id} - Par joueur"
echo ""

# Nettoyer
echo "Arrêt du serveur..."
kill $SERVER_PID 2>/dev/null

echo "Demo terminée !"
echo "Voir API_PERSONNAGES_SANS_BDD.md pour la documentation Angular"