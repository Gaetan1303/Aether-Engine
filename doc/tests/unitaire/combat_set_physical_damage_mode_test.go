package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_SetPhysicalDamageMode teste la méthode SetPhysicalDamageMode()
func TestCombat_SetPhysicalDamageMode(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	combat.SetMagicalDamageMode() // Changer d'abord en mode magique

	// Act
	combat.SetPhysicalDamageMode()

	// Assert
	calculator := combat.GetDamageCalculator()
	assert.Equal(t, "Physical", calculator.GetType(), "Le type devrait être physique après SetPhysicalDamageMode")
}
