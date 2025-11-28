# Test unitaire – Status (Version rigoureuse)

## Objectif
Gérer les effets temporaires (buffs/debuffs) appliqués aux unités en combat : poison, silence, hâte, bouclier, etc. Ces effets modifient le comportement des unités via des hooks dans le pipeline de résolution (Phase 3).

---

## Règles métier

### Contraintes
- **Durée limitée** : Chaque status a une durée en tours (≥ 0)
- **Intensité** : Certains status ont une intensité (ex: Poison inflige X dégâts/tour)
- **Unicité par type** : Une unité ne peut avoir qu'un status de chaque type actif
- **Expiration** : Durée décrémentée à chaque tour, status supprimé à 0
- **Immutabilité du type** : Le `StatusType` ne change jamais

### Types de Status (Exemples tactiques)

| Type | Catégorie | Effet | Intensité | Durée typique |
|------|-----------|-------|-----------|---------------|
| **Poison** | Debuff | Dégâts à chaque tour | Oui (dégâts/tour) | 3-5 tours |
| **Regen** | Buff | Soins à chaque tour | Oui (HP/tour) | 3-5 tours |
| **Silence** | Debuff | Bloque les sorts | Non | 2-3 tours |
| **Haste** | Buff | +50% ATB speed | Non | 3 tours |
| **Slow** | Debuff | -50% ATB speed | Non | 3 tours |
| **Shield** | Buff | Absorbe X dégâts | Oui (HP absorbés) | Jusqu'à rupture |
| **Berserk** | Mixed | +ATK mais perte contrôle | Oui (+% ATK) | 4 tours |
| **Stun** | Debuff | Skip le tour | Non | 1 tour |
| **Blind** | Debuff | Réduit la précision | Oui (-% accuracy) | 3 tours |
| **Protect** | Buff | Réduit dégâts physiques | Oui (-% dégâts) | 5 tours |

### Hooks pour le Pipeline (Phase 3)
Les Status s'intègrent dans le pipeline de résolution via des hooks :

1. **OnTurnStart(unit)** : Déclenché en début de tour (ex: Poison inflige dégâts)
2. **OnTurnEnd(unit)** : Déclenché en fin de tour (décrémente durée)
3. **OnIncomingDamage(damage, source)** : Modifie les dégâts reçus (ex: Shield absorbe)
4. **OnOutgoingDamage(damage, target)** : Modifie les dégâts infligés (ex: Berserk boost)
5. **OnActionAttempt(action)** : Valide si l'action est autorisée (ex: Silence bloque sorts)
6. **OnExpire(unit)** : Nettoyage à l'expiration

---

## Structure proposée

### 1. StatusType (Value Object - Enum)

```go
package domain

// StatusType représente le type d'effet de statut
type StatusType int

const (
    StatusPoison StatusType = iota
    StatusRegen
    StatusSilence
    StatusHaste
    StatusSlow
    StatusShield
    StatusBerserk
    StatusStun
    StatusBlind
    StatusProtect
)

var statusTypeNames = map[StatusType]string{
    StatusPoison:  "Poison",
    StatusRegen:   "Regen",
    StatusSilence: "Silence",
    StatusHaste:   "Haste",
    StatusSlow:    "Slow",
    StatusShield:  "Shield",
    StatusBerserk: "Berserk",
    StatusStun:    "Stun",
    StatusBlind:   "Blind",
    StatusProtect: "Protect",
}

// String retourne le nom du status
func (st StatusType) String() string {
    if name, ok := statusTypeNames[st]; ok {
        return name
    }
    return "Unknown"
}

// IsDebuff indique si le status est négatif
func (st StatusType) IsDebuff() bool {
    return st == StatusPoison || st == StatusSilence || 
           st == StatusSlow || st == StatusStun || st == StatusBlind
}

// IsBuff indique si le status est positif
func (st StatusType) IsBuff() bool {
    return st == StatusRegen || st == StatusHaste || 
           st == StatusShield || st == StatusProtect
}
```

### 2. Status (Entity)

```go
package domain

import (
    "errors"
    "fmt"
)

// Status représente un effet temporaire sur une unité
type Status struct {
    statusType StatusType
    duration   int // Tours restants
    intensity  int // Puissance de l'effet (optionnel)
    sourceID   UnitID // Qui a appliqué le status (pour tracking)
}

// NewStatus crée un nouveau status validé
func NewStatus(statusType StatusType, duration, intensity int, sourceID UnitID) (Status, error) {
    if duration < 0 {
        return Status{}, errors.New("duration cannot be negative")
    }
    if duration == 0 {
        return Status{}, errors.New("duration must be positive (use Remove for instant expiration)")
    }
    if intensity < 0 {
        return Status{}, errors.New("intensity cannot be negative")
    }
    
    return Status{
        statusType: statusType,
        duration:   duration,
        intensity:  intensity,
        sourceID:   sourceID,
    }, nil
}

// Getters
func (s Status) Type() StatusType { return s.statusType }
func (s Status) Duration() int    { return s.duration }
func (s Status) Intensity() int   { return s.intensity }
func (s Status) SourceID() UnitID { return s.sourceID }

// DecrementDuration réduit la durée de 1 tour
func (s *Status) DecrementDuration() {
    if s.duration > 0 {
        s.duration--
    }
}

// IsExpired vérifie si le status a expiré
func (s Status) IsExpired() bool {
    return s.duration == 0
}

// String implémente fmt.Stringer pour le logging
func (s Status) String() string {
    if s.intensity > 0 {
        return fmt.Sprintf("%s(duration=%d, intensity=%d)", 
            s.statusType, s.duration, s.intensity)
    }
    return fmt.Sprintf("%s(duration=%d)", s.statusType, s.duration)
}

// Equals compare deux status (même type et source)
func (s Status) Equals(other Status) bool {
    return s.statusType == other.statusType && 
           s.sourceID.Equals(other.sourceID)
}
```

### 3. StatusCollection (Agrégat helper pour Unit)

```go
package domain

import "errors"

// StatusCollection gère les statuts actifs d'une unité
type StatusCollection struct {
    statuses map[StatusType]Status
}

// NewStatusCollection crée une collection vide
func NewStatusCollection() *StatusCollection {
    return &StatusCollection{
        statuses: make(map[StatusType]Status),
    }
}

// Add ajoute ou remplace un status
func (sc *StatusCollection) Add(status Status) error {
    if status.IsExpired() {
        return errors.New("cannot add expired status")
    }
    
    // Si existe déjà, on remplace (refresh)
    sc.statuses[status.Type()] = status
    return nil
}

// Remove supprime un status par type
func (sc *StatusCollection) Remove(statusType StatusType) {
    delete(sc.statuses, statusType)
}

// Get récupère un status par type
func (sc *StatusCollection) Get(statusType StatusType) (Status, bool) {
    status, exists := sc.statuses[statusType]
    return status, exists
}

// Has vérifie si un status est actif
func (sc *StatusCollection) Has(statusType StatusType) bool {
    _, exists := sc.statuses[statusType]
    return exists
}

// All retourne tous les status actifs
func (sc *StatusCollection) All() []Status {
    result := make([]Status, 0, len(sc.statuses))
    for _, status := range sc.statuses {
        result = append(result, status)
    }
    return result
}

// DecrementAll réduit la durée de tous les status
func (sc *StatusCollection) DecrementAll() []StatusType {
    expired := make([]StatusType, 0)
    
    for statusType, status := range sc.statuses {
        status.DecrementDuration()
        
        if status.IsExpired() {
            expired = append(expired, statusType)
            delete(sc.statuses, statusType)
        } else {
            sc.statuses[statusType] = status // Mise à jour
        }
    }
    
    return expired
}

// Count retourne le nombre de status actifs
func (sc *StatusCollection) Count() int {
    return len(sc.statuses)
}

// Clear supprime tous les status (mort, purge magique)
func (sc *StatusCollection) Clear() {
    sc.statuses = make(map[StatusType]Status)
}
```

---

## Couverture des tests

| Catégorie | Tests |
|-----------|-------|
| StatusType | 3 |
| Création Status | 5 |
| Durée | 3 |
| String/Logging | 2 |
| Égalité | 3 |
| Collection - Add/Remove | 4 |
| Collection - Get/Has | 3 |
| Collection - DecrementAll | 1 |
| Collection - Clear/Count | 2 |
| Intégration | 2 |
| **TOTAL** | **28 tests** |

---
## Résultats obtenus :

=== RUN   TestStatusType_String
--- PASS: TestStatusType_String (0.00s)
=== RUN   TestStatusType_IsDebuff
--- PASS: TestStatusType_IsDebuff (0.00s)
=== RUN   TestStatusType_IsBuff
--- PASS: TestStatusType_IsBuff (0.00s)
=== RUN   TestNewStatus_Valid
--- PASS: TestNewStatus_Valid (0.00s)
=== RUN   TestNewStatus_ZeroDuration
--- PASS: TestNewStatus_ZeroDuration (0.00s)
=== RUN   TestNewStatus_NegativeDuration
--- PASS: TestNewStatus_NegativeDuration (0.00s)
=== RUN   TestNewStatus_NegativeIntensity
--- PASS: TestNewStatus_NegativeIntensity (0.00s)
=== RUN   TestNewStatus_ZeroIntensity
--- PASS: TestNewStatus_ZeroIntensity (0.00s)
=== RUN   TestStatus_DecrementDuration
--- PASS: TestStatus_DecrementDuration (0.00s)
=== RUN   TestStatus_DecrementDuration_AtZero
--- PASS: TestStatus_DecrementDuration_AtZero (0.00s)
=== RUN   TestStatus_IsExpired
--- PASS: TestStatus_IsExpired (0.00s)
=== RUN   TestStatus_String_WithIntensity
--- PASS: TestStatus_String_WithIntensity (0.00s)
=== RUN   TestStatus_String_WithoutIntensity
--- PASS: TestStatus_String_WithoutIntensity (0.00s)
=== RUN   TestStatus_Equals_SameTypeAndSource
--- PASS: TestStatus_Equals_SameTypeAndSource (0.00s)
=== RUN   TestStatus_Equals_DifferentType
--- PASS: TestStatus_Equals_DifferentType (0.00s)
=== RUN   TestStatus_Equals_DifferentSource
--- PASS: TestStatus_Equals_DifferentSource (0.00s)
=== RUN   TestStatusCollection_Add
--- PASS: TestStatusCollection_Add (0.00s)
=== RUN   TestStatusCollection_Add_ReplaceExisting
--- PASS: TestStatusCollection_Add_ReplaceExisting (0.00s)
=== RUN   TestStatusCollection_Add_ExpiredStatus
--- PASS: TestStatusCollection_Add_ExpiredStatus (0.00s)
=== RUN   TestStatusCollection_Remove
--- PASS: TestStatusCollection_Remove (0.00s)
=== RUN   TestStatusCollection_Get
--- PASS: TestStatusCollection_Get (0.00s)
=== RUN   TestStatusCollection_Get_NotFound
--- PASS: TestStatusCollection_Get_NotFound (0.00s)
=== RUN   TestStatusCollection_Has
--- PASS: TestStatusCollection_Has (0.00s)
=== RUN   TestStatusCollection_All
--- PASS: TestStatusCollection_All (0.00s)
=== RUN   TestStatusCollection_DecrementAll
--- PASS: TestStatusCollection_DecrementAll (0.00s)
=== RUN   TestStatusCollection_Clear
--- PASS: TestStatusCollection_Clear (0.00s)
=== RUN   TestStatusCollection_Count
--- PASS: TestStatusCollection_Count (0.00s)
=== RUN   TestStatusCollection_MultipleStatusLifecycle
--- PASS: TestStatusCollection_MultipleStatusLifecycle (0.00s)
=== RUN   TestStatusCollection_RefreshDuration
--- PASS: TestStatusCollection_RefreshDuration (0.00s)
PASS
ok      aether-engine-server/internal/combat/domain     0.003s
