package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_JoueurID teste la méthode JoueurID()
func TestEquipe_JoueurID(t *testing.T) {
	// Arrange
	expectedJoueurID := "player-xyz-789"
	equipe, _ := domain.NewEquipe(domain.TeamID("team-1"), "Joueurs", "#0000FF", false, &expectedJoueurID)

	// Act
	actualJoueurID := equipe.JoueurID()

	// Assert
	assert.NotNil(t, actualJoueurID, "Le joueurID ne devrait pas être nil")
	assert.Equal(t, expectedJoueurID, *actualJoueurID, "Le joueurID retourné devrait correspondre")
}
