package states

import (
	"fmt"

	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
)

// TurnBeginState gère le début du tour d'une unité
// - Get CurrentUnit from Queue
// - Trigger OnTurnStart hooks
// - Apply Status effects (Poison, Regen)
// - Check if Unit can act (Stun, Sleep)
type TurnBeginState struct {
	BaseState
	currentUnit *domain.Unite
}

// NewTurnBeginState crée un nouvel état TurnBegin
func NewTurnBeginState() *TurnBeginState {
	return &TurnBeginState{
		BaseState: BaseState{
			name: "TurnBegin",
			allowedTransitions: map[string]bool{
				"Stunned":         true,
				"ActionSelection": true,
			},
		},
	}
}

// Enter initialise le début du tour
func (s *TurnBeginState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	// 1. Obtenir l'unité avec le plus haut ATB
	readyUnits := ctx.ATBSystem.GetReadyUnits()
	if len(readyUnits) == 0 {
		return fmt.Errorf("aucune unité prête pour agir")
	}

	// Prendre la première unité prête (TODO: ordre d'initiative)
	unitID := readyUnits[0]
	s.currentUnit = ctx.Combat.TrouverUnite(unitID)
	if s.currentUnit == nil {
		return fmt.Errorf("unité %s non trouvée", unitID)
	}

	fmt.Printf("[State] Tour de l'unité: %s\n", s.currentUnit.Nom())

	// 2. Déclencher OnTurnStart hooks
	s.currentUnit.NouveauTour()

	// 3. Appliquer les effets de statut (Poison, Regen, etc.)
	effets := s.currentUnit.TraiterStatuts()
	for _, effet := range effets {
		fmt.Printf("[State] Effet de statut appliqué: %+v\n", effet)
	}

	// 4. Réinitialiser la jauge ATB de cette unité
	ctx.ATBSystem.ResetGauge(unitID)

	return nil
}

// Exit est appelé lors de la sortie
func (s *TurnBeginState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans TurnBegin
func (s *TurnBeginState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventUnitCannotAct:
		// Unité stunned/sleep, passer le tour
		return NewStunnedState(), nil

	case EventUnitCanAct:
		// Unité peut agir, aller en sélection d'action
		return NewActionSelectionState(s.currentUnit), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}

// CurrentUnit retourne l'unité active
func (s *TurnBeginState) CurrentUnit() *domain.Unite {
	return s.currentUnit
}
