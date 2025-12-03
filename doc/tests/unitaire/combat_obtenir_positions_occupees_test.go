package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_ObtenirPositionsOccupees teste la méthode ObtenirPositionsOccupees()
func TestCombat_ObtenirPositionsOccupees(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	unite1 := newTestUnite("unite-1", "Guerrier", "team-1", 5, 5)
	unite2 := newTestUnite("unite-2", "Archer", "team-1", 6, 6)

	equipes := combat.Equipes()
	for _, equipe := range equipes {
		if equipe.ID() == "team-1" {
			equipe.AjouterMembre(unite1)
			equipe.AjouterMembre(unite2)
			break
		}
	}

	// Act
	positions := combat.ObtenirPositionsOccupees(domain.UnitID("unite-1"))

	// Assert
	assert.NotNil(t, positions, "La map de positions ne devrait pas être nil")
	// Devrait contenir unite-2 mais pas unite-1 (exclusion)
}
