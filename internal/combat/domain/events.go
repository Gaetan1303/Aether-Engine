package domain

import (
	"time"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// Evenement est l'interface de base pour tous les événements
type Evenement interface {
	EventType() string
	AggregateID() string
	AggregateVersion() int
	Timestamp() time.Time
	SetAggregateID(id string)
	SetAggregateVersion(version int)
	SetTimestamp(t time.Time)
}

// BaseEvent implémente les méthodes communes à tous les événements
type BaseEvent struct {
	eventType        string
	aggregateID      string
	aggregateVersion int
	timestamp        time.Time
}

func (e *BaseEvent) EventType() string         { return e.eventType }
func (e *BaseEvent) AggregateID() string       { return e.aggregateID }
func (e *BaseEvent) AggregateVersion() int     { return e.aggregateVersion }
func (e *BaseEvent) Timestamp() time.Time      { return e.timestamp }
func (e *BaseEvent) SetAggregateID(id string)  { e.aggregateID = id }
func (e *BaseEvent) SetAggregateVersion(v int) { e.aggregateVersion = v }
func (e *BaseEvent) SetTimestamp(t time.Time)  { e.timestamp = t }

// CombatDemarreEvent - Le combat a démarré
type CombatDemarreEvent struct {
	BaseEvent
	Tour            int
	OrdreInitiative []UnitID
}

func NewCombatDemarreEvent(combatID string, tour int, ordre []UnitID) *CombatDemarreEvent {
	return &CombatDemarreEvent{
		BaseEvent:       BaseEvent{eventType: "CombatDemarre"},
		Tour:            tour,
		OrdreInitiative: ordre,
	}
}

// TourDemarreEvent - Un nouveau tour démarre
type TourDemarreEvent struct {
	BaseEvent
	Tour int
}

func NewTourDemarreEvent(combatID string, tour int) *TourDemarreEvent {
	return &TourDemarreEvent{
		BaseEvent: BaseEvent{eventType: "TourDemarre"},
		Tour:      tour,
	}
}

// ActionExecuteeEvent - Une action a été exécutée
type ActionExecuteeEvent struct {
	BaseEvent
	Tour     int
	Action   *ActionCombat
	Resultat *ResultatAction
}

func NewActionExecuteeEvent(combatID string, tour int, action *ActionCombat, resultat *ResultatAction) *ActionExecuteeEvent {
	return &ActionExecuteeEvent{
		BaseEvent: BaseEvent{eventType: "ActionExecutee"},
		Tour:      tour,
		Action:    action,
		Resultat:  resultat,
	}
}

// DegatsInfligesEvent - Des dégâts ont été infligés
type DegatsInfligesEvent struct {
	BaseEvent
	Tour     int
	ActeurID UnitID
	CibleID  UnitID
	Degats   int
}

func NewDegatsInfligesEvent(combatID string, tour int, acteurID, cibleID UnitID, degats int) *DegatsInfligesEvent {
	return &DegatsInfligesEvent{
		BaseEvent: BaseEvent{eventType: "DegatsInfliges"},
		Tour:      tour,
		ActeurID:  acteurID,
		CibleID:   cibleID,
		Degats:    degats,
	}
}

// SoinApliqueEvent - Un soin a été appliqué
type SoinApliqueEvent struct {
	BaseEvent
	Tour     int
	ActeurID UnitID
	CibleID  UnitID
	Soin     int
}

func NewSoinApliqueEvent(combatID string, tour int, acteurID, cibleID UnitID, soin int) *SoinApliqueEvent {
	return &SoinApliqueEvent{
		BaseEvent: BaseEvent{eventType: "SoinApplique"},
		Tour:      tour,
		ActeurID:  acteurID,
		CibleID:   cibleID,
		Soin:      soin,
	}
}

// StatutAppliqueEvent - Un statut a été appliqué
type StatutAppliqueEvent struct {
	BaseEvent
	Tour     int
	ActeurID UnitID
	CibleID  UnitID
	Statut   *shared.Statut
}

func NewStatutAppliqueEvent(combatID string, tour int, acteurID, cibleID UnitID, statut *shared.Statut) *StatutAppliqueEvent {
	return &StatutAppliqueEvent{
		BaseEvent: BaseEvent{eventType: "StatutApplique"},
		Tour:      tour,
		ActeurID:  acteurID,
		CibleID:   cibleID,
		Statut:    statut,
	}
}

// UniteElimineeEvent - Une unité a été éliminée
type UniteElimineeEvent struct {
	BaseEvent
	Tour    int
	UniteID UnitID
}

func NewUniteElimineeEvent(combatID string, tour int, uniteID UnitID) *UniteElimineeEvent {
	return &UniteElimineeEvent{
		BaseEvent: BaseEvent{eventType: "UniteEliminee"},
		Tour:      tour,
		UniteID:   uniteID,
	}
}

// UniteDeplaceeEvent - Une unité s'est déplacée
type UniteDeplaceeEvent struct {
	BaseEvent
	Tour            int
	UniteID         UnitID
	PositionDepart  *shared.Position
	PositionArrivee *shared.Position
	CoutDeplacement int
}

func NewUniteDeplaceeEvent(combatID string, tour int, uniteID UnitID, depart, arrivee *shared.Position, cout int) *UniteDeplaceeEvent {
	return &UniteDeplaceeEvent{
		BaseEvent:       BaseEvent{eventType: "UniteDeplacee"},
		Tour:            tour,
		UniteID:         uniteID,
		PositionDepart:  depart,
		PositionArrivee: arrivee,
		CoutDeplacement: cout,
	}
}

// CompetenceUtiliseeEvent - Une compétence a été utilisée
type CompetenceUtiliseeEvent struct {
	BaseEvent
	Tour         int
	ActeurID     UnitID
	CompetenceID CompetenceID
	Cibles       []UnitID
}

func NewCompetenceUtiliseeEvent(combatID string, tour int, acteurID UnitID, compID CompetenceID, cibles []UnitID) *CompetenceUtiliseeEvent {
	return &CompetenceUtiliseeEvent{
		BaseEvent:    BaseEvent{eventType: "CompetenceUtilisee"},
		Tour:         tour,
		ActeurID:     acteurID,
		CompetenceID: compID,
		Cibles:       cibles,
	}
}

// CombatTermineEvent - Le combat est terminé
type CombatTermineEvent struct {
	BaseEvent
	Tour      int
	Vainqueur *TeamID
}

func NewCombatTermineEvent(combatID string, tour int, vainqueur *TeamID) *CombatTermineEvent {
	return &CombatTermineEvent{
		BaseEvent: BaseEvent{eventType: "CombatTermine"},
		Tour:      tour,
		Vainqueur: vainqueur,
	}
}

// ActionCombat représente une action à exécuter
type ActionCombat struct {
	Type          TypeAction
	ActeurID      UnitID
	CibleID       *UnitID // Peut être nil pour actions sans cible
	PositionCible *shared.Position
	CompetenceID  *CompetenceID
	ObjetID       *shared.ObjetID
}

// TypeAction énumère les types d'actions possibles
type TypeAction int

const (
	TypeActionAttaque TypeAction = iota
	TypeActionCompetence
	TypeActionDeplacement
	TypeActionObjet
	TypeActionPasser
)

// ResultatAction représente le résultat d'une action
type ResultatAction struct {
	ActeurID   UnitID
	TypeAction TypeAction
	Succes     bool
	Message    string
	Effets     []EffetAction
}

// EffetAction représente un effet produit par une action
type EffetAction struct {
	Type     TypeEffetAction
	CibleID  UnitID
	Valeur   int            // Dégâts/soin
	Statut   *shared.Statut // Pour effets de statut
	Position *shared.Position
}

// TypeEffetAction énumère les types d'effets d'actions
type TypeEffetAction int

const (
	TypeEffetActionDegats TypeEffetAction = iota
	TypeEffetActionSoin
	TypeEffetActionStatut
	TypeEffetActionDeplacement
	TypeEffetActionModificationStat
)
