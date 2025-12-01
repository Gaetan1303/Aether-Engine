
# Phase 1 : Fondations & Contrats

> **Note de synchronisation** :
> Les concepts d'agrégats, Value Objects, etc. sont centralisés dans `/doc/agregats.md`.
> Les diagrammes et la documentation utilisent le nommage français, sauf pour les termes internationalement utilisés (item, Tank, DPS, Heal, etc.).

## Objectif du Jalon

La Phase 1, dite des "Fondations", établit l’hygiène du projet et formalise les contrats conceptuels essentiels. Elle garantit que les phases suivantes (Combat Core et Engine Logic) seront construites sur une architecture modulaire, déterministe et alignée sur les principes du Domain-Driven Design (DDD).

**Résultats attendus :**
- Une structure de dépôt claire basée sur Go.
- Les Value Objects fondamentaux du domaine implémentés et testés.
- Les contrats conceptuels documentés dans `/docs`.
- Aucune logique de combat encore implémentée : uniquement les bases structurelles.

---

# 1. Architecture du Dépôt (Go + DDD)

Le dépôt suit une structure propre au langage Go et inspirée des architectures hexagonales/DDD modernes.

| Dossier | Responsabilité | Contenu |
|--------|----------------|---------|
| `/cmd` | Point d'entrée de l'application | `main.go`, initialisation HTTP/GRPC |
| `/internal` | Logique métier non exportable | Bounded Contexts, agrégats, services |
| `/internal/combat` | Contexte de combat (Moteur de Combat) | Agrégats Combat, services du domaine, repositories |
| `/internal/shared` | Domaines transverses | `IdentifiantUnite`, `Position`, `Statistiques`, abstractions temporelles |
| `/pkg` | Code réutilisable et public | utilitaires, math, helpers |
| `/tests` | Tests d'intégration / end-to-end | scénarios complets de combat |
| `/docs` | Documentation | contrats, modèles, diagrammes, pipelines |

---

# 2. Définition des Contrats et Modélisation Domaine

L'objectif principal est la formalisation des **contrats conceptuels** décrivant le domaine du combat, avant toute logique métier.

## 2.1 Value Objects Fondamentaux

Trois Value Objects clés ont été définis pour la Phase 1 :

| Artefact | Type DDD | Rôle | Statut |
|---------|----------|------|--------|
| `Position` | Value Object | Coordonnées X, Y, Z et leurs opérations (distances, voisinage). | Validé (3D stable) |
| `IdentifiantUnite` | Value Object | Identifiant unique typé d'une unité. | Validé |
| `Statistiques` | Semi-mutable Value Object | Points de vie, mana, attaque, défense. | Validé |

### Clarification sur la dimension Z

Le champ Z représente la verticalité :
- différence de hauteur (bonus/malus, couverture),
- élévation (marches, plateformes),
- extensible vers des grilles multi-niveaux.

La distance exploite une métrique adaptée (Manhattan 3D ou Euclidienne selon les besoins).

---

## 2.2 Services de Domaine — Rôles explicites

Les services sont divisés selon trois catégories (important pour éviter les dérives) :

### Services Applicatifs
- orchestrent des use cases,
- ne contiennent pas de logique métier pure,
- appellent les repositories, agrégats et services de domaine.

### Services de Domaine
- contiennent les règles du combat (résolution de dégâts, portée, etc.),
- strictement déterministes.

### Services Techniques
- Ticker / scheduler,
- générateur aléatoire injectable,
- serializers, IO abstraits.

---

# 3. Contrats de Déterminisme

Le moteur de combat doit être **strictement déterministe**.

Règles définies dès la Phase 1 :

1. Aucun accès à l’horloge système dans le domaine.
2. Aucun IO dans les agrégats ou services de domaine.
3. Le random est uniquement fourni par une interface `RandomProvider`, seedable.
4. Les évolutions temporelles sont orchestrées *uniquement* par `TickerCombat`.
5. Les entités doivent être intégralement sérialisables pour permettre le replay.

Ces contrats préparent la Phase 2 (Combat Core Déterministe).

---

# 4. Diagrammes de Flux (Mermaid)

Quatre diagrammes cadrent la compréhension du moteur.

## 4.1 Architecture globale (Frontend → Engine)

Client → API Gateway → Aether-Engine (Go) → Base de données.


## 4.2 Machine d'états du Combat

États principaux :
```
Attente → EnCours → AttenteAction → Résolution → Terminé
```

## 4.3 Cycle de vie d'une instance de combat

- Création via "Monde/Carte".
- Construction de l'Agrégat Combat.
- Enregistrement du TickerCombat.
- Nettoyage en fin de combat.

## 4.4 Pipeline de règles Fabric (esquisse)

Entrée → Validation → Résolution → Effets → Post-Traitement → Evénements

Le pipeline sera finalisé en Phase 3.

---

# 5. Setup Toolchain & Hygiène

| Tâche | Description | Outils |
|-------|-------------|--------|
| Environnement Go | Setup Go 1.24+ | Go |
| Tests unitaires | Tests sur `internal/combat/domain` | go test |
| Qualité | Linter Go | golangci-lint |
| Client | Initialisation minimal Angular | Angular CLI |

---

# 6. Mini-Cartographie des Bounded Contexts

Liste des contextes prévus :
- Combat (Phase 1 + 2 + 3)
- Personnage (classes, progression)
- Inventory & Items
- World & Map
- Events & Triggers
- Progression

Pour la Phase 1, **seul le Contexte Combat est concerné**.

---

# Conclusion de Phase 1

La Phase 1 est complète lorsque :

1. La structure Go est en place.
2. Les Value Objects `Position`, `UnitID`, `Stats` sont implémentés et testés.
3. Les contrats conceptuels sont dans `/docs`.
4. Les règles de déterminisme sont formalisées.
5. Les diagrammes de flux Merlin sont validés.

La Phase 2 pourra démarrer avec :
- `BattleTicker`,
- le Command Pattern,
- la file d’actions.

