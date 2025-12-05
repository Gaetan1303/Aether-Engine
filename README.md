
# Aether-Engine – Combat Tactique Tour par Tour

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Tests](https://img.shields.io/badge/Tests-190%20passing-success)](doc/tests/)

**Aether Engine** est un moteur de combat tactique inspiré de *Final Fantasy Tactics Advance*, avec architecture Domain-Driven Design (DDD) et 190 tests automatisés.

## Démo Jouable (30 secondes)

```bash
# Démarrage rapide
./start-demo.sh

# Ou manuellement
go build -o bin/demo cmd/demo/main.go
./bin/demo
```

**Combat 2 Héros vs 2 Gobelins** avec :
-   Attaques physiques et magiques
- Compétences avec MP et cooldowns
-  Déplacement tactique sur grille 8x8
-  IA ennemie fonctionnelle
-  Système de dégâts ATK/DEF et MATK/MDEF

![Demo Screenshot](cmd/demo/screenshot.png)

**[ Guide complet de la démo](cmd/demo/README.md)**

---

##  Quick Start

### Prérequis
- Go 1.21+
- (Optionnel) Docker pour PostgreSQL/Kafka

### Installation

```bash
git clone https://github.com/Gaetan1303/Aether-Engine.git
cd Aether-Engine

# Lancer la démo CLI
./start-demo.sh

# Ou lancer les tests
go test ./... -v
```

---

##  Documentation Technique

**Aether Engine** implémente les règles métier d'un système de combat tactique au tour par tour avec architecture Domain-Driven Design.

---

## Vision du projet

### Qu'est-ce que le Serveur Fabric ?

Dans l'architecture MMO de Fantasy Tower, le **Fabric** est le service responsable de :

1. **Validation autoritaire** des actions de combat (portée, coûts, cibles)
2. **Résolution déterministe** des mécaniques de jeu (dégâts, soins, effets)
3. **Application des règles métier** via un pipeline modulaire (hooks, buffs, statuts)
4. **Persistance événementielle** (Event Sourcing) pour traçabilité et résilience
5. **Publication d'événements** vers les autres services (Kafka/Event Bus)

Le Fabric **ne gère pas** :
- L'interface utilisateur (client Angular séparé)
- La synchronisation temps réel clients (API Observer)
- Le chat et les échanges (API Chat)
- L'authentification (API Gateway)
- Les analytics (API Big Data)

---

## Architecture & Principes

### Domain-Driven Design (DDD)

Agrégats principaux documentés dans [`doc/agregats.md`](doc/agregats.md) :
- **Combat** (agrégat racine) : Gère le cycle de vie d'une instance de combat
- **Unite** : Représente un participant (joueur ou PNJ)
- **Equipe** : Regroupe plusieurs unités
- **GrilleDeCombat** : Grille tactique 3D (X, Y, Z)
- **Competence** (Value Object) : Définition immuable d'une compétence

### Event Sourcing / CQRS

Architecture documentée dans [`doc/bases_donnees/README.md`](doc/bases_donnees/README.md) :

```
Command (POST /actions) → Agrégat → Événements → Event Store (PostgreSQL)
                                         ↓
                                    Event Bus (Kafka)
                                         ↓
                              ┌──────────┴──────────┐
                              ↓                     ↓
                        Projections           Autres Services
                     (PostgreSQL + Redis)    (Observer, BigData)
```

- **Event Store** : Source de vérité immuable (append-only)
- **Projections** : Modèles de lecture optimisés (dénormalisés)
- **Cache Redis** : État temps réel des combats actifs

---

## État actuel du projet

### Ce qui est fait

| Composant | État | Documentation |
|-----------|------|---------------|
| **Architecture DDD** | Documentée + Implémentée | [`doc/agregats.md`](doc/agregats.md) |
| **Value Objects** | Implémentés + Testés | [`internal/shared/domain/`](internal/shared/domain/) |
| - Position (3D) | 100% | [`doc/tests/position/`](doc/tests/position/) |
| - Statistiques | 100% | [`doc/tests/stats/`](doc/tests/stats/) |
| - UnitID | 100% | [`doc/tests/unitID/`](doc/tests/unitID/) |
| - Statut | 100% | [`doc/tests/statut/`](doc/tests/statut/) |
| **Agrégats Domain** | Implémentés | [`internal/combat/domain/`](internal/combat/domain/) |
| - Combat (Aggregate Root) | 100% | [`internal/combat/domain/combat.go`](internal/combat/domain/combat.go) |
| - Unite (Entity) | 100% | [`internal/combat/domain/unite.go`](internal/combat/domain/unite.go) |
| - Equipe (Entity) | 100% | [`internal/combat/domain/equipe.go`](internal/combat/domain/equipe.go) |
| - Competence (Value Object) | 100% | [`internal/combat/domain/competence.go`](internal/combat/domain/competence.go) |
| **Design Patterns** | 11/12 (92%) | [`doc/tests/domain/`](doc/tests/domain/) |
| - Strategy Pattern | 100% (6 strategies) | [`doc/tests/domain/STRATEGY_PATTERN_IMPLEMENTED.md`](doc/tests/domain/STRATEGY_PATTERN_IMPLEMENTED.md) |
| - Singleton Pattern | 100% (ID Generator) | [`doc/tests/domain/SINGLETON_PATTERN_IMPLEMENTED.md`](doc/tests/domain/SINGLETON_PATTERN_IMPLEMENTED.md) |
| **Event Store** | Implémenté | [`internal/combat/infrastructure/event_store.go`](internal/combat/infrastructure/event_store.go) |
| **Combat Engine** | Implémenté | [`internal/combat/application/combat_engine.go`](internal/combat/application/combat_engine.go) |
| **API REST** | Implémentée | [`api/handlers/combat_handler.go`](api/handlers/combat_handler.go) |
| **Kafka Publisher** | Implémenté | [`pkg/eventbus/kafka_publisher.go`](pkg/eventbus/kafka_publisher.go) |
| **Tests PostgreSQL** | 14/14 passed | [`doc/tests/bases_donnees/`](doc/tests/bases_donnees/) |
| **Machines d'états** | Documentées | [`doc/machines_etats/`](doc/machines_etats/) |
| **Hooks Fabric** | Documentés | [`doc/tour_unite_hooks_integres.md`](doc/tour_unite_hooks_integres.md) |
| **40+ Types d'événements** | Spécifiés | [`doc/matrice_evenements.md`](doc/matrice_evenements.md) |

### En cours / À faire (Phase actuelle : P2 → P3)

| Composant | Priorité | État |
|-----------|----------|------|
| **Agrégats Go** (Combat, Unite, Equipe) | P0 | FAIT |
| **Event Store (implémentation)** | P0 | FAIT |
| **Use Cases** (DemarrerCombat, ExecuterAction) | P0 | FAIT |
| **API REST** (endpoints combat) | P0 | FAIT |
| **Kafka Publisher** | P1 | FAIT |
| **Design Patterns GoF** | P1 | 11/12 (92%) |
| **Projections (handlers)** | P1 | EN COURS |
| **Pipeline Fabric** (hooks + effets) | P1 | EN COURS |
| **Pathfinding A*** | P2 | À FAIRE |
| **Turn Manager** | P2 | À FAIRE |
| **Redis Cache** | P2 | À FAIRE |
| **Builder Pattern** | P3 | OPTIONNEL |

---

## Responsabilités du Fabric



### Ce que fait le Fabric

1. **Validation déterministe des actions**
   - Portée de compétence (Manhattan/Euclidienne 3D)
   - Coûts en MP/Stamina
   - Cibles valides (Single, AoE, Row)
   - État de l'unité (silencée, morte, étourdie)

2. **Résolution des actions**
   - Calculs de dégâts (formules ATK/DEF/SPD)
   - Application des effets (Poison, Haste, Shield)
   - Gestion des statuts (durée, stack, immunité)
   - Système ATB (Active Time Battle)

3. **Persistance événementielle**
   - Event Store PostgreSQL (source de vérité immuable)
   - Snapshots (optimisation reconstruction)
   - Projections read-only (état combat courant)

4. **Publication d'événements**
   - Kafka publisher (`CombatDemarre`, `ActionExecutee`, `DegatsInfliges`, etc.)
   - Contract: 40+ types d'événements JSON
   - Permet aux autres services de réagir (Observer, BigData)

5. **API REST pour commandes**
   - `POST /api/v1/combats` (démarrer combat)
   - `POST /api/v1/combats/:id/actions` (exécuter action)
   - `GET /api/v1/combats/:id` (état combat via projection)

### Ce que le Fabric NE fait PAS

- Interface utilisateur → Client Angular séparé
- Synchronisation temps réel → API Observer (écoute Kafka → WebSocket)
- Authentification → API Gateway
- Chat/Échanges → API Chat
- Analytics → API Big Data

---

## Stack Technique

| Composant | Technologie | Justification |
|-----------|-------------|---------------|
| **Backend** | Go 1.23+ | Performance, concurrence, typage fort |
| **Framework Web** | Gin | Léger, rapide, idiomatique Go |
| **Event Store** | PostgreSQL 15 (pgx/v5) | ACID, requêtes temporelles, robuste |
| **Cache** | Redis 7 | Latence sub-ms, pub/sub natif |
| **Event Bus** | Kafka (à implémenter) | Découplage, scalabilité, replay events |
| **Tests** | Testify + pgx/v5 | Assertions idiomatiques + tests PostgreSQL |
| **Logging** | Zap (prévu) | Structured logging, performance |
| **Metrics** | Prometheus (prévu) | Standard Cloud Native |
| **Deployment** | Kubernetes + Helm (prévu) | Scalabilité, rolling updates |

---

## Installation & Configuration

### Prérequis

- Go 1.23+
- PostgreSQL 15+
- Redis 7+ (optionnel pour cache)
- Make (optionnel)

### Installation locale

```bash
# Cloner le repository
git clone https://github.com/Gaetan1303/Aether-Engine.git
cd Aether-Engine

# Installer les dépendances Go
go mod download

# Configurer PostgreSQL de test
sudo -u postgres createdb aether_test
sudo -u postgres psql -c "CREATE USER test WITH PASSWORD 'test';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE aether_test TO test;"

# Lancer les tests unitaires (Design Patterns)
go test ./doc/tests/domain -v

# Lancer les tests d'intégration (PostgreSQL)
go test ./doc/tests/bases_donnees -v

# Compiler le serveur
go build -o bin/fabric ./cmd/fabric

# Lancer le serveur
./bin/fabric
# Ou avec variables d'environnement personnalisées
DATABASE_URL=postgres://test:test@localhost:5432/aether_test PORT=8080 ./bin/fabric
```

### Variables d'environnement

```env
# Serveur
PORT=8080                                    # Port du serveur (défaut: 8080)

# PostgreSQL (Event Store + Projections)
DATABASE_URL=postgres://test:test@localhost:5432/aether_test?sslmode=disable

# Kafka (Event Bus)
KAFKA_BROKERS=localhost:9092                 # Brokers Kafka (défaut: localhost:9092)
KAFKA_TOPIC=combat-events                    # Topic des événements (défaut: combat-events)

# Redis (Cache) - À venir
REDIS_HOST=localhost
REDIS_PORT=6379

# Observabilité - À venir
LOG_LEVEL=info
METRICS_PORT=9090
```

---

## Structure du Projet (Architecture Hexagonale)

```
Aether-Engine/
├── cmd/
│       └─ main.go                  # Point d'entrée principal
├── go.mod                           # Dépendances Go
├── go.sum                           # Checksums des dépendances
├── internal/                        # Code non exportable
│   ├── combat/                      # Bounded Context Combat
│   │   ├── domain/                  # IMPLÉMENTÉ
│   │   │   ├── combat.go            # Agrégat racine Combat
│   │   │   ├── unite.go             # Entité Unite
│   │   │   ├── equipe.go            # Entité Equipe
│   │   │   ├── competence.go        # Value Object Competence
│   │   │   ├── damage_calculator.go # Strategy Pattern (6 calculateurs)
│   │   │   ├── events.go            # Événements domain
│   │   │   └── enums.go             # Énumérations
│   │   ├── application/             # IMPLÉMENTÉ
│   │   │   ├── combat_engine.go     # Moteur de combat (Use Cases)
│   │   │   └── commands.go          # Commandes CQRS
│   │   └── infrastructure/          # IMPLÉMENTÉ
│   │       └── event_store.go       # Repository Event Store PostgreSQL
│   └── shared/                      # Code partagé
│       └── domain/                  # IMPLÉMENTÉ
│           ├── value_objects.go     # Position, Stats, UnitID, Statut
│           ├── id_generator.go      # Singleton Pattern (génération IDs)
│           ├── types.go             # Types de base
│           └── interfaces.go        # Interfaces partagées
├── api/                             # API Layer
│   └── handlers/
│       └── combat_handler.go        # Endpoints REST (IMPLÉMENTÉ)
├── pkg/                             # Packages exportables
│   └── eventbus/
│       └── kafka_publisher.go       # Publisher Kafka (IMPLÉMENTÉ)
├── bin/                             # Binaires compilés
│   └── fabric                       # Exécutable serveur
├── doc/                             # Documentation complète
│   ├── agregats.md                  # Définition des agrégats
│   ├── bases_donnees/               # Schémas Event Store + Projections
│   ├── machines_etats/              # Machines d'états du combat
│   ├── diagrammes_*/                # Diagrammes Mermaid
│   └── tests/                       # Tests et documentation
│       ├── domain/                  # Tests Design Patterns
│       │   ├── strategy_pattern_test.go        # Tests Strategy (8/8 PASS)
│       │   ├── id_generator_test.go            # Tests Singleton (14/14 PASS)
│       │   ├── STRATEGY_PATTERN_IMPLEMENTED.md
│       │   └── SINGLETON_PATTERN_IMPLEMENTED.md
│       ├── bases_donnees/           # Tests PostgreSQL (14/14 PASS)
│       ├── position/                # Tests Position 3D
│       ├── stats/                   # Tests Statistiques
│       ├── unitID/                  # Tests UnitID
│       └── statut/                  # Tests Statut
└── Plan/                            # Documentation planning
    ├── semaine_1.md
    └── architecture_globale.md
```

---

## Tests

### Tests Design Patterns

**22/22** tests de patterns réussis :

```bash
# Strategy Pattern (6 calculateurs de dégâts)
go test -v ./doc/tests/domain/strategy_pattern_test.go
# Résultat: 8/8 PASS

# Singleton Pattern (ID Generator)
go test -v ./doc/tests/domain/id_generator_test.go
# Résultat: 14/14 PASS

# Benchmarks ID Generator
go test -bench=. -benchmem ./doc/tests/domain/id_generator_test.go
# Performance: 3.2M IDs/sec, 371ns/op
```

### Tests Value Objects

**100%** des Value Objects testés :

```bash
# Position 3D
go test -v ./doc/tests/position/position_test.go

# Statistiques
go test -v ./doc/tests/stats/stats_test.go

# UnitID
go test -v ./doc/tests/unitID/unitid_test.go

# Statut
go test -v ./doc/tests/statut/statut_test.go
```

### Tests d'Intégration (PostgreSQL)

**14/14** tests Event Store + Projections :

```bash
# Tous les tests PostgreSQL
go test ./doc/tests/bases_donnees -v

# Event Store uniquement
go test ./doc/tests/bases_donnees -v -run "TestInsert|TestOptimistic|TestSnapshot|TestReconstruct|TestQuery|TestTransactional"

# Projections uniquement
go test ./doc/tests/bases_donnees -v -run "TestCombat.*Projection|TestProjectionIdempotence"
```

Documentation détaillée : [`doc/tests/bases_donnees/README.md`](doc/tests/bases_donnees/README.md)

---

## Documentation

### Documentation Centrale

- **[`doc/agregats.md`](doc/agregats.md)** : Définition des agrégats DDD
- **[`doc/presentation.md`](doc/presentation.md)** : Vision globale du Fabric
- **[`doc/feuille_de_route.md`](doc/feuille_de_route.md)** : Roadmap P1 → P6
- **[`doc/phase_1_domaine_metier.md`](doc/phase_1_domaine_metier.md)** : Phase actuelle

### Architecture Event Sourcing

- **[`doc/bases_donnees/README.md`](doc/bases_donnees/README.md)** : Vue d'ensemble
- **[`doc/bases_donnees/event_store.md`](doc/bases_donnees/event_store.md)** : Event Store
- **[`doc/bases_donnees/projections_combat.md`](doc/bases_donnees/projections_combat.md)** : Projections Combat
- **[`doc/matrice_evenements.md`](doc/matrice_evenements.md)** : 40+ types d'événements

### Machines d'États

- **[`doc/machines_etats/combat_core_p2.md`](doc/machines_etats/combat_core_p2.md)** : Machine d'états Combat
- **[`doc/machines_etats/tour.md`](doc/machines_etats/tour.md)** : Machine d'états Tour
- **[`doc/machines_etats/instance_combat.md`](doc/machines_etats/instance_combat.md)** : Instance de Combat

### Hooks & Pipeline

- **[`doc/tour_unite_hooks_integres.md`](doc/tour_unite_hooks_integres.md)** : Système de hooks Fabric

---

## Roadmap (Phases DDD)

| Phase | Objectif | État | Progrès |
|-------|----------|------|--------|
| **P1** | Fondations & Contrats | TERMINÉE | 100% |
| **P2** | Cœur Combat Déterministe | EN COURS | 75% |
| **P3** | Fabric & Résolution | EN COURS | 40% |
| **P4** | Résilience & Event Sourcing | IMPLÉMENTÉE | 80% |
| **P5** | API & Scalabilité | IMPLÉMENTÉE | 70% |
| **P6** | Production-Ready | EN COURS | 30% |

### Réalisations majeures
- Architecture DDD complète (Agrégats, Entités, Value Objects)
- 11/12 Design Patterns GoF (92%)
- Event Store PostgreSQL fonctionnel
- API REST opérationnelle (5 endpoints)
- Strategy Pattern (6 calculateurs de dégâts)
- Singleton Pattern (génération d'IDs thread-safe)
- Kafka Event Publisher intégré

Détails : [`doc/feuille_de_route.md`](doc/feuille_de_route.md)

---

## API Endpoints

Le serveur Fabric expose les endpoints REST suivants :

### Combat Management

```http
# Démarrer un nouveau combat
POST /api/v1/combats
Content-Type: application/json

{
  "equipes": [...],
  "grille": {...}
}

# Obtenir l'état d'un combat
GET /api/v1/combats/:id

# Exécuter une action
POST /api/v1/combats/:id/actions
Content-Type: application/json

{
  "uniteID": "unit_123",
  "typeAction": "ATTAQUE",
  "cibleID": "unit_456"
}

# Passer au tour suivant
POST /api/v1/combats/:id/tour-suivant

# Terminer un combat
POST /api/v1/combats/:id/terminer
```

### Health Check

```http
# Vérifier l'état du serveur
GET /ping

Response: {"message": "pong"}
```

### Démarrage du serveur

```bash
# Avec configuration par défaut
./bin/fabric

# Le serveur démarre sur http://localhost:8080
# Logs attendus:
# Connexion PostgreSQL établie
# Event Publisher Kafka créé
# Serveur Fabric démarré sur le port 8080
```

---

## Contribution

Ce projet suit les principes **Domain-Driven Design (DDD)** et **Event Sourcing**.

### Règles de contribution

1. **Déterminisme strict** : Pas d'horloge système, pas de random non seedé
2. **Event Sourcing** : Toute modification passe par un événement
3. **Tests obligatoires** : Chaque agrégat/use case doit avoir ses tests
4. **Documentation à jour** : Mettre à jour `/doc` si modification du domaine

---

## Licence

Projet sous licence de El Miminette 

---

## Écosystème Fantasy Tower

Le **Serveur Fabric (Aether Engine)** fait partie d'une architecture MMO plus large :

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │    │     API     │    │     API     │
│   Angular   │◄───┤   Observer  │◄───┤   Gateway   │
└─────────────┘    └─────────────┘    └─────────────┘
                          ▲                   ▲
                          │ Kafka Events      │ REST
                          │                   │
                   ┌──────┴───────────────────┴──────┐
                   │   AETHER ENGINE (Fabric)  │
                   │   - Validation autoritaire      │
                   │   - Résolution déterministe     │
                   │   - Event Store PostgreSQL      │
                   │   - Projections + Cache Redis   │
                   └─────────────────────────────────┘
```

**Services connexes** (hors scope Fabric) :
- **API Observer** : Synchronisation état temps réel
- **API Gateway** : Authentification, rate limiting, routing
- **API Chat** : Messages entre joueurs
- **API Big Data** : Analytics et métriques


