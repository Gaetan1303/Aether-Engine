# Aether-Engine
Ce serveur est le coeur du projet Tower fantasy, gère l'intégralité de la logique métier du jeu : combats tactiques au tour par tour, progression des personnages, système de quêtes, économie et règles complexes de gameplay. Conçu selon une architecture microservices scalable, Aether Engine garantit la cohérence des règles.


# Aether Engine

Le moteur de logique métier et de règles de jeu pour Fantasy Tower - Un système de combat tactique MMO RPG tour par tour

[![Node.js Version](https://img.shields.io/badge/node-%3E%3D18.0.0-brightgreen)](https://nodejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-blue)](https://www.typescriptlang.org/)

---

## Vue d'ensemble

**Aether Engine** est le cerveau de Fantasy Tower, gérant toute la logique métier, les règles de jeu et les mécaniques de combat. Inspiré des RPG tactiques comme *Final Fantasy Tactics Advance*, il traite les combats au tour par tour, la progression des quêtes, les systèmes économiques et l'évolution des personnages à travers une architecture microservices.

Ce moteur opère en tant que **Serveur Fabric** au sein de l'écosystème Fantasy Tower, garantissant la cohérence du jeu, appliquant les règles métier et calculant les résultats de gameplay complexes en temps réel.

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
- Échanges entre joueurs (via API Chat)

---

## Technologies utilisées

### Stack principal

- **Runtime** : Node.js 18+
- **Langage** : TypeScript 5.0+
- **Architecture** : Microservices, Domain-Driven Design (DDD)
- **Message Broker** : Kafka / RabbitMQ
- **Cache** : Redis
- **Base de données** : PostgreSQL (données relationnelles), MongoDB (données flexibles)

### Frameworks et librairies

- **Framework** : NestJS ou Express
- **Validation** : Zod / Joi
- **ORM** : TypeORM / Prisma
- **Testing** : Jest, Supertest
- **Logging** : Winston / Pino
- **Documentation** : Swagger / OpenAPI

### DevOps

- **Containerisation** : Docker
- **Orchestration** : Kubernetes
- **CI/CD** : GitHub Actions / GitLab CI
- **Monitoring** : Prometheus, Grafana

---

## Installation

### Prérequis

- Node.js >= 18.0.0
- npm >= 9.0.0 ou yarn >= 1.22.0
- Docker et Docker Compose (pour l'environnement de développement)
- PostgreSQL 14+
- Redis 7+

### Installation locale

```bash
# Cloner le repository
git clone https://github.com/votre-organisation/aether-engine.git
cd aether-engine

# Installer les dépendances
npm install

# Copier le fichier de configuration
cp .env.example .env

# Configurer les variables d'environnement
# Éditer .env avec vos paramètres

# Démarrer les services via Docker Compose
docker-compose up -d

# Lancer les migrations de base de données
npm run migration:run

# Démarrer le serveur en mode développement
npm run dev
```

### Installation avec Docker

```bash
# Build de l'image
docker build -t aether-engine:latest .

# Lancer le container
docker run -p 3000:3000 --env-file .env aether-engine:latest
```

---

## Configuration

### Variables d'environnement

```env
# Application
NODE_ENV=development
PORT=3000
API_VERSION=v1

# Base de données
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=fantasy_tower
DATABASE_USER=postgres
DATABASE_PASSWORD=your_password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_CLIENT_ID=aether-engine
KAFKA_GROUP_ID=aether-engine-group

# Services externes
OBSERVER_API_URL=http://localhost:3001
CHAT_API_URL=http://localhost:3002
MIDDLEWARE_API_URL=http://localhost:3003

# Sécurité
JWT_SECRET=your_jwt_secret
API_KEY=your_api_key

# Logs
LOG_LEVEL=debug
```

---

## Utilisation

### Démarrer le serveur

```bash
# Mode développement avec hot-reload
npm run dev

# Mode production
npm run build
npm run start:prod

# Mode debug
npm run start:debug
```

### Lancer les tests

```bash
# Tests unitaires
npm run test

# Tests d'intégration
npm run test:e2e

# Coverage
npm run test:cov

# Tests en mode watch
npm run test:watch
```

### Scripts utiles

```bash
# Linter
npm run lint

# Formatter
npm run format

# Migrations
npm run migration:generate
npm run migration:run
npm run migration:revert

# Seeds
npm run seed
```

---

## Structure du projet

```
aether-engine/
├── src/
│   ├── combat/              # Module de combat tactique
│   │   ├── domain/          # Entités, agrégats, value objects
│   │   ├── application/     # Use cases, services applicatifs
│   │   ├── infrastructure/  # Repositories, adapters
│   │   └── presentation/    # Controllers, DTOs
│   ├── character/           # Module de gestion des personnages
│   ├── quest/               # Module de quêtes
│   ├── economy/             # Module économique
│   ├── shared/              # Code partagé (utils, interfaces)
│   │   ├── domain/          # Domain primitives
│   │   ├── infrastructure/  # Database, messaging
│   │   └── application/     # Services transverses
│   └── main.ts              # Point d'entrée
├── test/                    # Tests e2e
├── migrations/              # Migrations de BDD
├── docs/                    # Documentation
├── docker/                  # Dockerfiles et configs
├── .github/                 # CI/CD workflows
├── docker-compose.yml
├── tsconfig.json
├── package.json
└── README.md
```

---

## API Documentation

Une fois le serveur démarré, la documentation Swagger est accessible à :

```
http://localhost:3000/api/docs
```

### Endpoints principaux

#### Combat

- `POST /api/v1/combat/initialize` - Initialiser un combat
- `POST /api/v1/combat/action` - Exécuter une action de combat
- `GET /api/v1/combat/:id/state` - Récupérer l'état d'un combat
- `POST /api/v1/combat/:id/end` - Terminer un combat

#### Personnages

- `GET /api/v1/characters/:id` - Récupérer un personnage
- `PATCH /api/v1/characters/:id` - Mettre à jour un personnage
- `POST /api/v1/characters/:id/level-up` - Faire monter de niveau
- `POST /api/v1/characters/:id/equip` - Équiper un objet

#### Quêtes

- `GET /api/v1/quests/available` - Lister les quêtes disponibles
- `POST /api/v1/quests/:id/accept` - Accepter une quête
- `POST /api/v1/quests/:id/complete` - Compléter une quête
- `GET /api/v1/quests/:id/progress` - Progression d'une quête

---

## Domain-Driven Design

Aether Engine suit les principes du **Domain-Driven Design (DDD)** :

### Agrégats principaux

1. **Combattant** : Gère l'état vital et les capacités d'une unité de combat
2. **Grille de Combat** : Garantit la cohérence spatiale du champ de bataille
3. **Compétence** : Définit les actions disponibles et leurs règles
4. **Tour de Combat** : Orchestre l'exécution atomique d'une action

### Bounded Contexts

- **Combat Context** : Tout ce qui concerne les batailles tactiques
- **Character Context** : Progression et gestion des personnages
- **Quest Context** : Système de quêtes et objectifs
- **Economy Context** : Craft, échanges, récompenses

### Event Sourcing

Les événements métier critiques sont capturés et publiés via Kafka :

- `CombatInitialized`
- `ActionExecuted`
- `CharacterLeveledUp`
- `QuestCompleted`
- `ItemCrafted`

---

### Guidelines

- Suivre les conventions de code (ESLint + Prettier)
- Écrire des tests pour les nouvelles fonctionnalités
- Documenter les changements dans le CHANGELOG.md
- Respecter les principes SOLID et DDD

---


## Licence

Ce projet est sous licence de El miminette

---


## Écosystème Fantasy Tower

Aether Engine fait partie de l'écosystème Fantasy Tower :

- **Front-End** : Interface utilisateur Angular
- **API Observer** : Maintien de l'état joueur/monde
- **Aether Engine** : Logique métier et règles (vous êtes ici)
- **Middleware API** : Gestion des instances et sessions
- **Chat API** : Communication entre joueurs
- **Big Data API** : Analyse et reporting
- **API Adapter** : Passerelle entre services

---

**Développé avec passion pour créer l'expérience MMO tactique ultime**