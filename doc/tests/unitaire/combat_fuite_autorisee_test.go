package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_FuiteAutorisee teste la méthode FuiteAutorisee()
func TestCombat_FuiteAutorisee(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	fuiteAutorisee := combat.FuiteAutorisee()

	// Assert
	assert.True(t, fuiteAutorisee, "La fuite devrait être autorisée par défaut")
}
