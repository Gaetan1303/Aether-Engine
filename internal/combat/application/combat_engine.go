package application

import (
	"errors"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
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

// ExecuterAction exécute une action dans un combat
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

	// Construire l'action
	action := &domain.ActionCombat{
		Type:     parseTypeAction(cmd.TypeAction),
		ActeurID: domain.UnitID(cmd.ActeurID),
	}

	if cmd.CibleID != nil {
		cibleID := domain.UnitID(*cmd.CibleID)
		action.CibleID = &cibleID
	}

	if cmd.PositionCible != nil {
		pos, err := cmd.PositionCible.ToPosition()
		if err != nil {
			return nil, err
		}
		action.PositionCible = pos
	}

	if cmd.CompetenceID != nil {
		compID := domain.CompetenceID(*cmd.CompetenceID)
		action.CompetenceID = &compID
	}

	if cmd.ObjetID != nil {
		objID := shared.ObjetID(*cmd.ObjetID)
		action.ObjetID = &objID
	}

	// Exécuter l'action
	if err := combat.ExecuterAction(action); err != nil {
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

	// Clear uncommitted events
	combat.ClearUncommittedEvents()

	// Retourner le résultat
	// TODO: Extraire le résultat réel
	resultat := &ResultatActionDTO{
		Succes:  true,
		Message: "Action exécutée avec succès",
		Effets:  make([]EffetDTO, 0),
	}

	return resultat, nil
}

// PasserTour passe au tour suivant
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

	// Passer au tour suivant
	if err := combat.TourSuivant(); err != nil {
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

// TerminerCombat termine un combat
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

	// Terminer le combat
	if err := combat.Terminer(); err != nil {
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
