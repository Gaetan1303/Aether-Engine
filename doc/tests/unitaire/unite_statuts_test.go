package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_Statuts teste la méthode Statuts()
func TestUnite_Statuts(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Archer", "team-1", 5, 5)

	// Act
	statuts := unite.Statuts()

	// Assert
	assert.NotNil(t, statuts, "La liste de statuts ne devrait pas être nil")
	assert.Len(t, statuts, 0, "Devrait avoir 0 statuts initialement")
}
