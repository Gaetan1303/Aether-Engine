package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCompetence_ObtenirPositionsDansZone teste la méthode ObtenirPositionsDansZone()
func TestCompetence_ObtenirPositionsDansZone(t *testing.T) {
	// Arrange
	comp := newTestCompetence("skill-1", "Zone", 1)
	centre := newTestPosition(5, 5)
	grille := newTestGrille(10, 10)

	// Act
	positions := comp.ObtenirPositionsDansZone(centre, grille)

	// Assert
	assert.NotNil(t, positions, "La liste de positions ne devrait pas être nil")
	assert.GreaterOrEqual(t, len(positions), 1, "Devrait avoir au moins 1 position (la cible)")
}
