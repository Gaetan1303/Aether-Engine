package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_GetValidationChain teste la méthode GetValidationChain()
func TestCombat_GetValidationChain(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	chain := combat.GetValidationChain()

	// Assert - Initialement nil
	assert.Nil(t, chain, "La validation chain devrait être nil initialement")
}
