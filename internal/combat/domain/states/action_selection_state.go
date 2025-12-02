package states

import (
	"fmt"

	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/aether-engine/aether-engine/internal/combat/domain/commands"
)

// ActionSelectionState permet au joueur de choisir une action
// En attente d'input: Move, Attack, Skill, Item, Flee, Wait
type ActionSelectionState struct {
	BaseState
	currentUnit *domain.Unite
}

// NewActionSelectionState crée un nouvel état ActionSelection
func NewActionSelectionState(currentUnit *domain.Unite) *ActionSelectionState {
	return &ActionSelectionState{
		BaseState: BaseState{
			name: "ActionSelection",
			allowedTransitions: map[string]bool{
				"Validating": true,
				"TurnEnd":    true, // Wait command passe directement à TurnEnd
			},
		},
		currentUnit: currentUnit,
	}
}

// Enter initialise la sélection d'action
func (s *ActionSelectionState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s pour %s\n", s.Name(), s.currentUnit.Nom())

	// Si c'est une IA, déclencher automatiquement la sélection d'action
	if s.currentUnit.EstIA() {
		// L'IA sélectionnera automatiquement une action via EventCommandSelected
		s.currentUnit.IAChoisirAction(ctx.Combat)
	}

	return nil
}

// Exit est appelé lors de la sortie
func (s *ActionSelectionState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans ActionSelection
func (s *ActionSelectionState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventCommandSelected:
		// Le joueur ou l'IA a choisi une commande
		// Stocker la commande dans le contexte
		if cmd, ok := event.Data.(commands.Command); ok {
			ctx.PendingCommand = cmd

			// Si c'est Wait, passer directement à TurnEnd
			if cmd.GetType() == commands.CommandTypeWait {
				return NewTurnEndState(), nil
			}

			// Sinon, passer à la validation
			return NewValidatingState(), nil
		}
		return nil, fmt.Errorf("données de commande invalides")

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}

// CurrentUnit retourne l'unité active
func (s *ActionSelectionState) CurrentUnit() *domain.Unite {
	return s.currentUnit
}
