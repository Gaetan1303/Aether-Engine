package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_CooldownActuel teste la méthode CooldownActuel()
func TestCompetence_CooldownActuel(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("skill"), "Compétence", domain.CompetenceAttaque)

	// Act
	cooldownActuel := comp.CooldownActuel()

	// Assert
	assert.Equal(t, 0, cooldownActuel, "Le cooldown actuel devrait être 0 au départ")
}
