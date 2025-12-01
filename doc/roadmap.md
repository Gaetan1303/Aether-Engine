

# Roadmap Technique — Aether-Engine (Serveur Autoritatif)

> **Note de synchronisation** :
> Les concepts d'agrégats, Value Objects, etc. sont centralisés dans `/doc/Agrégats.md`.
> Les diagrammes et la documentation utilisent le nommage français, sauf pour les termes internationalement utilisés (item, Tank, DPS, Heal, etc.).

> **Objectif** : Développer un moteur de combat tactique déterministe, résilient aux pannes (Event Sourcing) et modulaire (DDD/Fabric).

## Vision Générale

| Phase | Nom du Jalon                | Focus Principal                        | Livrables Clés                                                        |
|-------|-----------------------------|----------------------------------------|------------------------------------------------------------------------|
| P1    | Fondations (Structure)      | Hygiène du projet, Contrats & Diagrammes| Structure du dépôt Go, Contrats (Markdown), Diagrammes Mermaid          |
| P2    | Cœur du Combat (Squelette)  | Déterminisme, Temps & Intention         | Position3D, TickerCombat (ATB), Agrégat Combat (racine), Command Pattern|
| P3    | Fabric & Résolution         | Logique Métier, Pipeline de règles      | RésolveurDégâts (Pipeline), Interfaces EffetStatut (Hooks), Tests de résolution complexes |
| P4    | Résilience & Mémoire        | Event Sourcing & Récupération d'état    | Event Store, Snapshotting, Rehydrator, Implémentation de l'Idempotence  |
| P5    | Monde & API                 | Interface externe, Scalabilité          | MoteurMonde, API REST/WebSocket, Mapping DTO, Verrous de Concurrence (Redis) |
| P6    | Finalisation & Portfolio    | Qualité et Présentation                 | Documentation finale, Métriques (Prometheus), Vidéos Portfolio          |

---

## Phase 1 : Foundations (Semaine 1)

| Tâche | Description | Dépendance |
|-------|-------------|------------|
| 1.1. Architecture du Repo | Créer la structure DDD (e.g., /internal/combat/domain, /internal/shared). | N/A |
| 1.2. Définition des Contrats | Documenter le Combat Flow, la State Machine et les principes du Fabric (markdown). | N/A |
| 1.3. Diagrammes Majeurs | Produire les diagrammes d'architecture et de la Battle State Machine (Mermaid). | N/A |
| 1.4. Setup Toolchain | Configuration go.mod (Go 1.24+), tests unitaires (Go test), linter. | N/A |

---

## Phase 2 : Combat Core Déterministe (Semaine 2-3)

> **Objectif** : Rendre la boucle de jeu fonctionnelle et entièrement déterministe en mémoire.

| Tâche | Description | Dépendance |
|-------|-------------|------------|
| 2.1. Position 3D | Refactor du Value Object Position pour inclure l'axe Z. | P1 |
| 2.2. TickerCombat | Implémenter le moteur ATB déterministe (Tick(), ResetATB()). | 2.1 |
| 2.3. Command Pattern | Définir l'interface Command et implémenter MoveUnitCommand, AttackCommand. | P1 |
| 2.4. Agrégat Combat | Créer l'Agrégat racine Combat, intégrant le Ticker et la file d'acteurs. | 2.2, 2.3 |
| 2.5. Tests Core | Validation unitaire que le Turn Order est toujours le même pour les mêmes SPD. | 2.4 |

---

## Phase 3 : Fabric & Résolution (Semaine 4-5)

> **Objectif** : Implémenter la logique métier complexe des règles de combat de manière modulaire et extensible (le Fabric).

| Tâche | Description | Dépendance |
|-------|-------------|------------|
| 3.1. Ruleset & Data | Définir les Value Objects pour Skill (coût, portée, type) et les formules de base (Dégâts/Soins). | P2 |
| 3.2. Pipeline de Résolution | Implémenter le RésolveurDégâts (Chain of Responsibility) utilisant l'objet DamageSnapshot. | 3.1, P2.4 |
| 3.3. Hooks EffetStatut | Transformer l'EffetStatut en interface avec des hooks (OnIncomingDamage, OnTurnStart) qui s'insèrent dans le Pipeline. | 3.2 |
| 3.4. Cibles & Portée | Logique de validation de la portée (Manhattan/Euclidienne) et des cibles (Single, AOE, Row). | 2.1 |
| 3.5. Tests Pipeline | Scénario complexe : Unité A (Rage) attaque Unité B (Shielded, Poisoned). S'assurer que tous les hooks sont appliqués. | 3.3 |

---

## Phase 4 : Résilience & Mémoire (Semaine 6-7)

> **Objectif** : Rendre l'état du combat persistant et résilient aux pannes via l'Event Sourcing.

| Tâche | Description | Dépendance |
|-------|-------------|------------|
| 4.1. Event Store Setup | Configurer la base de données (PostgreSQL ou Mongo) et créer le schéma pour stocker les événements. | P3 |
| 4.2. Event Sourcing | Implémenter la méthode Battle.ApplyEvent() et le mécanisme Repository.SaveEvents(). | P3 |
| 4.3. State Rehydrator | Coder Repository.LoadAggregate() qui reconstruit l'état du Battle à partir de tous les événements passés. | 4.2 |
| 4.4. Snapshotting | Implémenter la logique pour sauvegarder l'état complet du Battle tous les N événements (accélération du Replay). | 4.3 |
| 4.5. Idempotence Check | Middleware/Logique dans ReceiveCommand pour rejeter les actions ayant un OperationID déjà traité. | 2.3 |

---

## Phase 5 : World & API (Semaine 8-9)

> **Objectif** : Exposer le moteur aux clients et gérer le contexte de jeu plus large.

| Tâche | Description | Dépendance |
|-------|-------------|------------|
| 5.1. API Gateway (Gin) | Créer les endpoints StartBattle (POST), SendAction (POST/WebSocket). | P4 |
| 5.2. DTO Mapping | Créer la couche de translation entre Command/Event (Domaine) et le format JSON/HTTP (Client). | P4 |
| 5.3. MoteurMonde (Contexte) | Implémenter le MoteurMonde (contexte supérieur) gérant les instances de Combat et les joueurs hors combat. | P4 |
| 5.4. Communication Temps Réel | Mettre en place les WebSockets (ou SSE) pour pusher les Events (P4) au client en temps réel. | 5.2 |
| 5.5. Verrous de Concurrence | Utilisation de Redis pour verrouiller chaque CombatID si le serveur doit être horizontalement scalable. | 5.3 |

---

## Phase 6 : Finalisation & Portfolio (Semaine 10)

> **Objectif** : Mettre le projet en état de production et de présentation professionnelle.

| Tâche | Description | Dépendance |
|-------|-------------|------------|
| 6.1. Observabilité | Intégrer Prometheus/Grafana pour les métriques critiques (temps de tick, latence des commandes). | P5 |
| 6.2. Documentation Finale | Finaliser le README.md, l'Architecture Overview, et le guide du Fabric Scripting. | P5 |
| 6.3. Tests End-to-End | Écrire des tests d'intégration complets simulant un client envoyant des commandes via l'API. | P5 |
| 6.4. Portfolio Vidéo | Capturer une séquence de combat tour par tour pour le portfolio. | P5 |

Phase 1 : Foundations (Semaine 1)

Tâche	Description	Dépendance
1.1. Architecture du Repo	Créer la structure DDD (e.g., /internal/combat/domain, /internal/shared).	N/A
1.2. Définition des Contrats	Documenter le Combat Flow, la State Machine et les principes du Fabric (markdown).	N/A
1.3. Diagrammes Majeurs	Produire les diagrammes d'architecture et de la Battle State Machine (Mermaid).	N/A
1.4. Setup Toolchain	Configuration go.mod (Go 1.24+), tests unitaires (Go test), linter.	N/A

Phase 2 : Combat Core Déterministe (Semaine 2-3)

Objectif : Rendre la boucle de jeu fonctionnelle et entièrement déterministe en mémoire.
Tâche	Description	Dépendance
2.1. Position 3D	Refactor du Value Object Position pour inclure l'axe Z.	P1
2.2. Battle Ticker	Implémenter le moteur ATB déterministe (Tick(), ResetATB()).	2.1
2.3. Command Pattern	Définir l'interface Command et implémenter MoveUnitCommand, AttackCommand.	P1
2.4. Battle Aggregate	Créer l'Agrégat Root Battle, intégrant le Ticker et la Queue d'acteurs.	2.2, 2.3
2.5. Tests Core	Validation unitaire que le Turn Order est toujours le même pour les mêmes SPD.	2.4

Phase 3 : Fabric & Résolution (Semaine 4-5)

Objectif : Implémenter la logique métier complexe des règles de combat de manière modulaire et extensible (le Fabric).
Tâche	Description	Dépendance
3.1. Ruleset & Data	Définir les Value Objects pour Skill (coût, portée, type) et les formules de base (Dégâts/Soins).	P2
3.2. Resolution Pipeline	Implémenter le DamageResolver (Chain of Responsibility) utilisant l'objet DamageSnapshot.	3.1, P2.4
3.3. Status Hooks	Transformer le Status en interface avec des hooks (OnIncomingDamage, OnTurnStart) qui s'insèrent dans le Pipeline.	3.2
3.4. Cibles & Portée	Logique de validation de la portée (Manhattan/Euclidienne) et des cibles (Single, AOE, Row).	2.1
3.5. Tests Pipeline	Scénario complexe : Unité A (Rage) attaque Unité B (Shielded, Poisoned). S'assurer que tous les hooks sont appliqués.	3.3

Phase 4 :  Résilience & Mémoire (Semaine 6-7)

Objectif : Rendre l'état du combat persistant et résilient aux pannes via l'Event Sourcing.
Tâche	Description	Dépendance
4.1. Event Store Setup	Configurer la base de données (PostgreSQL ou Mongo) et créer le schéma pour stocker les événements.	P3
4.2. Event Sourcing	Implémenter la méthode Battle.ApplyEvent() et le mécanisme Repository.SaveEvents().	P3
4.3. State Rehydrator	Coder Repository.LoadAggregate() qui reconstruit l'état du Battle à partir de tous les événements passés.	4.2
4.4. Snapshotting	Implémenter la logique pour sauvegarder l'état complet du Battle tous les N événements (accélération du Replay).	4.3
4.5. Idempotence Check	Middleware/Logique dans ReceiveCommand pour rejeter les actions ayant un OperationID déjà traité.	2.3

Phase 5 :  World & API (Semaine 8-9)

Objectif : Exposer le moteur aux clients et gérer le contexte de jeu plus large.
Tâche	Description	Dépendance
5.1. API Gateway (Gin)	Créer les endpoints StartBattle (POST), SendAction (POST/WebSocket).	P4
5.2. DTO Mapping	Créer la couche de translation entre Command/Event (Domaine) et le format JSON/HTTP (Client).	P4
5.3. World Engine (Context)	Implémenter le WorldEngine (contexte supérieur) gérant les instances de Battle et les joueurs hors combat.	P4
5.4. Communication Temps Réel	Mettre en place les WebSockets (ou SSE) pour pusher les Events (P4) au client en temps réel.	5.2
5.5. Concurrency Locks	Utilisation de Redis pour verrouiller chaque BattleID si le serveur doit être horizontalement scalable.	5.3

Phase 6 :  Finalisation & Portfolio (Semaine 10)

Objectif : Mettre le projet en état de production et de présentation professionnelle.
Tâche	Description	Dépendance
6.1. Observabilité	Intégrer Prometheus/Grafana pour les métriques critiques (temps de tick, latence des commandes).	P5
6.2. Documentation Finale	Finaliser le README.md, l'Architecture Overview, et le guide du Fabric Scripting.	P5
6.3. Tests End-to-End	Écrire des tests d'intégration complets simulant un client envoyant des commandes via l'API.	P5
6.4. Portfolio Vidéo	Capturer une séquence de combat tour par tour pour le portfolio.	P5