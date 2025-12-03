package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_SetCooldownActuel teste la méthode SetCooldownActuel()
func TestCompetence_SetCooldownActuel(t *testing.T) {
	// Arrange
	comp := newTestCompetence("skill-1", "Attaque", domain.CompetenceAttaque)

	// Act
	comp.SetCooldownActuel(5)

	// Assert
	assert.Equal(t, 5, comp.CooldownActuel(), "Le cooldown actuel devrait être défini à 5")
}
