package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_AnnulerFuite teste la méthode AnnulerFuite()
func TestCombat_AnnulerFuite(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	combat.MarquerEquipeFuite(domain.TeamID("team-1"))

	// Act
	combat.AnnulerFuite(domain.TeamID("team-1"))

	// Assert
	resultat := combat.VerifierConditionsVictoire()
	assert.NotEqual(t, "FLED", resultat, "Ne devrait plus être en fuite après annulation")
}
