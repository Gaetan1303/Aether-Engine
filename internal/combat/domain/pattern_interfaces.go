package domain

// pattern_interfaces.go
// Interfaces typées pour éviter interface{} dans l'agrégat Combat
// Respecte Interface Segregation Principle (SOLID)

// StateMachineProvider fournit l'accès à la machine à états
// Interface minimale pour éviter les cycles d'imports
type StateMachineProvider interface {
	// GetCurrentState retourne le nom de l'état courant
	GetCurrentState() string

	// CanTransitionTo vérifie si une transition est possible
	CanTransitionTo(stateName string) bool
}

// CommandSystemProvider fournit l'accès au système de commandes
type CommandSystemProvider interface {
	// CanUndo vérifie si un undo est possible
	CanUndo() bool

	// Undo annule la dernière commande
	Undo() error

	// History retourne l'historique des commandes
	History() []string
} // CommandFactoryProvider crée des commandes
// Interface marker - les méthodes concrètes sont définies dans commands.CommandFactory
type CommandFactoryProvider interface {
	// Marker interface - implémentée par CommandFactory
}

// ObserverProvider fournit l'accès au système d'observateurs
type ObserverProvider interface {
	// Notify notifie tous les observateurs d'un événement
	Notify(eventType string, data interface{})

	// AttachObserver ajoute un observateur
	AttachObserver(observer interface{})

	// DetachObserver retire un observateur
	DetachObserver(observerName string)

	// ObserverCount retourne le nombre d'observateurs
	ObserverCount() int
}

// ValidationProvider fournit l'accès à la chaîne de validation
type ValidationProvider interface {
	// Validate valide une action via la chaîne de responsabilité
	Validate(action interface{}) error

	// AddValidator ajoute un validateur à la chaîne
	AddValidator(validator interface{})
}
