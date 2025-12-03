package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_DegatsBase teste la méthode DegatsBase()
func TestCompetence_DegatsBase(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("attack"), "Attaque", domain.CompetenceAttaque)

	// Act
	degatsBase := comp.DegatsBase()

	// Assert
	assert.Equal(t, 20, degatsBase, "Les dégâts de base devraient être 20")
}
