package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestNewCombat teste la création d'un nouveau combat
func TestNewCombat(t *testing.T) {
	// Arrange
	id := "combat-1"
	joueur1 := "player-1"
	joueur2 := "player-2"
	equipe1, _ := domain.NewEquipe(domain.TeamID("team-1"), "Héros", "#0000FF", false, &joueur1)
	equipe2, _ := domain.NewEquipe(domain.TeamID("team-2"), "Ennemis", "#FF0000", true, &joueur2)
	grille := newTestGrille(10, 10)

	// Act
	combat, err := domain.NewCombat(id, []*domain.Equipe{equipe1, equipe2}, grille)

	// Assert
	assert.NoError(t, err, "La création du combat ne devrait pas retourner d'erreur")
	assert.NotNil(t, combat, "Le combat ne devrait pas être nil")
	assert.Equal(t, id, combat.ID(), "L'ID devrait correspondre")
	assert.Equal(t, domain.EtatAttente, combat.Etat(), "L'état initial devrait être EtatAttente")
	assert.Equal(t, 0, combat.TourActuel(), "Le tour actuel devrait être 0 au départ")
	assert.Equal(t, 2, len(combat.Equipes()), "Le combat devrait avoir 2 équipes")
}
