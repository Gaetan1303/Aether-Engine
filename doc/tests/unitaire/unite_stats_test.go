package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_Stats teste la méthode Stats()
func TestUnite_Stats(t *testing.T) {
	// Arrange
	expectedStats := &shared.Stats{
		HP:      150,
		MP:      80,
		Stamina: 100,
		ATK:     40,
		DEF:     25,
		MATK:    20,
		MDEF:    18,
		SPD:     15,
		MOV:     6,
	}
	unite := domain.NewUnite(domain.UnitID("u1"), "Tank", domain.TeamID("team-1"), expectedStats, newTestPosition(0, 0))

	// Act
	actualStats := unite.Stats()

	// Assert
	assert.NotNil(t, actualStats, "Les stats ne devraient pas être nil")
	assert.Equal(t, expectedStats.HP, actualStats.HP, "HP devrait correspondre")
	assert.Equal(t, expectedStats.MP, actualStats.MP, "MP devrait correspondre")
	assert.Equal(t, expectedStats.Stamina, actualStats.Stamina, "Stamina devrait correspondre")
	assert.Equal(t, expectedStats.ATK, actualStats.ATK, "ATK devrait correspondre")
	assert.Equal(t, expectedStats.DEF, actualStats.DEF, "DEF devrait correspondre")
	assert.Equal(t, expectedStats.MATK, actualStats.MATK, "MATK devrait correspondre")
	assert.Equal(t, expectedStats.MDEF, actualStats.MDEF, "MDEF devrait correspondre")
	assert.Equal(t, expectedStats.SPD, actualStats.SPD, "SPD devrait correspondre")
	assert.Equal(t, expectedStats.MOV, actualStats.MOV, "MOV devrait correspondre")
}
