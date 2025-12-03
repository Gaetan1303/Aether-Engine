package commands

import (
	"fmt"
	"math/rand"

	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
)

// FleeCommand représente la tentative de fuite du combat
// Probabilité basée sur la différence de niveau et SPD
type FleeCommand struct {
	*BaseCommand
	fleeSuccess bool
}

// NewFleeCommand crée une nouvelle commande de fuite
func NewFleeCommand(actor *domain.Unite, combat *domain.Combat) *FleeCommand {
	return &FleeCommand{
		BaseCommand: NewBaseCommand(actor, combat, CommandTypeFlee),
	}
}

// Validate vérifie si la fuite est possible
func (c *FleeCommand) Validate() error {
	// 1. Vérifier que l'acteur peut agir
	if !c.actor.PeutAgir() {
		return fmt.Errorf("l'unité %s ne peut pas agir", c.actor.Nom())
	}

	// 2. Vérifier si le combat autorise la fuite
	if !c.combat.FuiteAutorisee() {
		return fmt.Errorf("la fuite n'est pas autorisée dans ce combat (boss/arène)")
	}

	// 3. Vérifier le statut (Root empêche la fuite)
	if c.actor.EstRoot() {
		return fmt.Errorf("l'unité est enracinée, impossible de fuir")
	}

	return nil
}

// Execute tente de fuir
func (c *FleeCommand) Execute() (*CommandResult, error) {
	// Calculer la probabilité de fuite
	// Base: 50% + (SPD acteur - SPD moyenne ennemis) / 10
	probability := c.calculateFleeProbability()

	// Roll pour déterminer le succès
	roll := rand.Float64() * 100
	c.fleeSuccess = roll < probability

	result := &CommandResult{
		Success: c.fleeSuccess,
		Effects: []CommandEffect{},
	}

	if c.fleeSuccess {
		result.Message = fmt.Sprintf("%s a réussi à fuir! (probabilité: %.1f%%)", c.actor.Nom(), probability)

		// Marquer l'équipe comme ayant fui
		c.combat.MarquerEquipeFuite(c.actor.TeamID())
	} else {
		result.Message = fmt.Sprintf("%s n'a pas réussi à fuir (probabilité: %.1f%%)", c.actor.Nom(), probability)
	}

	return result, nil
}

// calculateFleeProbability calcule la probabilité de fuite
func (c *FleeCommand) calculateFleeProbability() float64 {
	baseProbability := 50.0

	// Obtenir la SPD moyenne des ennemis
	equipeActeur := c.actor.TeamID()
	ennemis := c.combat.ObtenirEnnemis(equipeActeur)

	if len(ennemis) == 0 {
		return 100.0 // Pas d'ennemi = fuite garantie
	}

	totalSPD := 0
	for _, ennemi := range ennemis {
		totalSPD += ennemi.Stats().SPD
	}
	spdMoyenne := float64(totalSPD) / float64(len(ennemis))

	// Ajuster selon la différence de SPD
	spdDiff := float64(c.actor.Stats().SPD) - spdMoyenne
	adjustment := spdDiff / 10.0

	probability := baseProbability + adjustment

	// Clamp entre 10% et 95%
	if probability < domain.FuiteProbabiliteMin {
		probability = domain.FuiteProbabiliteMin
	}
	if probability > domain.FuiteProbabiliteMax {
		probability = domain.FuiteProbabiliteMax
	}

	return probability
}

// Rollback annule la fuite (remet l'équipe en jeu)
func (c *FleeCommand) Rollback() error {
	if c.fleeSuccess {
		c.combat.AnnulerFuite(c.actor.TeamID())
	}
	return nil
}
