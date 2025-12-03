package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_GetTimestamp teste la méthode GetTimestamp()
func TestCombat_GetTimestamp(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	timestamp := combat.GetTimestamp()

	// Assert
	assert.Greater(t, timestamp, int64(0), "Le timestamp devrait être supérieur à 0")
}
