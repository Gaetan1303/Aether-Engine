# Présentation des tests unitaires – Value Objects du moteur Tactical RPG

## Pourquoi écrire des tests unitaires ?

Les tests unitaires sont essentiels pour garantir la fiabilité, la robustesse et l’évolutivité du code, en particulier dans un projet orienté DDD (Domain-Driven Design) et jeu tactique. Ils permettent de :

- **Valider chaque brique métier** indépendamment (Value Object, Entité, Agrégat)
- **Détecter rapidement les régressions** lors des évolutions
- **Documenter le comportement attendu** (spécification vivante)
- **Favoriser la conception SOLID** (responsabilités claires, code testable)
- **Sécuriser les invariants** (ex : immutabilité, égalité, contraintes)

---

## Organisation des tests

Chaque Value Object clé possède :
- Un dossier dédié (`doc/tests/UnitID`, `doc/tests/position`, `doc/tests/Stats`)
- Un fichier de test unitaire principal (`UnitID.md`, `Position.md`, `Stats.md`)
- Des diagrammes Mermaid pour la structure et les interactions
- Des exemples d’utilisation et d’extension

---

## Exemples de Value Objects testés

### 1. UnitID
- Garantit l’unicité et l’immutabilité d’un identifiant d’unité
- Teste la création, l’égalité, la gestion des erreurs
- Voir : `doc/tests/UnitID/UnitID.md`, `Diagramme de Sequence.md`, `Extension possible.md`, `Utilisation.md`

### 2. Position
- Représente une position sur la grille (x, y)
- Teste l’égalité, la validité, les cas limites
- Voir : `doc/tests/position/Position.md`, `Diagramme de classe Position.md`, `Utilisation.md`

### 3. Stats
- Encapsule les statistiques d’une unité (HP, ATK, DEF, etc.)
- Teste l’égalité, la validité, la gestion des valeurs extrêmes
- Voir : `doc/tests/Stats/Stats.md`, `Diagramme de classe Stats.md`, `Relation.md`, `Utilisation.md`

---

## Bénéfices pour le projet

- **Sécurité métier** : chaque règle ou contrainte est vérifiée automatiquement
- **Refactoring serein** : les tests protègent contre les effets de bord
- **Documentation vivante** : les tests servent d’exemples d’usage
- **Base solide pour l’intégration** : des Value Objects fiables facilitent la construction des entités et agrégats complexes

---

> Les tests unitaires sont la première ligne de défense pour garantir la qualité et la pérennité du moteur de jeu. Ils permettent d’itérer rapidement, d’oser refactorer et d’assurer la cohérence métier à chaque évolution.
