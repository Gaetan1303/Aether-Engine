package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_TourActuel teste la méthode TourActuel()
func TestCombat_TourActuel(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	tourActuel := combat.TourActuel()

	// Assert
	assert.Equal(t, 0, tourActuel, "Le tour actuel devrait être 0 au départ")
}
