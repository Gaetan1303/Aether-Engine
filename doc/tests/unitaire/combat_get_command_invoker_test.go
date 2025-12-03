package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_GetCommandInvoker teste la méthode GetCommandInvoker()
func TestCombat_GetCommandInvoker(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	invoker := combat.GetCommandInvoker()

	// Assert - Initialement nil
	assert.Nil(t, invoker, "Le command invoker devrait être nil initialement")
}
