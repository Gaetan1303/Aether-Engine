package commands

import (
	"fmt"

	domain "github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// Command représente une action que peut effectuer un joueur
// Command Pattern - Interface pour encapsuler les actions
// Chaque commande implémente Validate, Execute, Rollback
type Command interface {
	// Validate vérifie si la commande peut être exécutée
	Validate() error

	// Execute exécute la commande
	Execute() (*CommandResult, error)

	// Rollback annule les effets de la commande en cas d'erreur
	Rollback() error

	// GetType retourne le type de commande
	GetType() CommandType

	// GetActor retourne l'acteur qui exécute la commande
	GetActor() *domain.Unite
}

// CommandType énumère les types de commandes
type CommandType string

const (
	CommandTypeMove   CommandType = "MOVE"
	CommandTypeAttack CommandType = "ATTACK"
	CommandTypeSkill  CommandType = "SKILL"
	CommandTypeItem   CommandType = "ITEM"
	CommandTypeFlee   CommandType = "FLEE"
	CommandTypeWait   CommandType = "WAIT"
)

// CommandResult représente le résultat de l'exécution d'une commande
type CommandResult struct {
	Success       bool
	Message       string
	Effects       []CommandEffect
	CostMP        int
	CostStamina   int
	CostMovement  int
	DamageDealt   int
	HealingDone   int
	StatusApplied []*shared.Statut
}

// CommandEffect représente un effet produit par une commande
type CommandEffect struct {
	Type     EffectType
	TargetID domain.UnitID
	Value    int
	Position *shared.Position
	Status   *shared.Statut
}

// EffectType énumère les types d'effets
type EffectType string

const (
	EffectTypeDamage     EffectType = "DAMAGE"
	EffectTypeHealing    EffectType = "HEALING"
	EffectTypeStatus     EffectType = "STATUS"
	EffectTypeMovement   EffectType = "MOVEMENT"
	EffectTypeStatChange EffectType = "STAT_CHANGE"
)

// BaseCommand fournit une implémentation de base pour les commandes
// Template Method Pattern - implémentation par défaut
type BaseCommand struct {
	actor       *domain.Unite
	combat      *domain.Combat
	commandType CommandType

	// Snapshot pour rollback
	snapshot *CommandSnapshot
}

// CommandSnapshot sauvegarde l'état avant exécution pour rollback
type CommandSnapshot struct {
	ActorHP       int
	ActorMP       int
	ActorStamina  int
	ActorPosition *shared.Position
	TargetStates  map[domain.UnitID]*UnitSnapshot
}

// UnitSnapshot sauvegarde l'état d'une unité
type UnitSnapshot struct {
	HP       int
	MP       int
	Stamina  int
	Position *shared.Position
	Statuses []*shared.Statut
}

// NewBaseCommand crée une nouvelle commande de base
func NewBaseCommand(actor *domain.Unite, combat *domain.Combat, commandType CommandType) *BaseCommand {
	return &BaseCommand{
		actor:       actor,
		combat:      combat,
		commandType: commandType,
	}
}

// GetType retourne le type de commande
func (c *BaseCommand) GetType() CommandType {
	return c.commandType
}

// GetActor retourne l'acteur
func (c *BaseCommand) GetActor() *domain.Unite {
	return c.actor
}

// CreateSnapshot crée un snapshot de l'état actuel
func (c *BaseCommand) CreateSnapshot() {
	c.snapshot = &CommandSnapshot{
		ActorHP:       c.actor.HPActuels(),
		ActorMP:       c.actor.StatsActuelles().MP,
		ActorStamina:  c.actor.StatsActuelles().Stamina,
		ActorPosition: c.actor.Position(),
		TargetStates:  make(map[domain.UnitID]*UnitSnapshot),
	}
}

// Rollback implémentation par défaut (ne fait rien)
func (c *BaseCommand) Rollback() error {
	// Les commandes concrètes peuvent override
	return nil
}

// CommandInvoker gère l'historique et l'exécution des commandes
// Invoker du Command Pattern - stocke et exécute les commandes
type CommandInvoker struct {
	history    []Command
	maxHistory int
}

// NewCommandInvoker crée un nouvel invoker
func NewCommandInvoker(maxHistory int) *CommandInvoker {
	return &CommandInvoker{
		history:    make([]Command, 0),
		maxHistory: maxHistory,
	}
}

// Execute valide et exécute une commande
func (inv *CommandInvoker) Execute(cmd Command) (*CommandResult, error) {
	// 1. Valider la commande
	if err := cmd.Validate(); err != nil {
		return nil, fmt.Errorf("validation échouée: %w", err)
	}

	// 2. Exécuter la commande
	result, err := cmd.Execute()
	if err != nil {
		// En cas d'erreur, tenter un rollback
		if rollbackErr := cmd.Rollback(); rollbackErr != nil {
			return nil, fmt.Errorf("erreur lors de l'exécution ET du rollback: %w, %v", err, rollbackErr)
		}
		return nil, fmt.Errorf("erreur lors de l'exécution: %w", err)
	}

	// 3. Ajouter à l'historique
	inv.addToHistory(cmd)

	return result, nil
}

// addToHistory ajoute une commande à l'historique
func (inv *CommandInvoker) addToHistory(cmd Command) {
	inv.history = append(inv.history, cmd)

	// Limiter la taille de l'historique
	if len(inv.history) > inv.maxHistory {
		inv.history = inv.history[1:]
	}
}

// GetHistory retourne l'historique des commandes
func (inv *CommandInvoker) GetHistory() []Command {
	return inv.history
}

// Clear vide l'historique
func (inv *CommandInvoker) Clear() {
	inv.history = make([]Command, 0)
}

// Implémentation de CommandSystemProvider interface

// CanUndo implémente CommandSystemProvider
func (inv *CommandInvoker) CanUndo() bool {
	return len(inv.history) > 0
}

// Undo implémente CommandSystemProvider
func (inv *CommandInvoker) Undo() error {
	if !inv.CanUndo() {
		return fmt.Errorf("aucune commande à annuler")
	}

	lastCmd := inv.history[len(inv.history)-1]
	if err := lastCmd.Rollback(); err != nil {
		return fmt.Errorf("échec du rollback: %w", err)
	}

	inv.history = inv.history[:len(inv.history)-1]
	return nil
}

// History implémente CommandSystemProvider
func (inv *CommandInvoker) History() []string {
	result := make([]string, len(inv.history))
	for i, cmd := range inv.history {
		result[i] = fmt.Sprintf("%s by %s", cmd.GetType(), cmd.GetActor().ID())
	}
	return result
}
