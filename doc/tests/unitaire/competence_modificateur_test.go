package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_Modificateur teste la méthode Modificateur()
func TestCompetence_Modificateur(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("skill"), "Compétence", domain.CompetenceAttaque)

	// Act
	modificateur := comp.Modificateur()

	// Assert
	assert.Equal(t, 0.5, modificateur, "Le modificateur devrait être 0.5")
}
