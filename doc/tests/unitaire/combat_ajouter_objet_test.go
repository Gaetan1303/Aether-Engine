package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_AjouterObjet teste la méthode AjouterObjet()
func TestCombat_AjouterObjet(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act - Méthode void, ne devrait pas crasher
	combat.AjouterObjet("potion-2", 5)

	// Assert
	assert.NotNil(t, combat, "Le combat devrait toujours exister")
}
