package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_SeDeplacer teste la méthode SeDeplacer()
func TestUnite_SeDeplacer(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Voleur", "team-1", 5, 5)
	nouvellePosition := newTestPosition(6, 6)

	// Act
	err := unite.SeDeplacer(nouvellePosition, 1)

	// Assert
	assert.NoError(t, err, "Le déplacement devrait réussir")
	assert.Equal(t, 6, unite.Position().X(), "Position X devrait être mise à jour")
	assert.Equal(t, 6, unite.Position().Y(), "Position Y devrait être mise à jour")
}
