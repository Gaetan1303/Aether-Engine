package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_ID teste la méthode ID()
func TestCombat_ID(t *testing.T) {
	// Arrange
	expectedID := "combat-123"
	combat := newTestCombat(expectedID)

	// Act
	actualID := combat.ID()

	// Assert
	assert.Equal(t, expectedID, actualID, "L'ID retourné devrait correspondre à l'ID fourni")
}
