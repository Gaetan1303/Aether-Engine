package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_ObtenirQuantiteObjet teste la méthode ObtenirQuantiteObjet()
func TestCombat_ObtenirQuantiteObjet(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	quantite := combat.ObtenirQuantiteObjet("potion-1")

	// Assert - Placeholder retourne 1
	assert.Equal(t, 1, quantite, "Le placeholder retourne 1")
}
