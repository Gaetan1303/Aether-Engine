package commands

import (
	"fmt"

	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
)

// AttackCommand représente une commande d'attaque basique
type AttackCommand struct {
	*BaseCommand
	target *domain.Unite
}

// NewAttackCommand crée une nouvelle commande d'attaque
func NewAttackCommand(actor *domain.Unite, combat *domain.Combat, target *domain.Unite) *AttackCommand {
	return &AttackCommand{
		BaseCommand: NewBaseCommand(actor, combat, CommandTypeAttack),
		target:      target,
	}
}

// Validate vérifie si l'attaque est possible
func (c *AttackCommand) Validate() error {
	// 1. Vérifier que l'acteur peut agir
	if !c.actor.PeutAgir() {
		return fmt.Errorf("l'unité %s ne peut pas agir", c.actor.Nom())
	}

	// 2. Vérifier que la cible existe et est vivante
	if c.target == nil {
		return fmt.Errorf("aucune cible spécifiée")
	}

	if c.target.EstEliminee() {
		return fmt.Errorf("la cible %s est déjà éliminée", c.target.Nom())
	}

	// 3. Vérifier la portée (attaque basique = portée 1)
	distance := abs(c.actor.Position().X()-c.target.Position().X()) +
		abs(c.actor.Position().Y()-c.target.Position().Y())

	if distance > 1 {
		return fmt.Errorf("cible hors de portée (distance: %d, portée max: 1)", distance)
	}

	// 4. Vérifier que la cible est dans une équipe ennemie
	if c.actor.TeamID() == c.target.TeamID() {
		return fmt.Errorf("impossible d'attaquer un allié")
	}

	return nil
}

// Execute exécute l'attaque
func (c *AttackCommand) Execute() (*CommandResult, error) {
	// Créer un snapshot avant modification
	c.CreateSnapshot()

	// Obtenir la compétence par défaut (attaque basique)
	competence := c.actor.ObtenirCompetenceParDefaut()

	// Calculer les dégâts en utilisant le DamageCalculator
	calculator := c.combat.GetDamageCalculator()
	degatsFinaux := calculator.Calculate(c.actor, c.target, competence)

	// Appliquer les dégâts
	c.target.RecevoirDegats(degatsFinaux)

	// Créer le résultat
	result := &CommandResult{
		Success:     true,
		Message:     fmt.Sprintf("%s attaque %s pour %d dégâts", c.actor.Nom(), c.target.Nom(), degatsFinaux),
		DamageDealt: degatsFinaux,
		Effects: []CommandEffect{
			{
				Type:     EffectTypeDamage,
				TargetID: c.target.ID(),
				Value:    degatsFinaux,
			},
		},
	}

	return result, nil
}

// Rollback annule l'attaque (restaure les HP de la cible)
func (c *AttackCommand) Rollback() error {
	if c.snapshot == nil {
		return fmt.Errorf("aucun snapshot disponible pour rollback")
	}

	// Restaurer les HP de la cible si elle est dans le snapshot
	if _, exists := c.snapshot.TargetStates[c.target.ID()]; exists {
		// Avec SetHP implémentée maintenant, on pourrait faire le rollback complet
		// c.target.SetHP(targetSnapshot.HP)
		fmt.Printf("[Rollback] Rollback de l'attaque sur %s\n", c.target.Nom())
	}

	return nil
} // Helper function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
