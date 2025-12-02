package combatfacade

import (
	"errors"
	"fmt"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/aether-engine/aether-engine/internal/combat/domain/commands"
	"github.com/aether-engine/aether-engine/internal/combat/domain/observers"
	"github.com/aether-engine/aether-engine/internal/combat/domain/states"
)

// Step C - Fonctions de Facade pour éviter les cycles d'imports
// Ces fonctions acceptent *domain.Combat et utilisent les types concrets des patterns Step C

// SetupStateMachine initialise la state machine (appelé depuis combatinitializer)
func SetupStateMachine(c *domain.Combat, sm *states.CombatStateMachine) error {
	if sm == nil {
		return errors.New("state machine cannot be nil")
	}
	if c.GetStateMachineRaw() != nil {
		return errors.New("state machine déjà initialisée")
	}
	c.SetStateMachine(sm)
	return nil
}

// SetupCommandSystem initialise le système de commandes (appelé depuis combatinitializer)
func SetupCommandSystem(c *domain.Combat, invoker *commands.CommandInvoker, factory *commands.CommandFactory) {
	c.SetCommandInvoker(invoker)
	c.SetCommandFactory(factory)
}

// SetupObservers initialise le système d'observateurs (appelé depuis combatinitializer)
func SetupObservers(c *domain.Combat, subject *observers.CombatSubject) {
	c.SetObserverSubject(subject)
}


// GetStateMachine retourne la state machine (type concret pour les tests)
func GetStateMachine(c *domain.Combat) *states.CombatStateMachine {
	if sm, ok := c.GetStateMachineRaw().(*states.CombatStateMachine); ok {
		return sm
	}
	return nil
}

// GetCommandInvoker retourne le command invoker
func GetCommandInvoker(c *domain.Combat) *commands.CommandInvoker {
	if inv, ok := c.GetCommandInvokerRaw().(*commands.CommandInvoker); ok {
		return inv
	}
	return nil
}

// GetCommandFactory retourne la command factory
func GetCommandFactory(c *domain.Combat) *commands.CommandFactory {
	if factory, ok := c.GetCommandFactoryRaw().(*commands.CommandFactory); ok {
		return factory
	}
	return nil
}

// GetObserverSubject retourne l'observer subject
func GetObserverSubject(c *domain.Combat) *observers.CombatSubject {
	if subject, ok := c.GetObserverSubjectRaw().(*observers.CombatSubject); ok {
		return subject
	}
	return nil
}

// InitializeCombatWithStateMachine initialise le combat avec la state machine
// Configure les listeners pour les événements de transition d'état
func InitializeCombatWithStateMachine(c *domain.Combat) error {
	sm := GetStateMachine(c)
	if sm == nil {
		return errors.New("state machine non initialisée")
	}

	// La state machine est déjà initialisée lors de sa création
	// avec NewCombatStateMachine qui crée le contexte et l'état initial

	return nil
}

// CommandType représente le type d'action à exécuter
type CommandType string

const (
	CommandTypeMove   CommandType = "move"
	CommandTypeAttack CommandType = "attack"
	CommandTypeSkill  CommandType = "skill"
	CommandTypeItem   CommandType = "item"
	CommandTypeFlee   CommandType = "flee"
	CommandTypeWait   CommandType = "wait"
)


// ExecutePlayerAction exécute une action de joueur via le système de commandes
// SOLID Principles appliqués:
// - SRP: Fonction focalisée sur la création et l'exécution de commandes
// - OCP: Extensible via ajout de nouveaux CommandType sans modifier le code existant
// - LSP: Utilise l'interface Command, toutes les implémentations sont substituables
// - ISP: Interface Command minimaliste (Execute, Validate, Rollback, GetActor)
// - DIP: Dépend des abstractions (Command interface) pas des implémentations concrètes
func ExecutePlayerAction(c *domain.Combat, actorID domain.UnitID, actionType CommandType, params map[string]interface{}) (*commands.CommandResult, error) {
factory := GetCommandFactory(c)
if factory == nil {
return nil, errors.New("command factory non initialisée")
}

actor := c.TrouverUnite(actorID)
if actor == nil {
return nil, fmt.Errorf("acteur %s introuvable", actorID)
}

// Factory Pattern: création polymorphique de commandes
var cmd commands.Command
var err error

switch actionType {
case CommandTypeMove:
targetX, okX := params["targetX"].(int)
targetY, okY := params["targetY"].(int)
if !okX || !okY {
return nil, errors.New("coordonnées cibles invalides pour Move")
}
cmd, err = factory.CreateMoveCommand(actor, targetX, targetY)

case CommandTypeAttack:
targetID, ok := params["targetID"].(domain.UnitID)
if !ok {
return nil, errors.New("ID cible invalide pour Attack")
}
cmd, err = factory.CreateAttackCommand(actor, targetID)

case CommandTypeSkill:
skillID, ok := params["skillID"].(string)
if !ok {
return nil, errors.New("ID compétence invalide pour Skill")
}
targetIDs, ok := params["targetIDs"].([]domain.UnitID)
if !ok {
if targetID, ok := params["targetID"].(domain.UnitID); ok {
targetIDs = []domain.UnitID{targetID}
} else {
return nil, errors.New("IDs cibles invalides pour Skill")
}
}
cmd, err = factory.CreateSkillCommand(actor, skillID, targetIDs)

case CommandTypeItem:
itemID, ok := params["itemID"].(string)
if !ok {
return nil, errors.New("ID item invalide pour Item")
}
targetID, ok := params["targetID"].(domain.UnitID)
if !ok {
return nil, errors.New("ID cible invalide pour Item")
}
cmd, err = factory.CreateItemCommand(actor, itemID, targetID)

case CommandTypeFlee:
cmd, err = factory.CreateFleeCommand(actor)

case CommandTypeWait:
cmd, err = factory.CreateWaitCommand(actor)

default:
return nil, fmt.Errorf("type d'action inconnu: %s", actionType)
}

if err != nil {
return nil, fmt.Errorf("échec de création de commande: %w", err)
}

// Command Pattern: exécution polymorphique
result, err := cmd.Execute()
if err != nil {
return nil, fmt.Errorf("échec d'exécution: %w", err)
}

return result, nil
}

// AttachObserver attache un observateur au système
func AttachObserver(c *domain.Combat, observer observers.CombatObserver) {
	subject := GetObserverSubject(c)
	if subject != nil {
		subject.Attach(observer)
	}
}

// DetachObserver détache un observateur du système
func DetachObserver(c *domain.Combat, observerName string) {
	subject := GetObserverSubject(c)
	if subject != nil {
		subject.Detach(observerName)
	}
}

// NotifyObservers notifie tous les observateurs d'un événement
func NotifyObservers(c *domain.Combat, eventType string, ctx *states.CombatContext) {
	subject := GetObserverSubject(c)
	if subject != nil {
		// La méthode Notify de CombatSubject prend une string, pas EventType
		// On utilise le pattern actuel du Subject
		for _, obs := range subject.GetObservers() {
			obs.OnNotify(eventType, ctx)
		}
	}
}

// GetCurrentStateString retourne l'état actuel sous forme de string
func GetCurrentStateString(c *domain.Combat) string {
	sm := GetStateMachine(c)
	if sm == nil || sm.Context() == nil || sm.Context().CurrentState == nil {
		return "UNKNOWN"
	}
	return sm.Context().CurrentState.Name()
}

// GetStateHistoryList retourne l'historique des transitions d'état
func GetStateHistoryList(c *domain.Combat) []states.StateTransition {
	sm := GetStateMachine(c)
	if sm == nil {
		return []states.StateTransition{}
	}
	return sm.GetStateHistory()
}

// GetCommandHistoryList retourne l'historique des commandes exécutées
func GetCommandHistoryList(c *domain.Combat) []commands.Command {
	invoker := GetCommandInvoker(c)
	if invoker == nil {
		return []commands.Command{}
	}
	return invoker.GetHistory()
}

// UndoLastCommand annule la dernière commande exécutée
func UndoLastCommand(c *domain.Combat) error {
	invoker := GetCommandInvoker(c)
	if invoker == nil {
		return errors.New("command invoker non initialisé")
	}
	return errors.New("undo non implémenté")
}

// RedoLastCommand refait la dernière commande annulée
func RedoLastCommand(c *domain.Combat) error {
	invoker := GetCommandInvoker(c)
	if invoker == nil {
		return errors.New("command invoker non initialisé")
	}
	return errors.New("redo non implémenté")
}

// CanUndo vérifie si une commande peut être annulée
func CanUndo(c *domain.Combat) bool {
	invoker := GetCommandInvoker(c)
	if invoker == nil {
		return false
	}
	return len(invoker.GetHistory()) > 0
}

// CanRedo vérifie si une commande peut être refaite
func CanRedo(c *domain.Combat) bool {
	invoker := GetCommandInvoker(c)
	if invoker == nil {
		return false
	}
	return false
}
