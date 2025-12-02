package states

import (
	"fmt"
)

// ConfirmedState représente une action confirmée prête à être exécutée
type ConfirmedState struct {
	BaseState
}

// NewConfirmedState crée un nouvel état Confirmed
func NewConfirmedState() *ConfirmedState {
	return &ConfirmedState{
		BaseState: BaseState{
			name: "Confirmed",
			allowedTransitions: map[string]bool{
				"Executing": true,
			},
		},
	}
}

// Enter prépare l'exécution
func (s *ConfirmedState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	// L'action est confirmée, prête à être exécutée
	// Notifier les observateurs
	s.notifyObservers(ctx, "ActionConfirmed")

	return nil
}

// Exit est appelé lors de la sortie
func (s *ConfirmedState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans Confirmed
func (s *ConfirmedState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventBeginExecution:
		// Commencer l'exécution
		return NewExecutingState(), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}

// notifyObservers notifie tous les observateurs
func (s *ConfirmedState) notifyObservers(ctx *CombatContext, eventType string) {
	// Les observers sont gérés via le CombatSubject dans combatfacade
	// Pour l'instant, on logue simplement
	fmt.Printf("[State] Notification: %s\n", eventType)
}

// ActionRejectedState représente une action rejetée par la validation
type ActionRejectedState struct {
	BaseState
}

// NewActionRejectedState crée un nouvel état ActionRejected
func NewActionRejectedState() *ActionRejectedState {
	return &ActionRejectedState{
		BaseState: BaseState{
			name: "ActionRejected",
			allowedTransitions: map[string]bool{
				"ActionSelection": true, // Retour à la sélection d'action
			},
		},
	}
}

// Enter affiche l'erreur de validation
func (s *ActionRejectedState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	if ctx.ValidationError != nil {
		fmt.Printf("[State] Action rejetée: %v\n", ctx.ValidationError)
	}

	// Notifier les observateurs de l'échec
	s.notifyObservers(ctx, "ActionRejected")

	return nil
}

// Exit est appelé lors de la sortie
func (s *ActionRejectedState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	// Nettoyer l'erreur
	ctx.ValidationError = nil
	ctx.PendingCommand = nil
	return nil
}

// Handle gère les événements dans ActionRejected
func (s *ActionRejectedState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventRetryAction:
		// Retourner à la sélection d'action
		// Récupérer l'acteur depuis Combat ou créer un nouvel état générique
		return NewActionSelectionState(nil), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}

// notifyObservers notifie tous les observateurs
func (s *ActionRejectedState) notifyObservers(ctx *CombatContext, eventType string) {
	// Les observers sont gérés via le CombatSubject dans combatfacade
	// Pour l'instant, on logue simplement
	fmt.Printf("[State] Notification: %s\n", eventType)
}

// StunnedState représente une unité qui ne peut pas agir
type StunnedState struct {
	BaseState
}

// NewStunnedState crée un nouvel état Stunned
func NewStunnedState() *StunnedState {
	return &StunnedState{
		BaseState: BaseState{
			name: "Stunned",
			allowedTransitions: map[string]bool{
				"TurnEnd": true,
			},
		},
	}
}

// Enter indique que l'unité passe son tour
func (s *StunnedState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())
	fmt.Printf("[State] L'unité ne peut pas agir (Stun/Sleep)\n")
	return nil
}

// Exit est appelé lors de la sortie
func (s *StunnedState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans Stunned
func (s *StunnedState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventSkipTurn:
		// Passer directement à la fin du tour
		return NewTurnEndState(), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}
