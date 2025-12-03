package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_Couleur teste la méthode Couleur()
func TestEquipe_Couleur(t *testing.T) {
	// Arrange
	expectedCouleur := "#00FF00"
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(domain.TeamID("team-1"), "Verts", expectedCouleur, false, &joueurID)

	// Act
	actualCouleur := equipe.Couleur()

	// Assert
	assert.Equal(t, expectedCouleur, actualCouleur, "La couleur retournée devrait correspondre à la couleur fournie")
}
