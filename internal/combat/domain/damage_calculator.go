package domain

import (
	"math"
)

// DamageCalculator est l'interface Strategy pour calculer les dégâts
// Strategy Pattern - permet d'avoir plusieurs algorithmes de calcul interchangeables
type DamageCalculator interface {
	Calculate(attacker *Unite, defender *Unite, competence *Competence) int
	CalculerDegats(attacker *Unite, defender *Unite, competence *Competence) int // Alias français
	GetType() string
}

// ===============================================
// Strategy 1: Dégâts Physiques
// ===============================================

// PhysicalDamageCalculator calcule les dégâts physiques
// Formule: (ATK de l'attaquant - DEF du défenseur) + bonus compétence
type PhysicalDamageCalculator struct{}

// NewPhysicalDamageCalculator crée une nouvelle instance
func NewPhysicalDamageCalculator() DamageCalculator {
	return &PhysicalDamageCalculator{}
}

func (c *PhysicalDamageCalculator) Calculate(attacker *Unite, defender *Unite, competence *Competence) int {
	// Stats de base
	attackerStats := attacker.Stats()
	defenderStats := defender.Stats()

	// Calcul de base: ATK - DEF
	baseDamage := attackerStats.ATK - defenderStats.DEF

	// Bonus de la compétence
	skillBonus := competence.DegatsBase()

	// Modificateur de la compétence (scaling)
	scaling := competence.Modificateur() * float64(attackerStats.ATK)

	// Dégâts totaux
	totalDamage := float64(baseDamage) + float64(skillBonus) + scaling

	// Variance aléatoire ±5% (simulé ici comme fixe, à randomiser plus tard)
	// totalDamage = totalDamage * (0.95 + rand.Float64()*0.1)

	// Minimum 1 dégât
	if totalDamage < 1 {
		totalDamage = 1
	}

	return int(math.Round(totalDamage))
}

func (c *PhysicalDamageCalculator) CalculerDegats(attacker *Unite, defender *Unite, competence *Competence) int {
	return c.Calculate(attacker, defender, competence)
}

func (c *PhysicalDamageCalculator) GetType() string {
	return "Physical"
}

// ===============================================
// Strategy 2: Dégâts Magiques
// ===============================================

// MagicalDamageCalculator calcule les dégâts magiques
// Formule: (MATK de l'attaquant - MDEF du défenseur) + bonus compétence + scaling MATK
type MagicalDamageCalculator struct{}

// NewMagicalDamageCalculator crée une nouvelle instance
func NewMagicalDamageCalculator() DamageCalculator {
	return &MagicalDamageCalculator{}
}

func (c *MagicalDamageCalculator) Calculate(attacker *Unite, defender *Unite, competence *Competence) int {
	// Stats de base
	attackerStats := attacker.Stats()
	defenderStats := defender.Stats()

	// Calcul de base: MATK - MDEF
	baseDamage := attackerStats.MATK - defenderStats.MDEF

	// Bonus de la compétence
	skillBonus := competence.DegatsBase()

	// Modificateur de la compétence (scaling sur MATK)
	scaling := competence.Modificateur() * float64(attackerStats.MATK)

	// Dégâts totaux
	totalDamage := float64(baseDamage) + float64(skillBonus) + scaling

	// Les sorts magiques ignorent 20% de la MDEF (pénétration magique)
	magicPenetration := float64(defenderStats.MDEF) * 0.2
	totalDamage += magicPenetration

	// Minimum 1 dégât
	if totalDamage < 1 {
		totalDamage = 1
	}

	return int(math.Round(totalDamage))
}

func (c *MagicalDamageCalculator) CalculerDegats(attacker *Unite, defender *Unite, competence *Competence) int {
	return c.Calculate(attacker, defender, competence)
}

func (c *MagicalDamageCalculator) GetType() string {
	return "Magical"
}

// ===============================================
// Strategy 3: Dégâts Fixes (Objets, Pièges)
// ===============================================

// FixedDamageCalculator inflige un montant fixe de dégâts
// Utilisé pour les objets (bombes, potions offensives) ou pièges
type FixedDamageCalculator struct {
	amount int
}

// NewFixedDamageCalculator crée une nouvelle instance
func NewFixedDamageCalculator(amount int) DamageCalculator {
	return &FixedDamageCalculator{amount: amount}
}

func (c *FixedDamageCalculator) Calculate(attacker *Unite, defender *Unite, competence *Competence) int {
	// Ignore complètement les stats - dégâts fixes
	return c.amount
}

func (c *FixedDamageCalculator) CalculerDegats(attacker *Unite, defender *Unite, competence *Competence) int {
	return c.Calculate(attacker, defender, competence)
}

func (c *FixedDamageCalculator) GetType() string {
	return "Fixed"
}

// ===============================================
// Strategy 4: Dégâts Hybrides (Physique + Magique)
// ===============================================

// HybridDamageCalculator combine dégâts physiques et magiques
// Utilisé pour les compétences hybrides (ex: lame élémentaire)
type HybridDamageCalculator struct {
	physicalRatio float64 // Ratio physique (0.0 à 1.0)
	magicalRatio  float64 // Ratio magique (0.0 à 1.0)
}

// NewHybridDamageCalculator crée une nouvelle instance
// Example: NewHybridDamageCalculator(0.6, 0.4) = 60% physique, 40% magique
func NewHybridDamageCalculator(physicalRatio, magicalRatio float64) DamageCalculator {
	return &HybridDamageCalculator{
		physicalRatio: physicalRatio,
		magicalRatio:  magicalRatio,
	}
}

func (c *HybridDamageCalculator) Calculate(attacker *Unite, defender *Unite, competence *Competence) int {
	attackerStats := attacker.Stats()
	defenderStats := defender.Stats()

	// Partie physique
	physicalDamage := float64(attackerStats.ATK-defenderStats.DEF) * c.physicalRatio

	// Partie magique
	magicalDamage := float64(attackerStats.MATK-defenderStats.MDEF) * c.magicalRatio

	// Bonus compétence
	skillBonus := competence.DegatsBase()

	// Scaling hybride
	scaling := competence.Modificateur() * (float64(attackerStats.ATK)*c.physicalRatio + float64(attackerStats.MATK)*c.magicalRatio)

	totalDamage := physicalDamage + magicalDamage + float64(skillBonus) + scaling

	// Minimum 1 dégât
	if totalDamage < 1 {
		totalDamage = 1
	}

	return int(math.Round(totalDamage))
}

func (c *HybridDamageCalculator) CalculerDegats(attacker *Unite, defender *Unite, competence *Competence) int {
	return c.Calculate(attacker, defender, competence)
}

func (c *HybridDamageCalculator) GetType() string {
	return "Hybrid"
}

// ===============================================
// Strategy 5: Dégâts Basés sur HP (Proportionnels)
// ===============================================

// ProportionalDamageCalculator inflige des dégâts basés sur les HP
// Utilisé pour les compétences "% HP" (ex: sacrifice, drain)
type ProportionalDamageCalculator struct {
	percentageOfTargetHP float64 // % des HP actuels de la cible
	useCurrentHP         bool    // true = HP actuels, false = HP max
}

// NewProportionalDamageCalculator crée une nouvelle instance
// Example: NewProportionalDamageCalculator(0.15, true) = 15% des HP actuels
func NewProportionalDamageCalculator(percentage float64, useCurrentHP bool) DamageCalculator {
	return &ProportionalDamageCalculator{
		percentageOfTargetHP: percentage,
		useCurrentHP:         useCurrentHP,
	}
}

func (c *ProportionalDamageCalculator) Calculate(attacker *Unite, defender *Unite, competence *Competence) int {
	defenderStats := defender.Stats()

	var targetHP int
	if c.useCurrentHP {
		targetHP = defender.HPActuels()
	} else {
		targetHP = defenderStats.HP // HP max
	}

	damage := float64(targetHP) * c.percentageOfTargetHP

	// Minimum 1 dégât
	if damage < 1 {
		damage = 1
	}

	return int(math.Round(damage))
}

func (c *ProportionalDamageCalculator) CalculerDegats(attacker *Unite, defender *Unite, competence *Competence) int {
	return c.Calculate(attacker, defender, competence)
}

func (c *ProportionalDamageCalculator) GetType() string {
	return "Proportional"
}

// ===============================================
// Strategy 6: Dégâts Critiques (Wrapper)
// ===============================================

// CriticalDamageCalculator est un décorateur qui ajoute la chance de critique
// Design Pattern: Decorator + Strategy
type CriticalDamageCalculator struct {
	baseCalculator DamageCalculator
	critChance     float64 // Chance de critique (0.0 à 1.0)
	critMultiplier float64 // Multiplicateur critique (ex: 1.5 = +50%)
}

// NewCriticalDamageCalculator wrap un calculator existant avec des critiques
func NewCriticalDamageCalculator(base DamageCalculator, critChance, critMultiplier float64) DamageCalculator {
	return &CriticalDamageCalculator{
		baseCalculator: base,
		critChance:     critChance,
		critMultiplier: critMultiplier,
	}
}

func (c *CriticalDamageCalculator) Calculate(attacker *Unite, defender *Unite, competence *Competence) int {
	// Calculer les dégâts de base avec la stratégie wrappée
	baseDamage := c.baseCalculator.Calculate(attacker, defender, competence)

	// TODO: Implémenter RNG pour critique (pour l'instant, pas de crit)
	// isCrit := rand.Float64() < c.critChance
	// if isCrit {
	//     return int(float64(baseDamage) * c.critMultiplier)
	// }

	return baseDamage
}

func (c *CriticalDamageCalculator) CalculerDegats(attacker *Unite, defender *Unite, competence *Competence) int {
	return c.Calculate(attacker, defender, competence)
}

func (c *CriticalDamageCalculator) GetType() string {
	return "Critical-" + c.baseCalculator.GetType()
}

// ===============================================
// Strategy Factory (Factory Pattern)
// ===============================================

// DamageCalculatorFactory crée les calculators selon le type de compétence
// Factory Pattern - centralise la création des strategies
type DamageCalculatorFactory struct{}

// NewDamageCalculatorFactory crée une nouvelle factory
func NewDamageCalculatorFactory() *DamageCalculatorFactory {
	return &DamageCalculatorFactory{}
}

// CreateCalculator crée le calculator approprié selon le type de compétence
func (f *DamageCalculatorFactory) CreateCalculator(competence *Competence) DamageCalculator {
	switch competence.Type() {
	case CompetenceAttaque:
		return NewPhysicalDamageCalculator()
	case CompetenceMagie:
		return NewMagicalDamageCalculator()
	case CompetenceSoin:
		// Les compétences de soin n'infligent pas de dégâts
		return NewFixedDamageCalculator(0)
	case CompetenceBuff, CompetenceDebuff:
		// Les buffs/debuffs n'infligent pas de dégâts directement
		return NewFixedDamageCalculator(0)
	default:
		// Par défaut, utiliser physique
		return NewPhysicalDamageCalculator()
	}
}

// CreateHybridCalculator crée un calculator hybride
func (f *DamageCalculatorFactory) CreateHybridCalculator(physicalRatio, magicalRatio float64) DamageCalculator {
	return NewHybridDamageCalculator(physicalRatio, magicalRatio)
}

// CreateProportionalCalculator crée un calculator proportionnel
func (f *DamageCalculatorFactory) CreateProportionalCalculator(percentage float64, useCurrentHP bool) DamageCalculator {
	return NewProportionalDamageCalculator(percentage, useCurrentHP)
}

// CreateWithCritical wrap un calculator avec des critiques
func (f *DamageCalculatorFactory) CreateWithCritical(base DamageCalculator, critChance, critMultiplier float64) DamageCalculator {
	return NewCriticalDamageCalculator(base, critChance, critMultiplier)
}

// ===============================================
// Helper: Integration dans Combat
// ===============================================

// CalculerDegats est une méthode helper pour Combat qui utilise le Strategy Pattern
func CalculerDegats(attacker *Unite, defender *Unite, competence *Competence, calculator DamageCalculator) int {
	if calculator == nil {
		// Factory par défaut
		factory := NewDamageCalculatorFactory()
		calculator = factory.CreateCalculator(competence)
	}

	return calculator.Calculate(attacker, defender, competence)
}
