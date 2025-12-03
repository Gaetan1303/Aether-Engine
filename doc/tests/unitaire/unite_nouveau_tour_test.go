package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_NouveauTour teste la méthode NouveauTour()
func TestUnite_NouveauTour(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Berserker", "team-1", 5, 5)

	// Act
	unite.NouveauTour()

	// Assert
	assert.True(t, unite.PeutAgir(), "Devrait pouvoir agir après un nouveau tour")
	assert.True(t, unite.PeutSeDeplacer(), "Devrait pouvoir se déplacer après un nouveau tour")
}
