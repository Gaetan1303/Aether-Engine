package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_SetDamageCalculator teste la méthode SetDamageCalculator()
func TestCombat_SetDamageCalculator(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	customCalculator := domain.NewMagicalDamageCalculator()

	// Act
	combat.SetDamageCalculator(customCalculator)

	// Assert
	calculator := combat.GetDamageCalculator()
	assert.Equal(t, "Magical", calculator.GetType(), "Le calculateur devrait être de type Magical")
}
