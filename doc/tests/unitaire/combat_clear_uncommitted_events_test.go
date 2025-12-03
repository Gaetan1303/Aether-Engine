package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_ClearUncommittedEvents teste la méthode ClearUncommittedEvents()
func TestCombat_ClearUncommittedEvents(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	event := domain.NewCombatDemarreEvent("combat-1", 1, []domain.UnitID{})
	combat.RaiseEvent(event)

	// Act
	combat.ClearUncommittedEvents()

	// Assert
	events := combat.GetUncommittedEvents()
	assert.Len(t, events, 0, "Les événements devraient être effacés")
}
