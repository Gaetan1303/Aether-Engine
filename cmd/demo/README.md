# 🎮 Aether Engine - Combat Demo CLI

## Démarrage Rapide

```bash
# Compiler
go build -o bin/demo cmd/demo/main.go

# Lancer
./bin/demo
```

## Description

Démo jouable en ligne de commande d'un combat tactique **2 Héros vs 2 Gobelins**.

**Fonctionnalités démontrées :**
- ✅ Combat tour par tour fonctionnel
- ✅ Système de dégâts avec ATK/DEF et MATK/MDEF
- ✅ Compétences avec coûts MP et cooldowns
- ✅ Déplacement sur grille 8x8
- ✅ IA simple pour les ennemis
- ✅ Conditions de victoire/défaite

## Unités

### Équipe Héros (Contrôlée par le joueur)
- **Guerrier** : Tank avec haute défense
  - HP: 120, ATK: 25, DEF: 15, MOV: 4, **ATH: 85%**
  
- **Mage** : DPS magique avec compétence Boule de Feu
  - HP: 80, MP: 100, MATK: 30, MDEF: 20, MOV: 3, **ATH: 90%**
  - **Compétence** : `fireball` (35 dégâts de base, coût 20 MP, cooldown 2 tours, **ATH effectif: 100%**)

### Équipe Gobelins (IA)
- **Gobelin Guerrier** : Attaquant rapide
  - HP: 70, ATK: 18, DEF: 10, SPD: 15, **ATH: 75%**
  
- **Gobelin Archer** : DPS physique
  - HP: 60, ATK: 22, DEF: 8, SPD: 18, **ATH: 80%**

## Commandes

### Actions de Combat
```
attack <cible-id>          # Attaque de base
skill <nom> <cible-id>     # Utiliser une compétence
move <x> <y>               # Se déplacer sur la grille
pass                       # Passer son tour
```

### Utilitaires
```
help                       # Afficher l'aide
quit                       # Quitter le jeu
```

## Exemples

```bash
# Tour du Guerrier
> attack gobelin-1
⚔️  Guerrier attaque Gobelin Guerrier et inflige 15 dégâts!

# Tour du Mage
> skill fireball gobelin-2
✨ Mage lance Boule de Feu sur Gobelin Archer et inflige 48 dégâts!

# Déplacement
> move 3 4
🏃 Guerrier se déplace en (3, 4)
```

## IDs des Cibles

- Héros : `hero-guerrier`, `hero-mage`
- Gobelins : `gobelin-1`, `gobelin-2`

## Grille de Combat

- Dimensions : **8x8**
- Coordonnées : **0-7** en X et Y
- Position initiale :
  - Héros à gauche (x=1)
  - Gobelins à droite (x=6)

## Mécanique de Jeu

### Système ATH (Attack Hit)

Chaque unité a une statistique **ATH** (chance de toucher en %) :
- **Guerrier** : 85% de chance de toucher
- **Mage** : 90% de chance (magie plus précise)
- **Gobelin Guerrier** : 75% de chance
- **Gobelin Archer** : 80% de chance (archer précis)

**Bonus pour les compétences magiques** : +10% ATH (ex: Boule de Feu = 100%)

**Quand une attaque rate** :
- Les dégâts ne sont pas infligés
- Les ressources (MP) sont quand même consommées pour les compétences
- Le message "RATE!" s'affiche avec le jet de dés

**Exemples** :
```
⚔️  Guerrier attaque Gobelin mais RATE! (ATH:85% vs jet:92)
✨ Mage lance Boule de Feu sur Gobelin et inflige 48 dégâts! (ATH:100%)
```

### Calcul des Dégâts
- **Attaque physique** : `(ATK - DEF) × modificateur`
- **Attaque magique** : `(MATK - MDEF) × modificateur`
- **Minimum** : 1 dégât garanti
- **Chance de toucher** : Jet 1-100 doit être ≤ ATH

### Déplacement
- Coût calculé en **distance Manhattan**
- `distance = |x2 - x1| + |y2 - y1|`
- Bloqué si distance > MOV

### Compétences
- Consomment des MP
- Ont un cooldown (tours de recharge)
- Plus puissantes que l'attaque de base

### Tour de Jeu
1. Tour de chaque héros (dans l'ordre)
2. Tour de chaque gobelin (IA)
3. Régénération automatique :
   - MP : +10%
   - Stamina : +20%
4. Décrémentation des cooldowns

## IA Ennemie

Comportement simple mais efficace :
1. Trouver le héros le plus proche
2. Si à portée (distance ≤ 1) → **Attaquer**
3. Sinon → **Se rapprocher**

## Conditions de Fin

- **Victoire** : Tous les gobelins sont vaincus
- **Défaite** : Tous les héros sont vaincus

## Architecture Utilisée

Cette démo utilise le domaine existant d'Aether Engine :
- `internal/combat/domain/combat.go` - Agrégat Combat
- `internal/combat/domain/unite.go` - Entité Unite
- `internal/combat/domain/equipe.go` - Entité Equipe
- `internal/combat/domain/competence.go` - Value Object Competence
- `internal/shared/domain/value_objects.go` - Position, Stats, GrilleCombat

**Patterns démontrés :**
- ✅ Strategy Pattern (DamageCalculator)
- ✅ Composition Pattern (Unite → UnitCombatBehavior)
- ✅ Domain-Driven Design (Agrégats, Entités, Value Objects)

## Prochaines Étapes

Pour améliorer la démo :
1. **Implémenter les statuts** (Poison, Stun, Root)
2. **Ajouter plus de compétences** (Soin, Buffs, Debuffs)
3. **Améliorer l'IA** (pathfinding, stratégie)
4. **Ajouter une UI graphique** (Raylib, Ebiten, ou web)

## Feedback

Cette démo prouve que l'architecture fonctionne et permet de jouer ! 🎉

**Points forts :**
- Code domain réutilisé tel quel
- Combat fonctionnel et équilibré
- Interface claire et responsive

**À améliorer :**
- Statuts non implémentés (Poison, etc.)
- HP non régénéré (par design actuel)
- Pas de système de récompenses

---

**Temps de développement de la démo :** ~30 minutes  
**Lignes de code :** ~600 (interface CLI uniquement)  
**Utilise le domaine existant :** ✅ 100%
