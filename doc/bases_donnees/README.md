# Vue d'Ensemble des Bases de Données - Aether Engine

## Architecture Globale

L'Aether Engine utilise une architecture Event Sourcing avec CQRS, séparant strictement l'écriture (Command) de la lecture (Query). Cette séparation permet une scalabilité optimale et une traçabilité complète.

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENTS                                  │
│              (WebSocket, REST API, GraphQL)                     │
└────────────┬────────────────────────────────────┬───────────────┘
             │                                    │
             ▼                                    ▼
┌────────────────────────┐            ┌─────────────────────────┐
│   COMMAND HANDLERS     │            │   QUERY HANDLERS        │
│  (Agrégats + Règles)   │            │  (Read Models)          │
└────────────┬───────────┘            └────────┬────────────────┘
             │                                 │
             │ Événements                      │ Projections
             ▼                                 ▼
┌────────────────────────┐            ┌─────────────────────────┐
│   EVENT STORE          │───────────>│   PROJECTIONS           │
│   PostgreSQL           │  Diffusion │   PostgreSQL + Redis    │
│   (Write Model)        │            │   (Read Models)         │
└────────────────────────┘            └─────────────────────────┘
         Immuable                          Optimisé pour lecture
```

## 1. Event Store (Write Model)

**Base de données**: PostgreSQL

**Responsabilité**: Stockage immuable de tous les événements du système.

**Fichier**: `event_store.md`

### Caractéristiques

- **Append-only**: Les événements ne sont jamais modifiés ou supprimés
- **Ordre total**: Séquence globale garantie (`event_sequence`)
- **Source de vérité**: Toutes les données du système dérivent de ces événements
- **Audit complet**: Traçabilité totale via `correlationId`, `causationId`
- **Snapshots**: Optimisation pour agrégats volumineux

### Tables Principales

- `evenements`: Tous les événements du système (40+ types)
- `snapshots`: États périodiques des agrégats pour performance

### Usage Typique

```sql
-- Reconstruction d'un agrégat (joueur, combat, etc.)
SELECT event_type, aggregate_version, payload, timestamp
FROM evenements
WHERE aggregate_id = 'uuid-joueur'
ORDER BY aggregate_version ASC;
```

## 2. Projections Combat (Read Model)

**Base de données**: PostgreSQL + Redis

**Responsabilité**: Modèles de lecture optimisés pour le système de combat temps-réel.

**Fichier**: `projections_combat.md`

### Caractéristiques

- **Dénormalisées**: Données agrégées pour performance
- **Eventual consistency**: Mises à jour asynchrones depuis l'Event Store
- **Cache Redis**: Combats actifs en mémoire pour ultra-performance
- **TTL adaptatif**: Nettoyage automatique des combats terminés

### Tables Principales

- `instances_combat`: État global de chaque combat
- `participants_combat`: État de chaque participant (HP, MP, stamina, effets)
- `actions_combat`: Historique des actions
- `effets_statut`: Effets actifs (poison, régénération, etc.)
- `historique_combat`: Journal complet pour replay

### Cache Redis

```redis
combat:{combat_id}:state               # État du combat
combat:{combat_id}:participant:{pid}   # État d'un participant
combat:{combat_id}:turn_order          # Ordre d'initiative
combat:{combat_id}:recent_actions      # 20 dernières actions
zone:{zone_id}:active_combats          # Index par zone
```

## 3. Projections Monde (Read Model)

**Base de données**: PostgreSQL

**Responsabilité**: Catalogues d'items, compétences, quêtes, économie et état du monde.

**Fichier**: `projections_monde.md`

### Caractéristiques

- **Catalogues statiques**: Items, compétences, quêtes (peu de modifications)
- **Données dynamiques**: Ordres économiques, transactions, prix du marché
- **État global**: Monde singleton (cycle jour/nuit, saisons, boss vaincus)
- **Historique complet**: Toutes les transactions du marché

### Tables Principales

- `items`: Catalogue complet des items
- `competences`: Catalogue des compétences
- `quetes`: Catalogue des quêtes
- `etat_monde`: État global du monde (singleton)
- `ordres_economie`: Ordres d'achat/vente actifs
- `transactions_economie`: Historique des transactions
- `prix_marche`: Statistiques et historique des prix

### Requêtes Typiques

```sql
-- Items filtrés pour le marché
SELECT * FROM items 
WHERE type_item = 'ARME' 
  AND rarete = 'EPIQUE'
  AND echangeable = true;

-- Meilleurs prix de vente
SELECT * FROM ordres_economie
WHERE type_ordre = 'VENTE' 
  AND item_id = 'uuid-item'
  AND statut = 'ACTIF'
ORDER BY prix_unitaire ASC
LIMIT 10;
```

## 4. Projections Joueur (Read Model)

**Base de données**: PostgreSQL + Redis

**Responsabilité**: Profils joueurs, inventaires, équipements, compétences et quêtes.

**Fichier**: `projections_joueur.md`

### Caractéristiques

- **Haute dénormalisation**: Stats totales précalculées
- **Relations complexes**: Inventaire ↔ Équipement ↔ Stats
- **Cache Redis**: Profils fréquemment accédés
- **Vues matérialisées**: Leaderboards globaux

### Tables Principales

- `joueurs`: Profil complet (stats, niveau, XP, or, réputation)
- `inventaires`: Container principal
- `items_inventaire`: Items individuels avec propriétés d'instance
- `sets_equipement`: Équipement porté (13 slots)
- `competences_joueur`: Compétences apprises avec cooldowns
- `quetes_joueur`: Quêtes actives avec progression
- `historique_niveau`: Historique des level-ups

### Cache Redis

```redis
player:{player_id}:profile             # Profil complet
player:{player_id}:stats               # Stats détaillées
player:{player_id}:inventory:meta      # Métadonnées inventaire
player:{player_id}:inventory:items     # Liste des items
player:{player_id}:cooldowns           # Cooldowns actifs
```

### Vues Matérialisées

```sql
-- Leaderboard global
CREATE MATERIALIZED VIEW vue_leaderboard_global AS
SELECT joueur_id, pseudo, classe, niveau, 
       ROW_NUMBER() OVER (ORDER BY niveau DESC, experience DESC) as rang
FROM joueurs
WHERE derniere_connexion >= NOW() - INTERVAL '30 days';
```

## 5. Cache Redis (Temps-Réel)

**Base de données**: Redis 7.x

**Responsabilité**: Cache ultra-rapide et pub/sub pour données temps-réel.

**Fichier**: `projections_cache_redis.md`

### Cas d'Usage

1. **Write-Through** (Combat, Cooldowns)
   - Données écrites simultanément Redis + PostgreSQL
   - Cohérence forte

2. **Cache Aside** (Profils, Inventaires)
   - Redis comme cache, PostgreSQL comme source
   - Invalidation explicite après update

3. **Temporaire** (Sessions, Verrous, Buffers)
   - Données éphémères avec TTL automatique
   - Pas de persistance PostgreSQL

### Structures Principales

- **Hash**: État d'entités (combat, profil, session)
- **Sorted Set**: Leaderboards (classements avec scores)
- **Set**: Index (combats actifs, sessions utilisateur)
- **List**: Files (actions récentes, buffer événements)
- **String**: Verrous distribués (lock avec NX)
- **Pub/Sub**: Notifications temps-réel (WebSocket)

### Leaderboards

```redis
# Classement par niveau (score = niveau*1000000 + xp)
ZADD leaderboard:level 25000500 "p1:Aventurier"
ZREVRANGE leaderboard:level 0 99 WITHSCORES  # Top 100
ZREVRANK leaderboard:level "p1:Aventurier"   # Rang du joueur
```

### Verrous Distribués

```redis
# Empêcher double-action dans un combat
SET lock:combat:{combat_id}:action {session_id} NX EX 5
```

### Pub/Sub Temps-Réel

```redis
# Publication d'événements
PUBLISH combat:{combat_id}:updates '{"type":"TourDebute",...}'
PUBLISH player:{player_id}:notifications '{"type":"ItemRecu",...}'
PUBLISH world:events '{"type":"BossVaincu",...}'

# Abonnement WebSocket
SUBSCRIBE combat:{combat_id}:updates
```

## Flux de Données

### 1. Commande → Événement → Projection

```
┌────────────┐
│  Client    │ POST /combat/action
└─────┬──────┘
      │
      ▼
┌────────────────────┐
│ Command Handler    │ Valider règles métier
│ (CombatAggregate)  │
└─────┬──────────────┘
      │
      ▼ Émettre ActionExecutee
┌────────────────────┐
│   Event Store      │ INSERT INTO evenements
│   (PostgreSQL)     │
└─────┬──────────────┘
      │
      ├──────────────────┐
      │                  │
      ▼                  ▼
┌──────────────┐  ┌─────────────────┐
│ Projection   │  │  Redis Cache    │
│ (PostgreSQL) │  │  (Write-Through)│
└──────────────┘  └─────────────────┘
      │                  │
      ▼                  ▼
┌────────────────────────────┐
│      Query API             │
│  (GraphQL, REST, WebSocket)│
└────────────────────────────┘
```

### 2. Requête (Query)

```
┌────────────┐
│  Client    │ GET /player/profile
└─────┬──────┘
      │
      ▼
┌────────────────────┐
│ Query Handler      │ 1. Check Redis cache
└─────┬──────────────┘
      │
      ├── Cache Hit ──────────┐
      │                       ▼
      │              ┌─────────────────┐
      │              │  Return data    │
      │              └─────────────────┘
      │
      └── Cache Miss
            │
            ▼
      ┌──────────────┐
      │  PostgreSQL  │ SELECT * FROM joueurs
      │  Projection  │
      └─────┬────────┘
            │
            ▼
      ┌──────────────┐
      │ Update Cache │ SET player:{id}:profile
      │   (Redis)    │
      └─────┬────────┘
            │
            ▼
      ┌──────────────┐
      │ Return data  │
      └──────────────┘
```

## Handlers de Projection

Définis dans `event_handlers.md`, les handlers consomment les événements pour maintenir les projections.

### Projections

| Handler | Agrégat Source | Tables Maintenues |
|---------|---------------|-------------------|
| CombatProjection | Combat | instances_combat, participants_combat |
| ParticipantProjection | Combat | participants_combat, effets_statut |
| PlayerProjection | Joueur | joueurs, historique_niveau |
| InventoryProjection | Inventaire | inventaires, items_inventaire |
| EquipmentProjection | Equipement | sets_equipement |
| ItemProjection | Item | items |
| SkillProjection | Competence | competences |
| QuestProjection | Quete | quetes |
| EconomyProjection | Economy | ordres_economie, transactions_economie, prix_marche |
| WorldProjection | WorldState | etat_monde |

### Exemple: CombatProjection

```python
class CombatProjection:
    async def handle(self, event: Event):
        if event.event_type == "CombatDemarre":
            await self.create_combat_instance(event)
        elif event.event_type == "TourDebute":
            await self.update_turn(event)
        elif event.event_type == "DegatsInfliges":
            await self.update_participant_hp(event)
        elif event.event_type == "CombatTermine":
            await self.finalize_combat(event)
            await self.cleanup_redis_cache(event.aggregate_id)
```

## Stratégies de Synchronisation

### Eventual Consistency

Les projections sont mises à jour asynchronément après l'écriture dans l'Event Store. Un délai de quelques millisecondes peut exister.

```
T0: Événement écrit dans Event Store
T0+10ms: Projection PostgreSQL mise à jour
T0+15ms: Cache Redis invalidé/mis à jour
T0+20ms: Client reçoit notification via WebSocket
```

### Garanties de Cohérence

- **Event Store**: Cohérence forte (ACID PostgreSQL)
- **Projections PostgreSQL**: Eventual consistency (généralement <100ms)
- **Cache Redis**: Eventual consistency avec invalidation explicite
- **Pub/Sub Redis**: At-most-once delivery (peut perdre des messages)

### Gestion des Conflits

- **Concurrence optimiste**: `aggregate_version` dans l'Event Store
- **Verrous distribués**: Redis pour actions critiques (combat, échange)
- **Idempotence**: Handlers de projection doivent être idempotents

## Performance et Scalabilité

### PostgreSQL

- **Partitionnement**: `evenements` par mois (time-series)
- **Index**: Optimisés pour les patterns de lecture fréquents
- **Connexions**: Pool de connexions (pgbouncer)
- **Réplication**: Read replicas pour les projections

### Redis

- **Clustering**: Redis Cluster pour scalabilité horizontale
- **Persistance**: RDB + AOF pour durabilité
- **Éviction**: LRU pour gérer la mémoire
- **Sentinel**: Haute disponibilité avec failover automatique

### Métriques Cibles

- **Write latency** (Event Store): <10ms (P99)
- **Projection lag**: <100ms (P95)
- **Read latency** (Cache hit): <5ms (P99)
- **Read latency** (Cache miss): <50ms (P99)
- **WebSocket notification**: <50ms après événement

## Maintenance et Opérations

### Archivage

```sql
-- Archiver les événements > 1 an
INSERT INTO evenements_archive
SELECT * FROM evenements
WHERE timestamp < NOW() - INTERVAL '1 year';

-- Archiver les combats terminés > 30 jours
INSERT INTO instances_combat_archive
SELECT * FROM instances_combat
WHERE etat = 'TERMINE' AND heure_fin < NOW() - INTERVAL '30 days';
```

### Reconstruction de Projection

```sql
-- 1. Trouver la dernière séquence traitée
SELECT MAX(last_event_sequence) FROM joueurs;

-- 2. Rejouer les événements depuis cette séquence
SELECT * FROM evenements
WHERE event_sequence > 12345678
  AND aggregate_type = 'Joueur'
ORDER BY event_sequence ASC;

-- 3. Appliquer chaque événement via le handler
```

### Monitoring

```sql
-- Lag de projection (écart avec Event Store)
SELECT 
    'joueurs' as projection,
    (SELECT MAX(event_sequence) FROM evenements) - 
    (SELECT MAX(last_event_sequence) FROM joueurs) as lag;

-- Combats actifs
SELECT COUNT(*) FROM instances_combat WHERE etat = 'EN_COURS';

-- Transactions économiques du jour
SELECT COUNT(*), SUM(montant_total) 
FROM transactions_economie
WHERE executee_a >= CURRENT_DATE;
```

```redis
# Redis: Utilisation mémoire
INFO memory

# Redis: Cache hit rate
INFO stats
# keyspace_hits / (keyspace_hits + keyspace_misses)
```

## Références

- **event_store.md**: Détails du stockage des événements
- **projections_combat.md**: Schémas des projections de combat
- **projections_monde.md**: Schémas des projections monde (items, économie, etc.)
- **projections_joueur.md**: Schémas des projections joueur
- **projections_cache_redis.md**: Structures Redis détaillées
- **event_handlers.md**: Implémentation des handlers de projection
- **matrice_evenements.md**: Définition de tous les types d'événements (40+)
- **architecture_generale.md**: Vue d'ensemble de l'architecture Event Sourcing
