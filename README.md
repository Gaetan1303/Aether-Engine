
# Aether-Engine – Serveur Fabric (Combat Tactique Déterministe)

> **Note de synchronisation** :
> Les concepts d'agrégats, Value Objects, etc. sont centralisés dans `/doc/agregats.md`.  
> Les diagrammes et la documentation utilisent le nommage français, sauf pour les termes internationalement utilisés (item, Tank, DPS, Heal, etc.).

**Aether Engine** est le **Serveur Fabric** du projet Fantasy Tower : un moteur de combat tactique autoritatif, déterministe et résilient, implémentant les règles métier du système de combat au tour par tour inspiré de *Final Fantasy Tactics Advance*.

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
| **Architecture DDD** |  Documentée | [`doc/agregats.md`](doc/agregats.md) |
| **Value Objects** |  Implémentés + Testés | [`server/internal/shared/domain/`](server/internal/shared/domain/) |
| - Position (3D) |  100% | [`doc/tests/position/`](doc/tests/position/) |
| - Statistiques |  100% | [`doc/tests/stats/`](doc/tests/stats/) |
| - UnitID |  100% | [`doc/tests/unitID/`](doc/tests/unitID/) |
| - Statut |  100% | [`doc/tests/statut/`](doc/tests/statut/) |
| **Event Store (schémas)** |  Documenté | [`doc/bases_donnees/event_store.md`](doc/bases_donnees/event_store.md) |
| **Projections (schémas)** |  Documentées | [`doc/bases_donnees/projections_combat.md`](doc/bases_donnees/projections_combat.md) |
| **Tests PostgreSQL** |  14/14 passed | [`doc/tests/bases_donnees/`](doc/tests/bases_donnees/) |
| **Machines d'états** |  Documentées | [`doc/machines_etats/`](doc/machines_etats/) |
| **Hooks Fabric** |  Documentés | [`doc/tour_unite_hooks_integres.md`](doc/tour_unite_hooks_integres.md) |
| **40+ Types d'événements** |  Spécifiés | [`doc/matrice_evenements.md`](doc/matrice_evenements.md) |

### En cours / À faire (Phase actuelle : P1 → P2)

| Composant | Priorité | Effort estimé |
|-----------|----------|---------------|
| **Agrégats Go** (Combat, Unite, Equipe) | 🔴 P0 | 3-4 jours |
| **Event Store (implémentation)** | 🔴 P0 | 2-3 jours |
| **Use Cases** (DemarrerCombat, ExecuterAction) | 🔴 P0 | 3-4 jours |
| **Projections (handlers)** | 🔴 P0 | 2-3 jours |
| **API REST** (endpoints combat) | 🔴 P0 | 2-3 jours |
| **Pipeline Fabric** (hooks + effets) | 🟠 P1 | 1 semaine |
| **Kafka Publisher** | 🟠 P1 | 2-3 jours |
| **Redis Cache** | 🟡 P2 | 2-3 jours |

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
cd server
go mod download

# Configurer PostgreSQL de test
sudo -u postgres createdb aether_test
sudo -u postgres psql -c "CREATE USER test WITH PASSWORD 'test';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE aether_test TO test;"

# Lancer les tests
go test ./tests/bases_donnees -v

# Lancer le serveur (actuellement minimal)
go run main.go
```

### Variables d'environnement (futures)

```env
# Serveur
GIN_MODE=release
PORT=8080

# PostgreSQL (Event Store + Projections)
DB_HOST=localhost
DB_PORT=5432
DB_NAME=aether_engine
DB_USER=aether
DB_PASSWORD=your_password

# Redis (Cache)
REDIS_HOST=localhost
REDIS_PORT=6379

# Kafka (Event Bus)
KAFKA_BROKERS=localhost:9092

# Observabilité
LOG_LEVEL=info
METRICS_PORT=9090
```

---

## Structure du Projet (Architecture Hexagonale)

```
Aether-Engine/
├── server/
│   ├── main.go                      # Point d'entrée (actuellement minimal)
│   ├── go.mod                       # Dépendances Go
│   └── internal/                    # Code non exportable
│       ├── combat/                  # Bounded Context Combat
│       │   ├── domain/              # À IMPLÉMENTER
│       │   │   ├── combat.go        # Agrégat racine
│       │   │   ├── unite.go         # Entité Unite
│       │   │   ├── equipe.go        # Entité Equipe
│       │   │   ├── competence.go    # Value Object
│       │   │   └── grille.go        # Grille tactique 3D
│       │   ├── application/         # À IMPLÉMENTER
│       │   │   ├── demarrer_combat.go
│       │   │   ├── executer_action.go
│       │   │   └── terminer_combat.go
│       │   ├── infrastructure/      # À IMPLÉMENTER
│       │   │   ├── event_store.go   # Repository Event Store
│       │   │   ├── projections.go   # Handlers projections
│       │   │   └── kafka.go         # Publisher Kafka
│       │   └── api/                 # À IMPLÉMENTER
│       │       └── handlers.go      # Endpoints REST
│       └── shared/                  # Code partagé
│           └── domain/              # FAIT
│               ├── position.go      # Value Object Position (3D)
│               ├── stats.go         # Value Object Statistiques
│               ├── unit_id.go       # Value Object UnitID
│               └── status.go        # Value Object Statut
├── doc/                             # Documentation complète
│   ├── agregats.md                  # Définition des agrégats
│   ├── bases_donnees/               # Schémas Event Store + Projections
│   ├── machines_etats/              # Machines d'états du combat
│   ├── diagrammes_*/                # Diagrammes Mermaid
│   └── tests/                       # Documentation des tests
└── tests/                           # Tests à migrer dans server/
    └── bases_donnees/               # Tests PostgreSQL (14/14 passed)
```

---

## Tests

### Tests Unitaires (Value Objects)

**100%** des Value Objects testés :

```bash
# Position 3D
go test -v server/internal/shared/domain/position_test.go

# Statistiques
go test -v server/internal/shared/domain/stats_test.go

# UnitID
go test -v server/internal/shared/domain/unit_id_test.go

# Statut
go test -v server/internal/shared/domain/status_test.go
```

### Tests d'Intégration (PostgreSQL)

**14/14** tests Event Store + Projections :

```bash
# Tous les tests PostgreSQL
go test ./tests/bases_donnees -v

# Event Store uniquement
go test ./tests/bases_donnees -v -run "TestInsert|TestOptimistic|TestSnapshot|TestReconstruct|TestQuery|TestTransactional"

# Projections uniquement
go test ./tests/bases_donnees -v -run "TestCombat.*Projection|TestProjectionIdempotence"
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

| Phase | Objectif | État | ETA |
|-------|----------|------|-----|
| **P1** | Fondations & Contrats | 80% | Actuelle |
| **P2** | Cœur Combat Déterministe | 20% | 2-3 sem |
| **P3** | Fabric & Résolution | 0% | 2-3 sem |
| **P4** | Résilience & Event Sourcing | 0% | 2 sem |
| **P5** | API & Scalabilité | 0% | 2 sem |
| **P6** | Production-Ready | 0% | 1 sem |

Détails : [`doc/feuille_de_route.md`](doc/feuille_de_route.md)

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


