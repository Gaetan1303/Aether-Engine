package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_EstComplete teste la méthode EstComplete()
func TestEquipe_EstComplete(t *testing.T) {
	// Arrange
	teamID := domain.TeamID("team-1")
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(teamID, "Héros", "#0000FF", false, &joueurID)

	// Act - équipe vide
	estCompleteVide := equipe.EstComplete()

	// Assert
	assert.False(t, estCompleteVide, "Une équipe vide ne devrait pas être complète")
}
