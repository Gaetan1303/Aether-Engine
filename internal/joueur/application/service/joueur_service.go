package application

import (
	"context"
	"fmt"
	"time"

	"github.com/aether-engine/aether-engine/internal/joueur/application/dto"
	"github.com/aether-engine/aether-engine/internal/joueur/domain"
	"github.com/aether-engine/aether-engine/internal/joueur/interfaces"
)

// JoueurService couche application pour orchestrer les opérations joueur
// Prêt pour l'intégration future avec l'API BDD externe
type JoueurService struct {
	repository interfaces.JoueurRepository
	jobManager *domain.JobManager
}

// NewJoueurService crée un nouveau service joueur
func NewJoueurService(repository interfaces.JoueurRepository) *JoueurService {
	return &JoueurService{
		repository: repository,
		jobManager: domain.NewJobManager(),
	}
}

// CreateJoueur crée un nouveau joueur avec validation métier
func (s *JoueurService) CreateJoueur(ctx context.Context, createDTO *dto.JoueurCreateDTO) (*dto.JoueurResponseDTO, error) {
	// Validation métier
	if err := s.validateCreation(ctx, createDTO); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Générer ID unique (temporairement timestamp, plus tard UUID via API BDD)
	joueurID := domain.JoueurID(fmt.Sprintf("joueur_%d", time.Now().UnixNano()))

	// Convertir DTO vers entité domain
	joueur := createDTO.ToJoueur(joueurID)

	// Initialiser les compétences de base selon le job
	if err := s.initializerCompetencesBase(joueur); err != nil {
		return nil, fmt.Errorf("failed to initialize base competences: %w", err)
	}

	// Sauvegarder (via API BDD future)
	if err := s.repository.Save(ctx, joueur); err != nil {
		return nil, fmt.Errorf("failed to save joueur: %w", err)
	}

	// Retourner DTO de réponse
	return dto.FromJoueur(joueur), nil
}

// GetJoueur récupère un joueur par ID
func (s *JoueurService) GetJoueur(ctx context.Context, id string) (*dto.JoueurResponseDTO, error) {
	joueur, err := s.repository.GetByID(ctx, domain.JoueurID(id))
	if err != nil {
		return nil, fmt.Errorf("joueur not found: %w", err)
	}

	return dto.FromJoueur(joueur), nil
}

// UpdateJoueur met à jour un joueur existant
func (s *JoueurService) UpdateJoueur(ctx context.Context, id string, updateData interface{}) (*dto.JoueurResponseDTO, error) {
	// TODO: Implémenter la logique de mise à jour
	// Sera implémenté avec l'API BDD externe
	return nil, fmt.Errorf("update not implemented yet - waiting for external BDD API")
}

// DeleteJoueur supprime un joueur
func (s *JoueurService) DeleteJoueur(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, domain.JoueurID(id))
}

// GetAllJoueurs récupère tous les joueurs
func (s *JoueurService) GetAllJoueurs(ctx context.Context) ([]*dto.JoueurResponseDTO, error) {
	joueurs, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responseDTOs := make([]*dto.JoueurResponseDTO, len(joueurs))
	for i, joueur := range joueurs {
		responseDTOs[i] = dto.FromJoueur(joueur)
	}

	return responseDTOs, nil
}

// AssignerCompetence assigne une compétence à un slot du joueur
func (s *JoueurService) AssignerCompetence(ctx context.Context, joueurID string, typeComp domain.TypeCompetence, slotIndex int, competenceID domain.CompetenceID) error {
	joueur, err := s.repository.GetByID(ctx, domain.JoueurID(joueurID))
	if err != nil {
		return fmt.Errorf("joueur not found: %w", err)
	}

	if err := joueur.SlotsCompetences.AssignerCompetence(typeComp, slotIndex, competenceID); err != nil {
		return fmt.Errorf("failed to assign competence: %w", err)
	}

	return s.repository.Update(ctx, joueur)
}

// ChangerJob change le job du joueur
func (s *JoueurService) ChangerJob(ctx context.Context, joueurID string, nouveauJob domain.JobID) error {
	joueur, err := s.repository.GetByID(ctx, domain.JoueurID(joueurID))
	if err != nil {
		return fmt.Errorf("joueur not found: %w", err)
	}

	if err := joueur.ChangerJob(nouveauJob); err != nil {
		return fmt.Errorf("failed to change job: %w", err)
	}

	return s.repository.Update(ctx, joueur)
}

// Validation métier pour la création
func (s *JoueurService) validateCreation(ctx context.Context, createDTO *dto.JoueurCreateDTO) error {
	// Vérifier que le nom n'existe pas déjà
	exists, err := s.repository.ExistsNom(ctx, createDTO.Nom)
	if err != nil {
		return fmt.Errorf("failed to check name uniqueness: %w", err)
	}
	if exists {
		return fmt.Errorf("nom '%s' already exists", createDTO.Nom)
	}

	// Vérifier que le job initial est valide
	if !s.jobManager.IsJobDepart(createDTO.JobInitial) {
		return fmt.Errorf("job initial '%s' is not a valid starter job", createDTO.JobInitial)
	}

	return nil
}

// Initialise les compétences de base selon le job
func (s *JoueurService) initializerCompetencesBase(joueur *domain.Joueur) error {
	// TODO: Récupérer les compétences de base du job via CompetenceRepository
	// Pour l'instant, ajouter quelques compétences fixes selon le job
	
	var competencesBase []domain.CompetenceID

	switch joueur.JobActuel {
	case domain.JobGuerrier:
		competencesBase = []domain.CompetenceID{
			"attaque_basique",
			"coup_puissant",
		}
	case domain.JobMage:
		competencesBase = []domain.CompetenceID{
			"boule_de_feu",
			"soin_mineur",
		}
	case domain.JobArcher:
		competencesBase = []domain.CompetenceID{
			"tir_precis",
			"tir_perforant",
		}
	case domain.JobVoleur:
		competencesBase = []domain.CompetenceID{
			"attaque_sournoise",
			"evasion",
		}
	case domain.JobClerc:
		competencesBase = []domain.CompetenceID{
			"soin_majeur",
			"benediction",
		}
	}

	// Ajouter les compétences disponibles
	for _, competence := range competencesBase {
		joueur.SlotsCompetences.AjouterCompetenceDisponible(competence)
	}

	// Auto-équiper la première compétence dans le slot actif
	if len(competencesBase) > 0 {
		return joueur.SlotsCompetences.AssignerCompetence(domain.TypeActif, 0, competencesBase[0])
	}

	return nil
}