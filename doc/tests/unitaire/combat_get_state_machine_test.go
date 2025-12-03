package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_GetStateMachine teste la méthode GetStateMachine()
func TestCombat_GetStateMachine(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	sm := combat.GetStateMachine()

	// Assert - Initialement nil
	assert.Nil(t, sm, "La state machine devrait être nil initialement")
}
