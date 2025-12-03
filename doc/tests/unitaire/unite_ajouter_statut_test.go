package unitaire

import (
	"testing"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_AjouterStatut teste la méthode AjouterStatut()
func TestUnite_AjouterStatut(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Ninja", "team-1", 5, 5)
	statut := &shared.Statut{}

	// Act
	err := unite.AjouterStatut(statut)

	// Assert
	assert.NoError(t, err, "L'ajout de statut ne devrait pas échouer")
}
