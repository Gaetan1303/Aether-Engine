package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_EstSilence teste la méthode EstSilence()
func TestUnite_EstSilence(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Barde", "team-1", 5, 5)

	// Act
	estSilence := unite.EstSilence()

	// Assert
	assert.False(t, estSilence, "L'unité ne devrait pas être silence initialement")
}
