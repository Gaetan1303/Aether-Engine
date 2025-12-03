package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_EstIA teste la méthode EstIA()
func TestUnite_EstIA(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Bot", "team-1", 5, 5)

	// Act
	estIA := unite.EstIA()

	// Assert
	assert.False(t, estIA, "Par défaut, l'unité ne devrait pas être contrôlée par l'IA (placeholder)")
}
