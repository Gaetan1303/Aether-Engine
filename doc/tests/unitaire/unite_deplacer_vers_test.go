package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_DeplacerVers teste la méthode DeplacerVers()
func TestUnite_DeplacerVers(t *testing.T) {
	// Arrange
	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	positionInitiale := newTestPosition(0, 0)
	unite := domain.NewUnite(domain.UnitID("u1"), "Scout", domain.TeamID("team-1"), stats, positionInitiale)

	nouvellePosition := newTestPosition(3, 4)

	// Act
	unite.DeplacerVers(nouvellePosition)

	// Assert
	assert.Equal(t, nouvellePosition, unite.Position(), "La position devrait être mise à jour")
	assert.Equal(t, 3, unite.Position().X(), "La coordonnée X devrait être 3")
	assert.Equal(t, 4, unite.Position().Y(), "La coordonnée Y devrait être 4")
}
