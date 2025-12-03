package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_Apply teste la méthode Apply()
func TestCombat_Apply(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	event := domain.NewCombatDemarreEvent("combat-1", 1, []domain.UnitID{})

	// Act
	err := combat.Apply(event)

	// Assert
	assert.NoError(t, err, "L'application de l'événement ne devrait pas échouer")
}
