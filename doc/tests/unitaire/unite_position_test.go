package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_Position teste la méthode Position()
func TestUnite_Position(t *testing.T) {
	// Arrange
	expectedPosition := newTestPosition(5, 10)
	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(domain.UnitID("u1"), "Archer", domain.TeamID("team-1"), stats, expectedPosition)

	// Act
	actualPosition := unite.Position()

	// Assert
	assert.Equal(t, expectedPosition, actualPosition, "La position retournée devrait correspondre à la position fournie")
	assert.Equal(t, 5, actualPosition.X(), "La coordonnée X devrait être 5")
	assert.Equal(t, 10, actualPosition.Y(), "La coordonnée Y devrait être 10")
}
