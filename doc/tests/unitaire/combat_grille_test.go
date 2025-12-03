package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_Grille teste la méthode Grille()
func TestCombat_Grille(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	grille := combat.Grille()

	// Assert
	assert.NotNil(t, grille, "La grille ne devrait pas être nil")
	assert.Equal(t, 10, grille.Largeur(), "La largeur devrait être 10")
	assert.Equal(t, 10, grille.Hauteur(), "La hauteur devrait être 10")
}
