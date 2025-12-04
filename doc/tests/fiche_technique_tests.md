# Fiche Technique des Tests - Aether Engine

##  Vue d'ensemble

- **Total des tests unitaires** : 130 tests
- **Total des tests d'intégration** : 60 tests
- **Total général** : 190 tests
- **Framework** : Go testing + testify/assert
- **Couverture** : Domaine Combat (Unite, Equipe, Combat, Competence)

---

##  Tests Unitaires (130 tests)

### 1. Combat (42 tests)

#### Gestion des Objets
- `TestCombat_AjouterObjet` : Ajout d'objets à l'inventaire du combat
- `TestCombat_ConsommerObjet` : Consommation d'objets avec décrémentation de quantité
- `TestCombat_ObtienirObjet` : Récupération d'un objet par ID
- `TestCombat_ObtienirQuantiteObjet` : Obtention de la quantité d'un objet
- `TestCombat_PossedeObjet` : Vérification de la présence d'un objet

#### Gestion des Fuites
- `TestCombat_AnnulerFuite` : Annulation du statut de fuite d'une équipe
- `TestCombat_FuiteAutorisee` : Vérification si la fuite est autorisée
- `TestCombat_MarquerEquipeFuite` : Marquage d'une équipe comme ayant fui
- `TestCombat_SetFuiteAutorisee` : Configuration de l'autorisation de fuite

#### Events & CQRS
- `TestCombat_Apply` : Application d'événements de domaine
- `TestCombat_ClearUncommittedEvents` : Nettoyage des événements non committés
- `TestCombat_GetUncommittedEvents` : Récupération des événements en attente
- `TestCombat_RaiseEvent` : Levée d'événements de domaine

#### Pattern Command
- `TestCombat_GetCommandFactory` : Récupération de la factory de commandes
- `TestCombat_GetCommandInvoker` : Récupération de l'invoker de commandes
- `TestCombat_SetCommandFactory` : Injection de la factory de commandes
- `TestCombat_SetCommandInvoker` : Injection de l'invoker de commandes

#### Pattern Strategy (Damage Calculator)
- `TestCombat_GetDamageCalculator` : Récupération du calculateur de dégâts
- `TestCombat_SetDamageCalculator` : Injection du calculateur personnalisé
- `TestCombat_SetHybridDamageMode` : Configuration mode hybride (physique + magique)
- `TestCombat_SetMagicalDamageMode` : Configuration mode dégâts magiques
- `TestCombat_SetPhysicalDamageMode` : Configuration mode dégâts physiques

#### Pattern Observer
- `TestCombat_GetObserverSubject` : Récupération du subject observable
- `TestCombat_SetObserverSubject` : Injection du subject observable

#### Pattern State Machine
- `TestCombat_GetStateMachine` : Récupération de la machine à états
- `TestCombat_SetStateMachine` : Injection de la machine à états

#### Pattern Chain of Responsibility
- `TestCombat_GetValidationChain` : Récupération de la chaîne de validation
- `TestCombat_SetValidationChain` : Injection de la chaîne de validation

#### Gestion du Combat
- `TestCombat_Demarrer` : Démarrage d'un combat
- `TestCombat_Equipes` : Accès à la map des équipes
- `TestCombat_Etat` : Vérification de l'état du combat
- `TestCombat_Grille` : Accès à la grille de combat
- `TestCombat_ID` : Vérification de l'ID du combat
- `TestCombat_ObtienirEnnemis` : Récupération des ennemis d'une équipe
- `TestCombat_ObtienirPositionsOccupees` : Récupération des positions occupées
- `TestCombat_ObtienirResultat` : Obtention du résultat (VICTORY/DEFEAT/CONTINUE)
- `TestCombat_TourActuel` : Récupération du numéro de tour actuel
- `TestCombat_TrouverUnite` : Recherche d'une unité par ID
- `TestCombat_VerifierConditionsVictoire` : Vérification des conditions de victoire
- `TestCombat_Version` : Gestion du versioning pour Event Sourcing
- `TestCombat_GetTimestamp` : Récupération du timestamp de création

#### Constructeur
- `TestNewCombat` : Création d'un combat avec validation

---

### 2. Competence (28 tests)

#### Propriétés de Base
- `TestCompetence_ID` : Vérification de l'ID unique
- `TestCompetence_Nom` : Vérification du nom
- `TestCompetence_Description` : Vérification de la description
- `TestCompetence_Type` : Vérification du type (Attaque, Magie, Soin, etc.)
- `TestCompetence_Portee` : Vérification de la portée
- `TestCompetence_DegatsBase` : Vérification des dégâts de base
- `TestCompetence_Modificateur` : Vérification du modificateur de dégâts

#### Coûts et Cooldown
- `TestCompetence_CoutMP` : Vérification du coût en MP
- `TestCompetence_CoutStamina` : Vérification du coût en Stamina
- `TestCompetence_Cooldown` : Vérification du cooldown max
- `TestCompetence_CooldownActuel` : Vérification du cooldown actuel
- `TestCompetence_EstEnCooldown` : Vérification si en cooldown
- `TestCompetence_ActiverCooldown` : Activation du cooldown après usage
- `TestCompetence_DecrementerCooldown` : Décrémentation du cooldown
- `TestCompetence_SetCooldownActuel` : Définition manuelle du cooldown

#### Zone d'Effet
- `TestCompetence_ObtienirPositionsDansZone` : Calcul des positions affectées

#### Ciblage
- `TestCompetence_Cibles` : Vérification du type de cibles (Ennemis, Alliés, etc.)
- `TestCompetence_EstCibleValide` : Validation d'une cible selon le type

#### Effets
- `TestCompetence_Effets` : Récupération de la liste d'effets
- `TestCompetence_AjouterEffet` : Ajout d'un effet à la compétence

#### Calculs
- `TestCompetence_CalculerDegats` : Calcul des dégâts en fonction des stats

#### Utilitaires
- `TestCompetence_Clone` : Clonage d'une compétence

#### Constructeur
- `TestNewCompetence` : Création d'une compétence avec validation

---

### 3. Equipe (18 tests)

#### Propriétés de Base
- `TestEquipe_ID` : Vérification de l'ID unique
- `TestEquipe_Nom` : Vérification du nom
- `TestEquipe_Couleur` : Vérification de la couleur (hex)
- `TestEquipe_IsIA` : Vérification si l'équipe est contrôlée par IA
- `TestEquipe_JoueurID` : Vérification de l'ID du joueur (null si IA)

#### Gestion des Membres
- `TestEquipe_AjouterMembre` : Ajout d'une unité à l'équipe
- `TestEquipe_RetirerMembre` : Retrait d'une unité de l'équipe
- `TestEquipe_NombreMembres` : Comptage du nombre de membres
- `TestEquipe_ObtienirMembre` : Récupération d'un membre par ID
- `TestEquipe_ContientUnite` : Vérification si une unité est dans l'équipe

#### État de l'Équipe
- `TestEquipe_MembresVivants` : Récupération des membres vivants
- `TestEquipe_MembresElimines` : Récupération des membres morts
- `TestEquipe_ADesMembresVivants` : Vérification s'il reste des membres vivants
- `TestEquipe_EstComplete` : Vérification si l'équipe est complète (4 membres)

#### Relations
- `TestEquipe_EstEnnemie` : Vérification si une autre équipe est ennemie

#### Statistiques Agrégées
- `TestEquipe_StatsMoyennes` : Calcul des stats moyennes de l'équipe
- `TestEquipe_PuissanceTotale` : Calcul de la puissance totale (somme ATK+MATK)

#### Constructeur
- `TestNewEquipe` : Création d'une équipe avec validation

---

### 4. Unite (42 tests)

#### Propriétés de Base
- `TestUnite_ID` : Vérification de l'ID unique
- `TestUnite_Nom` : Vérification du nom
- `TestUnite_TeamID` : Vérification de l'équipe d'appartenance
- `TestUnite_Position` : Vérification de la position sur la grille
- `TestUnite_IsIA` : Vérification si l'unité est contrôlée par IA

#### Statistiques
- `TestUnite_Stats` : Récupération des stats de base
- `TestUnite_StatsActuelles` : Récupération des stats modifiées
- `TestUnite_RecalculerStats` : Recalcul des stats avec modificateurs
- `TestUnite_AppliquerModificateurStat` : Application d'un modificateur temporaire
- `TestUnite_RetirerModificateurStat` : Retrait d'un modificateur

#### Points de Vie (HP)
- `TestUnite_HPActuels` : Récupération des HP actuels
- `TestUnite_SetHP` : Définition manuelle des HP
- `TestUnite_RecevoirDegats` : Application de dégâts avec défense
- `TestUnite_RecevoirSoin` : Application de soins avec limite max
- `TestUnite_Soigner` : Alias pour soigner
- `TestUnite_EstEliminee` : Vérification si HP ≤ 0
- `TestUnite_Ressusciter` : Résurrection avec restauration HP

#### Points de Magie (MP)
- `TestUnite_SetMP` : Définition manuelle des MP
- `TestUnite_ConsommerMP` : Consommation de MP
- `TestUnite_RestaurerMP` : Restauration de MP avec limite max

#### Stamina
- `TestUnite_ConsommerStamina` : Consommation de Stamina

#### Compétences
- `TestUnite_Competences` : Récupération de la liste des compétences
- `TestUnite_AjouterCompetence` : Ajout d'une compétence
- `TestUnite_ObtienirCompetence` : Récupération d'une compétence par ID
- `TestUnite_ObtienirCompetenceParDefaut` : Récupération de l'attaque de base
- `TestUnite_UtiliserCompetence` : Utilisation d'une compétence avec consommation
- `TestUnite_PeutUtiliserCompetence` : Vérification des conditions d'usage
- `TestUnite_SkillEstPret` : Vérification si une compétence est disponible
- `TestUnite_ActiverCooldown` : Activation du cooldown d'une compétence

#### Statuts
- `TestUnite_Statuts` : Récupération de la liste des statuts actifs
- `TestUnite_AjouterStatut` : Ajout d'un statut (Poison, Stun, etc.)
- `TestUnite_RetirerStatut` : Retrait d'un statut
- `TestUnite_TraiterStatuts` : Traitement des statuts (décrémentation durée)
- `TestUnite_EstEmpoisonne` : Vérification du statut Poison
- `TestUnite_EstStun` : Vérification du statut Stun
- `TestUnite_EstSilence` : Vérification du statut Silence
- `TestUnite_EstRoot` : Vérification du statut Root

#### Déplacement
- `TestUnite_SeDeplacer` : Déplacement avec consommation de points de mouvement
- `TestUnite_DeplacerVers` : Téléportation directe sans coût
- `TestUnite_PeutSeDeplacer` : Vérification si l'unité peut se déplacer
- `TestUnite_EstBloqueDeplacement` : Vérification des blocages (Stun, Root)

#### Actions
- `TestUnite_PeutAgir` : Vérification si l'unité peut agir (non Stun)

#### Buffs
- `TestUnite_AppliquerBuff` : Application d'un buff temporaire

#### Régénération & Tour
- `TestUnite_RegenerarStatut` : Régénération automatique (MP 10%, Stamina 20%)
- `TestUnite_NouveauTour` : Traitement du nouveau tour (regen, cooldowns, reset actions)

#### Constructeur
- `TestNewUnite` : Création d'une unité avec validation

---

##  Tests d'Intégration (60 tests)

### 1. Combat - Damage Calculator Integration (11 tests)

**Fichier** : `combat_damage_calculator_integration_test.go`

#### Tests des Stratégies de Dégâts
- `TestCombatDamageCalculator_Physical` : Mode dégâts physiques (ATK vs DEF)
- `TestCombatDamageCalculator_ChangerModePhysique` : Changement vers mode physique
- `TestCombatDamageCalculator_Magical` : Mode dégâts magiques (MATK vs MDEF)
- `TestCombatDamageCalculator_Hybrid` : Mode hybride 50/50 physique-magique
- `TestCombatDamageCalculator_HybridRatios` : Ratios hybrides personnalisés (80/20, 20/80)
- `TestCombatDamageCalculator_ComparePhysicalVsMagical` : Comparaison des stratégies
- `TestCombatDamageCalculator_SetCustomCalculator` : Injection calculateur personnalisé (fixed damage)
- `TestCombatDamageCalculator_MinimumDamage` : Vérification du minimum de 1 dégât
- `TestCombatDamageCalculator_AvecModificateurCompetence` : Dégâts avec modificateur de compétence
- `TestCombatDamageCalculator_InfligerDegatsIntegration` : Flux complet d'infliction de dégâts
- `TestCombatDamageCalculator_CalculerPourDifferentsTypes` : Calculs pour différents types de compétences

**Comportement vérifié** :
- Pattern Strategy appliqué au calcul de dégâts
- Injection de dépendances via `SetDamageCalculator`
- Formules : Physical = ATK - DEF, Magical = MATK - MDEF, Hybrid = mix des deux
- Minimum garanti de 1 dégât

---

### 2. Combat - Equipe - Unite Integration (14 tests)

**Fichier** : `combat_equipe_unite_integration_test.go`

#### Tests de Recherche et Relations
- `TestCombat_TrouverUnites` : Recherche d'unités par ID dans le combat
- `TestCombat_ObtenirEnnemis` : Récupération des ennemis d'une équipe donnée
- `TestCombat_PositionsOccupees` : Récupération des positions occupées
- `TestCombat_PositionsOccupeesAvecExclusion` : Positions occupées en excluant une unité

#### Tests des Conditions de Victoire
- `TestCombat_VerifierConditionsVictoire` : Retour "CONTINUE" quand les deux équipes sont vivantes
- `TestCombat_MarquerEquipeFuite` : Marquage d'une équipe en fuite → "FLED"
- `TestCombat_AnnulerFuite` : Annulation de fuite → retour "CONTINUE"
- `TestCombat_DesactiverFuite` : Désactivation de la fuite
- `TestCombat_ObtenirResultat` : Délégation à `VerifierConditionsVictoire()`

#### Tests Multi-Équipes
- `TestCombat_MultipleEquipes` : Combat avec 3 équipes
- `TestCombat_EquipesVides` : Combat avec équipes vides retourne "DEFEAT"

#### Tests Métadonnées
- `TestCombat_GetTimestamp` : Timestamp > 0 à la création
- `TestCombat_VersionIncrementation` : Versioning pour Event Sourcing
- `TestCombat_Equipes` : Accès à la map des équipes

**Comportement vérifié** :
- Conditions de victoire : 0 ou 1 équipe active → fin de combat
- Fuite : "FLED" si une équipe fuit, "CONTINUE" si annulé
- Équipes vides (0 membres vivants) → "DEFEAT"
- Multi-équipes supporté (3+ équipes)

---

### 3. Equipe - Unite Integration (11 tests)

**Fichier** : `equipe_unite_integration_test.go`

#### Tests de Gestion des Membres
- `TestEquipeUnite_AjouterPlusieursUnites` : Ajout de plusieurs unités
- `TestEquipeUnite_RetirerMembre` : Retrait d'un membre
- `TestEquipeUnite_MembresVivantsEtElimines` : Séparation vivants/morts après élimination
- `TestEquipeUnite_TousElimines` : Détection de tous les membres éliminés
- `TestEquipeUnite_ContientUnite` : Vérification de la présence d'une unité
- `TestEquipeUnite_EstComplete` : Équipe complète = 4 membres

#### Tests de Statistiques Agrégées
- `TestEquipeUnite_StatsMoyennes` : Calcul des stats moyennes
- `TestEquipeUnite_PuissanceTotale` : Calcul de la puissance totale

#### Tests de Relations
- `TestEquipeUnite_EquipeEnnemie` : Identification des équipes ennemies

#### Tests de Résurrection
- `TestEquipeUnite_RessusciterMembre` : Résurrection d'un membre mort

**Comportement vérifié** :
- Séparation correcte entre vivants et éliminés (HP ≤ 0)
- Résurrection fonctionne (HP restaurés, `EstEliminee()` retourne false)
- Stats agrégées calculées correctement (moyenne, puissance)
- Équipe complète = 4 membres maximum

---

### 4. Unite - Competence Integration (7 tests)

**Fichier** : `unite_competence_integration_test.go`

#### Tests d'Usage de Compétences
- `TestUniteUtiliserCompetence_FluxComplet` : Usage complet avec consommation MP/Stamina
- `TestUniteUtiliserCompetence_MPInsuffisant` : Erreur si MP insuffisants
- `TestUniteUtiliserCompetence_EnCooldown` : Erreur si compétence en cooldown

#### Tests de Cooldown
- `TestUniteCompetence_CooldownDecremente` : Cooldowns décrément avec `NouveauTour()`
  - **Note** : Cooldown 2 → 1 après 1 tour (bug corrigé)

#### Tests de Gestion des Compétences
- `TestUniteCompetence_ObtenirCompetenceParDefaut` : ID par défaut = "attaque-basique"
- `TestUniteCompetence_AjouterPlusieursCompetences` : Gestion de plusieurs compétences

**Comportement vérifié** :
- Consommation MP/Stamina lors de l'usage
- Blocage si MP insuffisants ou en cooldown
- Décrémentation des cooldowns à chaque tour
- Compétence par défaut toujours disponible

---

### 5. Unite - Deplacement Integration (12 tests)

**Fichier** : `unite_deplacement_integration_test.go`

#### Tests de Déplacement de Base
- `TestUniteDeplacement_Basique` : Déplacement basique avec changement de position
- `TestUniteDeplacement_CoutDeplacement` : Consommation des points de mouvement
- `TestUniteDeplacement_DeplacementInsuffisant` : Erreur si points insuffisants

#### Tests de Blocages
- `TestUniteDeplacement_AvecRoot` : Root ne bloque PAS (non implémenté)
- `TestUniteDeplacement_UniteEliminee` : Unité morte ne peut pas se déplacer

#### Tests de Restauration
- `TestUniteDeplacement_NouveauTourRestaure` : Points de mouvement restaurés à chaque tour

#### Tests Avancés
- `TestUniteDeplacement_DeplacerVers` : Téléportation directe (pas de coût)
- `TestUniteDeplacement_PositionsOccupeesDansCombat` : Vérification des positions occupées
- `TestUniteDeplacement_DeplacementMultiple` : Plusieurs déplacements successifs
- `TestUniteDeplacement_VerificationEstBloqueDeplacement` : Vérification des blocages
- `TestUniteDeplacement_GrilleValidation` : Validation des limites de la grille

**Comportement vérifié** :
- Coût de déplacement calculé par distance Manhattan
- Restauration des points à chaque tour
- Blocage si unité morte (HP ≤ 0)
- **Statuts de blocage (Root, Stun) NON implémentés**
- Validation des limites de grille

---

### 6. Tour Complet Integration (10 tests)

**Fichier** : `tour_complet_integration_test.go`

#### Tests de Régénération
- `TestTourComplet_NouveauTourRestaureRessources` : MP/Stamina restaurés (PAS HP)
- `TestTourComplet_RegenerationMP` : Régénération de 10% MP (5 MP sur 50)

#### Tests de Cooldown
- `TestTourComplet_CooldownsDecrementes` : Cooldowns décrément avec `NouveauTour()`
  - **Note** : Cooldown 2 → 1 après 1 tour (bug corrigé)

#### Tests de Cycle Complet
- `TestTourComplet_CycleComplet` : Actions + déplacement + cooldowns + regen
- `TestTourComplet_PlusieursUnites` : Cycle pour plusieurs unités
- `TestTourComplet_AvecCombat` : Tour dans contexte de combat

#### Tests de Restrictions
- `TestTourComplet_UniteElimineeNePasProceder` : Unité morte ne peut pas agir

#### Tests de Restauration
- `TestTourComplet_ActionsEtDeplacementRestaurees` : Compteurs d'actions restaurés

**Comportement vérifié** :
- Régénération : MP +10%, Stamina +20%, **HP non régénéré**
- Cooldowns décrémentés de 1 à chaque tour
- Actions et déplacement restaurés
- Unité morte (HP ≤ 0) ne peut pas agir
- Cycle complet : regen → cooldowns → reset actions

---

##  Infrastructure de Test

### Helpers Unitaires (`doc/tests/unitaire/helpers_test.go`)

```go
// Création de stats de test
newStats(hp, mp, atk, def, matk, mdef, spd, mov, sta int) *shared.Stats

// Création de position
newPosition(x, y int) shared.Position

// Création d'unité
newUnite(id, nom, teamID string, x, y int) *domain.Unite

// Création d'équipe
newEquipe(id, nom string) (*domain.Equipe, error)

// Création de compétence
newCompetence(id, nom string, typeComp domain.TypeCompetence) *domain.Competence

// Création de grille
newGrilleCombat(largeur, hauteur int) *domain.GrilleCombat

// Création de combat
newCombat(id string) (*domain.Combat, error)
```

### Helpers Intégration (`doc/tests/integration/helpers_test.go`)

```go
// Position
newTestPosition(x, y int) shared.Position

// Stats avec 9 paramètres
newTestStats(hp, mp, atk, def, matk, mdef, spd, mov, stamina int) shared.Stats

// Unité avec position
newTestUnite(id, nom, teamID string, x, y int) *domain.Unite

// Équipe (joueur)
newTestEquipe(id, nom string) *domain.Equipe
// Paramètres: couleur="#FF0000", isIA=false, joueurID="player-1"

// Équipe IA
newTestEquipeIA(id, nom string) *domain.Equipe
// Paramètres: couleur="#0000FF", isIA=true, joueurID=nil

// Compétence de base
newTestCompetence(id, nom string, typeComp domain.TypeCompetence) *domain.Competence

// Compétence avec coûts personnalisés
newTestCompetenceAvecCouts(id, nom string, coutMP, coutStamina, cooldown int) *domain.Competence

// Grille
newTestGrille(largeur, hauteur int) *domain.GrilleCombat

// Combat basique (2 équipes vides)
newTestCombat(id string) *domain.Combat

// Combat avec 2 unités (hero vs enemy)
newTestCombatAvecUnites(id string) (*domain.Combat, *domain.Unite, *domain.Unite)

// Statut
newTestStatut(type shared.TypeStatut, duree, degatsParTour int) *shared.Statut
```

---

##  Statistiques de Couverture

### Par Entité

| Entité | Tests Unitaires | Tests Intégration | Total |
|--------|----------------|-------------------|-------|
| Combat | 42 | 25 | 67 |
| Unite | 42 | 29 | 71 |
| Equipe | 18 | 11 | 29 |
| Competence | 28 | 7 | 35 |
| **TOTAL** | **130** | **60** | **190** |

### Par Catégorie de Tests

| Catégorie | Nombre de Tests | Description |
|-----------|----------------|-------------|
| Getters/Setters | 35 | Tests de propriétés simples |
| Gestion de Collections | 28 | Ajout/Retrait/Recherche |
| Logique Métier | 42 | Calculs, conditions, validations |
| Patterns GoF | 18 | Strategy, Command, Observer, State, Chain |
| Event Sourcing | 7 | Events, Apply, Version |
| Intégration Multi-Entités | 60 | Interactions entre entités |

---

##  Configuration de Test

### Dépendances

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

### Patterns de Nommage

**Unitaires** : `Test<Entité>_<Méthode>`
- Exemple : `TestUnite_RecevoirDegats`

**Intégration** : `Test<Scope>_<Scenario>`
- Exemple : `TestCombatDamageCalculator_Physical`

### Structure AAA (Arrange-Act-Assert)

```go
func TestExample(t *testing.T) {
    // Arrange - Préparation
    unite := newTestUnite("u1", "Hero", "team-1", 0, 0)
    
    // Act - Action
    unite.RecevoirDegats(50)
    
    // Assert - Vérification
    assert.Equal(t, 50, unite.HPActuels())
}
```

---

##  Comportements Notables Documentés

### 1. Régénération (NouveauTour)
- ✅ MP : +10% (minimum 1)
- ✅ Stamina : +20%
- ❌ HP : **NON régénéré** (non implémenté)

### 2. Cooldowns
- ✅ Décrémente de 1 par tour
- ✅ Bug de double décrémentation **corrigé**

### 3. Statuts de Blocage
- ❌ Silence : **NON implémenté** (pas d'effet)
- ❌ Stun : **NON implémenté** (pas d'effet)
- ❌ Root : **NON implémenté** (pas d'effet)
- ❌ Poison : **NON implémenté** (pas d'effet)

### 4. Conditions de Victoire
- Si 0 équipes actives : `"DEFEAT"`
- Si 1 équipe active : `"VICTORY"` (sauf si fuite)
- Si équipe fuit : `"FLED"`
- Sinon : `"CONTINUE"`

### 5. Compétence par Défaut
- ID : `"attaque-basique"`
- Toujours disponible (pas de cooldown)

---

##  Exécution des Tests

### Tous les tests unitaires
```bash
go test ./doc/tests/unitaire -v
```

### Tous les tests d'intégration
```bash
go test ./doc/tests/integration -v
```

### Test spécifique
```bash
go test ./doc/tests/unitaire -run TestUnite_RecevoirDegats -v
```

### Avec couverture
```bash
go test ./doc/tests/unitaire -cover
go test ./doc/tests/integration -cover
```

---

##  Notes Techniques

### Patterns de Conception Testés

1. **Strategy Pattern** : DamageCalculator (Physical, Magical, Hybrid, Fixed)
2. **Command Pattern** : CommandFactory, CommandInvoker
3. **Observer Pattern** : ObserverSubject pour notifications
4. **State Pattern** : StateMachine pour états de combat
5. **Chain of Responsibility** : ValidationChain pour validations
6. **Composition Pattern** : Unite délègue à UnitCombatBehavior, UnitStatusManager, UnitInventory

### Architecture CQRS/Event Sourcing

- `RaiseEvent()` : Levée d'événements
- `Apply()` : Application d'événements
- `GetUncommittedEvents()` : Récupération des événements non persistés
- `ClearUncommittedEvents()` : Nettoyage après commit
- `Version` : Versioning pour Event Sourcing

### Formules de Calcul

**Dégâts Physiques** :
```
damage = (ATK - DEF) * modificateur
damage = max(1, damage)
```

**Dégâts Magiques** :
```
damage = (MATK - MDEF) * modificateur
damage = max(1, damage)
```

**Dégâts Hybrides** :
```
physical = (ATK - DEF) * ratioPhysique
magical = (MATK - MDEF) * ratioMagique
damage = (physical + magical) * modificateur
damage = max(1, damage)
```

**Régénération** :
```
MP_regen = max(1, MP_max * 0.10)
Stamina_regen = Stamina_max * 0.20
HP_regen = 0  // Non implémenté
```

---
