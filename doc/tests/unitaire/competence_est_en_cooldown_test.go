package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_EstEnCooldown teste la méthode EstEnCooldown()
func TestCompetence_EstEnCooldown(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("skill"), "Compétence", domain.CompetenceAttaque)

	// Act
	estEnCooldown := comp.EstEnCooldown()

	// Assert
	assert.False(t, estEnCooldown, "La compétence ne devrait pas être en cooldown au départ")
}
