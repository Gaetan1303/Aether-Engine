package commands

import (
	"fmt"

	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
)

// WaitCommand représente l'action d'attendre (passer son tour)
// Réinitialise la jauge ATB à 0 et passe le tour
type WaitCommand struct {
	*BaseCommand
}

// NewWaitCommand crée une nouvelle commande d'attente
func NewWaitCommand(actor *domain.Unite, combat *domain.Combat) *WaitCommand {
	return &WaitCommand{
		BaseCommand: NewBaseCommand(actor, combat, CommandTypeWait),
	}
}

// Validate vérifie si l'unité peut attendre (toujours possible)
func (c *WaitCommand) Validate() error {
	return nil
}

// Execute passe le tour
func (c *WaitCommand) Execute() (*CommandResult, error) {
	result := &CommandResult{
		Success: true,
		Message: fmt.Sprintf("%s attend", c.actor.Nom()),
		Effects: []CommandEffect{},
	}

	return result, nil
}

// Rollback ne fait rien (attendre n'a pas d'effet à annuler)
func (c *WaitCommand) Rollback() error {
	return nil
}
