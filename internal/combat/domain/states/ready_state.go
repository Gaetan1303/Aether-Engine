package states

import "fmt"

// ReadyState représente l'état où le combat est prêt à démarrer
// En attente du premier acteur avec ATB >= 100
type ReadyState struct {
	BaseState
}

// NewReadyState crée un nouvel état Ready
func NewReadyState() *ReadyState {
	return &ReadyState{
		BaseState: BaseState{
			name: "Ready",
			allowedTransitions: map[string]bool{
				"TurnBegin": true,
			},
		},
	}
}

// Enter initialise l'état Ready
func (s *ReadyState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	// Commencer à faire progresser les jauges ATB
	// jusqu'à ce qu'une unité soit prête
	return nil
}

// Exit est appelé lors de la sortie
func (s *ReadyState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans Ready
func (s *ReadyState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventFirstUnitReady:
		// Une unité est prête, commencer le combat
		return NewTurnBeginState(), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}

// FailedState représente l'échec de l'initialisation
type FailedState struct {
	BaseState
	error error
}

// NewFailedState crée un nouvel état Failed
func NewFailedState() *FailedState {
	return &FailedState{
		BaseState: BaseState{
			name:               "Failed",
			allowedTransitions: map[string]bool{
				// État terminal, aucune transition autorisée
			},
		},
	}
}

// Enter enregistre l'erreur
func (s *FailedState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s (ERREUR)\n", s.Name())
	s.error = ctx.ValidationError
	return nil
}

// Exit est appelé lors de la sortie (ne devrait jamais arriver)
func (s *FailedState) Exit(ctx *CombatContext) error {
	return nil
}

// Handle ne gère aucun événement (état terminal)
func (s *FailedState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	return nil, fmt.Errorf("état Failed est terminal, aucune transition possible")
}
