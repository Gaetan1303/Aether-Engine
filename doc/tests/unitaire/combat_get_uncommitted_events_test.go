package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_GetUncommittedEvents teste la méthode GetUncommittedEvents()
func TestCombat_GetUncommittedEvents(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	events := combat.GetUncommittedEvents()

	// Assert
	assert.NotNil(t, events, "La liste d'événements ne devrait pas être nil")
	assert.Len(t, events, 0, "Devrait avoir 0 événements initialement")
}
