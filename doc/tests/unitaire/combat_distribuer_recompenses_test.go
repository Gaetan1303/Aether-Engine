package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_DistribuerRecompenses teste la méthode DistribuerRecompenses()
func TestCombat_DistribuerRecompenses(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act - Méthode void, ne devrait pas crasher
	combat.DistribuerRecompenses()

	// Assert
	assert.NotNil(t, combat, "Le combat devrait toujours exister après distribution")
}
