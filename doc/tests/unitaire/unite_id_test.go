package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_ID teste la méthode ID()
func TestUnite_ID(t *testing.T) {
	// Arrange
	expectedID := domain.UnitID("warrior-123")
	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(expectedID, "Guerrier", domain.TeamID("team-1"), stats, newTestPosition(0, 0))

	// Act
	actualID := unite.ID()

	// Assert
	assert.Equal(t, expectedID, actualID, "L'ID retourné devrait correspondre à l'ID fourni")
}
