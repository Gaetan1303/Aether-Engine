package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_GetObserverSubject teste la méthode GetObserverSubject()
func TestCombat_GetObserverSubject(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	subject := combat.GetObserverSubject()

	// Assert - Initialement nil
	assert.Nil(t, subject, "L'observer subject devrait être nil initialement")
}
