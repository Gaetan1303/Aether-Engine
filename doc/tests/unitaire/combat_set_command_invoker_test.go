package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_SetCommandInvoker teste la méthode SetCommandInvoker()
func TestCombat_SetCommandInvoker(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act - Méthode void avec paramètre nil, ne devrait pas crasher
	combat.SetCommandInvoker(nil)

	// Assert
	assert.NotNil(t, combat, "Le combat devrait toujours exister")
}
