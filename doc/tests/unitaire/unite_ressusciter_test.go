package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_Ressusciter teste la méthode Ressusciter()
func TestUnite_Ressusciter(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Phénix", "team-1", 5, 5)
	unite.RecevoirDegats(200) // Éliminer l'unité

	// Act
	unite.Ressusciter(50)

	// Assert
	assert.False(t, unite.EstEliminee(), "L'unité ne devrait plus être éliminée")
	assert.Equal(t, 50, unite.HPActuels(), "HP devrait être restauré à 50")
}
