package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_EstEliminee teste la méthode EstEliminee()
func TestUnite_EstEliminee(t *testing.T) {
	// Arrange
	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(domain.UnitID("u1"), "Guerrier", domain.TeamID("team-1"), stats, newTestPosition(0, 0))

	// Act
	estEliminee := unite.EstEliminee()

	// Assert
	assert.False(t, estEliminee, "L'unité ne devrait pas être éliminée au début")
}
