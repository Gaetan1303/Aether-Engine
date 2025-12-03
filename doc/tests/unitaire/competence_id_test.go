package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_ID teste la méthode ID()
func TestCompetence_ID(t *testing.T) {
	// Arrange
	expectedID := domain.CompetenceID("lightning-bolt")
	comp := newTestCompetence(expectedID, "Éclair", domain.CompetenceMagie)

	// Act
	actualID := comp.ID()

	// Assert
	assert.Equal(t, expectedID, actualID, "L'ID retourné devrait correspondre à l'ID fourni")
}
