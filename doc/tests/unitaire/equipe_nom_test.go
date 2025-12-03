package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_Nom teste la méthode Nom()
func TestEquipe_Nom(t *testing.T) {
	// Arrange
	expectedNom := "Chevaliers de la Table Ronde"
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(domain.TeamID("team-1"), expectedNom, "#FFD700", false, &joueurID)

	// Act
	actualNom := equipe.Nom()

	// Assert
	assert.Equal(t, expectedNom, actualNom, "Le nom retourné devrait correspondre au nom fourni")
}
