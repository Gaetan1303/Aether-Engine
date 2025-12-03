package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_SetHybridDamageMode teste la méthode SetHybridDamageMode()
func TestCombat_SetHybridDamageMode(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	combat.SetHybridDamageMode(0.6, 0.4)

	// Assert
	calculator := combat.GetDamageCalculator()
	assert.Equal(t, "Hybrid", calculator.GetType(), "Le calculateur devrait être de type Hybrid")
}
