package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_GetCommandFactory teste la méthode GetCommandFactory()
func TestCombat_GetCommandFactory(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	factory := combat.GetCommandFactory()

	// Assert - Initialement nil
	assert.Nil(t, factory, "La command factory devrait être nil initialement")
}
