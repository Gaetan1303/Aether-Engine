package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_SetCommandFactory teste la méthode SetCommandFactory()
func TestCombat_SetCommandFactory(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act - Méthode void avec paramètre nil, ne devrait pas crasher
	combat.SetCommandFactory(nil)

	// Assert
	assert.NotNil(t, combat, "Le combat devrait toujours exister")
}
