package states

import (
	"fmt"
)

// CheckVictoryState vérifie les conditions de victoire/défaite
type CheckVictoryState struct {
	BaseState
}

// NewCheckVictoryState crée un nouvel état CheckVictory
func NewCheckVictoryState() *CheckVictoryState {
	return &CheckVictoryState{
		BaseState: BaseState{
			name: "CheckVictory",
			allowedTransitions: map[string]bool{
				"TurnEnd":     true,
				"BattleEnded": true,
			},
		},
	}
}

// Enter vérifie les conditions de fin de combat
func (s *CheckVictoryState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	// Vérifier si une équipe a gagné ou perdu
	victoryCondition := ctx.Combat.VerifierConditionsVictoire()

	switch victoryCondition {
	case "CONTINUE":
		// Combat continue
		fmt.Printf("[State] Le combat continue\n")
		return nil

	case "VICTORY", "DEFEAT", "FLED":
		// Combat terminé
		fmt.Printf("[State] Combat terminé: %s\n", victoryCondition)
		return nil

	default:
		return fmt.Errorf("condition de victoire inconnue: %s", victoryCondition)
	}
}

// Exit est appelé lors de la sortie
func (s *CheckVictoryState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans CheckVictory
func (s *CheckVictoryState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventVictoryOrDefeat:
		// Une équipe a gagné/perdu, fin du combat
		return NewBattleEndedState(), nil

	case EventCombatContinue:
		// Combat continue, fin du tour
		return NewTurnEndState(), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}

// TurnEndState termine le tour de l'unité active
type TurnEndState struct {
	BaseState
}

// NewTurnEndState crée un nouvel état TurnEnd
func NewTurnEndState() *TurnEndState {
	return &TurnEndState{
		BaseState: BaseState{
			name: "TurnEnd",
			allowedTransitions: map[string]bool{
				"WaitingATB": true,
			},
		},
	}
}

// Enter termine le tour
func (s *TurnEndState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	// Nettoyer le contexte
	ctx.PendingCommand = nil
	ctx.PendingResult = nil
	ctx.ValidationError = nil

	// Notifier la fin du tour
	fmt.Printf("[State] Notification: TurnEnd\n")

	return nil
}

// Exit est appelé lors de la sortie
func (s *TurnEndState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans TurnEnd
func (s *TurnEndState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventTurnComplete:
		// Tour terminé, attendre la prochaine unité
		return NewWaitingATBState(), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}

// WaitingATBState attend qu'une unité soit prête (ATB >= 100)
type WaitingATBState struct {
	BaseState
}

// NewWaitingATBState crée un nouvel état WaitingATB
func NewWaitingATBState() *WaitingATBState {
	return &WaitingATBState{
		BaseState: BaseState{
			name: "WaitingATB",
			allowedTransitions: map[string]bool{
				"TurnBegin": true,
			},
		},
	}
}

// Enter commence à faire progresser les jauges ATB
func (s *WaitingATBState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	// Faire progresser les jauges jusqu'à ce qu'une unité soit prête
	for {
		ctx.ATBSystem.Tick()

		readyUnits := ctx.ATBSystem.GetReadyUnits()
		if len(readyUnits) > 0 {
			// Une unité est prête
			fmt.Printf("[State] Unité prête: %s\n", readyUnits[0])
			break
		}
	}

	return nil
}

// Exit est appelé lors de la sortie
func (s *WaitingATBState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans WaitingATB
func (s *WaitingATBState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventNextUnitReady:
		// Une unité est prête, commencer son tour
		return NewTurnBeginState(), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}

// BattleEndedState représente la fin du combat
type BattleEndedState struct {
	BaseState
}

// NewBattleEndedState crée un nouvel état BattleEnded
func NewBattleEndedState() *BattleEndedState {
	return &BattleEndedState{
		BaseState: BaseState{
			name: "BattleEnded",
			allowedTransitions: map[string]bool{
				"Finalizing": true,
			},
		},
	}
}

// Enter termine le combat
func (s *BattleEndedState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	// Afficher le résultat
	result := ctx.Combat.ObtenirResultat()
	fmt.Printf("[State] Résultat du combat: %s\n", result)

	// Notifier les observateurs
	fmt.Printf("[State] Notification: BattleEnded\n")

	return nil
}

// Exit est appelé lors de la sortie
func (s *BattleEndedState) Exit(ctx *CombatContext) error {
	fmt.Printf("[State] Sortie de l'état: %s\n", s.Name())
	return nil
}

// Handle gère les événements dans BattleEnded
func (s *BattleEndedState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	switch event.Type {
	case EventFinalizeCombat:
		// Finaliser le combat (récompenses, XP, etc.)
		return NewFinalizingState(), nil

	default:
		return nil, fmt.Errorf("événement %s non géré dans l'état %s", event.Type, s.Name())
	}
}

// FinalizingState finalise le combat (récompenses, XP)
type FinalizingState struct {
	BaseState
}

// NewFinalizingState crée un nouvel état Finalizing
func NewFinalizingState() *FinalizingState {
	return &FinalizingState{
		BaseState: BaseState{
			name:               "Finalizing",
			allowedTransitions: map[string]bool{
				// État terminal
			},
		},
	}
}

// Enter finalise le combat
func (s *FinalizingState) Enter(ctx *CombatContext) error {
	fmt.Printf("[State] Entrée dans état: %s\n", s.Name())

	// Distribuer XP, loots, etc.
	ctx.Combat.DistribuerRecompenses()

	// Notifier les observateurs
	fmt.Printf("[State] Notification: Finalizing\n")

	fmt.Printf("[State] Combat finalisé\n")
	return nil
}

// Exit est appelé lors de la sortie (ne devrait jamais arriver)
func (s *FinalizingState) Exit(ctx *CombatContext) error {
	return nil
}

// Handle ne gère aucun événement (état terminal)
func (s *FinalizingState) Handle(ctx *CombatContext, event StateEvent) (CombatState, error) {
	return nil, fmt.Errorf("état Finalizing est terminal, aucune transition possible")
}
