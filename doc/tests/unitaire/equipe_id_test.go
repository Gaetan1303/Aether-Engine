package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_ID teste la méthode ID()
func TestEquipe_ID(t *testing.T) {
	// Arrange
	expectedID := domain.TeamID("team-heroes")
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(expectedID, "Héros", "#0000FF", false, &joueurID)

	// Act
	actualID := equipe.ID()

	// Assert
	assert.Equal(t, expectedID, actualID, "L'ID retourné devrait correspondre à l'ID fourni")
}
