package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_EstBloqueDeplacement teste la méthode EstBloqueDeplacement()
func TestUnite_EstBloqueDeplacement(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Paladin", "team-1", 5, 5)

	// Act
	estBloque := unite.EstBloqueDeplacement()

	// Assert
	assert.False(t, estBloque, "Une unité saine ne devrait pas être bloquée")
}
