package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_EstRoot teste la méthode EstRoot()
func TestUnite_EstRoot(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Éclaireur", "team-1", 5, 5)

	// Act
	estRoot := unite.EstRoot()

	// Assert
	assert.False(t, estRoot, "L'unité ne devrait pas être root initialement")
}
