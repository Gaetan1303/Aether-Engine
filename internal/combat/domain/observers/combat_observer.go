package observers

import (
	"fmt"

	"github.com/aether-engine/aether-engine/internal/combat/domain/commands"
	"github.com/aether-engine/aether-engine/internal/combat/domain/states"
)

// CombatObserver représente un observateur du combat
// Observer Pattern - Interface pour surveiller les événements
type CombatObserver interface {
	// OnNotify est appelé lors d'un événement
	OnNotify(eventType string, context *states.CombatContext)

	// GetName retourne le nom de l'observateur
	GetName() string
}

// CombatSubject gère l'enregistrement et la notification des observateurs
// Subject du Observer Pattern
type CombatSubject struct {
	observers []CombatObserver
}

// NewCombatSubject crée un nouveau sujet
func NewCombatSubject() *CombatSubject {
	return &CombatSubject{
		observers: make([]CombatObserver, 0),
	}
}

// Attach ajoute un observateur
func (s *CombatSubject) Attach(observer CombatObserver) {
	s.observers = append(s.observers, observer)
	fmt.Printf("[Observer] Observateur ajouté: %s\n", observer.GetName())
}

// Detach retire un observateur
func (s *CombatSubject) Detach(observerName string) {
	for i, obs := range s.observers {
		if obs.GetName() == observerName {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			fmt.Printf("[Observer] Observateur retiré: %s\n", observerName)
			return
		}
	}
}

// NotifyAll notifie tous les observateurs
func (s *CombatSubject) NotifyAll(eventType string, context *states.CombatContext) {
	for _, observer := range s.observers {
		observer.OnNotify(eventType, context)
	}
}

// GetObservers retourne la liste des observateurs
func (s *CombatSubject) GetObservers() []CombatObserver {
	return s.observers
}

// StateObserver surveille les transitions d'états
type StateObserver struct {
	name string
}

// NewStateObserver crée un nouvel observateur d'états
func NewStateObserver(name string) *StateObserver {
	return &StateObserver{
		name: name,
	}
}

// OnNotify réagit aux événements d'états
func (o *StateObserver) OnNotify(eventType string, context *states.CombatContext) {
	switch eventType {
	case "StateTransition":
		fmt.Printf("[StateObserver] Transition vers: %s\n",
			context.CurrentState.Name())

	case "ActionConfirmed":
		if cmd, ok := context.PendingCommand.(commands.Command); ok {
			fmt.Printf("[StateObserver] Action confirmée: %s\n", cmd.GetType())
		}

	case "ActionRejected":
		fmt.Printf("[StateObserver] Action rejetée: %v\n", context.ValidationError)

	case "TurnEnd":
		fmt.Printf("[StateObserver] Fin du tour\n")

	case "BattleEnded":
		fmt.Printf("[StateObserver] Combat terminé\n")
	}
}

// GetName retourne le nom de l'observateur
func (o *StateObserver) GetName() string {
	return o.name
}

// UnitObserver surveille l'état des unités (HP, MP, statuts)
type UnitObserver struct {
	name string
}

// NewUnitObserver crée un nouvel observateur d'unités
func NewUnitObserver(name string) *UnitObserver {
	return &UnitObserver{
		name: name,
	}
}

// OnNotify réagit aux événements d'unités
func (o *UnitObserver) OnNotify(eventType string, context *states.CombatContext) {
	switch eventType {
	case "Effect_DAMAGE":
		// Un dégât a été infligé
		if result, ok := context.PendingResult.(*commands.CommandResult); ok {
			fmt.Printf("[UnitObserver] Dégâts infligés: %d\n", result.DamageDealt)
		}

	case "Effect_HEALING":
		// Des soins ont été effectués
		if result, ok := context.PendingResult.(*commands.CommandResult); ok {
			fmt.Printf("[UnitObserver] Soins effectués: %d\n", result.HealingDone)
		}

	case "Effect_STATUS":
		// Un statut a été appliqué
		if result, ok := context.PendingResult.(*commands.CommandResult); ok {
			for _, status := range result.StatusApplied {
				fmt.Printf("[UnitObserver] Statut appliqué: %d\n", status.Type())
			}
		}

	case "UnitDefeated":
		// Une unité a été vaincue
		fmt.Printf("[UnitObserver] Unité vaincue\n")

	case "TurnBegin":
		// Début du tour d'une unité
		if cmd, ok := context.PendingCommand.(commands.Command); ok {
			fmt.Printf("[UnitObserver] Tour de: %s\n", cmd.GetActor().Nom())
		}
	}
}

// GetName retourne le nom de l'observateur
func (o *UnitObserver) GetName() string {
	return o.name
}

// ConnectionObserver surveille les connexions joueur/serveur
type ConnectionObserver struct {
	name                string
	disconnectedPlayers map[string]bool
}

// NewConnectionObserver crée un nouvel observateur de connexions
func NewConnectionObserver(name string) *ConnectionObserver {
	return &ConnectionObserver{
		name:                name,
		disconnectedPlayers: make(map[string]bool),
	}
}

// OnNotify réagit aux événements de connexion
func (o *ConnectionObserver) OnNotify(eventType string, context *states.CombatContext) {
	switch eventType {
	case "PlayerDisconnected":
		// Un joueur s'est déconnecté
		// TODO: Implémenter le système de connexion
		fmt.Printf("[ConnectionObserver] Joueur déconnecté, activation de l'IA\n")

	case "PlayerReconnected":
		// Un joueur s'est reconnecté
		fmt.Printf("[ConnectionObserver] Joueur reconnecté, désactivation de l'IA\n")

	case "ServerTimeout":
		// Timeout du serveur
		fmt.Printf("[ConnectionObserver] Timeout du serveur\n")
	}
}

// GetName retourne le nom de l'observateur
func (o *ConnectionObserver) GetName() string {
	return o.name
}

// EventLogger enregistre tous les événements pour replay
type EventLogger struct {
	name string
	log  []LogEntry
}

// LogEntry représente une entrée dans le log
type LogEntry struct {
	Timestamp   int64
	EventType   string
	StateName   string
	CommandType string
	Details     string
}

// NewEventLogger crée un nouvel enregistreur d'événements
func NewEventLogger(name string) *EventLogger {
	return &EventLogger{
		name: name,
		log:  make([]LogEntry, 0),
	}
}

// OnNotify enregistre tous les événements
func (o *EventLogger) OnNotify(eventType string, context *states.CombatContext) {
	entry := LogEntry{
		Timestamp: context.Combat.GetTimestamp(),
		EventType: eventType,
		StateName: context.CurrentState.Name(),
	}

	if cmd, ok := context.PendingCommand.(commands.Command); ok {
		entry.CommandType = string(cmd.GetType())
	}

	o.log = append(o.log, entry)

	// Pour debug, afficher l'entrée
	fmt.Printf("[EventLogger] %s | État: %s | Commande: %s\n",
		eventType, entry.StateName, entry.CommandType)
}

// GetName retourne le nom de l'observateur
func (o *EventLogger) GetName() string {
	return o.name
}

// GetLog retourne le log complet
func (o *EventLogger) GetLog() []LogEntry {
	return o.log
}

// ExportLog exporte le log pour replay
func (o *EventLogger) ExportLog() string {
	// TODO: Implémenter l'export JSON
	return fmt.Sprintf("Log avec %d entrées", len(o.log))
}
