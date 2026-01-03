package dto

import (
	"time"

	"github.com/aether-engine/aether-engine/internal/joueur/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// JoueurCreateDTO données pour créer un joueur via API
type JoueurCreateDTO struct {
	Nom       string                   `json:"nom" binding:"required,min=3,max=20"`
	Apparence ApparencePhysiqueDTO     `json:"apparence" binding:"required"`
	JobInitial domain.JobID            `json:"job_initial" binding:"required"`
}

// JoueurResponseDTO données retournées par l'API
type JoueurResponseDTO struct {
	ID                  string                   `json:"id"`
	Nom                 string                   `json:"nom"`
	Apparence           ApparencePhysiqueDTO     `json:"apparence"`
	Niveau              int                      `json:"niveau"`
	Experience          int64                    `json:"experience"`
	ExperienceMax       int64                    `json:"experience_max"`
	JobActuel           string                   `json:"job_actuel"`
	JobsDebloquees      []string                 `json:"jobs_debloquees"`
	SlotsCompetences    CompetenceSlotsDTO       `json:"slots_competences"`
	StatsBase           shared.Stats             `json:"stats_base"`
	DateCreation        time.Time                `json:"date_creation"`
	DernierLogin        time.Time                `json:"dernier_login"`
	TempsJeu            string                   `json:"temps_jeu"` // Duration formatée
	ZoneActuelle        string                   `json:"zone_actuelle"`
	SousZone            string                   `json:"sous_zone"`
}

// ApparencePhysiqueDTO version simplifiée pour l'API (MVP)
type ApparencePhysiqueDTO struct {
	Sexe           string            `json:"sexe" binding:"required,oneof=masculin feminin autre"`
	Taille         string            `json:"taille" binding:"required,oneof=petite moyenne grande"`
	CouleurPeau    CouleurDTO        `json:"couleur_peau" binding:"required"`
	CouleurCheveux CouleurDTO        `json:"couleur_cheveux" binding:"required"`
	CouleurYeux    CouleurDTO        `json:"couleur_yeux" binding:"required"`
	StyleCheveux   string            `json:"style_cheveux" binding:"required"`
	StyleBarbe     string            `json:"style_barbe"`
	// Tatouages, Cicatrices, Maquillage retirés pour MVP
}

// CouleurDTO représentation RGB pour l'API
type CouleurDTO struct {
	R uint8 `json:"r" binding:"max=255"`
	G uint8 `json:"g" binding:"max=255"`
	B uint8 `json:"b" binding:"max=255"`
}

// CompetenceSlotsDTO version simplifiée pour l'API
type CompetenceSlotsDTO struct {
	SlotsActifs            []SlotCompetenceDTO `json:"slots_actifs"`
	MaxSlotsActifs         int                 `json:"max_slots_actifs"`
	CompetencesDisponibles []string            `json:"competences_disponibles"`
	// Passif, Réaction, Support retirés pour MVP
}

// SlotCompetenceDTO représente un slot via API
type SlotCompetenceDTO struct {
	Actif        bool   `json:"actif"`
	CompetenceID string `json:"competence_id,omitempty"`
}

// ToJoueur convertit le DTO en entité domain
func (dto *JoueurCreateDTO) ToJoueur(id domain.JoueurID) *domain.Joueur {
	apparence := domain.ApparencePhysique{
		Sexe:           domain.Sexe(dto.Apparence.Sexe),
		Taille:         domain.TaillePersonnage(dto.Apparence.Taille),
		CouleurPeau:    domain.CouleurPersonnage(dto.Apparence.CouleurPeau),
		CouleurCheveux: domain.CouleurPersonnage(dto.Apparence.CouleurCheveux),
		CouleurYeux:    domain.CouleurPersonnage(dto.Apparence.CouleurYeux),
		StyleCheveux:   dto.Apparence.StyleCheveux,
		StyleBarbe:     dto.Apparence.StyleBarbe,
		// Pas de tatouages, cicatrices, maquillage pour MVP
		Tatouages:  []domain.Tatouage{},
		Cicatrices: []domain.Cicatrice{},
		Maquillage: nil,
	}

	return domain.NewJoueur(id, dto.Nom, apparence, dto.JobInitial)
}

// FromJoueur convertit l'entité domain en DTO de réponse
func FromJoueur(joueur *domain.Joueur) *JoueurResponseDTO {
	return &JoueurResponseDTO{
		ID:            string(joueur.ID),
		Nom:           joueur.Nom,
		Apparence:     fromApparence(joueur.Apparence),
		Niveau:        joueur.Niveau,
		Experience:    joueur.Experience,
		ExperienceMax: joueur.ExperienceMax,
		JobActuel:     string(joueur.JobActuel),
		JobsDebloquees: func() []string {
			jobs := make([]string, len(joueur.JobsDebloquees))
			for i, job := range joueur.JobsDebloquees {
				jobs[i] = string(job)
			}
			return jobs
		}(),
		SlotsCompetences: fromCompetenceSlots(joueur.SlotsCompetences),
		StatsBase:        joueur.StatsBase,
		DateCreation:     joueur.DateCreation,
		DernierLogin:     joueur.DernierLogin,
		TempsJeu:         joueur.TempsJeu.String(),
		ZoneActuelle:     joueur.ZoneActuelle,
		SousZone:         joueur.SousZone,
	}
}

func fromApparence(apparence domain.ApparencePhysique) ApparencePhysiqueDTO {
	return ApparencePhysiqueDTO{
		Sexe:           string(apparence.Sexe),
		Taille:         string(apparence.Taille),
		CouleurPeau:    CouleurDTO(apparence.CouleurPeau),
		CouleurCheveux: CouleurDTO(apparence.CouleurCheveux),
		CouleurYeux:    CouleurDTO(apparence.CouleurYeux),
		StyleCheveux:   apparence.StyleCheveux,
		StyleBarbe:     apparence.StyleBarbe,
	}
}

func fromCompetenceSlots(slots domain.CompetenceSlots) CompetenceSlotsDTO {
	slotsActifs := make([]SlotCompetenceDTO, len(slots.SlotsActifs))
	for i, slot := range slots.SlotsActifs {
		slotsActifs[i] = SlotCompetenceDTO{
			Actif:        slot.Actif,
			CompetenceID: string(slot.CompetenceID),
		}
	}

	competencesDisponibles := make([]string, len(slots.CompetencesDisponibles))
	for i, comp := range slots.CompetencesDisponibles {
		competencesDisponibles[i] = string(comp)
	}

	return CompetenceSlotsDTO{
		SlotsActifs:            slotsActifs,
		MaxSlotsActifs:         slots.MaxSlotsActifs,
		CompetencesDisponibles: competencesDisponibles,
	}
}