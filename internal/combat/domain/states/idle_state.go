package states

import "fmt"

// IdleState représente l'état initial où aucun combat n'est actif
// État initial de la machine d'états
type IdleState struct {
	BaseState
}

func init() {
	// Pas d'initialisation nécessaire pour l'instant
}

// NewIdleState crée un nouvel état Idle
func NewIdleState() *IdleState {
	return &IdleState{
		BaseState: BaseState{
			name: "Idle",
			allowedTransitions: map[string]bool{
				"Initializing": true,
			},
		},
	}
}

// Enter est appelé lors de l'entrée dans l'état Idle
func (s *IdleState) Enter(ctx *CombatContext) error {
	// Log de l'entrée dans l'état
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	// Aucune initialisation nécessaire en Idle
	return nil
}

// Exit est appelé lors de la sortie de l'état Idle
func (s *IdleState) Exit(ctx *CombatContext) error {
	// Log de la sortie
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans l'état Idle
func (s *IdleState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventStartBattle:
		// Transition vers Initializing
		return NewInitializingState(), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}
