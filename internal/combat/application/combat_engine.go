package application

import (
	"errors"

	"github.com/aether-engine/aether-engine/internal/combat/combatfacade"
	"github.com/aether-engine/aether-engine/internal/combat/combatinitializer"
	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/aether-engine/aether-engine/internal/combat/domain/states"
)

// CombatEngine est le service applicatif qui orchestre les combats
type CombatEngine interface {
	// DemarrerCombat démarre un nouveau combat
	DemarrerCombat(cmd CommandeDemarrerCombat) (*CombatDTO, error)

	// ExecuterAction exécute une action dans un combat
	ExecuterAction(cmd CommandeExecuterAction) (*ResultatActionDTO, error)

	// PasserTour passe au tour suivant
	PasserTour(cmd CommandePasserTour) error

	// TerminerCombat termine un combat
	TerminerCombat(cmd CommandeTerminerCombat) error

	// ObtenirCombat récupère l'état d'un combat
	ObtenirCombat(query QueryObtenirCombat) (*CombatDTO, error)
}

// CombatEngineImpl implémente CombatEngine
type CombatEngineImpl struct {
	eventStore EventStore
	publisher  EventPublisher
}

// NewCombatEngine crée une nouvelle instance du moteur de combat
func NewCombatEngine(eventStore EventStore, publisher EventPublisher) CombatEngine {
	return &CombatEngineImpl{
		eventStore: eventStore,
		publisher:  publisher,
	}
}

// DemarrerCombat démarre un nouveau combat
func (e *CombatEngineImpl) DemarrerCombat(cmd CommandeDemarrerCombat) (*CombatDTO, error) {
	// Valider la commande
	if cmd.CombatID == "" {
		return nil, errors.New("CombatID requis")
	}

	if len(cmd.Equipes) < 2 {
		return nil, errors.New("au moins 2 équipes requises")
	}

	// Construire le domain model
	equipes := make([]*domain.Equipe, 0)
	for _, equipeDTO := range cmd.Equipes {
		equipe, err := equipeDTO.ToEquipe()
		if err != nil {
			return nil, err
		}
		equipes = append(equipes, equipe)
	}

	grille, err := cmd.Grille.ToGrilleCombat()
	if err != nil {
		return nil, err
	}

	// Créer l'agrégat Combat
	combat, err := domain.NewCombat(cmd.CombatID, equipes, grille)
	if err != nil {
		return nil, err
	}

	// Démarrer le combat
	if err := combat.Demarrer(); err != nil {
		return nil, err
	}

	// Initialiser les patterns Step C (State Machine, Commands, Observers, Validation)
	initializer := combatinitializer.NewCombatInitializer(combat)
	if err := initializer.InitializeAll(); err != nil {
		return nil, err
	}

	// Sauvegarder les événements
	events := combat.GetUncommittedEvents()
	if err := e.eventStore.AppendEvents(cmd.CombatID, events, 0); err != nil {
		return nil, err
	}

	// Publier les événements
	for _, evt := range events {
		if err := e.publisher.Publish(evt); err != nil {
			// Log error mais ne pas échouer
		}
	}

	// Clear uncommitted events
	combat.ClearUncommittedEvents()

	// Retourner le DTO
	dto := FromCombat(combat)
	return &dto, nil
}

// ExecuterAction exécute une action dans un combat via la State Machine
func (e *CombatEngineImpl) ExecuterAction(cmd CommandeExecuterAction) (*ResultatActionDTO, error) {
	// Charger le combat depuis l'Event Store
	events, err := e.eventStore.LoadEvents(cmd.CombatID)
	if err != nil {
		return nil, err
	}

	combat, err := domain.ReconstruireDepuisEvenements(events)
	if err != nil {
		return nil, err
	}

	// Convertir le type d'action vers CommandType de la facade
	var actionType combatfacade.CommandType
	switch parseTypeAction(cmd.TypeAction) {
	case domain.TypeActionAttaque:
		actionType = combatfacade.CommandTypeAttack
	case domain.TypeActionCompetence:
		actionType = combatfacade.CommandTypeSkill
	case domain.TypeActionDeplacement:
		actionType = combatfacade.CommandTypeMove
	case domain.TypeActionObjet:
		actionType = combatfacade.CommandTypeItem
	case domain.TypeActionPasser:
		actionType = combatfacade.CommandTypeWait
	default:
		actionType = combatfacade.CommandTypeWait
	}

	// Préparer les paramètres pour ExecutePlayerAction
	params := make(map[string]interface{})

	if cmd.CibleID != nil {
		params["targetID"] = domain.UnitID(*cmd.CibleID)
	}

	if cmd.PositionCible != nil {
		pos, err := cmd.PositionCible.ToPosition()
		if err != nil {
			return nil, err
		}
		params["targetX"] = pos.X()
		params["targetY"] = pos.Y()
	}

	if cmd.CompetenceID != nil {
		params["skillID"] = string(*cmd.CompetenceID)
		if cmd.CibleID != nil {
			params["targetID"] = domain.UnitID(*cmd.CibleID)
		}
	}

	if cmd.ObjetID != nil {
		params["itemID"] = string(*cmd.ObjetID)
		if cmd.CibleID != nil {
			params["targetID"] = domain.UnitID(*cmd.CibleID)
		}
	}

	// Exécuter via combatfacade (architecture Step C)
	result, err := combatfacade.ExecutePlayerAction(
		combat,
		domain.UnitID(cmd.ActeurID),
		actionType,
		params,
	)
	if err != nil {
		return nil, err
	}

	// Sauvegarder les nouveaux événements
	newEvents := combat.GetUncommittedEvents()
	if err := e.eventStore.AppendEvents(cmd.CombatID, newEvents, combat.Version()); err != nil {
		return nil, err
	}

	// Publier les événements
	for _, evt := range newEvents {
		if err := e.publisher.Publish(evt); err != nil {
			// Log error mais ne pas échouer
		}
	}

	combat.ClearUncommittedEvents()

	// Convertir le résultat vers DTO
	resultat := &ResultatActionDTO{
		Succes:  result.Success,
		Message: result.Message,
		Effets:  []EffetDTO{}, // TODO: implémenter convertEffectsToDTO si nécessaire
	}

	return resultat, nil
}

// PasserTour passe au tour suivant via la State Machine
func (e *CombatEngineImpl) PasserTour(cmd CommandePasserTour) error {
	// Charger le combat
	events, err := e.eventStore.LoadEvents(cmd.CombatID)
	if err != nil {
		return err
	}

	combat, err := domain.ReconstruireDepuisEvenements(events)
	if err != nil {
		return err
	}

	// Récupérer la State Machine
	sm := combatfacade.GetStateMachine(combat)
	if sm == nil {
		return errors.New("state machine non initialisée")
	}

	// Déclencher l'événement de fin de tour
	event := states.StateEvent{Type: states.EventTurnComplete}
	if err := sm.HandleEvent(event); err != nil {
		return err
	}

	// Sauvegarder les événements
	newEvents := combat.GetUncommittedEvents()
	if err := e.eventStore.AppendEvents(cmd.CombatID, newEvents, combat.Version()); err != nil {
		return err
	}

	// Publier les événements
	for _, evt := range newEvents {
		if err := e.publisher.Publish(evt); err != nil {
			// Log error
		}
	}

	combat.ClearUncommittedEvents()

	return nil
}

// TerminerCombat termine un combat via la State Machine
func (e *CombatEngineImpl) TerminerCombat(cmd CommandeTerminerCombat) error {
	// Charger le combat
	events, err := e.eventStore.LoadEvents(cmd.CombatID)
	if err != nil {
		return err
	}

	combat, err := domain.ReconstruireDepuisEvenements(events)
	if err != nil {
		return err
	}

	// Récupérer la State Machine
	sm := combatfacade.GetStateMachine(combat)
	if sm == nil {
		return errors.New("state machine non initialisée")
	}

	// Déclencher l'événement de finalisation
	event := states.StateEvent{Type: states.EventFinalizeCombat}
	if err := sm.HandleEvent(event); err != nil {
		return err
	}

	// Distribuer les récompenses
	combat.DistribuerRecompenses()

	// Sauvegarder les événements
	newEvents := combat.GetUncommittedEvents()
	if err := e.eventStore.AppendEvents(cmd.CombatID, newEvents, combat.Version()); err != nil {
		return err
	}

	// Publier les événements
	for _, evt := range newEvents {
		if err := e.publisher.Publish(evt); err != nil {
			// Log error
		}
	}

	combat.ClearUncommittedEvents()

	return nil
}

// ObtenirCombat récupère l'état d'un combat
func (e *CombatEngineImpl) ObtenirCombat(query QueryObtenirCombat) (*CombatDTO, error) {
	// TODO: Lire depuis une projection (read model) plutôt que reconstruire
	// Pour l'instant, reconstruire depuis les événements
	events, err := e.eventStore.LoadEvents(query.CombatID)
	if err != nil {
		return nil, err
	}

	combat, err := domain.ReconstruireDepuisEvenements(events)
	if err != nil {
		return nil, err
	}

	dto := FromCombat(combat)
	return &dto, nil
}

// parseTypeAction convertit une string en TypeAction
func parseTypeAction(typeStr string) domain.TypeAction {
	switch typeStr {
	case "attaque":
		return domain.TypeActionAttaque
	case "competence":
		return domain.TypeActionCompetence
	case "deplacement":
		return domain.TypeActionDeplacement
	case "objet":
		return domain.TypeActionObjet
	case "passer":
		return domain.TypeActionPasser
	default:
		return domain.TypeActionPasser
	}
}

// EventStore interface pour persister les événements
type EventStore interface {
	// AppendEvents ajoute des événements à un agrégat
	AppendEvents(aggregateID string, events []domain.Evenement, expectedVersion int) error

	// LoadEvents charge tous les événements d'un agrégat
	LoadEvents(aggregateID string) ([]domain.Evenement, error)

	// LoadEventsFromVersion charge les événements depuis une version
	LoadEventsFromVersion(aggregateID string, fromVersion int) ([]domain.Evenement, error)

	// SaveSnapshot sauvegarde un snapshot
	SaveSnapshot(aggregateID string, version int, data []byte) error

	// LoadSnapshot charge le dernier snapshot
	LoadSnapshot(aggregateID string) (version int, data []byte, err error)
}

// EventPublisher interface pour publier les événements
type EventPublisher interface {
	// Publish publie un événement
	Publish(event domain.Evenement) error
}
