package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_DecrementerCooldown teste la méthode DécrémenterCooldown()
func TestCompetence_DecrementerCooldown(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("skill"), "Compétence", domain.CompetenceAttaque)
	comp.ActiverCooldown() // Cooldown à 2

	// Act
	comp.DecrémenterCooldown()

	// Assert
	assert.Equal(t, 1, comp.CooldownActuel(), "Le cooldown actuel devrait être décrémenté de 1 (2 -> 1)")
	assert.True(t, comp.EstEnCooldown(), "La compétence devrait toujours être en cooldown")
}
