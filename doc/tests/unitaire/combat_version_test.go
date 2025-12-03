package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_Version teste la méthode Version()
func TestCombat_Version(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	version := combat.Version()

	// Assert
	assert.Equal(t, 0, version, "La version devrait être 0 au départ")
}
