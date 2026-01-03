package infrastructure

import (
	"context"
	"fmt"
	"sync"

	"github.com/aether-engine/aether-engine/internal/joueur/domain"
	"github.com/aether-engine/aether-engine/internal/joueur/interfaces"
)

// InMemoryJoueurRepository implémentation temporaire en mémoire
// Sera remplacée par l'implémentation API BDD externe
type InMemoryJoueurRepository struct {
	joueurs map[domain.JoueurID]*domain.Joueur
	mu      sync.RWMutex
}

// NewInMemoryJoueurRepository crée un nouveau repository en mémoire
func NewInMemoryJoueurRepository() interfaces.JoueurRepository {
	return &InMemoryJoueurRepository{
		joueurs: make(map[domain.JoueurID]*domain.Joueur),
	}
}

// Save sauvegarde un joueur en mémoire
func (r *InMemoryJoueurRepository) Save(ctx context.Context, joueur *domain.Joueur) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Clone pour éviter les modifications externes
	clonedJoueur := *joueur
	r.joueurs[joueur.ID] = &clonedJoueur

	return nil
}

// GetByID récupère un joueur par ID
func (r *InMemoryJoueurRepository) GetByID(ctx context.Context, id domain.JoueurID) (*domain.Joueur, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	joueur, exists := r.joueurs[id]
	if !exists {
		return nil, fmt.Errorf("joueur with ID %s not found", id)
	}

	// Clone pour éviter les modifications externes
	clonedJoueur := *joueur
	return &clonedJoueur, nil
}

// GetAll récupère tous les joueurs
func (r *InMemoryJoueurRepository) GetAll(ctx context.Context) ([]*domain.Joueur, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	joueurs := make([]*domain.Joueur, 0, len(r.joueurs))
	for _, joueur := range r.joueurs {
		// Clone pour éviter les modifications externes
		clonedJoueur := *joueur
		joueurs = append(joueurs, &clonedJoueur)
	}

	return joueurs, nil
}

// Update met à jour un joueur existant
func (r *InMemoryJoueurRepository) Update(ctx context.Context, joueur *domain.Joueur) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.joueurs[joueur.ID]; !exists {
		return fmt.Errorf("joueur with ID %s not found for update", joueur.ID)
	}

	// Clone pour éviter les modifications externes
	clonedJoueur := *joueur
	r.joueurs[joueur.ID] = &clonedJoueur

	return nil
}

// Delete supprime un joueur
func (r *InMemoryJoueurRepository) Delete(ctx context.Context, id domain.JoueurID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.joueurs[id]; !exists {
		return fmt.Errorf("joueur with ID %s not found for deletion", id)
	}

	delete(r.joueurs, id)
	return nil
}

// GetByNom récupère un joueur par nom
func (r *InMemoryJoueurRepository) GetByNom(ctx context.Context, nom string) (*domain.Joueur, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, joueur := range r.joueurs {
		if joueur.Nom == nom {
			// Clone pour éviter les modifications externes
			clonedJoueur := *joueur
			return &clonedJoueur, nil
		}
	}

	return nil, fmt.Errorf("joueur with nom '%s' not found", nom)
}

// GetByZone récupère les joueurs dans une zone
func (r *InMemoryJoueurRepository) GetByZone(ctx context.Context, zone string) ([]*domain.Joueur, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var joueurs []*domain.Joueur
	for _, joueur := range r.joueurs {
		if joueur.ZoneActuelle == zone {
			// Clone pour éviter les modifications externes
			clonedJoueur := *joueur
			joueurs = append(joueurs, &clonedJoueur)
		}
	}

	return joueurs, nil
}

// ExistsNom vérifie si un nom existe déjà
func (r *InMemoryJoueurRepository) ExistsNom(ctx context.Context, nom string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, joueur := range r.joueurs {
		if joueur.Nom == nom {
			return true, nil
		}
	}

	return false, nil
}