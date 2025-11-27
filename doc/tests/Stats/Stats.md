Test unitaire – Stats (Version rigoureuse)
Objectif
Gérer les statistiques d'une unité de combat : points de vie, points de magie, attributs offensifs/défensifs, et leur modification dynamique en combat.

Règles métier
Statistiques de base

HP (Health Points) : Points de vie maximum (> 0)
MP (Magic Points) : Points de magie maximum (≥ 0, certaines classes n'ont pas de MP)
ATK (Attack) : Puissance physique (≥ 1)
DEF (Defense) : Résistance physique (≥ 0)
SPD (Speed) : Vitesse (≥ 1, détermine l'initiative CT/ATB)
MAG (Magic) : Puissance magique (≥ 0)
RES (Resistance) : Résistance magique (≥ 0)

Statistiques dynamiques

CurrentHP : HP actuels (0 ≤ CurrentHP ≤ MaxHP)
CurrentMP : MP actuels (0 ≤ CurrentMP ≤ MaxMP)

Contraintes

Les stats de base ne peuvent pas être négatives
CurrentHP/CurrentMP ne peuvent jamais dépasser les max
Une unité avec CurrentHP = 0 est KO
Les modificateurs temporaires (buffs/debuffs) sont gérés ailleurs (Status)

Comportements attendus

Création : Validation des stats de base
Dégâts : Réduction de CurrentHP avec plancher à 0
Soins : Augmentation de CurrentHP avec plafond à MaxHP
Coût MP : Réduction de CurrentMP pour les compétences
Régénération : Récupération de MP
État KO : Détection si CurrentHP = 0


#### Structure proposée : 
```go
package domain

import "errors"

// Stats représente les statistiques d'une unité
type Stats struct {
    // Stats de base (immutables une fois créées)
    maxHP int
    maxMP int
    atk   int
    def   int
    spd   int
    mag   int
    res   int
    
    // Stats dynamiques (mutables en combat)
    currentHP int
    currentMP int
}

// NewStats crée des stats validées
func NewStats(maxHP, maxMP, atk, def, spd, mag, res int) (Stats, error) {
    if maxHP <= 0 {
        return Stats{}, errors.New("maxHP must be positive")
    }
    if maxMP < 0 {
        return Stats{}, errors.New("maxMP cannot be negative")
    }
    if atk < 1 {
        return Stats{}, errors.New("atk must be at least 1")
    }
    if def < 0 {
        return Stats{}, errors.New("def cannot be negative")
    }
    if spd < 1 {
        return Stats{}, errors.New("spd must be at least 1")
    }
    if mag < 0 {
        return Stats{}, errors.New("mag cannot be negative")
    }
    if res < 0 {
        return Stats{}, errors.New("res cannot be negative")
    }
    
    return Stats{
        maxHP:     maxHP,
        maxMP:     maxMP,
        atk:       atk,
        def:       def,
        spd:       spd,
        mag:       mag,
        res:       res,
        currentHP: maxHP, // Commence au max
        currentMP: maxMP,
    }, nil
}

// Getters
func (s Stats) MaxHP() int     { return s.maxHP }
func (s Stats) MaxMP() int     { return s.maxMP }
func (s Stats) ATK() int       { return s.atk }
func (s Stats) DEF() int       { return s.def }
func (s Stats) SPD() int       { return s.spd }
func (s Stats) MAG() int       { return s.mag }
func (s Stats) RES() int       { return s.res }
func (s Stats) CurrentHP() int { return s.currentHP }
func (s Stats) CurrentMP() int { return s.currentMP }

// TakeDamage réduit les HP actuels (plancher à 0)
func (s *Stats) TakeDamage(amount int) {
    if amount < 0 {
        return // Ignorer les valeurs négatives
    }
    s.currentHP -= amount
    if s.currentHP < 0 {
        s.currentHP = 0
    }
}

// Heal restaure les HP (plafond à MaxHP)
func (s *Stats) Heal(amount int) {
    if amount < 0 {
        return
    }
    s.currentHP += amount
    if s.currentHP > s.maxHP {
        s.currentHP = s.maxHP
    }
}

// ConsumeMP réduit les MP actuels
func (s *Stats) ConsumeMP(amount int) error {
    if amount < 0 {
        return errors.New("MP consumption cannot be negative")
    }
    if s.currentMP < amount {
        return errors.New("insufficient MP")
    }
    s.currentMP -= amount
    return nil
}

// RestoreMP restaure les MP (plafond à MaxMP)
func (s *Stats) RestoreMP(amount int) {
    if amount < 0 {
        return
    }
    s.currentMP += amount
    if s.currentMP > s.maxMP {
        s.currentMP = s.maxMP
    }
}

// IsKO vérifie si l'unité est KO
func (s Stats) IsKO() bool {
    return s.currentHP == 0
}

// HPPercentage retourne le % de HP restants
func (s Stats) HPPercentage() float64 {
    if s.maxHP == 0 {
        return 0
    }
    return float64(s.currentHP) / float64(s.maxHP) * 100
}

// MPPercentage retourne le % de MP restants
func (s Stats) MPPercentage() float64 {
    if s.maxMP == 0 {
        return 0
    }
    return float64(s.currentMP) / float64(s.maxMP) * 100
}

// FullRestore restaure HP et MP au maximum
func (s *Stats) FullRestore() {
    s.currentHP = s.maxHP
    s.currentMP = s.maxMP
}
```

Tests unitaires complets

```go
package domain_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "aether-engine-server/internal/combat/domain"
)

// ========== Tests de création ==========

func TestNewStats_Valid(t *testing.T) {
    stats, err := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    assert.NoError(t, err)
    assert.Equal(t, 100, stats.MaxHP())
    assert.Equal(t, 50, stats.MaxMP())
    assert.Equal(t, 10, stats.ATK())
    assert.Equal(t, 5, stats.DEF())
    assert.Equal(t, 8, stats.SPD())
    assert.Equal(t, 12, stats.MAG())
    assert.Equal(t, 6, stats.RES())
    assert.Equal(t, 100, stats.CurrentHP())
    assert.Equal(t, 50, stats.CurrentMP())
}

func TestNewStats_InvalidMaxHP_Zero(t *testing.T) {
    _, err := domain.NewStats(0, 50, 10, 5, 8, 12, 6)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "maxHP must be positive")
}

func TestNewStats_InvalidMaxHP_Negative(t *testing.T) {
    _, err := domain.NewStats(-10, 50, 10, 5, 8, 12, 6)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "maxHP must be positive")
}

func TestNewStats_InvalidMaxMP_Negative(t *testing.T) {
    _, err := domain.NewStats(100, -5, 10, 5, 8, 12, 6)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "maxMP cannot be negative")
}

func TestNewStats_ValidMaxMP_Zero(t *testing.T) {
    // Certaines classes n'ont pas de MP (guerriers pure physique)
    stats, err := domain.NewStats(100, 0, 10, 5, 8, 0, 6)
    
    assert.NoError(t, err)
    assert.Equal(t, 0, stats.MaxMP())
}

func TestNewStats_InvalidATK_Zero(t *testing.T) {
    _, err := domain.NewStats(100, 50, 0, 5, 8, 12, 6)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "atk must be at least 1")
}

func TestNewStats_InvalidATK_Negative(t *testing.T) {
    _, err := domain.NewStats(100, 50, -5, 5, 8, 12, 6)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "atk must be at least 1")
}

func TestNewStats_InvalidDEF_Negative(t *testing.T) {
    _, err := domain.NewStats(100, 50, 10, -2, 8, 12, 6)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "def cannot be negative")
}

func TestNewStats_InvalidSPD_Zero(t *testing.T) {
    _, err := domain.NewStats(100, 50, 10, 5, 0, 12, 6)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "spd must be at least 1")
}

func TestNewStats_InvalidMAG_Negative(t *testing.T) {
    _, err := domain.NewStats(100, 50, 10, 5, 8, -3, 6)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "mag cannot be negative")
}

func TestNewStats_InvalidRES_Negative(t *testing.T) {
    _, err := domain.NewStats(100, 50, 10, 5, 8, 12, -1)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "res cannot be negative")
}

// ========== Tests de dégâts ==========

func TestStats_TakeDamage_Normal(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    stats.TakeDamage(30)
    
    assert.Equal(t, 70, stats.CurrentHP())
}

func TestStats_TakeDamage_Multiple(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    stats.TakeDamage(20)
    stats.TakeDamage(15)
    
    assert.Equal(t, 65, stats.CurrentHP())
}

func TestStats_TakeDamage_Overkill(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    stats.TakeDamage(150)
    
    assert.Equal(t, 0, stats.CurrentHP())
    assert.True(t, stats.IsKO())
}

func TestStats_TakeDamage_ExactlyLethal(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    stats.TakeDamage(100)
    
    assert.Equal(t, 0, stats.CurrentHP())
    assert.True(t, stats.IsKO())
}

func TestStats_TakeDamage_Negative(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    stats.TakeDamage(-10)
    
    // Ne doit pas soigner
    assert.Equal(t, 100, stats.CurrentHP())
}

func TestStats_TakeDamage_Zero(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    stats.TakeDamage(0)
    
    assert.Equal(t, 100, stats.CurrentHP())
}

// ========== Tests de soins ==========

func TestStats_Heal_Normal(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.TakeDamage(40)
    
    stats.Heal(20)
    
    assert.Equal(t, 80, stats.CurrentHP())
}

func TestStats_Heal_Overheal(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.TakeDamage(20)
    
    stats.Heal(50)
    
    // Ne peut pas dépasser MaxHP
    assert.Equal(t, 100, stats.CurrentHP())
}

func TestStats_Heal_FullHP(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    stats.Heal(30)
    
    // Déjà au max, reste à 100
    assert.Equal(t, 100, stats.CurrentHP())
}

func TestStats_Heal_FromKO(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.TakeDamage(100)
    
    stats.Heal(50)
    
    assert.Equal(t, 50, stats.CurrentHP())
    assert.False(t, stats.IsKO())
}

func TestStats_Heal_Negative(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.TakeDamage(30)
    
    stats.Heal(-10)
    
    // Ne doit pas infliger de dégâts
    assert.Equal(t, 70, stats.CurrentHP())
}

// ========== Tests de consommation MP ==========

func TestStats_ConsumeMP_Valid(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    err := stats.ConsumeMP(20)
    
    assert.NoError(t, err)
    assert.Equal(t, 30, stats.CurrentMP())
}

func TestStats_ConsumeMP_ExactlyAll(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    err := stats.ConsumeMP(50)
    
    assert.NoError(t, err)
    assert.Equal(t, 0, stats.CurrentMP())
}

func TestStats_ConsumeMP_Insufficient(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    err := stats.ConsumeMP(60)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "insufficient MP")
    assert.Equal(t, 50, stats.CurrentMP()) // Ne change pas
}

func TestStats_ConsumeMP_AfterPartialUse(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.ConsumeMP(30)
    
    err := stats.ConsumeMP(25)
    
    assert.Error(t, err)
    assert.Equal(t, 20, stats.CurrentMP())
}

func TestStats_ConsumeMP_Negative(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    err := stats.ConsumeMP(-10)
    
    assert.Error(t, err)
    assert.Equal(t, 50, stats.CurrentMP())
}

func TestStats_ConsumeMP_Zero(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    err := stats.ConsumeMP(0)
    
    assert.NoError(t, err)
    assert.Equal(t, 50, stats.CurrentMP())
}

// ========== Tests de restauration MP ==========

func TestStats_RestoreMP_Normal(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.ConsumeMP(30)
    
    stats.RestoreMP(15)
    
    assert.Equal(t, 35, stats.CurrentMP())
}

func TestStats_RestoreMP_Overflow(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.ConsumeMP(10)
    
    stats.RestoreMP(50)
    
    assert.Equal(t, 50, stats.CurrentMP())
}

func TestStats_RestoreMP_FromZero(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.ConsumeMP(50)
    
    stats.RestoreMP(25)
    
    assert.Equal(t, 25, stats.CurrentMP())
}

func TestStats_RestoreMP_Negative(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.ConsumeMP(10)
    
    stats.RestoreMP(-5)
    
    assert.Equal(t, 40, stats.CurrentMP())
}

// ========== Tests d'état KO ==========

func TestStats_IsKO_Alive(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    assert.False(t, stats.IsKO())
}

func TestStats_IsKO_AfterDamage(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.TakeDamage(50)
    
    assert.False(t, stats.IsKO())
}

func TestStats_IsKO_Dead(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.TakeDamage(100)
    
    assert.True(t, stats.IsKO())
}

// ========== Tests de pourcentages ==========

func TestStats_HPPercentage_Full(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    percentage := stats.HPPercentage()
    
    assert.Equal(t, 100.0, percentage)
}

func TestStats_HPPercentage_Half(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.TakeDamage(50)
    
    percentage := stats.HPPercentage()
    
    assert.Equal(t, 50.0, percentage)
}

func TestStats_HPPercentage_Zero(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.TakeDamage(100)
    
    percentage := stats.HPPercentage()
    
    assert.Equal(t, 0.0, percentage)
}

func TestStats_MPPercentage_Full(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    
    percentage := stats.MPPercentage()
    
    assert.Equal(t, 100.0, percentage)
}

func TestStats_MPPercentage_Quarter(t *testing.T) {
    stats, _ := domain.NewStats(100, 100, 10, 5, 8, 12, 6)
    stats.ConsumeMP(75)
    
    percentage := stats.MPPercentage()
    
    assert.Equal(t, 25.0, percentage)
}

func TestStats_MPPercentage_NoMP(t *testing.T) {
    stats, _ := domain.NewStats(100, 0, 10, 5, 8, 12, 6)
    
    percentage := stats.MPPercentage()
    
    assert.Equal(t, 0.0, percentage)
}

// ========== Tests de restauration complète ==========

func TestStats_FullRestore(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.TakeDamage(60)
    stats.ConsumeMP(30)
    
    stats.FullRestore()
    
    assert.Equal(t, 100, stats.CurrentHP())
    assert.Equal(t, 50, stats.CurrentMP())
}

func TestStats_FullRestore_FromKO(t *testing.T) {
    stats, _ := domain.NewStats(100, 50, 10, 5, 8, 12, 6)
    stats.TakeDamage(100)
    stats.ConsumeMP(50)
    
    stats.FullRestore()
    
    assert.Equal(t, 100, stats.CurrentHP())
    assert.Equal(t, 50, stats.CurrentMP())
    assert.False(t, stats.IsKO())
}

// ========== Tests de cas limites ==========

func TestStats_MinimalWarrior(t *testing.T) {
    // Guerrier sans magie
    stats, err := domain.NewStats(50, 0, 15, 10, 5, 0, 3)
    
    assert.NoError(t, err)
    assert.Equal(t, 0, stats.MaxMP())
    assert.Equal(t, 0, stats.CurrentMP())
}

func TestStats_GlassCannon(t *testing.T) {
    // Mage avec peu de défense
    stats, err := domain.NewStats(30, 100, 1, 0, 6, 25, 0)
    
    assert.NoError(t, err)
    assert.Equal(t, 1, stats.ATK())
    assert.Equal(t, 0, stats.DEF())
    assert.Equal(t, 25, stats.MAG())
}

func TestStats_Tank(t *testing.T) {
    // Tank avec haute défense
    stats, err := domain.NewStats(200, 20, 8, 40, 3, 5, 30)
    
    assert.NoError(t, err)
    assert.Equal(t, 200, stats.MaxHP())
    assert.Equal(t, 40, stats.DEF())
    assert.Equal(t, 30, stats.RES())
}
```