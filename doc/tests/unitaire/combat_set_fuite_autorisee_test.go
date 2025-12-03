package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_SetFuiteAutorisee teste la méthode SetFuiteAutorisee()
func TestCombat_SetFuiteAutorisee(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	combat.SetFuiteAutorisee(false)

	// Assert
	assert.False(t, combat.FuiteAutorisee(), "La fuite devrait être désactivée après SetFuiteAutorisee(false)")
}
