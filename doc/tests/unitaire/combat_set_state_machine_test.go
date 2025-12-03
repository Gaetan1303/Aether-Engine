package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_SetStateMachine teste la méthode SetStateMachine()
func TestCombat_SetStateMachine(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act - Méthode void avec paramètre nil, ne devrait pas crasher
	combat.SetStateMachine(nil)

	// Assert
	assert.NotNil(t, combat, "Le combat devrait toujours exister")
}
