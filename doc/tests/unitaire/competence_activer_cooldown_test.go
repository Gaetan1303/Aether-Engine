package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_ActiverCooldown teste la méthode ActiverCooldown()
func TestCompetence_ActiverCooldown(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("skill"), "Compétence", domain.CompetenceAttaque)

	// Act
	comp.ActiverCooldown()

	// Assert
	assert.Equal(t, 2, comp.CooldownActuel(), "Le cooldown actuel devrait être égal au cooldown max (2)")
	assert.True(t, comp.EstEnCooldown(), "La compétence devrait être en cooldown après activation")
}
