# Projections Cache Redis - Données Temps-Réel

## Vue d'ensemble

Les projections cache Redis contiennent les données critiques nécessitant des accès ultra-rapides et des opérations temps-réel. Redis sert de couche de cache entre l'Event Store PostgreSQL et les clients, permettant des performances optimales pour les opérations fréquentes.

## Architecture

- **Base de données**: Redis 7.x
- **Pattern**: Cache aside / Write-through selon les cas d'usage
- **Persistance**: RDB + AOF pour durabilité
- **Expiration**: TTL adaptatif selon le type de données
- **Éviction**: LRU (Least Recently Used) si mémoire saturée
- **Pub/Sub**: Notifications temps-réel pour les clients WebSocket

## Cas d'Usage Redis

### Données Critiques (Write-Through)
- État des combats actifs
- Sessions utilisateur
- Cooldowns de compétences
- Verrous distribués (combat, échange)

### Données Fréquemment Accédées (Cache Aside)
- Profils joueur
- Inventaires
- Leaderboards
- État du monde

### Données Temporaires
- Actions en attente de résolution
- Buffers d'événements (avant batch insert PostgreSQL)
- Statistiques temps-réel

## Structures de Données

### 1. Combat Temps-Réel

#### État Global d'un Combat

```redis
# Hash: combat:{combat_id}:state
HSET combat:a1b2c3d4:state combat_id "a1b2c3d4-..."
HSET combat:a1b2c3d4:state etat "EN_COURS"
HSET combat:a1b2c3d4:state tour_actuel "5"
HSET combat:a1b2c3d4:state participant_actif_id "p1..."
HSET combat:a1b2c3d4:state prochain_timeout "1737820800"
HSET combat:a1b2c3d4:state zone_id "zone_123"
EXPIRE combat:a1b2c3d4:state 3600
```

**Opérations**:
```redis
# Lire l'état complet
HGETALL combat:{combat_id}:state

# Mettre à jour le tour
HINCRBY combat:{combat_id}:state tour_actuel 1

# Changer le participant actif
HSET combat:{combat_id}:state participant_actif_id {new_participant_id}

# Vérifier l'existence
EXISTS combat:{combat_id}:state
```

#### Participants d'un Combat

```redis
# Hash: combat:{combat_id}:participant:{participant_id}
HSET combat:a1b2c3d4:participant:p1 participant_id "p1..."
HSET combat:a1b2c3d4:participant:p1 nom "Aventurier"
HSET combat:a1b2c3d4:participant:p1 hp_actuel "85"
HSET combat:a1b2c3d4:participant:p1 hp_max "100"
HSET combat:a1b2c3d4:participant:p1 mp_actuel "40"
HSET combat:a1b2c3d4:participant:p1 mp_max "50"
HSET combat:a1b2c3d4:participant:p1 etat "ACTIF"
HSET combat:a1b2c3d4:participant:p1 peut_agir "1"
EXPIRE combat:a1b2c3d4:participant:p1 3600
```

**Opérations**:
```redis
# Appliquer des dégâts
HINCRBY combat:{combat_id}:participant:{pid} hp_actuel -25

# Soigner
HINCRBY combat:{combat_id}:participant:{pid} hp_actuel 15

# Liste de tous les participants
KEYS combat:{combat_id}:participant:*

# Marquer comme KO
HSET combat:{combat_id}:participant:{pid} etat "KO"
HSET combat:{combat_id}:participant:{pid} peut_agir "0"
```

#### Ordre des Tours

```redis
# Sorted Set: combat:{combat_id}:turn_order (score = initiative)
ZADD combat:a1b2c3d4:turn_order 18 "p1..."
ZADD combat:a1b2c3d4:turn_order 15 "p2..."
ZADD combat:a1b2c3d4:turn_order 12 "p3..."
EXPIRE combat:a1b2c3d4:turn_order 3600
```

**Opérations**:
```redis
# Obtenir l'ordre complet (initiative décroissante)
ZREVRANGE combat:{combat_id}:turn_order 0 -1 WITHSCORES

# Prochain participant
ZREVRANGE combat:{combat_id}:turn_order 0 0

# Retirer un participant (mort/fuite)
ZREM combat:{combat_id}:turn_order {participant_id}
```

#### Effets Actifs

```redis
# Set: combat:{combat_id}:participant:{pid}:effects
SADD combat:a1b2c3d4:participant:p1:effects "effect:poison:e1"
SADD combat:a1b2c3d4:participant:p1:effects "effect:regen:e2"
EXPIRE combat:a1b2c3d4:participant:p1:effects 3600

# Hash pour chaque effet: combat:{combat_id}:effect:{effect_id}
HSET combat:a1b2c3d4:effect:e1 type_effet "POISON"
HSET combat:a1b2c3d4:effect:e1 puissance "5"
HSET combat:a1b2c3d4:effect:e1 duree_tours_restants "3"
EXPIRE combat:a1b2c3d4:effect:e1 3600
```

**Opérations**:
```redis
# Tous les effets d'un participant
SMEMBERS combat:{combat_id}:participant:{pid}:effects

# Décrémenter durée d'un effet
HINCRBY combat:{combat_id}:effect:{eid} duree_tours_restants -1

# Retirer un effet expiré
DEL combat:{combat_id}:effect:{eid}
SREM combat:{combat_id}:participant:{pid}:effects "effect:poison:{eid}"
```

#### Actions Récentes

```redis
# List: combat:{combat_id}:recent_actions (FIFO, max 20)
LPUSH combat:a1b2c3d4:recent_actions '{"tour":5,"participant":"p1","action":"ATTAQUE","resultat":"CRITIQUE","degats":35}'
LTRIM combat:a1b2c3d4:recent_actions 0 19
EXPIRE combat:a1b2c3d4:recent_actions 3600
```

**Opérations**:
```redis
# 10 dernières actions
LRANGE combat:{combat_id}:recent_actions 0 9

# Ajouter une action
LPUSH combat:{combat_id}:recent_actions {action_json}
LTRIM combat:{combat_id}:recent_actions 0 19
```

#### Index: Combats Actifs par Zone

```redis
# Set: zone:{zone_id}:active_combats
SADD zone:zone_123:active_combats "a1b2c3d4"
SADD zone:zone_123:active_combats "b2c3d4e5"
EXPIRE zone:zone_123:active_combats 7200
```

**Opérations**:
```redis
# Tous les combats actifs d'une zone
SMEMBERS zone:{zone_id}:active_combats

# Ajouter/retirer
SADD zone:{zone_id}:active_combats {combat_id}
SREM zone:{zone_id}:active_combats {combat_id}

# Nombre de combats
SCARD zone:{zone_id}:active_combats
```

### 2. Sessions Utilisateur

```redis
# Hash: session:{session_id}
HSET session:sess_abc123 user_id "u1..."
HSET session:sess_abc123 username "Aventurier"
HSET session:sess_abc123 ip_address "192.168.1.100"
HSET session:sess_abc123 created_at "1737820800"
HSET session:sess_abc123 last_activity "1737820900"
HSET session:sess_abc123 websocket_connection_id "ws_xyz"
EXPIRE session:sess_abc123 3600

# Index inversé: user:{user_id}:sessions
SADD user:u1:sessions "sess_abc123"
EXPIRE user:u1:sessions 3600
```

**Opérations**:
```redis
# Valider une session
EXISTS session:{session_id}

# Mettre à jour l'activité
HSET session:{session_id} last_activity {timestamp}
EXPIRE session:{session_id} 3600

# Toutes les sessions d'un utilisateur
SMEMBERS user:{user_id}:sessions

# Déconnexion
DEL session:{session_id}
SREM user:{user_id}:sessions {session_id}
```

### 3. Profils Joueur (Cache)

```redis
# Hash: player:{player_id}:profile
HSET player:p1:profile joueur_id "p1..."
HSET player:p1:profile pseudo "Aventurier"
HSET player:p1:profile classe "GUERRIER"
HSET player:p1:profile niveau "25"
HSET player:p1:profile hp_actuel "250"
HSET player:p1:profile hp_max "300"
HSET player:p1:profile or "15000"
EXPIRE player:p1:profile 600

# Stats complètes (JSON compressé)
SET player:p1:stats '{"force":50,"dexterite":30,...}'
EXPIRE player:p1:stats 600
```

**Opérations**:
```redis
# Charger le profil
HGETALL player:{player_id}:profile

# Mise à jour HP
HINCRBY player:{player_id}:profile hp_actuel -20

# Mise à jour Or
HINCRBY player:{player_id}:profile or 500

# Invalidation (après update PostgreSQL)
DEL player:{player_id}:profile
DEL player:{player_id}:stats
```

### 4. Inventaire (Cache)

```redis
# Hash: player:{player_id}:inventory:meta
HSET player:p1:inventory:meta capacite "50"
HSET player:p1:inventory:meta slots_utilises "32"
HSET player:p1:inventory:meta poids_actuel "850"
HSET player:p1:inventory:meta poids_max "1000"
EXPIRE player:p1:inventory:meta 600

# Set: player:{player_id}:inventory:items
SADD player:p1:inventory:items "item:slot1"
SADD player:p1:inventory:items "item:slot2"
EXPIRE player:p1:inventory:items 600

# Hash par item: player:{player_id}:inventory:item:{slot_id}
HSET player:p1:inventory:item:slot1 item_id "i123"
HSET player:p1:inventory:item:slot1 nom "Épée Longue"
HSET player:p1:inventory:item:slot1 quantite "1"
HSET player:p1:inventory:item:slot1 est_equipe "1"
EXPIRE player:p1:inventory:item:slot1 600
```

**Opérations**:
```redis
# Inventaire complet
SMEMBERS player:{player_id}:inventory:items
# Puis pour chaque item:
HGETALL player:{player_id}:inventory:item:{slot_id}

# Ajouter un item
SADD player:{player_id}:inventory:items "item:{slot_id}"
HSET player:{player_id}:inventory:item:{slot_id} ...
HINCRBY player:{player_id}:inventory:meta slots_utilises 1

# Retirer un item
SREM player:{player_id}:inventory:items "item:{slot_id}"
DEL player:{player_id}:inventory:item:{slot_id}
HINCRBY player:{player_id}:inventory:meta slots_utilises -1

# Invalidation
DEL player:{player_id}:inventory:meta
DEL player:{player_id}:inventory:items
# + DEL pour chaque item
```

### 5. Cooldowns de Compétences

```redis
# Hash: player:{player_id}:cooldowns
HSET player:p1:cooldowns skill:s1 "15"
HSET player:p1:cooldowns skill:s2 "0"
HSET player:p1:cooldowns skill:s3 "5"
EXPIRE player:p1:cooldowns 300
```

**Opérations**:
```redis
# Démarrer un cooldown
HSET player:{player_id}:cooldowns skill:{skill_id} {seconds}

# Décrémenter (job périodique chaque seconde)
HINCRBY player:{player_id}:cooldowns skill:{skill_id} -1

# Vérifier disponibilité
HGET player:{player_id}:cooldowns skill:{skill_id}
# Si <= 0 ou NULL, compétence disponible

# Tous les cooldowns
HGETALL player:{player_id}:cooldowns
```

### 6. Leaderboards

#### Classement par Niveau

```redis
# Sorted Set: leaderboard:level (score = niveau * 1000000 + xp)
ZADD leaderboard:level 25000500 "p1:Aventurier"
ZADD leaderboard:level 24998000 "p2:Mage"
ZADD leaderboard:level 23500000 "p3:Archer"
```

**Opérations**:
```redis
# Top 100
ZREVRANGE leaderboard:level 0 99 WITHSCORES

# Rang d'un joueur
ZREVRANK leaderboard:level "p1:Aventurier"

# Joueurs autour d'un joueur (+/- 10)
ZREVRANGE leaderboard:level {rang-10} {rang+10} WITHSCORES

# Mise à jour
ZADD leaderboard:level {new_score} "{player_id}:{pseudo}"
```

#### Classement Arène

```redis
# Sorted Set: leaderboard:arena (score = rating ELO)
ZADD leaderboard:arena 1850 "p1:Aventurier"
ZADD leaderboard:arena 1720 "p2:Champion"
```

#### Classement Guilde

```redis
# Sorted Set: leaderboard:guild (score = points totaux)
ZADD leaderboard:guild 150000 "guild:g1:LesLégendes"
ZADD leaderboard:guild 125000 "guild:g2:Templiers"
```

### 7. État du Monde (Cache)

```redis
# Hash: world:state
HSET world:state cycle_jour_nuit "14"
HSET world:state saison "ETE"
HSET world:state jour_saison "15"
EXPIRE world:state 60

# Set: world:active_events
SADD world:active_events "event:invasion_dragon"
SADD world:active_events "event:marche_special"
EXPIRE world:active_events 300

# Hash par événement: world:event:{event_id}
HSET world:event:invasion_dragon nom "Invasion du Dragon Noir"
HSET world:event:invasion_dragon zone_id "zone_boss"
HSET world:event:invasion_dragon expire_a "1737824400"
EXPIRE world:event:invasion_dragon 3600
```

### 8. Verrous Distribués

```redis
# Pour empêcher les actions concurrentes

# Verrou de combat (empêcher double-action)
SET lock:combat:{combat_id}:action {session_id} NX EX 5
# NX = Only set if not exists
# EX 5 = Expire in 5 seconds

# Verrou d'inventaire (empêcher échange simultané)
SET lock:player:{player_id}:inventory {transaction_id} NX EX 10

# Verrou d'ordre économique
SET lock:economy:order:{order_id} {user_id} NX EX 30
```

**Opérations**:
```redis
# Acquérir un verrou
SET lock:combat:{combat_id}:action {session_id} NX EX 5
# Retourne OK si acquis, NULL si déjà verrouillé

# Libérer un verrou (avec vérification propriétaire)
# Utiliser un script Lua pour atomicité:
EVAL "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end" 1 lock:combat:{combat_id}:action {session_id}

# Étendre un verrou
EXPIRE lock:combat:{combat_id}:action 5
```

### 9. File d'Événements (Buffer)

```redis
# List: event:buffer (avant insertion batch PostgreSQL)
RPUSH event:buffer '{"eventType":"DegatsInfliges",...}'
RPUSH event:buffer '{"eventType":"SoinsRecus",...}'

# Stream: event:stream (alternative avec consumer groups)
XADD event:stream * event_type DegatsInfliges payload {...}
```

**Opérations**:
```redis
# Ajouter un événement
RPUSH event:buffer {event_json}

# Lire un batch (par worker)
LRANGE event:buffer 0 99
# Puis supprimer après insertion PostgreSQL
LTRIM event:buffer 100 -1

# Avec Streams (plus robuste)
XADD event:stream * event_type {type} payload {json}
XREADGROUP GROUP workers consumer1 COUNT 100 STREAMS event:stream >
```

### 10. Statistiques Temps-Réel

```redis
# Compteurs
INCR stats:combats:today
INCR stats:transactions:today
INCR stats:players:online

# Expiration automatique à minuit
EXPIREAT stats:combats:today {timestamp_minuit}

# HyperLogLog pour comptage unique (joueurs actifs)
PFADD stats:active_players:today "p1"
PFADD stats:active_players:today "p2"
PFCOUNT stats:active_players:today
```

## Pub/Sub pour Temps-Réel

### Canaux de Publication

```redis
# Combat updates
PUBLISH combat:{combat_id}:updates '{"type":"TourDebute","tour":6,...}'

# Notifications joueur
PUBLISH player:{player_id}:notifications '{"type":"ItemRecu","item":"Épée",...}'

# Chat de zone
PUBLISH zone:{zone_id}:chat '{"pseudo":"Aventurier","message":"Bonjour!"}'

# Événements globaux
PUBLISH world:events '{"type":"BossVaincu","boss":"Dragon Noir"}'
```

### Abonnements Clients WebSocket

```redis
# Client WebSocket s'abonne
SUBSCRIBE combat:{combat_id}:updates
SUBSCRIBE player:{player_id}:notifications
SUBSCRIBE zone:{zone_id}:chat
SUBSCRIBE world:events
```

## Stratégies d'Expiration

| Type de Données | TTL | Stratégie |
|----------------|-----|-----------|
| Combat actif | 1h | Prolongé à chaque action |
| Session utilisateur | 1h | Prolongé à chaque activité |
| Profil joueur | 10min | Cache aside, invalidation explicite |
| Inventaire | 10min | Cache aside, invalidation explicite |
| Cooldowns | 5min | Write-through, auto-expiration |
| Leaderboards | Pas d'expiration | Mise à jour incrémentale |
| État du monde | 1min | Rafraîchi périodiquement |
| Verrous | 5-30s | Expiration automatique (safety) |
| Buffer événements | Pas d'expiration | Vidé par workers |

## Patterns d'Utilisation

### Write-Through (Combat, Cooldowns)

```python
# Écriture
await redis.hset(f"combat:{combat_id}:state", "tour_actuel", 5)
await postgres.update("UPDATE instances_combat SET tour_actuel = 5 WHERE combat_id = $1", combat_id)

# Lecture
data = await redis.hgetall(f"combat:{combat_id}:state")
```

### Cache Aside (Profils, Inventaires)

```python
# Lecture
data = await redis.hgetall(f"player:{player_id}:profile")
if not data:
    # Cache miss
    data = await postgres.fetch("SELECT * FROM joueurs WHERE joueur_id = $1", player_id)
    await redis.hset(f"player:{player_id}:profile", mapping=data)
    await redis.expire(f"player:{player_id}:profile", 600)
return data
```

### Invalidation Explicite (Après Update)

```python
# Après modification PostgreSQL
await postgres.update("UPDATE joueurs SET niveau = $1 WHERE joueur_id = $2", new_level, player_id)

# Invalider le cache
await redis.delete(f"player:{player_id}:profile")
await redis.delete(f"player:{player_id}:stats")

# Publier notification
await redis.publish(f"player:{player_id}:notifications", json.dumps({
    "type": "NiveauAtteint",
    "niveau": new_level
}))
```

## Surveillance et Maintenance

### Métriques Clés

```redis
# Utilisation mémoire
INFO memory

# Nombre de clés
DBSIZE

# Hit rate du cache
INFO stats
# Regarder keyspace_hits vs keyspace_misses

# Clients connectés
CLIENT LIST

# Opérations par seconde
INFO stats
# Regarder instantaneous_ops_per_sec
```

### Nettoyage

```redis
# Trouver les clés par pattern
SCAN 0 MATCH combat:* COUNT 1000

# Supprimer les combats terminés
DEL combat:{combat_id}:state
DEL combat:{combat_id}:turn_order
# etc.

# Flush d'une base (DEV uniquement!)
FLUSHDB
```

### Persistance

Configuration `redis.conf`:
```conf
# RDB snapshots
save 900 1
save 300 10
save 60 10000

# AOF (Append-Only File)
appendonly yes
appendfsync everysec

# Compression
rdbcompression yes
```

## Références

- **event_store.md**: Source de vérité PostgreSQL
- **projections_combat.md**: Structure des données de combat
- **projections_joueur.md**: Structure des données joueur
- **projections_monde.md**: État du monde et économie
- **flux_reseaux.md**: Architecture WebSocket et synchronisation temps-réel
