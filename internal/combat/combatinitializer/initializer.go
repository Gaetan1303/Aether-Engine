package combatinitializer

import (
	"github.com/aether-engine/aether-engine/internal/combat/combatfacade"
	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/aether-engine/aether-engine/internal/combat/domain/commands"
	"github.com/aether-engine/aether-engine/internal/combat/domain/observers"
	"github.com/aether-engine/aether-engine/internal/combat/domain/states"
)

// CombatInitializer permet d'initialiser un Combat avec les patterns Step C
// sans créer de cycle d'imports
type CombatInitializer struct {
	combat *domain.Combat
}

// NewCombatInitializer crée un nouvel initialiseur
func NewCombatInitializer(combat *domain.Combat) *CombatInitializer {
	return &CombatInitializer{
		combat: combat,
	}
}

// InitializeStateMachine initialise la state machine
func (ci *CombatInitializer) InitializeStateMachine() error {
	return combatfacade.SetupStateMachine(ci.combat, states.NewCombatStateMachine(ci.combat))
}

// InitializeCommandSystem initialise le système de commandes
func (ci *CombatInitializer) InitializeCommandSystem() {
	combatfacade.SetupCommandSystem(
		ci.combat,
		commands.NewCommandInvoker(100),
		commands.NewCommandFactory(ci.combat),
	)
}

// InitializeObservers initialise le système d'observateurs
func (ci *CombatInitializer) InitializeObservers() {
	combatfacade.SetupObservers(ci.combat, observers.NewCombatSubject())
}

// InitializeValidation initialise la chaîne de validation
func (ci *CombatInitializer) InitializeValidation() {
}

// InitializeAll initialise tous les systèmes Step C
func (ci *CombatInitializer) InitializeAll() error {
	ci.InitializeCommandSystem()
	ci.InitializeObservers()
	ci.InitializeValidation()
	return ci.InitializeStateMachine()
}

// Helper functions pour créer des observateurs prêts à l'emploi

// AddStateObserver ajoute un observateur d'états
func (ci *CombatInitializer) AddStateObserver(name string) {
	observer := observers.NewStateObserver(name)
	combatfacade.AttachObserver(ci.combat, observer)
}

// AddUnitObserver ajoute un observateur d'unités
func (ci *CombatInitializer) AddUnitObserver(name string) {
	observer := observers.NewUnitObserver(name)
	combatfacade.AttachObserver(ci.combat, observer)
}

// AddConnectionObserver ajoute un observateur de connexions
func (ci *CombatInitializer) AddConnectionObserver(name string) {
	observer := observers.NewConnectionObserver(name)
	combatfacade.AttachObserver(ci.combat, observer)
}

// AddEventLogger ajoute un logger d'événements
func (ci *CombatInitializer) AddEventLogger(name string) *observers.EventLogger {
	logger := observers.NewEventLogger(name)
	combatfacade.AttachObserver(ci.combat, logger)
	return logger
}
