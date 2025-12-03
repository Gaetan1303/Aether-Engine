package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_PeutSeDeplacer teste la méthode PeutSeDeplacer()
func TestUnite_PeutSeDeplacer(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Eclaireur", "team-1", 5, 5)

	// Act
	peutSeDeplacer := unite.PeutSeDeplacer()

	// Assert
	assert.True(t, peutSeDeplacer, "Une unité saine devrait pouvoir se déplacer")
}
