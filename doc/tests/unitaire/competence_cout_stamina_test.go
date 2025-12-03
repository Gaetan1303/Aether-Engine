package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_CoutStamina teste la méthode CoutStamina()
func TestCompetence_CoutStamina(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("skill"), "Compétence", domain.CompetenceAttaque)

	// Act
	coutStamina := comp.CoutStamina()

	// Assert
	assert.Equal(t, 5, coutStamina, "Le coût Stamina devrait être 5")
}
