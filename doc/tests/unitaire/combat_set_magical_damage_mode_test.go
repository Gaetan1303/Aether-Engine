package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_SetMagicalDamageMode teste la méthode SetMagicalDamageMode()
func TestCombat_SetMagicalDamageMode(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	combat.SetMagicalDamageMode()

	// Assert
	calculator := combat.GetDamageCalculator()
	assert.Equal(t, "Magical", calculator.GetType(), "Le type devrait être magique après SetMagicalDamageMode")
}
