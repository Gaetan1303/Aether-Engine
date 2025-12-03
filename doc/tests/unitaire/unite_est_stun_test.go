package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_EstStun teste la méthode EstStun()
func TestUnite_EstStun(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Sentinelle", "team-1", 5, 5)

	// Act
	estStun := unite.EstStun()

	// Assert
	assert.False(t, estStun, "L'unité ne devrait pas être stun initialement")
}
