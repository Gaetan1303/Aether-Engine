package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_RecevoirSoin teste la méthode RecevoirSoin()
func TestUnite_RecevoirSoin(t *testing.T) {
	// Arrange
	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(domain.UnitID("u1"), "Guerrier", domain.TeamID("team-1"), stats, newTestPosition(0, 0))

	// Infliger des dégâts d'abord
	unite.RecevoirDegats(40)

	// Act
	soin := 20
	unite.RecevoirSoin(soin)

	// Assert
	assert.Equal(t, 80, unite.HPActuels(), "Les HP devraient être restaurés de 20 (60 + 20 = 80)")
}
