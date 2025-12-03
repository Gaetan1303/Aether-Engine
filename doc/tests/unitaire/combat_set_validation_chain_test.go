package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_SetValidationChain teste la méthode SetValidationChain()
func TestCombat_SetValidationChain(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act - Méthode void avec paramètre nil, ne devrait pas crasher
	combat.SetValidationChain(nil)

	// Assert
	assert.NotNil(t, combat, "Le combat devrait toujours exister")
}
