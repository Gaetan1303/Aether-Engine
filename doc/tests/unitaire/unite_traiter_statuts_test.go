package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_TraiterStatuts teste la méthode TraiterStatuts()
func TestUnite_TraiterStatuts(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Clerc", "team-1", 5, 5)

	// Act
	effets := unite.TraiterStatuts()

	// Assert
	assert.NotNil(t, effets, "La liste d'effets ne devrait pas être nil")
	assert.Len(t, effets, 0, "Devrait avoir 0 effets sans statuts actifs")
}
