package states

import (
	"fmt"

	"github.com/aether-engine/aether-engine/internal/combat/domain/commands"
)

// ValidatingState valide la commande sélectionnée
// Utilise la Chain of Responsibility pour les validations
type ValidatingState struct {
	BaseState
}

// NewValidatingState crée un nouvel état Validating
func NewValidatingState() *ValidatingState {
	return &ValidatingState{
		BaseState: BaseState{
			name: "Validating",
			allowedTransitions: map[string]bool{
				"Confirmed":      true,
				"ActionRejected": true,
			},
		},
	}
}

// Enter lance la validation de la commande
func (s *ValidatingState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	if ctx.PendingCommand == nil {
		return fmt.Errorf("aucune commande à valider")
	}

	// Type assertion pour la commande
	cmd, ok := ctx.PendingCommand.(commands.Command)
	if !ok {
		return fmt.Errorf("type de commande invalide")
	}

	// Valider la commande via son interface Command
	if err := cmd.Validate(); err != nil {
		// Stocker l'erreur de validation
		ctx.ValidationError = err
		return err
	}

	return nil
}

// Exit est appelé lors de la sortie
func (s *ValidatingState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans Validating
func (s *ValidatingState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventValidationSuccess:
		// Validation réussie, confirmer l'action
		return NewConfirmedState(), nil

	case EventValidationFailed:
		// Validation échouée, rejeter l'action
		return NewActionRejectedState(), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}
