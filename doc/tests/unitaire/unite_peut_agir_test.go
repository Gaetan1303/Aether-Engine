package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_PeutAgir teste la méthode PeutAgir()
func TestUnite_PeutAgir(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Chevalier", "team-1", 5, 5)

	// Act
	peutAgir := unite.PeutAgir()

	// Assert
	assert.True(t, peutAgir, "Une unité saine devrait pouvoir agir")
}
