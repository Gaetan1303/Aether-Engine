package states

import (
	"fmt"

	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
)

// CombatState représente un état dans la machine d'états du combat
// State Pattern - Interface pour tous les états concrets
// Chaque état implémente son propre comportement pour les transitions
type CombatState interface {
	// Enter est appelé lors de l'entrée dans cet état
	Enter(ctx *CombatContext) error

	// Exit est appelé lors de la sortie de cet état
	Exit(ctx *CombatContext) error

	// Handle gère les événements dans cet état et retourne le prochain état
	Handle(ctx *CombatContext, event StateEvent) (CombatState, error)

	// Name retourne le nom de l'état pour le logging et debugging
	Name() string

	// CanTransitionTo vérifie si une transition vers un état est valide
	CanTransitionTo(targetState string) bool
}

// StateEvent représente un événement qui déclenche une transition
type StateEvent struct {
	Type      EventType
	Data      interface{}
	Timestamp int64
}

// EventType énumère les types d'événements possibles
type EventType string

const (
	// Événements d'initialisation
	EventStartBattle     EventType = "START_BATTLE"
	EventSetupComplete   EventType = "SETUP_COMPLETE"
	EventValidationError EventType = "VALIDATION_ERROR"
	EventFirstUnitReady  EventType = "FIRST_UNIT_READY"

	// Événements de tour
	EventUnitCannotAct  EventType = "UNIT_CANNOT_ACT"
	EventUnitCanAct     EventType = "UNIT_CAN_ACT"
	EventActionReceived EventType = "ACTION_RECEIVED"
	EventTimeout        EventType = "TIMEOUT"
	EventWait           EventType = "WAIT"

	// Événements de validation
	EventCommandSelected  EventType = "COMMAND_SELECTED"
	EventValidationFailed EventType = "VALIDATION_FAILED"
	EventValidationOK     EventType = "VALIDATION_OK"
	EventExecuteCommand   EventType = "EXECUTE_COMMAND"

	// Événements d'exécution
	EventBeginExecution       EventType = "BEGIN_EXECUTION"
	EventActionCompleted      EventType = "ACTION_COMPLETED"
	EventErrorDuringExecution EventType = "ERROR_DURING_EXECUTION"
	EventEffectsApplied       EventType = "EFFECTS_APPLIED"
	EventRetryAction          EventType = "RETRY_ACTION"

	// Événements de victoire
	EventVictoryOrDefeat EventType = "VICTORY_OR_DEFEAT"
	EventBattleContinues EventType = "BATTLE_CONTINUES"

	// Événements de fin de tour
	EventNextUnitInQueue EventType = "NEXT_UNIT_IN_QUEUE"
	EventNoUnitReady     EventType = "NO_UNIT_READY"
	EventUnitReady       EventType = "UNIT_READY"
	EventSkipTurn        EventType = "SKIP_TURN"
	EventTurnComplete    EventType = "TURN_COMPLETE"
	EventNextUnitReady   EventType = "NEXT_UNIT_READY"

	// Événements de victoire/défaite
	EventCombatContinue EventType = "COMBAT_CONTINUE"

	// Événements d'exécution additionnels
	EventExecutionSuccess EventType = "EXECUTION_SUCCESS"
	EventExecutionError   EventType = "EXECUTION_ERROR"

	// Événements de validation
	EventValidationSuccess EventType = "VALIDATION_SUCCESS"

	// Événements de finalisation
	EventSaveResults    EventType = "SAVE_RESULTS"
	EventBattleClosed   EventType = "BATTLE_CLOSED"
	EventFinalizeCombat EventType = "FINALIZE_COMBAT"

	// Événements d'erreur
	EventCriticalError EventType = "CRITICAL_ERROR"
	EventErrorHandled  EventType = "ERROR_HANDLED"
)

// CombatContext contient toutes les données nécessaires pour les états
// Context du State Pattern - Maintient la référence à l'état actuel
type CombatContext struct {
	// Combat aggregate
	Combat *domain.Combat

	// État actuel
	CurrentState CombatState

	// Historique des états (pour debugging/rollback)
	StateHistory []StateTransition

	// Données temporaires pour l'état actuel
	PendingCommand  interface{} // *commands.Command
	PendingAction   *domain.ActionCombat
	PendingResult   interface{} // *commands.CommandResult
	ValidationError error

	// ATB System
	ATBSystem *ATBSystem

	// Observers (sera complété avec le Observer Pattern)
	Observers []interface{}
}

// StateTransition représente une transition d'état pour l'historique
type StateTransition struct {
	FromState string
	ToState   string
	Event     EventType
	Timestamp int64
}

// CombatStateMachine gère les transitions d'états
// Context du State Pattern - Délègue le comportement à l'état actuel
type CombatStateMachine struct {
	context *CombatContext
}

// NewCombatStateMachine crée une nouvelle machine d'états
func NewCombatStateMachine(combat *domain.Combat) *CombatStateMachine {
	ctx := &CombatContext{
		Combat:       combat,
		StateHistory: make([]StateTransition, 0),
		ATBSystem:    NewATBSystem(),
	}

	sm := &CombatStateMachine{
		context: ctx,
	}

	// État initial : Idle
	sm.TransitionTo(&IdleState{})

	return sm
}

// Context retourne le contexte de la machine d'états
func (sm *CombatStateMachine) Context() *CombatContext {
	return sm.context
}

// CurrentState retourne l'état actuel
func (sm *CombatStateMachine) CurrentState() CombatState {
	return sm.context.CurrentState
}

// TransitionTo effectue une transition vers un nouvel état
// Respect du Open/Closed Principle - pas besoin de modifier pour ajouter des états
func (sm *CombatStateMachine) TransitionTo(newState CombatState) error {
	if sm.context.CurrentState != nil {
		// Vérifier si la transition est autorisée
		if !sm.context.CurrentState.CanTransitionTo(newState.Name()) {
			return fmt.Errorf("transition invalide de %s vers %s",
				sm.context.CurrentState.Name(), newState.Name())
		}

		// Appeler Exit sur l'état actuel
		if err := sm.context.CurrentState.Exit(sm.context); err != nil {
			return fmt.Errorf("erreur lors de Exit de %s: %w",
				sm.context.CurrentState.Name(), err)
		}
	}

	// Sauvegarder la transition dans l'historique
	if sm.context.CurrentState != nil {
		transition := StateTransition{
			FromState: sm.context.CurrentState.Name(),
			ToState:   newState.Name(),
			Timestamp: sm.context.Combat.GetTimestamp(),
		}
		sm.context.StateHistory = append(sm.context.StateHistory, transition)
	}

	// Changer d'état
	oldState := sm.context.CurrentState
	sm.context.CurrentState = newState

	// Appeler Enter sur le nouvel état
	if err := newState.Enter(sm.context); err != nil {
		// Rollback en cas d'erreur
		sm.context.CurrentState = oldState
		return fmt.Errorf("erreur lors de Enter de %s: %w", newState.Name(), err)
	}

	return nil
}

// HandleEvent traite un événement et effectue la transition appropriée
func (sm *CombatStateMachine) HandleEvent(event StateEvent) error {
	if sm.context.CurrentState == nil {
		return fmt.Errorf("aucun état actuel défini")
	}

	// Laisser l'état actuel gérer l'événement
	nextState, err := sm.context.CurrentState.Handle(sm.context, event)
	if err != nil {
		return fmt.Errorf("erreur lors du traitement de l'événement %s dans %s: %w",
			event.Type, sm.context.CurrentState.Name(), err)
	}

	// Si un nouvel état est retourné, effectuer la transition
	if nextState != nil && nextState != sm.context.CurrentState {
		if err := sm.TransitionTo(nextState); err != nil {
			return err
		}
	}

	return nil
}

// GetStateHistory retourne l'historique des transitions
func (sm *CombatStateMachine) GetStateHistory() []StateTransition {
	return sm.context.StateHistory
}

// BaseState fournit une implémentation par défaut pour les méthodes communes
// Template Method Pattern - méthodes par défaut que les états peuvent override
type BaseState struct {
	name               string
	allowedTransitions map[string]bool
}

// Name retourne le nom de l'état
func (b *BaseState) Name() string {
	return b.name
}

// CanTransitionTo vérifie si une transition est autorisée
func (b *BaseState) CanTransitionTo(targetState string) bool {
	if b.allowedTransitions == nil {
		return true // Si pas de restrictions, tout est autorisé
	}
	return b.allowedTransitions[targetState]
}

// Enter implémentation par défaut (ne fait rien)
func (b *BaseState) Enter(ctx *CombatContext) error {
	return nil
}

// Exit implémentation par défaut (ne fait rien)
func (b *BaseState) Exit(ctx *CombatContext) error {
	return nil
}
