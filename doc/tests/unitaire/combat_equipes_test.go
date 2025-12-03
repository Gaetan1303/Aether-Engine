package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_Equipes teste la méthode Equipes()
func TestCombat_Equipes(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	equipes := combat.Equipes()

	// Assert
	assert.NotNil(t, equipes, "Les équipes ne devraient pas être nil")
	assert.Equal(t, 2, len(equipes), "Le combat devrait avoir 2 équipes")
	assert.NotNil(t, equipes[domain.TeamID("team-1")], "L'équipe 1 devrait exister")
	assert.NotNil(t, equipes[domain.TeamID("team-2")], "L'équipe 2 devrait exister")
}
