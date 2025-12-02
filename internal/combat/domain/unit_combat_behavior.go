package domain

import (
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// UnitCombatBehavior gère la logique de combat d'une unité
// Responsabilités: Dégâts, soins, état de vie/mort
// Single Responsibility Principle - Une seule raison de changer: logique de combat
type UnitCombatBehavior struct {
	baseStats    *shared.Stats
	currentStats *shared.Stats
	isEliminated bool
}

// NewUnitCombatBehavior crée un nouveau comportement de combat
func NewUnitCombatBehavior(baseStats *shared.Stats) *UnitCombatBehavior {
	return &UnitCombatBehavior{
		baseStats:    baseStats,
		currentStats: baseStats.Clone(),
		isEliminated: false,
	}
}

// BaseStats retourne les stats de base
func (c *UnitCombatBehavior) BaseStats() *shared.Stats {
	return c.baseStats
}

// CurrentStats retourne les stats actuelles
func (c *UnitCombatBehavior) CurrentStats() *shared.Stats {
	return c.currentStats
}

// IsEliminated vérifie si l'unité est éliminée
func (c *UnitCombatBehavior) IsEliminated() bool {
	return c.isEliminated
}

// TakeDamage applique des dégâts à l'unité
func (c *UnitCombatBehavior) TakeDamage(damage int) {
	if c.isEliminated {
		return // Unité déjà morte
	}

	c.currentStats.HP -= damage

	if c.currentStats.HP <= 0 {
		c.currentStats.HP = 0
		c.isEliminated = true
	}
}

// Heal applique un soin à l'unité
func (c *UnitCombatBehavior) Heal(healing int) {
	if c.isEliminated {
		return // Pas de soin si mort
	}

	c.currentStats.HP += healing

	// Cap aux HP max
	if c.currentStats.HP > c.baseStats.HP {
		c.currentStats.HP = c.baseStats.HP
	}
}

// Revive ressuscite l'unité avec un montant de HP
func (c *UnitCombatBehavior) Revive(hp int) {
	if !c.isEliminated {
		return // Déjà vivant
	}

	c.isEliminated = false
	c.currentStats.HP = hp

	// Cap aux HP max
	if c.currentStats.HP > c.baseStats.HP {
		c.currentStats.HP = c.baseStats.HP
	}
}

// RestoreMP restaure des points de mana
func (c *UnitCombatBehavior) RestoreMP(mp int) {
	c.currentStats.MP += mp

	if c.currentStats.MP > c.baseStats.MP {
		c.currentStats.MP = c.baseStats.MP
	}
}

// ConsumeMP consomme des points de mana
func (c *UnitCombatBehavior) ConsumeMP(mp int) error {
	if c.currentStats.MP < mp {
		return shared.NewDomainError("MP insuffisants", "INSUFFICIENT_MP")
	}

	c.currentStats.MP -= mp
	return nil
}

// ConsumeStamina consomme de l'endurance
func (c *UnitCombatBehavior) ConsumeStamina(stamina int) error {
	if c.currentStats.Stamina < stamina {
		return shared.NewDomainError("Stamina insuffisante", "INSUFFICIENT_STAMINA")
	}

	c.currentStats.Stamina -= stamina
	return nil
}

// FullRestore restaure complètement HP et MP
func (c *UnitCombatBehavior) FullRestore() {
	c.currentStats.HP = c.baseStats.HP
	c.currentStats.MP = c.baseStats.MP
	c.currentStats.Stamina = c.baseStats.Stamina
}

// CurrentHP retourne les HP actuels
func (c *UnitCombatBehavior) CurrentHP() int {
	return c.currentStats.HP
}

// CurrentMP retourne les MP actuels
func (c *UnitCombatBehavior) CurrentMP() int {
	return c.currentStats.MP
}

// HPPercentage retourne le pourcentage de HP restants
func (c *UnitCombatBehavior) HPPercentage() float64 {
	if c.baseStats.HP == 0 {
		return 0
	}
	return float64(c.currentStats.HP) / float64(c.baseStats.HP) * 100
}
