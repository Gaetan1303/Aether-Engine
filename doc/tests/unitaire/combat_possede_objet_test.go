package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_PossedeObjet teste la méthode PossedeObjet()
func TestCombat_PossedeObjet(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	possede := combat.PossedeObjet("potion-1")

	// Assert - Placeholder retourne toujours true
	assert.True(t, possede, "Le placeholder retourne toujours true")
}
