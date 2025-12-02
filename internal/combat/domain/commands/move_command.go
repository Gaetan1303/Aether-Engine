package commands

import (
	"fmt"

	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// MoveCommand représente une commande de déplacement
// Utilise le PathfindingService du Step B
type MoveCommand struct {
	*BaseCommand
	targetPosition *shared.Position
	path           []*shared.Position
	cost           int
}

// NewMoveCommand crée une nouvelle commande de déplacement
func NewMoveCommand(actor *domain.Unite, combat *domain.Combat, targetPosition *shared.Position) *MoveCommand {
	return &MoveCommand{
		BaseCommand:    NewBaseCommand(actor, combat, CommandTypeMove),
		targetPosition: targetPosition,
	}
}

// Validate vérifie si le déplacement est possible
func (c *MoveCommand) Validate() error {
	// 1. Vérifier que l'acteur peut se déplacer (pas Root/Stun)
	if c.actor.EstBloqueDeplacement() {
		return fmt.Errorf("l'unité %s ne peut pas se déplacer (statut bloquant)", c.actor.Nom())
	}

	// 2. Vérifier que la position cible est valide
	if c.targetPosition == nil {
		return fmt.Errorf("position cible non spécifiée")
	}

	grille := c.combat.Grille()
	if !grille.EstDansLimites(c.targetPosition) {
		return fmt.Errorf("position cible hors limites")
	}

	// 3. Calculer le chemin avec pathfinding
	pathfindingService := domain.NewPathfindingService()
	pathfindingService.SetStrategyType("manhattan")

	// Créer la map des positions occupées (excluant l'acteur)
	unitesOccupees := c.combat.ObtenirPositionsOccupees(c.actor.ID())

	// Calculer le chemin avec portée
	porteeMax := c.actor.Stats().MOV
	path, cost, err := pathfindingService.TrouverCheminAvecPortee(
		grille,
		c.actor.Position(),
		c.targetPosition,
		unitesOccupees,
		porteeMax,
	)

	if err != nil {
		return fmt.Errorf("impossible de trouver un chemin: %w", err)
	}

	c.path = path
	c.cost = cost

	return nil
}

// Execute déplace l'unité
func (c *MoveCommand) Execute() (*CommandResult, error) {
	// Créer un snapshot avant modification
	c.CreateSnapshot()

	// Déplacer l'unité
	c.actor.DeplacerVers(c.targetPosition)

	// Créer le résultat
	result := &CommandResult{
		Success:      true,
		Message:      fmt.Sprintf("%s se déplace vers (%d,%d)", c.actor.Nom(), c.targetPosition.X(), c.targetPosition.Y()),
		CostMovement: c.cost,
		Effects: []CommandEffect{
			{
				Type:     EffectTypeMovement,
				TargetID: c.actor.ID(),
				Position: c.targetPosition,
			},
		},
	}

	return result, nil
}

// Rollback annule le déplacement
func (c *MoveCommand) Rollback() error {
	if c.snapshot == nil {
		return fmt.Errorf("aucun snapshot disponible pour rollback")
	}

	// Restaurer la position précédente
	c.actor.DeplacerVers(c.snapshot.ActorPosition)
	return nil
}
