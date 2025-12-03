package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_Cooldown teste la méthode Cooldown()
func TestCompetence_Cooldown(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("ultimate"), "Ultime", domain.CompetenceAttaque)

	// Act
	cooldown := comp.Cooldown()

	// Assert
	assert.Equal(t, 2, cooldown, "Le cooldown devrait être 2")
}
