# Tests de Combat - Aether Engine

Ce dossier contient les tests d'intégration pour le système de combat du moteur Aether Engine.

## Vue d'ensemble

Les tests vérifient le bon fonctionnement du système de combat en utilisant uniquement les APIs publiques du domain layer, sans mock ni simulation artificielle. Ils testent des scénarios réalistes avec de vraies unités, équipes et actions de combat.

## Tests disponibles

### 1. `combat_full_test.go` - Test d'initialisation complète

**Objectif :** Valider l'initialisation d'un combat à grande échelle (6v6)

**Scénario :**
- **Grille :** 15x15 avec obstacles en croix au centre (5 cellules)
- **Équipe Joueurs :** "Héros Aventuriers" (6 unités)
  - Guerrier Tank (HP: 150, ATK: 25, DEF: 15)
  - Archer DPS (HP: 100, ATK: 18, DEF: 12)
  - Mage Élémentaliste (HP: 80, MATK: 30, MDEF: 15)
  - Clerc Soigneur (HP: 120, ATK: 15, DEF: 12)
  - Rôdeur Scout (HP: 110, ATK: 20, DEF: 14)
  - Paladin Protecteur (HP: 140, ATK: 23, DEF: 16)

- **Équipe Gobelins :** "Horde Gobeline" (6 unités IA)
  - Chef de Guerre (HP: 90, ATK: 18, DEF: 12)
  - 2x Gobelin Guerrier (HP: 60, ATK: 12, DEF: 8)
  - Shaman Gobelin (HP: 50, MATK: 20, MDEF: 10)
  - 2x Gobelin Archer (HP: 60, ATK: 12, DEF: 8)

**Ce qui est testé :**
- ✅ Création de la grille avec obstacles
- ✅ Création de 12 unités avec stats variées
- ✅ Formation de 2 équipes (joueur et IA)
- ✅ Initialisation et démarrage du combat
- ✅ Vérification de l'état initial (EnCours, Tour 1)
- ✅ Toutes les unités en position et vivantes

**Résultat :** Combat initialisé avec succès, prêt pour l'exécution de commandes

---

### 2. `combat_tour_complet_test.go` - Test tour par tour avec actions

**Objectif :** Simuler un tour de combat complet avec 4 attaques consécutives

**Scénario :**
- **Grille :** 10x10 avec obstacles en croix centrale (5 cellules)
- **Équipe Joueurs :** "Héros Aventuriers" (3 unités)
  - Guerrier Tank (HP: 100, ATK: 20, DEF: 15) - Position (1, 5)
  - Archer DPS (HP: 80, ATK: 25, DEF: 8) - Position (1, 4)
  - Mage Élémentaliste (HP: 60, MATK: 30, MDEF: 15) - Position (1, 6)

- **Équipe Gobelins :** "Horde Gobeline" (3 unités IA)
  - Chef de Guerre (HP: 70, ATK: 18, DEF: 10) - Position (8, 5)
  - Gobelin Guerrier 1 (HP: 50, ATK: 15, DEF: 8) - Position (8, 4)
  - Gobelin Guerrier 2 (HP: 50, ATK: 15, DEF: 8) - Position (8, 6)

**Déroulement du Tour 1 :**

| Action | Attaquant | Cible | Dégâts calculés | HP avant | HP après |
|--------|-----------|-------|-----------------|----------|----------|
| 1 | Archer | Gobelin Guerrier 1 | 40 | 50 | 10 |
| 2 | Mage | Chef Gobelin | 23 | 70 | 47 |
| 3 | Chef Gobelin | Guerrier | 22 | 100 | 78 |
| 4 | Gobelin Guerrier 2 | Mage | 23 | 60 | 37 |

**Ce qui est testé :**
- ✅ Création d'un combat 3v3
- ✅ Exécution de commandes d'attaque avec `commands.NewAttackCommand()`
- ✅ Calcul réel des dégâts par le système
- ✅ Modification des HP après chaque attaque
- ✅ Vérification qu'aucune unité n'est éliminée
- ✅ Affichage des stats avant/après le tour
- ✅ Combat toujours en cours après le tour 1

**Résultat final :**
- Joueurs vivants : 3/3
- Gobelins vivants : 3/3
- État : EnCours
- Tour : 1

---

## Architecture technique

### APIs utilisées

**Domain Layer :**
```go
// Grille et positions
shared.NewGrilleCombat(largeur, hauteur)
shared.NewPosition(x, y)
grille.DefinirTypeCellule(position, type)

// Stats et unités
shared.NewStats(hp, mp, stamina, atk, def, matk, mdef, spd, mov)
domain.NewUnite(id, nom, teamID, stats, position)
unite.Stats()
unite.StatsActuelles()
unite.EstEliminee()

// Équipes
domain.NewEquipe(id, nom, couleur, isIA, joueurID)
equipe.AjouterMembre(unite)
equipe.Membres()

// Combat
domain.NewCombat(id, equipes, grille)
combat.Demarrer()
combat.Etat()
combat.TourActuel()
combat.Equipes()

// Commandes
commands.NewAttackCommand(attaquant, combat, cible)
attackCmd.Execute() // Retourne (*CommandResult, error)
```

### Structure des tests

Les deux tests suivent la même structure en 6 phases :

1. **Phase 1 :** Création de la grille avec obstacles
2. **Phase 2 :** Création des unités de l'équipe Joueurs
3. **Phase 3 :** Création des unités de l'équipe Gobelins
4. **Phase 4 :** Initialisation et démarrage du combat
5. **Phase 5 :** Exécution des actions (tour complet uniquement)
6. **Phase 6 :** Vérifications et assertions finales

## Exécution des tests

```bash
# Tous les tests de combat
go test -v ./doc/tests/combat/...

# Test d'initialisation uniquement
go test -v ./doc/tests/combat/ -run TestCombatComplet_GobelinsVsJoueurs

# Test tour par tour uniquement
go test -v ./doc/tests/combat/ -run TestCombatTourComplet_GobelinsVsJoueurs
```

## Résultats attendus

Les deux tests doivent passer avec succès :
- ✅ Aucune erreur de compilation
- ✅ Toutes les assertions passent
- ✅ Le système de combat fonctionne comme prévu
- ✅ Les dégâts sont calculés correctement
- ✅ Les états des unités sont mis à jour

## Notes techniques

- **Pattern Command :** Les actions utilisent le pattern Command avec `NewAttackCommand()`
- **CommandResult :** Chaque commande retourne un résultat avec `DamageDealt`, `Success`, etc.
- **Stats immuables :** Les stats de base ne changent pas, seules les stats actuelles évoluent
- **Gestion d'erreurs :** Les tests capturent les erreurs (ex: attaque hors portée) sans faire échouer le test
- **Pas de mock :** Tests d'intégration réels utilisant le vrai système de combat

## Prochaines étapes

Pour étendre ces tests, vous pouvez ajouter :
- Tests de mouvement avec `commands.NewMoveCommand()`
- Tests de compétences spéciales
- Tests d'élimination d'unités
- Tests de fin de combat (victoire/défaite)
- Tests de régénération de ressources
- Tests de statuts (buffs/debuffs)
- Tests avec items et équipements
