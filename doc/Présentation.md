# Présentation du Projet Tactical RPG – Serveur avec un design pattern de type Fabric

## 1. Vision Générale

- **Serveur autoritatif** : Source de vérité, toutes les règles et validations sont côté serveur.
- **Tour par tour dynamique** : Système hybride avec gestion dynamique du CT/ATB.
- **Scalabilité** : Architecture pensée pour héberger de nombreuses instances de combat simultanées.
- **Client** : Angular 20, simple interface graphique reflétant l’état serveur.
- **Technos** : Go (serveur Fabric), API Gateway (Go/NestJS), Kafka/Redis, PostgreSQL.

---

---

## 2. Domain-Driven Design (DDD)

- **Agrégats principaux** :
  - Unit
  - Team
  - BattleGrid
  - Battle (agrégat maître)
- **Skill** : Value Object (immuable)

---

## 3. State Machine du Combat

Cycle : Idle → Preparing → Running → WaitingForAction → Executing → ApplyingEffects → PostProcessing → Ended

```mermaid
stateDiagram-v2
[*] --> Idle
Idle --> Preparing : startBattle()
Preparing --> Running : initTurnOrder()
Running --> WaitingForAction : nextActorReady
WaitingForAction --> ExecutingAction : receiveAction
ExecutingAction --> ApplyingEffects : resolveAction
ApplyingEffects --> PostProcessing : applyEffects, updateATB
PostProcessing --> Running : continue / checkEnd
Running --> Ended : victory OR defeat
Ended --> [*]
```

---

## 4. Diagramme de Classes (Mermaid)

```mermaid
classDiagram
class Unit {
  +UnitID id
  +Stats stats
  +Status[] statuses
  +Position pos
  +ApplyDamage()
  +ApplyStatus()
}
class Team {
  +TeamID id
  +Unit[] members
}
class BattleGrid {
  +width
  +height
  +cells
}
class Battle {
  +BattleID id
  +Team[] teams
  +BattleGrid grid
  +TurnOrder order
}
Unit -- Team
Team -- Battle
Battle -- BattleGrid
```

---

## 5. Diagramme de Séquence (Action de compétence)

```mermaid
sequenceDiagram
participant Client
participant API_GW
participant Fabric
participant Battle
participant DB
Client->>API_GW: castSkill()
API_GW->>Fabric: forward(action)
Fabric->>Battle: enqueue(action)
Battle->>Battle: validate
Battle->>DB: persistEvent
Battle-->>Fabric: finalEvent
Fabric-->>Client: pushEvent
```

---

## 6. Architecture & Concurrence

- 1 goroutine Go par instance de combat
- Channels pour la synchronisation
- Snapshots réguliers de l’état

---

## 7. CQRS & Event Sourcing

- Émission d’événements de jeu
- Projections ReadModel pour statistiques et analyse

---

## 8. Stratégie de Tests

- Tests unitaires DDD
- Tests d’intégration
- Tests end-to-end
- Injection de simulateurs (Fake RNG)

---

## 9. Roadmap Technique

- Sprint 0 : Conception et tests
- Sprint 1 : Base Fabric
- Sprint 2 : Skills et résolution
- Sprint 3 : API Gateway + events
- Sprint 4 : Scalabilité et observabilité

---
