
# Aether-Engine (Go + Gin)

> **Note de synchronisation** :
> Les concepts d'agrégats, Value Objects, etc. sont centralisés dans `/doc/agregats.md`.
> Les diagrammes et la documentation utilisent le nommage français, sauf pour les termes internationalement utilisés (item, Tank, DPS, Heal, etc.).

Ce serveur est le cœur du projet Fantasy Tower, gérant toute la logique métier du jeu : combats tactiques au tour par tour, progression des personnages, système de quêtes, économie et règles complexes de gameplay. Conçu selon une architecture microservices scalable, Aether Engine garantit la cohérence des règles et la robustesse du gameplay.

---

## Vue d'ensemble

**Aether Engine** orchestre la logique métier, les règles de jeu et les mécaniques de combat, inspiré des RPG tactiques comme *Final Fantasy Tactics Advance*. Il traite les combats au tour par tour, la progression des quêtes, l'économie et l'évolution des personnages à travers une architecture microservices.

Ce moteur fonctionne comme **serveur autoritatif** au sein de l'écosystème Fantasy Tower, garantissant la cohérence du jeu, l'application des règles métier et le calcul des résultats de gameplay en temps réel.

---

### Responsabilités principales



## Fonctionnalités principales

### Système de combat tactique

### Gestion des personnages

### Système de quêtes

### Économie et craft


## Technologies utilisées

### Stack principal

### Librairies & outils


## Installation

### Prérequis

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


## Structure du projet (prévisionnelle)

server/
├── cmd/
│   └── server/
│       └── main.go              # Point d'entrée
├── internal/
│   ├── combat/                  # Bounded Context Combat
│   │   ├── domain/              # Agrégats, Entities, Value Objects
│   │   │   ├── battle.go
│   │   │   ├── unit.go
│   │   │   └── skill.go
│   │   ├── application/         # Use cases
│   │   │   └── battle_service.go
│   │   ├── infrastructure/      # Repositories, Adapters
│   │   │   ├── battle_repository.go
│   │   │   └── event_publisher.go
│   │   └── presentation/        # Handlers HTTP
│   │       └── battle_handler.go
│   ├── character/               # Bounded Context Personnage
│   ├── quest/
│   └── shared/                  # Code partagé
│       ├── domain/
│       └── infrastructure/
├── pkg/                         # Code public réutilisable
├── tests/
└── go.mod

## API Documentation

Une fois le serveur démarré, la documentation Swagger est accessible à :

```
http://localhost:8080/swagger/index.html
```


## Domain-Driven Design



## Bonnes pratiques


## Licence

Ce projet est sous licence El miminette


## Écosystème Fantasy Tower


