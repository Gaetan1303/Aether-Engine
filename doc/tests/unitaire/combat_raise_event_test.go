package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_RaiseEvent teste la méthode RaiseEvent()
func TestCombat_RaiseEvent(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	event := domain.NewCombatDemarreEvent("combat-1", 1, []domain.UnitID{})

	// Act - Méthode void, ne devrait pas crasher
	combat.RaiseEvent(event)

	// Assert
	events := combat.GetUncommittedEvents()
	assert.Len(t, events, 1, "Devrait avoir 1 événement non commité")
}
