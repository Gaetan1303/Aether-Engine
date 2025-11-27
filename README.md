# Aether-Engine (Go + Gin)

Ce serveur est le cœur du projet Fantasy Tower, gérant toute la logique métier du jeu : combats tactiques au tour par tour, progression des personnages, système de quêtes, économie et règles complexes de gameplay. Conçu selon une architecture microservices scalable, Aether Engine garantit la cohérence des règles et la robustesse du gameplay.

---

## Vue d'ensemble

**Aether Engine** orchestre la logique métier, les règles de jeu et les mécaniques de combat, inspiré des RPG tactiques comme *Final Fantasy Tactics Advance*. Il traite les combats au tour par tour, la progression des quêtes, l'économie et l'évolution des personnages à travers une architecture microservices.

Ce moteur fonctionne comme **serveur autoritatif** au sein de l'écosystème Fantasy Tower, garantissant la cohérence du jeu, l'application des règles métier et le calcul des résultats de gameplay en temps réel.

---

### Responsabilités principales

- **Logique de combat** : Calculs de dégâts, résolution des actions, gestion des effets de statut
- **Règles de jeu** : Validation des actions, application des contraintes métier
- **Progression** : Gestion de l'expérience, montée de niveau, déverrouillage de compétences
- **Économie** : Système de craft, gestion des récompenses, équilibrage des ressources
- **Quêtes** : Orchestration des objectifs, validation des conditions de complétion

---

## Fonctionnalités principales

### Système de combat tactique
- Combats au tour par tour sur grille tactique
- Gestion du terrain avec altitude et types de cases
- Système d'initiative basé sur la vitesse (CT/ATB)
- Calcul des dégâts avec formules complexes (ATK, DEF, éléments, critique)
- Effets de statut (poison, silence, hâte, berserk, etc.)
- Zones d'effet et compétences à portée variable
- Système de réactions et contre-attaques

### Gestion des personnages
- Système de classes avec progression et spécialisations
- Gestion des compétences et capacités
- Équipement et inventaire avec contraintes par classe
- Statistiques dynamiques (HP, MP, ATK, DEF, SPD, MAG, RES)
- Talents passifs et arbres de compétences

### Système de quêtes
- Quêtes principales et secondaires
- Objectifs dynamiques et conditions de succès
- Récompenses paramétrables (XP, or, objets)
- Système de progression narrative

### Économie et craft
- Gestion de la monnaie et des ressources
- Système de craft avec recettes
- Boutiques et marchands
- Échanges entre joueurs (API Chat)

---

## Technologies utilisées

### Stack principal
- **Langage** : Go 1.21+
- **Framework HTTP** : Gin
- **Architecture** : Microservices, Domain-Driven Design (DDD)
- **Message Broker** : Kafka / RabbitMQ
- **Cache** : Redis
- **Base de données** : PostgreSQL (relationnel), MongoDB (flexible)

### Librairies & outils
- **Validation** : go-playground/validator
- **Tests** : Go test, testify
- **Logging** : zap, logrus
- **Documentation API** : Swagger (swaggo)
- **Containerisation** : Docker
- **Orchestration** : Kubernetes
- **CI/CD** : GitHub Actions / GitLab CI
- **Monitoring** : Prometheus, Grafana

---

## Installation

### Prérequis
- Go >= 1.21
- Docker & Docker Compose
- PostgreSQL 14+
- Redis 7+

### Installation locale

```bash
# Cloner le repository
git clone https://github.com/votre-organisation/aether-engine.git
cd aether-engine/server

# Installer les dépendances
go mod tidy

# Lancer le serveur Gin
go run main.go
```

### Installation avec Docker

```bash
cd server
docker build -t aether-engine-go:latest .
docker run -p 8080:8080 --env-file .env aether-engine-go:latest
```

---

## Configuration

### Variables d'environnement (exemple)

```env
GIN_MODE=release
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME=fantasy_tower
DB_USER=postgres
DB_PASSWORD=your_password
REDIS_HOST=localhost
REDIS_PORT=6379
KAFKA_BROKERS=localhost:9092
JWT_SECRET=your_jwt_secret
API_KEY=your_api_key
LOG_LEVEL=debug
```

---

## Utilisation

### Démarrer le serveur

```bash
# Mode développement
go run main.go

# Mode production
GIN_MODE=release go run main.go
```

### Lancer les tests

```bash
go test ./...
```

---

## Structure du projet (exemple)

```
server/
├── main.go                # Point d'entrée Gin
├── handlers/              # Handlers HTTP (routes)
├── models/                # Modèles métier (entités, agrégats)
├── services/              # Logique métier
├── repository/            # Accès BDD
├── middleware/            # Middlewares Gin
├── config/                # Chargement config/env
├── tests/                 # Tests unitaires et d'intégration
├── go.mod
└── README.md
```

---

## API Documentation

Une fois le serveur démarré, la documentation Swagger est accessible à :

```
http://localhost:8080/swagger/index.html
```

---

## Domain-Driven Design

- **Agrégats principaux** :
  - Combattant
  - Grille de Combat
  - Compétence
  - Tour de Combat
- **Bounded Contexts** : Combat, Personnage, Quête, Économie
- **Event Sourcing** : Événements métier via Kafka

---

## Bonnes pratiques
- Convention de code (golangci-lint)
- Tests pour chaque fonctionnalité
- Documentation des changements (CHANGELOG.md)
- Respect des principes SOLID et DDD

---

## Licence

Ce projet est sous licence El miminette

---

## Écosystème Fantasy Tower

- **Front-End** : Angular
- **API Observer** : État joueur/monde
- **Aether Engine** : Logique métier (ce projet)
- **Middleware API** : Instances/sessions
- **Chat API** : Communication joueurs
- **Big Data API** : Analyse/reporting
- **API Adapter** : Passerelle services

---
