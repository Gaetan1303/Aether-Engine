package bases_donnees

import (
	"encoding/json"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

// Event représente un événement générique
type Event struct {
	EventID          uuid.UUID              `json:"event_id"`
	AggregateType    string                 `json:"aggregate_type"`
	AggregateID      uuid.UUID              `json:"aggregate_id"`
	AggregateVersion int                    `json:"aggregate_version"`
	EventType        string                 `json:"event_type"`
	EventData        map[string]interface{} `json:"event_data"`
	Metadata         map[string]interface{} `json:"metadata"`
	TimestampUTC     time.Time              `json:"timestamp_utc"`
}

// EventBuilder aide à construire des événements de test
type EventBuilder struct {
	event Event
}

// NewEventBuilder crée un nouveau builder d'événements
func NewEventBuilder() *EventBuilder {
	return &EventBuilder{
		event: Event{
			EventID:      uuid.New(),
			AggregateID:  uuid.New(),
			EventData:    make(map[string]interface{}),
			Metadata:     make(map[string]interface{}),
			TimestampUTC: time.Now().UTC(),
		},
	}
}

func (b *EventBuilder) WithEventID(id uuid.UUID) *EventBuilder {
	b.event.EventID = id
	return b
}

func (b *EventBuilder) WithAggregateType(t string) *EventBuilder {
	b.event.AggregateType = t
	return b
}

func (b *EventBuilder) WithAggregateID(id uuid.UUID) *EventBuilder {
	b.event.AggregateID = id
	return b
}

func (b *EventBuilder) WithAggregateVersion(v int) *EventBuilder {
	b.event.AggregateVersion = v
	return b
}

func (b *EventBuilder) WithEventType(t string) *EventBuilder {
	b.event.EventType = t
	return b
}

func (b *EventBuilder) WithEventData(data map[string]interface{}) *EventBuilder {
	b.event.EventData = data
	return b
}

func (b *EventBuilder) WithMetadata(metadata map[string]interface{}) *EventBuilder {
	b.event.Metadata = metadata
	return b
}

func (b *EventBuilder) WithTimestamp(t time.Time) *EventBuilder {
	b.event.TimestampUTC = t
	return b
}

func (b *EventBuilder) Build() Event {
	return b.event
}

// --- Événements Joueur ---

// NewJoueurCreeEvent crée un événement JoueurCree pour les tests
func NewJoueurCreeEvent() Event {
	joueurID := uuid.New()
	username := gofakeit.Username()

	return NewEventBuilder().
		WithAggregateType("Joueur").
		WithAggregateID(joueurID).
		WithAggregateVersion(1).
		WithEventType("JoueurCree").
		WithEventData(map[string]interface{}{
			"joueur_id": joueurID.String(),
			"username":  username,
			"email":     gofakeit.Email(),
			"classe":    "GUERRIER",
			"race":      "HUMAIN",
		}).
		WithMetadata(map[string]interface{}{
			"ip_address": gofakeit.IPv4Address(),
			"user_agent": gofakeit.UserAgent(),
		}).
		Build()
}

// NewNiveauGagneEvent crée un événement NiveauGagne
func NewNiveauGagneEvent(joueurID uuid.UUID, version int, nouveauNiveau int) Event {
	return NewEventBuilder().
		WithAggregateType("Joueur").
		WithAggregateID(joueurID).
		WithAggregateVersion(version).
		WithEventType("NiveauGagne").
		WithEventData(map[string]interface{}{
			"joueur_id":      joueurID.String(),
			"ancien_niveau":  nouveauNiveau - 1,
			"nouveau_niveau": nouveauNiveau,
			"bonus_hp":       10,
			"bonus_mana":     5,
			"points_stat":    5,
		}).
		Build()
}

// NewExperienceGagneeEvent crée un événement ExperienceGagnee
func NewExperienceGagneeEvent(joueurID uuid.UUID, version int, xpGagne int) Event {
	return NewEventBuilder().
		WithAggregateType("Joueur").
		WithAggregateID(joueurID).
		WithAggregateVersion(version).
		WithEventType("ExperienceGagnee").
		WithEventData(map[string]interface{}{
			"joueur_id":           joueurID.String(),
			"experience_gagnee":   xpGagne,
			"experience_actuelle": xpGagne,
			"source":              "COMBAT",
		}).
		Build()
}

// --- Événements Combat ---

// NewCombatDemarreEvent crée un événement CombatDemarre
func NewCombatDemarreEvent() Event {
	combatID := uuid.New()
	joueur1ID := uuid.New()
	joueur2ID := uuid.New()

	return NewEventBuilder().
		WithAggregateType("Combat").
		WithAggregateID(combatID).
		WithAggregateVersion(1).
		WithEventType("CombatDemarre").
		WithEventData(map[string]interface{}{
			"combat_id":   combatID.String(),
			"type_combat": "PVP",
			"participants": []map[string]interface{}{
				{
					"joueur_id":   joueur1ID.String(),
					"initiative":  15,
					"hp_actuel":   100,
					"hp_max":      100,
					"mana_actuel": 50,
					"mana_max":    50,
				},
				{
					"joueur_id":   joueur2ID.String(),
					"initiative":  12,
					"hp_actuel":   100,
					"hp_max":      100,
					"mana_actuel": 50,
					"mana_max":    50,
				},
			},
			"ordre_tours": []string{joueur1ID.String(), joueur2ID.String()},
		}).
		Build()
}

// NewActionCombatExecuteeEvent crée un événement ActionCombatExecutee
func NewActionCombatExecuteeEvent(combatID uuid.UUID, version int, acteurID, cibleID uuid.UUID) Event {
	degats := gofakeit.Number(10, 30)

	return NewEventBuilder().
		WithAggregateType("Combat").
		WithAggregateID(combatID).
		WithAggregateVersion(version).
		WithEventType("ActionCombatExecutee").
		WithEventData(map[string]interface{}{
			"combat_id":   combatID.String(),
			"tour":        version,
			"acteur_id":   acteurID.String(),
			"type_action": "ATTAQUE",
			"cible_id":    cibleID.String(),
			"reussi":      true,
			"critique":    false,
			"degats":      degats,
		}).
		Build()
}

// NewDegatsInfligesEvent crée un événement DegatsInfliges
func NewDegatsInfligesEvent(combatID uuid.UUID, version int, joueurID uuid.UUID, degats int) Event {
	hpRestant := 100 - degats

	return NewEventBuilder().
		WithAggregateType("Combat").
		WithAggregateID(combatID).
		WithAggregateVersion(version).
		WithEventType("DegatsInfliges").
		WithEventData(map[string]interface{}{
			"combat_id":   combatID.String(),
			"joueur_id":   joueurID.String(),
			"degats":      degats,
			"hp_avant":    100,
			"hp_apres":    hpRestant,
			"type_degats": "PHYSIQUE",
		}).
		Build()
}

// NewEffetStatutAppliqueEvent crée un événement EffetStatutApplique
func NewEffetStatutAppliqueEvent(combatID uuid.UUID, version int, joueurID uuid.UUID) Event {
	return NewEventBuilder().
		WithAggregateType("Combat").
		WithAggregateID(combatID).
		WithAggregateVersion(version).
		WithEventType("EffetStatutApplique").
		WithEventData(map[string]interface{}{
			"combat_id":    combatID.String(),
			"joueur_id":    joueurID.String(),
			"type_effet":   "POISON",
			"puissance":    5,
			"duree_tours":  3,
			"applique_par": uuid.New().String(),
		}).
		Build()
}

// NewCombatTermineEvent crée un événement CombatTermine
func NewCombatTermineEvent(combatID uuid.UUID, version int, vainqueurID uuid.UUID) Event {
	return NewEventBuilder().
		WithAggregateType("Combat").
		WithAggregateID(combatID).
		WithAggregateVersion(version).
		WithEventType("CombatTermine").
		WithEventData(map[string]interface{}{
			"combat_id":      combatID.String(),
			"vainqueur_id":   vainqueurID.String(),
			"raison_fin":     "KO",
			"duree_secondes": 120,
			"recompenses": map[string]interface{}{
				"experience": 100,
				"or":         50,
			},
		}).
		Build()
}

// --- Événements Inventaire ---

// NewItemAjouteEvent crée un événement ItemAjoute
func NewItemAjouteEvent(joueurID uuid.UUID, version int) Event {
	return NewEventBuilder().
		WithAggregateType("Inventaire").
		WithAggregateID(joueurID).
		WithAggregateVersion(version).
		WithEventType("ItemAjoute").
		WithEventData(map[string]interface{}{
			"joueur_id": joueurID.String(),
			"item_id":   "EPEE_FER",
			"quantite":  1,
			"slot":      1,
			"poids":     2.5,
			"source":    "LOOT",
		}).
		Build()
}

// NewItemEquipeEvent crée un événement ItemEquipe
func NewItemEquipeEvent(joueurID uuid.UUID, version int, itemID string, slot string) Event {
	return NewEventBuilder().
		WithAggregateType("Inventaire").
		WithAggregateID(joueurID).
		WithAggregateVersion(version).
		WithEventType("ItemEquipe").
		WithEventData(map[string]interface{}{
			"joueur_id":       joueurID.String(),
			"item_id":         itemID,
			"slot_equipement": slot,
			"ancien_item_id":  nil,
		}).
		Build()
}

// NewItemUtiliseEvent crée un événement ItemUtilise
func NewItemUtiliseEvent(joueurID uuid.UUID, version int) Event {
	return NewEventBuilder().
		WithAggregateType("Inventaire").
		WithAggregateID(joueurID).
		WithAggregateVersion(version).
		WithEventType("ItemUtilise").
		WithEventData(map[string]interface{}{
			"joueur_id":         joueurID.String(),
			"item_id":           "POTION_VIE",
			"quantite_utilisee": 1,
			"quantite_restante": 2,
			"effet": map[string]interface{}{
				"type":   "SOIN",
				"valeur": 50,
			},
		}).
		Build()
}

// NewItemEchangeEvent crée un événement ItemEchange
func NewItemEchangeEvent(joueur1ID, joueur2ID uuid.UUID, version int) Event {
	return NewEventBuilder().
		WithAggregateType("Inventaire").
		WithAggregateID(joueur1ID).
		WithAggregateVersion(version).
		WithEventType("ItemEchange").
		WithEventData(map[string]interface{}{
			"joueur_source_id": joueur1ID.String(),
			"joueur_cible_id":  joueur2ID.String(),
			"item_id":          "EPEE_FER",
			"quantite":         1,
			"prix_or":          100,
		}).
		Build()
}

// --- Événements Économie ---

// NewOrdreEconomieCreeEvent crée un événement OrdreEconomieCree
func NewOrdreEconomieCreeEvent(joueurID uuid.UUID, version int) Event {
	ordreID := uuid.New()

	return NewEventBuilder().
		WithAggregateType("Economie").
		WithAggregateID(ordreID).
		WithAggregateVersion(version).
		WithEventType("OrdreEconomieCree").
		WithEventData(map[string]interface{}{
			"ordre_id":   ordreID.String(),
			"joueur_id":  joueurID.String(),
			"type_ordre": "VENTE",
			"item_id":    "EPEE_FER",
			"quantite":   5,
			"prix_unite": 100,
			"expire_a":   time.Now().Add(24 * time.Hour).Unix(),
		}).
		Build()
}

// NewTransactionEconomieExecuteeEvent crée un événement TransactionEconomieExecutee
func NewTransactionEconomieExecuteeEvent(acheteurID, vendeurID uuid.UUID, version int) Event {
	transactionID := uuid.New()

	return NewEventBuilder().
		WithAggregateType("Economie").
		WithAggregateID(transactionID).
		WithAggregateVersion(version).
		WithEventType("TransactionEconomieExecutee").
		WithEventData(map[string]interface{}{
			"transaction_id": transactionID.String(),
			"acheteur_id":    acheteurID.String(),
			"vendeur_id":     vendeurID.String(),
			"item_id":        "EPEE_FER",
			"quantite":       2,
			"prix_unite":     100,
			"prix_total":     200,
		}).
		Build()
}

// --- Événements Quête ---

// NewQueteAccepteeEvent crée un événement QueteAcceptee
func NewQueteAccepteeEvent(joueurID uuid.UUID, version int) Event {
	return NewEventBuilder().
		WithAggregateType("Joueur").
		WithAggregateID(joueurID).
		WithAggregateVersion(version).
		WithEventType("QueteAcceptee").
		WithEventData(map[string]interface{}{
			"joueur_id": joueurID.String(),
			"quete_id":  "QUETE_DEBUTANT_01",
			"objectifs": []map[string]interface{}{
				{
					"objectif_id": "TUE_10_SLIMES",
					"type":        "TUER",
					"cible":       "SLIME",
					"progression": 0,
					"requis":      10,
				},
			},
		}).
		Build()
}

// NewObjectifQueteProgresseEvent crée un événement ObjectifQueteProgresse
func NewObjectifQueteProgresseEvent(joueurID uuid.UUID, version int) Event {
	return NewEventBuilder().
		WithAggregateType("Joueur").
		WithAggregateID(joueurID).
		WithAggregateVersion(version).
		WithEventType("ObjectifQueteProgresse").
		WithEventData(map[string]interface{}{
			"joueur_id":         joueurID.String(),
			"quete_id":          "QUETE_DEBUTANT_01",
			"objectif_id":       "TUE_10_SLIMES",
			"progression_avant": 5,
			"progression_apres": 6,
			"complete":          false,
		}).
		Build()
}

// --- Helpers pour créer des agrégats complets ---

// JoueurAggregate représente l'état complet d'un joueur
type JoueurAggregate struct {
	ID                 uuid.UUID
	Username           string
	Niveau             int
	ExperienceActuelle int64
	HPActuel           int
	HPMax              int
	ManaActuel         int
	ManaMax            int
	Or                 int64
	Version            int
}

// NewJoueurAggregate crée un nouvel agrégat joueur pour les tests
func NewJoueurAggregate() JoueurAggregate {
	return JoueurAggregate{
		ID:                 uuid.New(),
		Username:           gofakeit.Username(),
		Niveau:             1,
		ExperienceActuelle: 0,
		HPActuel:           100,
		HPMax:              100,
		ManaActuel:         50,
		ManaMax:            50,
		Or:                 100,
		Version:            1,
	}
}

// ToJSON convertit l'agrégat en JSON
func (j JoueurAggregate) ToJSON() ([]byte, error) {
	return json.Marshal(j)
}

// CombatAggregate représente l'état complet d'un combat
type CombatAggregate struct {
	ID           uuid.UUID
	TypeCombat   string
	Etat         string
	TourActuel   int
	Participants []ParticipantCombat
	OrdreTours   []uuid.UUID
	Version      int
}

type ParticipantCombat struct {
	JoueurID   uuid.UUID
	Initiative int
	HPActuel   int
	HPMax      int
	ManaActuel int
	ManaMax    int
	Etat       string
}

// NewCombatAggregate crée un nouvel agrégat combat
func NewCombatAggregate(joueur1, joueur2 JoueurAggregate) CombatAggregate {
	return CombatAggregate{
		ID:         uuid.New(),
		TypeCombat: "PVP",
		Etat:       "EN_COURS",
		TourActuel: 1,
		Participants: []ParticipantCombat{
			{
				JoueurID:   joueur1.ID,
				Initiative: 15,
				HPActuel:   joueur1.HPActuel,
				HPMax:      joueur1.HPMax,
				ManaActuel: joueur1.ManaActuel,
				ManaMax:    joueur1.ManaMax,
				Etat:       "ACTIF",
			},
			{
				JoueurID:   joueur2.ID,
				Initiative: 12,
				HPActuel:   joueur2.HPActuel,
				HPMax:      joueur2.HPMax,
				ManaActuel: joueur2.ManaActuel,
				ManaMax:    joueur2.ManaMax,
				Etat:       "ACTIF",
			},
		},
		OrdreTours: []uuid.UUID{joueur1.ID, joueur2.ID},
		Version:    1,
	}
}

// ToJSON convertit l'agrégat en JSON
func (c CombatAggregate) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

// InventaireAggregate représente l'état complet d'un inventaire
type InventaireAggregate struct {
	JoueurID    uuid.UUID
	CapaciteMax int
	PoidsActuel float64
	Items       []ItemInventaire
	Version     int
}

type ItemInventaire struct {
	Slot     int
	ItemID   string
	Quantite int
	Equipe   bool
}

// NewInventaireAggregate crée un nouvel agrégat inventaire
func NewInventaireAggregate(joueurID uuid.UUID) InventaireAggregate {
	return InventaireAggregate{
		JoueurID:    joueurID,
		CapaciteMax: 50,
		PoidsActuel: 0,
		Items:       []ItemInventaire{},
		Version:     1,
	}
}

// ToJSON convertit l'agrégat en JSON
func (i InventaireAggregate) ToJSON() ([]byte, error) {
	return json.Marshal(i)
}
