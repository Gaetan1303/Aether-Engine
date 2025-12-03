package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_Soigner teste la méthode Soigner() (alias de RecevoirSoin)
func TestUnite_Soigner(t *testing.T) {
	// Arrange
	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(domain.UnitID("u1"), "Prêtre", domain.TeamID("team-1"), stats, newTestPosition(0, 0))

	// Infliger des dégâts d'abord
	unite.RecevoirDegats(50)

	// Act
	soin := 25
	unite.Soigner(soin)

	// Assert
	assert.Equal(t, 75, unite.HPActuels(), "Les HP devraient être restaurés de 25 (50 + 25 = 75)")
}
