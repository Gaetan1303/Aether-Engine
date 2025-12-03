package application

import (
	"errors"

	"github.com/aether-engine/aether-engine/internal/combat/combatfacade"
	"github.com/aether-engine/aether-engine/internal/combat/combatinitializer"
	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/aether-engine/aether-engine/internal/combat/domain/states"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
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
// Refactoré avec Extract Method Pattern pour réduire la complexité
func (e *CombatEngineImpl) DemarrerCombat(cmd CommandeDemarrerCombat) (*CombatDTO, error) {
	// Valider la commande
	if err := validateDemarrerCommand(cmd); err != nil {
		return nil, err
	}

	// Construire le domain model
	equipes, grille, err := buildDomainModel(cmd)
	if err != nil {
		return nil, err
	}

	// Créer et démarrer l'agrégat Combat
	combat, err := domain.NewCombat(cmd.CombatID, equipes, grille)
	if err != nil {
		return nil, err
	}

	if err := combat.Demarrer(); err != nil {
		return nil, err
	}

	// Initialiser les patterns Step C (State Machine, Commands, Observers, Validation)
	initializer := combatinitializer.NewCombatInitializer(combat)
	if err := initializer.InitializeAll(); err != nil {
		return nil, err
	}

	// Sauvegarder et publier les événements
	if err := e.saveAndPublishEvents(cmd.CombatID, combat); err != nil {
		return nil, err
	}

	// Retourner le DTO
	dto := FromCombat(combat)
	return &dto, nil
}

// ExecuterAction exécute une action dans un combat via la State Machine
// Refactoré avec Extract Method Pattern pour réduire la complexité
func (e *CombatEngineImpl) ExecuterAction(cmd CommandeExecuterAction) (*ResultatActionDTO, error) {
	// Charger le combat depuis l'Event Store
	combat, err := e.loadCombatFromEvents(cmd.CombatID)
	if err != nil {
		return nil, err
	}

	// Convertir le type d'action vers CommandType de la facade
	actionType := convertToCommandType(parseTypeAction(cmd.TypeAction))

	// Préparer les paramètres pour ExecutePlayerAction
	params, err := buildActionParameters(cmd)
	if err != nil {
		return nil, err
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

	// Sauvegarder et publier les événements
	if err := e.saveAndPublishEvents(cmd.CombatID, combat); err != nil {
		return nil, err
	}

	// Convertir le résultat vers DTO
	resultat := &ResultatActionDTO{
		Succes:  result.Success,
		Message: result.Message,
		Effets:  []EffetDTO{}, // TODO: implémenter convertEffectsToDTO si nécessaire
	}

	return resultat, nil
}

// PasserTour passe au tour suivant via la State Machine
// Refactoré avec Extract Method Pattern pour réduire la duplication
func (e *CombatEngineImpl) PasserTour(cmd CommandePasserTour) error {
	// Charger le combat
	combat, err := e.loadCombatFromEvents(cmd.CombatID)
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

	// Sauvegarder et publier les événements
	return e.saveAndPublishEvents(cmd.CombatID, combat)
}

// TerminerCombat termine un combat via la State Machine
// Refactoré avec Extract Method Pattern pour réduire la duplication
func (e *CombatEngineImpl) TerminerCombat(cmd CommandeTerminerCombat) error {
	// Charger le combat
	combat, err := e.loadCombatFromEvents(cmd.CombatID)
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

	// Sauvegarder et publier les événements
	return e.saveAndPublishEvents(cmd.CombatID, combat)
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

// --- Helper Methods (Extract Method Pattern) ---

// saveAndPublishEvents sauvegarde et publie les événements (DRY pattern)
func (e *CombatEngineImpl) saveAndPublishEvents(combatID string, combat *domain.Combat) error {
	newEvents := combat.GetUncommittedEvents()
	if len(newEvents) == 0 {
		return nil
	}

	// Sauvegarder les événements
	if err := e.eventStore.AppendEvents(combatID, newEvents, combat.Version()); err != nil {
		return err
	}

	// Publier les événements
	for _, evt := range newEvents {
		if err := e.publisher.Publish(evt); err != nil {
			// Log error mais ne pas échouer la transaction
			// TODO: implémenter proper logging
		}
	}

	// Clear uncommitted events
	combat.ClearUncommittedEvents()
	return nil
}

// loadCombatFromEvents charge un combat depuis l'Event Store
func (e *CombatEngineImpl) loadCombatFromEvents(combatID string) (*domain.Combat, error) {
	events, err := e.eventStore.LoadEvents(combatID)
	if err != nil {
		return nil, err
	}

	combat, err := domain.ReconstruireDepuisEvenements(events)
	if err != nil {
		return nil, err
	}

	return combat, nil
}

// validateDemarrerCommand valide la commande DemarrerCombat
func validateDemarrerCommand(cmd CommandeDemarrerCombat) error {
	if cmd.CombatID == "" {
		return errors.New("CombatID requis")
	}

	if len(cmd.Equipes) < 2 {
		return errors.New("au moins 2 équipes requises")
	}

	return nil
}

// buildDomainModel construit le modèle du domaine depuis les DTOs
func buildDomainModel(cmd CommandeDemarrerCombat) ([]*domain.Equipe, *shared.GrilleCombat, error) {
	// Construire les équipes
	equipes := make([]*domain.Equipe, 0)
	for _, equipeDTO := range cmd.Equipes {
		equipe, err := equipeDTO.ToEquipe()
		if err != nil {
			return nil, nil, err
		}
		equipes = append(equipes, equipe)
	}

	// Construire la grille
	grille, err := cmd.Grille.ToGrilleCombat()
	if err != nil {
		return nil, nil, err
	}

	return equipes, grille, nil
}

// convertToCommandType convertit TypeAction vers CommandType de la facade
func convertToCommandType(typeAction domain.TypeAction) combatfacade.CommandType {
	switch typeAction {
	case domain.TypeActionAttaque:
		return combatfacade.CommandTypeAttack
	case domain.TypeActionCompetence:
		return combatfacade.CommandTypeSkill
	case domain.TypeActionDeplacement:
		return combatfacade.CommandTypeMove
	case domain.TypeActionObjet:
		return combatfacade.CommandTypeItem
	case domain.TypeActionPasser:
		return combatfacade.CommandTypeWait
	default:
		return combatfacade.CommandTypeWait
	}
}

// buildActionParameters construit les paramètres pour ExecutePlayerAction
func buildActionParameters(cmd CommandeExecuterAction) (map[string]interface{}, error) {
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

	return params, nil
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
