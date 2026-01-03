package interfaces

import (
	"context"

	"github.com/aether-engine/aether-engine/internal/joueur/domain"
)

// JoueurRepository interface pour l'accès aux données des joueurs
// Sera implémentée par l'API BDD externe dans le futur
type JoueurRepository interface {
	// CRUD de base
	Save(ctx context.Context, joueur *domain.Joueur) error
	GetByID(ctx context.Context, id domain.JoueurID) (*domain.Joueur, error)
	GetAll(ctx context.Context) ([]*domain.Joueur, error)
	Update(ctx context.Context, joueur *domain.Joueur) error
	Delete(ctx context.Context, id domain.JoueurID) error

	// Requêtes spécifiques
	GetByNom(ctx context.Context, nom string) (*domain.Joueur, error)
	GetByZone(ctx context.Context, zone string) ([]*domain.Joueur, error)
	ExistsNom(ctx context.Context, nom string) (bool, error)
}

// CompetenceRepository interface pour les compétences
// Gère les compétences disponibles et leurs déblocages
type CompetenceRepository interface {
	GetCompetencesByJob(ctx context.Context, jobID domain.JobID) ([]domain.CompetenceID, error)
	GetCompetenceDetails(ctx context.Context, competenceID domain.CompetenceID) (*CompetenceDTO, error)
	GetCompetencesUnlocked(ctx context.Context, joueurID domain.JoueurID) ([]domain.CompetenceID, error)
}

// CompetenceDTO structure pour transfert de données des compétences
// Bridge entre joueur et combat via API externe
type CompetenceDTO struct {
	ID               domain.CompetenceID `json:"id"`
	Nom              string              `json:"nom"`
	Description      string              `json:"description"`
	Type             string              `json:"type"`
	JobOrigine       domain.JobID        `json:"job_origine"`
	NiveauRequis     int                 `json:"niveau_requis"`
	
	// Données pour le combat (proviennent de combat/domain)
	CoutMP           int                 `json:"cout_mp"`
	CoutAP           int                 `json:"cout_ap"`
	Portee           int                 `json:"portee"`
	Zone             string              `json:"zone"`
	Cooldown         int                 `json:"cooldown"`
}