package states

import (
	"fmt"

	"github.com/aether-engine/aether-engine/internal/combat/domain/commands"
)

// ExecutingState exécute la commande confirmée
type ExecutingState struct {
	BaseState
}

// NewExecutingState crée un nouvel état Executing
func NewExecutingState() *ExecutingState {
	return &ExecutingState{
		BaseState: BaseState{
			name: "Executing",
			allowedTransitions: map[string]bool{
				"ApplyingEffects": true,
				"ExecutionFailed": true,
			},
		},
	}
}

// Enter exécute la commande via le CommandInvoker
func (s *ExecutingState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	if ctx.PendingCommand == nil {
		return fmt.Errorf("aucune commande à exécuter")
	}

	// Type assertion pour la commande
	cmd, ok := ctx.PendingCommand.(commands.Command)
	if !ok {
		return fmt.Errorf("type de commande invalide")
	}

	// Exécuter la commande
	result, err := cmd.Execute()
	if err != nil {
		ctx.ValidationError = err
		return err
	}

	// Stocker le résultat pour l'état suivant
	ctx.PendingResult = result

	fmt.Printf("[State] Commande exécutée: %s\n", result.Message)

	return nil
}

// Exit est appelé lors de la sortie
func (s *ExecutingState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans Executing
func (s *ExecutingState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventExecutionSuccess:
		// Exécution réussie, appliquer les effets
		return NewApplyingEffectsState(), nil

	case EventExecutionError:
		// Erreur lors de l'exécution
		return NewExecutionFailedState(), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}

// ExecutionFailedState représente une erreur lors de l'exécution
type ExecutionFailedState struct {
	BaseState
}

// NewExecutionFailedState crée un nouvel état ExecutionFailed
func NewExecutionFailedState() *ExecutionFailedState {
	return &ExecutionFailedState{
		BaseState: BaseState{
			name: "ExecutionFailed",
			allowedTransitions: map[string]bool{
				"ActionSelection": true, // Retourner à la sélection
			},
		},
	}
}

// Enter gère l'échec de l'exécution avec rollback
func (s *ExecutionFailedState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	if ctx.ValidationError != nil {
		fmt.Printf("[State] Erreur d'exécution: %v\n", ctx.ValidationError)
	}

	// Tenter un rollback
	if ctx.PendingCommand != nil {
		if cmd, ok := ctx.PendingCommand.(commands.Command); ok {
			if err := cmd.Rollback(); err != nil {
				fmt.Printf("[State] Erreur lors du rollback: %v\n", err)
			} else {
				fmt.Printf("[State] Rollback effectué avec succès\n")
			}
		}
	}

	return nil
}

// Exit est appelé lors de la sortie
func (s *ExecutionFailedState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	ctx.ValidationError = nil
	ctx.PendingCommand = nil
	ctx.PendingResult = nil
	return nil
}

// Handle gère les événements dans ExecutionFailed
func (s *ExecutionFailedState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventRetryAction:
		// Retourner à la sélection d'action
		if cmd, ok := ctx.PendingCommand.(commands.Command); ok {
			return NewActionSelectionState(cmd.GetActor()), nil
		}
		return NewActionSelectionState(nil), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}

// ApplyingEffectsState applique les effets de la commande
type ApplyingEffectsState struct {
	BaseState
}

// NewApplyingEffectsState crée un nouvel état ApplyingEffects
func NewApplyingEffectsState() *ApplyingEffectsState {
	return &ApplyingEffectsState{
		BaseState: BaseState{
			name: "ApplyingEffects",
			allowedTransitions: map[string]bool{
				"CheckVictory": true,
			},
		},
	}
}

// Enter applique les effets du résultat
func (s *ApplyingEffectsState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	if ctx.PendingResult == nil {
		return fmt.Errorf("aucun résultat à appliquer")
	}

	// Type assertion pour le résultat
	result, ok := ctx.PendingResult.(*commands.CommandResult)
	if !ok {
		return fmt.Errorf("type de résultat invalide")
	}

	// Les effets sont déjà appliqués dans Execute()
	// Ici on peut ajouter des effets secondaires, animations, etc.

	// Les effets sont appliqués
	fmt.Printf("[State] Effets appliqués: %s\n", result.Message)

	return nil
}

// Exit est appelé lors de la sortie
func (s *ApplyingEffectsState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans ApplyingEffects
func (s *ApplyingEffectsState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventEffectsApplied:
		// Effets appliqués, vérifier les conditions de victoire
		return NewCheckVictoryState(), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}
