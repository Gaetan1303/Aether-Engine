package application

import (
	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// CommandeDemarrerCombat - Commande pour démarrer un nouveau combat
type CommandeDemarrerCombat struct {
	CombatID string
	Equipes  []EquipeDTO
	Grille   GrilleDTO
}

// EquipeDTO représente une équipe dans les commandes
type EquipeDTO struct {
	ID       string
	Nom      string
	Couleur  string
	IsIA     bool
	JoueurID *string
	Membres  []UniteDTO
}

// UniteDTO représente une unité dans les commandes
type UniteDTO struct {
	ID       string
	Nom      string
	TeamID   string
	Stats    StatsDTO
	Position PositionDTO
}

// StatsDTO représente des stats dans les commandes
type StatsDTO struct {
	HP      int
	MP      int
	Stamina int
	ATK     int
	DEF     int
	MATK    int
	MDEF    int
	SPD     int
	MOV     int
}

// PositionDTO représente une position dans les commandes
type PositionDTO struct {
	X int
	Y int
}

// GrilleDTO représente une grille dans les commandes
type GrilleDTO struct {
	Largeur int
	Hauteur int
}

// CommandeExecuterAction - Commande pour exécuter une action
type CommandeExecuterAction struct {
	CombatID      string
	ActeurID      string
	TypeAction    string // "attaque", "competence", "deplacement", "objet", "passer"
	CibleID       *string
	PositionCible *PositionDTO
	CompetenceID  *string
	ObjetID       *string
}

// CommandePasserTour - Commande pour passer au tour suivant
type CommandePasserTour struct {
	CombatID string
}

// CommandeTerminerCombat - Commande pour terminer un combat
type CommandeTerminerCombat struct {
	CombatID string
}

// QueryObtenirCombat - Query pour obtenir l'état d'un combat
type QueryObtenirCombat struct {
	CombatID string
}

// CombatDTO représente l'état d'un combat (Read Model)
type CombatDTO struct {
	ID          string
	Etat        string
	Equipes     []EquipeDTO
	Grille      GrilleDTO
	TourActuel  int
	UniteActive string
	Phase       string
	Version     int
}

// ResultatActionDTO représente le résultat d'une action
type ResultatActionDTO struct {
	Succes  bool
	Message string
	Effets  []EffetDTO
}

// EffetDTO représente un effet d'une action
type EffetDTO struct {
	Type    string
	CibleID string
	Valeur  int
}

// Conversion helpers

// ToUnite convertit UniteDTO vers domain.Unite
func (dto UniteDTO) ToUnite() (*domain.Unite, error) {
	stats, err := dto.Stats.ToStats()
	if err != nil {
		return nil, err
	}

	position, err := dto.Position.ToPosition()
	if err != nil {
		return nil, err
	}

	unite := domain.NewUnite(
		domain.UnitID(dto.ID),
		dto.Nom,
		domain.TeamID(dto.TeamID),
		stats,
		position,
	)

	return unite, nil
}

// ToStats convertit StatsDTO vers shared.Stats
func (dto StatsDTO) ToStats() (*shared.Stats, error) {
	return shared.NewStats(
		dto.HP,
		dto.MP,
		dto.Stamina,
		dto.ATK,
		dto.DEF,
		dto.MATK,
		dto.MDEF,
		dto.SPD,
		dto.MOV,
	)
}

// ToPosition convertit PositionDTO vers shared.Position
func (dto PositionDTO) ToPosition() (*shared.Position, error) {
	return shared.NewPosition(dto.X, dto.Y)
}

// ToGrilleCombat convertit GrilleDTO vers shared.GrilleCombat
func (dto GrilleDTO) ToGrilleCombat() (*shared.GrilleCombat, error) {
	return shared.NewGrilleCombat(dto.Largeur, dto.Hauteur)
}

// ToEquipe convertit EquipeDTO vers domain.Equipe
func (dto EquipeDTO) ToEquipe() (*domain.Equipe, error) {
	equipe, err := domain.NewEquipe(
		domain.TeamID(dto.ID),
		dto.Nom,
		dto.Couleur,
		dto.IsIA,
		dto.JoueurID,
	)
	if err != nil {
		return nil, err
	}

	// Ajouter les membres
	for _, membreDTO := range dto.Membres {
		membre, err := membreDTO.ToUnite()
		if err != nil {
			return nil, err
		}
		if err := equipe.AjouterMembre(membre); err != nil {
			return nil, err
		}
	}

	return equipe, nil
}

// FromUnite convertit domain.Unite vers UniteDTO
func FromUnite(unite *domain.Unite) UniteDTO {
	return UniteDTO{
		ID:     string(unite.ID()),
		Nom:    unite.Nom(),
		TeamID: string(unite.TeamID()),
		Stats:  FromStats(unite.Stats()),
		Position: PositionDTO{
			X: unite.Position().X(),
			Y: unite.Position().Y(),
		},
	}
}

// FromStats convertit shared.Stats vers StatsDTO
func FromStats(stats *shared.Stats) StatsDTO {
	return StatsDTO{
		HP:      stats.HP,
		MP:      stats.MP,
		Stamina: stats.Stamina,
		ATK:     stats.ATK,
		DEF:     stats.DEF,
		MATK:    stats.MATK,
		MDEF:    stats.MDEF,
		SPD:     stats.SPD,
		MOV:     stats.MOV,
	}
}

// FromEquipe convertit domain.Equipe vers EquipeDTO
func FromEquipe(equipe *domain.Equipe) EquipeDTO {
	membres := make([]UniteDTO, 0)
	for _, membre := range equipe.Membres() {
		membres = append(membres, FromUnite(membre))
	}

	return EquipeDTO{
		ID:       string(equipe.ID()),
		Nom:      equipe.Nom(),
		Couleur:  equipe.Couleur(),
		IsIA:     equipe.IsIA(),
		JoueurID: equipe.JoueurID(),
		Membres:  membres,
	}
}

// FromCombat convertit domain.Combat vers CombatDTO
func FromCombat(combat *domain.Combat) CombatDTO {
	equipes := make([]EquipeDTO, 0)
	// TODO: Récupérer les équipes du combat

	return CombatDTO{
		ID:          combat.ID(),
		Etat:        combat.Etat().String(),
		Equipes:     equipes,
		TourActuel:  combat.TourActuel(),
		UniteActive: "", // LEGACY - Géré par State Machine maintenant
		Phase:       "", // LEGACY - Géré par State Machine maintenant
		Version:     combat.Version(),
	}
}
