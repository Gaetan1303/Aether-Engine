package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_GetDamageCalculator teste la méthode GetDamageCalculator()
func TestCombat_GetDamageCalculator(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	calculator := combat.GetDamageCalculator()

	// Assert
	assert.NotNil(t, calculator, "Le calculateur de dégâts ne devrait pas être nil")
	assert.Equal(t, "Physical", calculator.GetType(), "Le type par défaut devrait être physique")
}
