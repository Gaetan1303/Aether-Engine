# Binaire Fabric - Serveur Aether Engine

Ce dossier contient le binaire compilé du **Serveur Fabric**, le moteur de combat tactique autoritatif d'Aether Engine.

## Qu'est-ce que Fabric ?

**Fabric** est l'exécutable du serveur de combat qui implémente toute la logique métier du système de combat au tour par tour. C'est le cœur du projet Aether Engine.

### Responsabilités du serveur

Le binaire `fabric` gère :

1. **Validation autoritaire des actions**
   - Portée des compétences (distance 3D)
   - Coûts en MP/Stamina
   - Cibles valides (Single, AoE, Row)
   - État des unités (silencée, morte, étourdie)

2. **Résolution déterministe du combat**
   - Calculs de dégâts (6 stratégies : Physical, Magical, Fixed, Hybrid, Proportional, Critical)
   - Application des effets (Poison, Haste, Shield, etc.)
   - Gestion des statuts (durée, stack, immunité)
   - Système de tour par tour

3. **Persistance événementielle (Event Sourcing)**
   - Event Store PostgreSQL (source de vérité immuable)
   - Événements : `CombatDemarre`, `ActionExecutee`, `DegatsInfliges`, etc.
   - Reconstruction de l'état à partir des événements

4. **Publication d'événements**
   - Event Bus Kafka pour notifier les autres services
   - 40+ types d'événements JSON
   - Permet aux services Observer, BigData, etc. de réagir

5. **API REST**
   - `POST /api/v1/combats` - Démarrer un combat
   - `GET /api/v1/combats/:id` - Obtenir l'état d'un combat
   - `POST /api/v1/combats/:id/actions` - Exécuter une action
   - `POST /api/v1/combats/:id/tour-suivant` - Passer au tour suivant
   - `POST /api/v1/combats/:id/terminer` - Terminer un combat

## Contenu du binaire

Le fichier `fabric` (37 Mo) contient :

- **Architecture DDD complète**
  - Agrégats : Combat, Unite, Equipe
  - Value Objects : Position, Stats, UnitID, Statut, Competence
  - 11/12 Design Patterns GoF (92%)

- **Couche Application**
  - CombatEngine (orchestration des use cases)
  - Commandes CQRS
  - Event Handlers

- **Infrastructure**
  - PostgreSQL Event Store
  - Kafka Event Publisher
  - Gin HTTP Router

- **API Layer**
  - REST Handlers
  - Validation des requêtes
  - Gestion des erreurs

## Utilisation

### Lancer le serveur

```bash
# Avec configuration par défaut
./bin/fabric

# Avec variables d'environnement personnalisées
DATABASE_URL=postgres://user:pass@localhost:5432/aether_db \
KAFKA_BROKERS=localhost:9092 \
PORT=8080 \
./bin/fabric
```

### Configuration requise

**Variables d'environnement :**

```env
# Serveur
PORT=8080                                    # Port HTTP (défaut: 8080)

# PostgreSQL (Event Store)
DATABASE_URL=postgres://test:test@localhost:5432/aether_test?sslmode=disable

# Kafka (Event Bus)
KAFKA_BROKERS=localhost:9092                 # Brokers Kafka (défaut: localhost:9092)
KAFKA_TOPIC=combat-events                    # Topic événements (défaut: combat-events)
```

**Services requis :**
- PostgreSQL 15+ (pour l'Event Store)
- Kafka (pour la publication d'événements)

### Logs au démarrage

Quand le serveur démarre correctement, vous verrez :

```
Connexion PostgreSQL établie
Event Publisher Kafka créé
Serveur Fabric démarré sur le port 8080
```

### Tester le serveur

```bash
# Health check
curl http://localhost:8080/ping
# Réponse: {"message":"pong"}

# Démarrer un combat (exemple)
curl -X POST http://localhost:8080/api/v1/combats \
  -H "Content-Type: application/json" \
  -d '{
    "equipes": [...],
    "grille": {...}
  }'
```

## Recompilation

Si vous modifiez le code source, recompilez avec :

```bash
# Depuis la racine du projet
go build -o bin/fabric ./cmd/fabric

# Vérifier le binaire
ls -lh bin/fabric
file bin/fabric
```

## Architecture technique

**Langage :** Go 1.24+

**Frameworks & Dépendances principales :**
- `gin-gonic/gin` - Framework web HTTP
- `jackc/pgx/v5` - Driver PostgreSQL
- `segmentio/kafka-go` - Client Kafka
- `google/uuid` - Génération UUID
- `stretchr/testify` - Framework de tests

**Patterns implémentés :**
- Domain-Driven Design (DDD)
- Event Sourcing / CQRS
- Strategy Pattern (calculs de dégâts)
- Singleton Pattern (génération d'IDs)
- Factory Pattern (création de calculateurs)
- Repository Pattern (Event Store)

## Notes importantes

- **Déterminisme** : Le serveur est 100% déterministe (pas de random non seedé)
- **Thread-safe** : Toutes les opérations sont thread-safe
- **Performance** : 
  - Génération d'IDs : 3.2M ops/sec
  - Calculs de dégâts : optimisés avec Strategy Pattern
- **Production-ready** : Event Store, validation stricte, gestion d'erreurs complète

## Déploiement

Le binaire `fabric` peut être déployé sur n'importe quel serveur Linux 64-bit sans installer Go :

```bash
# Copier le binaire sur le serveur
scp bin/fabric user@server:/opt/aether-engine/

# Sur le serveur
chmod +x /opt/aether-engine/fabric
./fabric
```

Pour un déploiement en production, utilisez :
- Docker (avec image Go alpine)
- Kubernetes + Helm
- Systemd service pour gestion automatique

## Plus d'informations

- Documentation complète : `/doc/`
- Code source : `/cmd/fabric/main.go`
- Tests : `go test ./...`
- Architecture : `/doc/agregats.md`
- API : Voir README principal

---

**Version :** Voir `go.mod` pour la version Go
**Dernière compilation :** Voir `ls -l bin/fabric`
