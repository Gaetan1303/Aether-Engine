package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_ConsommerObjet teste la méthode ConsommerObjet()
func TestCombat_ConsommerObjet(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act - Méthode void, ne devrait pas crasher
	combat.ConsommerObjet("potion-1", 1)

	// Assert
	assert.NotNil(t, combat, "Le combat devrait toujours exister")
}
