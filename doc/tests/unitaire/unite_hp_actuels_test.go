package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_HPActuels teste la méthode HPActuels()
func TestUnite_HPActuels(t *testing.T) {
	// Arrange
	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(domain.UnitID("u1"), "Guerrier", domain.TeamID("team-1"), stats, newTestPosition(0, 0))

	// Act
	hpActuels := unite.HPActuels()

	// Assert
	assert.Equal(t, 100, hpActuels, "Les HP actuels devraient être égaux aux HP de base au début")
}
