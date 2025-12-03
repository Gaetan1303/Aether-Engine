package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_MarquerEquipeFuite teste la méthode MarquerEquipeFuite()
func TestCombat_MarquerEquipeFuite(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	teamID := domain.TeamID("team-1")

	// Act
	combat.MarquerEquipeFuite(teamID)

	// Assert
	resultat := combat.VerifierConditionsVictoire()
	assert.Equal(t, "FLED", resultat, "Le résultat devrait être FLED après marquage de fuite")
}
