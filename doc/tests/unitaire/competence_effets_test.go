package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCompetence_Effets teste la méthode Effets()
func TestCompetence_Effets(t *testing.T) {
	// Arrange
	comp := newTestCompetence("skill-1", "Boule de Feu", 1) // CompetenceMagie = 1

	// Act
	effets := comp.Effets()

	// Assert
	assert.NotNil(t, effets, "La liste d'effets ne devrait pas être nil")
	assert.Len(t, effets, 0, "Devrait avoir 0 effets initialement")
}
