#  Aether Engine - Démo Avancée



## Grille de combat 
   0 1 2 3 4 5 6 7 8 9
 0 · · · · · · · · · ·
 1 · · · · P · · · · ·  P = Paladin
 2 · · A · · · · · B ·  A = Archer
 3 · · · · · · · · · ·  M = Mage
 4 · · · · · · · · C ·  C = Chef
 5 · · · · · · · · · ·  B = Berserker
 6 · · M · · · · · S ·  S = Shaman
 7 · · · · · · · · · ·
 8 · · · · · · · · · ·
 9 · · · · · · · · · ·


##  Démarrage Ultra-Rapide

```bash
# Compiler
go build -o bin/demo-advanced cmd/demo-advanced/main.go

# Lancer
./bin/demo-advanced
```

##  Nouvelles Fonctionnalités

Cette démo **AVANCÉE** exploite tout ce qui a été développé :

###  Système ATH (Attack Hit) Complet
- Chaque unité a une chance de toucher réaliste
- Jets de dés visibles (1-100)
- Bonus +10% pour les compétences magiques
- Messages détaillés avec ATH affiché

###  Combat 3v3 Épique
**Héros de Lumière :**
- **Paladin** 🛡️ - Tank robuste (HP:150, DEF:25, ATH:80%)
  - Compétence 1: Provocation (MP:15, CD:3)
  - Compétence 2: Soin Divin (MP:20, CD:3) 💚 **NOUVEAU**
- **Archer** 🏹 - Sniper précis (ATK:28, ATH:95%, Portée:4)
  - Compétence: Tir de Précision (portée 6, Stamina:15, CD:2)
- **Mage** ✨ - Élémentaliste (MATK:35, ATH:92%)
  - Compétence 1: Boule de Feu (MP:25, CD:2)
  - Compétence 2: Éclair (MP:20, CD:1)
  - Compétence 3: Sommeil (MP:15, CD:4) 😴 **NOUVEAU**
  - Compétence 4: Boost Magique (MP:18, CD:3) ⚡ **NOUVEAU**

**Horde Gobeline :**
- **Chef Gobelin** 👑 - Leader équilibré (HP:100, ATH:82%)
  - Compétence: Cri de Guerre (MP:10, CD:4)
- **Berserker** ⚡ - Brutal mais imprécis (ATK:30, ATH:72%)
  - Aucune compétence (pure force brute)
- **Shaman** 🔮 - Sorcier dangereux (MATK:28, ATH:88%)
  - Compétence 1: Éclair Sombre (MP:18, CD:2)
  - Compétence 2: Paralysie (MP:15, CD:3) ⚡ **NOUVEAU**
  - Compétence 3: Poison (MP:12, CD:4) ☠️ **NOUVEAU**

###  Statistiques de Combat
- Tracker en temps réel des performances
- Précision calculée (attaques réussies/total)
- Dégâts infligés/subis par unité
- Compétences utilisées
- MVP (Most Valuable Player) à la fin

###  IA Améliorée
- **Priorisation intelligente** : cible les héros faibles
- **Utilisation tactique** des compétences
- **Déplacements optimisés** vers les cibles
- Comportement adaptatif selon la situation

###  Interface Enrichie
- Bannière épique au démarrage
- Descriptions de rôles pour chaque unité
- Barres HP colorées selon le statut (vert/jaune/rouge)
- Effets visuels avec couleurs ANSI
- Messages de combat dramatiques

##  Commandes

### Actions de Combat
```bash
attack <cible-id>          # Attaque de base (portée selon l'unité)
skill <nom> <cible-id>     # Utiliser une compétence
move <x> <y>               # Se déplacer (grille 10x10)
map                        # Afficher la grille tactique 🗺️ **NOUVEAU**
stats                      # Voir statistiques détaillées
pass                       # Passer son tour
```

###  Portées d'Attaque
- **Mêlée** (Paladin, Mage) : Portée 1 (adjacent) ⚔️
- **À distance** (Archer) : Portée 4 (tir d'arc) 🏹
- Les compétences ont leurs propres portées individuelles

###  Nouvelles Compétences

**🛡️ Soin Divin (Paladin)**
- Restaure 50 HP à un allié
- MP: 20, Cooldown: 3 tours
- Portée: 4 cases
- Usage: `skill heal hero-archer`

**😴 Sommeil (Mage)**
- Immobilise un ennemi pendant 2 tours
- MP: 15, Cooldown: 4 tours
- Portée: 5 cases
- Usage: `skill sleep ennemi-berserker`

**⚡ Boost Magique (Mage)**
- Augmente MATK d'un allié (+15 pendant 3 tours)
- MP: 18, Cooldown: 3 tours
- Portée: 4 cases
- Usage: `skill boost hero-mage`

**⚡ Paralysie (Shaman)**
- Paralyse un héros pendant 2 tours
- MP: 15, Cooldown: 3 tours
- Portée: 4 cases

**☠️ Poison (Shaman)**
- Empoisonne un héros (10 dégâts/tour pendant 3 tours)
- MP: 12, Cooldown: 4 tours
- Portée: 5 cases

### Utilitaires
```bash
help                       # Afficher l'aide complète
quit                       # Quitter
```

##  IDs des Unités

**Héros :**
- `hero-paladin` - Paladin (Tank)
- `hero-archer` - Archer (Sniper)
- `hero-mage` - Mage (Élémentaliste)

**Ennemis :**
- `ennemi-chef` - Chef Gobelin
- `ennemi-berserker` - Gobelin Berserker
- `ennemi-shaman` - Shaman Gobelin

##  Exemples de Gameplay

### Combat Tactique
```bash
# Tour 1 - Paladin attaque le chef (corps-à-corps)
> attack ennemi-chef
⚔️  Paladin attaque Chef Gobelin et inflige 12 dégâts! (ATH:80%, portée:1)

# Tour 2 - Archer tire à distance
> attack ennemi-shaman
🏹 Archer attaque Shaman Gobelin et inflige 28 dégâts! (ATH:95%, portée:4)

# Tour 3 - Archer utilise Tir de Précision (portée 6)
> skill precision-shot ennemi-shaman
✨ Archer lance Tir de Précision sur Shaman Gobelin et inflige 52 dégâts! (ATH:100%)

# Tour 4 - Mage lance Boule de Feu
> skill fireball ennemi-berserker
🔥 Mage lance Boule de Feu sur Gobelin Berserker et inflige 48 dégâts! (ATH:100%)
💀 Gobelin Berserker a été vaincu!
```

### Voir les Stats
```bash
> stats

 STATISTIQUES DE COMBAT

  HÉROS

  Paladin (Vivant)
    Tours joués: 3
    Attaques: 3 (✓2 ✗1) - Précision: 66.7%
    Dégâts infligés: 24
    Dégâts subis: 18

  Archer (Vivant)
    Tours joués: 3
    Attaques: 2 (✓2 ✗0) - Précision: 100.0%
    Compétences utilisées: 1
    Dégâts infligés: 67
```

##  Système ATH Détaillé

### Valeurs ATH par Unité

| Unité | ATH Base | ATH Magie | Rôle |
|-------|----------|-----------|------|
| Paladin | 80% | 90% | Tank solide mais moins précis |
| Archer | 95% | 100% | Sniper ultra-précis |
| Mage | 92% | 100% | Magie quasi-infaillible |
| Chef | 82% | 92% | Leader compétent |
| Berserker | 72% | 82% | Brutal mais imprécis |
| Shaman | 88% | 98% | Sorcier précis |

### Calcul des Chances

**Attaque de base :**
```
Jet de dés: 1-100
Réussite si: jet ≤ ATH
```

**Compétence magique :**
```
ATH effectif = ATH + 10% (max 100%)
Jet de dés: 1-100
Réussite si: jet ≤ ATH effectif
```

### Impact Gameplay

**Scénario 1 - Archer vs Berserker:**
```
Archer (ATH 95%) attaque
Jet: 23 ≤ 95 → ✅ TOUCHÉ (52 dégâts)

Berserker (ATH 72%) contre-attaque
Jet: 88 > 72 → ❌ RATÉ!
```

**Résultat:** L'Archer gagne grâce à sa précision supérieure

**Scénario 2 - Mage avec Éclair:**
```
Mage (ATH 92% → 100% avec bonus magie) lance Éclair
Jet: 99 ≤ 100 → ✅ TOUCHÉ (45 dégâts)
```

**Résultat:** La magie ne rate (presque) jamais!

##  Statistiques Finales

À la fin du combat, tu obtiens un récapitulatif complet :

```
═══════════════════════════════════════════════

        🎉 VICTOIRE HÉROÏQUE! 🎉

Les héros ont triomphé de la horde gobeline!
Le royaume peut dormir tranquille cette nuit.

═══════════════════════════════════════════════

 STATISTIQUES FINALES
─────────────────────────────────────
Tours total: 8

🏆 MVP (Most Valuable Player)
   Mage - 142 dégâts infligés

HÉROS:
  ⚔️  Paladin - 98/150 HP
      Dégâts: 48 | Précision: 75.0% | Compétences: 1
  ⚔️  Archer - 65/90 HP
      Dégâts: 89 | Précision: 100.0% | Compétences: 3
  ⚔️  Mage - 42/70 HP
      Dégâts: 142 | Précision: 88.9% | Compétences: 5

ENNEMIS:
  💀 Chef Gobelin - VAINCU
      Dégâts: 36 | Précision: 66.7%
  💀 Gobelin Berserker - VAINCU
      Dégâts: 52 | Précision: 60.0%
  💀 Shaman Gobelin - VAINCU
      Dégâts: 41 | Précision: 75.0%
```

##  Stratégies Gagnantes

### 1. Exploitation des Faiblesses ATH
```
Cible prioritaire: Berserker (ATH 72%)
→ Il rate 28% de ses attaques
→ Facile à esquiver
```

### 2. Fiabilité Magique
```
Compétences magiques: ATH effectif 90-100%
→ Utiliser pour finir les ennemis
→ Ne jamais gaspiller
```

### 3. Positionnement Tactique
```
Paladin en avant (Tank)
Archer en position élevée (Portée 6)
Mage en arrière (Protection)
```

### 4. Focus Fire
```
Tour 1: Tout le monde attaque le Shaman
Tour 2: Éliminer le Berserker
Tour 3: Finir le Chef
```

##  Comparaison Démo Simple vs Avancée

| Feature | Démo Simple | Démo Avancée |
|---------|-------------|--------------|
| Unités | 2v2 | 3v3 |
| Compétences | 1 | 5+ |
| Système ATH | ✅ | ✅ Amélioré |
| Statistiques | ❌ | ✅ Complètes |
| IA | Simple | Intelligente |
| Grille | 8x8 | 10x10 |
| Tracking | ❌ | ✅ Temps réel |
| MVP | ❌ | ✅ |

##  Architecture Utilisée

**100% du code domaine réutilisé :**
- ✅ `Combat` (agrégat)
- ✅ `Unite`, `Equipe` (entités)
- ✅ `Competence` (value object)
- ✅ `Stats` avec ATH
- ✅ `GrilleCombat` 10x10
- ✅ `DamageCalculator` (Strategy pattern)

**Nouveautés dans la démo :**
- `CombatStats` - Tracking des performances
- IA avec priorisation des cibles
- Système de statistiques en temps réel
- Interface enrichie

##  Prochaines Améliorations

### Court Terme
- [ ] Système de Statuts (Poison, Stun, Root)
- [ ] Coups critiques (stat CRT)
- [ ] Terrain avec modificateurs ATH
- [ ] Plus de compétences (soin, buffs)

### Moyen Terme
- [ ] Mode histoire avec 5-10 combats
- [ ] Système de progression (XP, levels)
- [ ] Équipement (armes, armures)
- [ ] Sauvegarde/Chargement

### Long Terme
- [ ] UI graphique (Raylib/Ebiten)
- [ ] Multijoueur local
- [ ] Éditeur de campagnes

## 💬 Feedback

Cette démo démontre que **l'architecture fonctionne** et permet de créer un jeu **complet et équilibré** ! 🎉

**Ce qui fonctionne :**
- ✅ Combat tactique profond
- ✅ Système ATH réaliste
- ✅ IA compétitive
- ✅ Statistiques complètes
- ✅ Équilibrage solide

**Améliorations possibles :**
- Statuts manquants (Poison, etc.)
- Interface pourrait être graphique
- Plus de variété de compétences

---


