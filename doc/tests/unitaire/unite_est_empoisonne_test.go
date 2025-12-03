package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_EstEmpoisonne teste la méthode EstEmpoisonne()
func TestUnite_EstEmpoisonne(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Rôdeur", "team-1", 5, 5)

	// Act
	estEmpoisonne := unite.EstEmpoisonne()

	// Assert
	assert.False(t, estEmpoisonne, "L'unité ne devrait pas être empoisonnée initialement")
}
