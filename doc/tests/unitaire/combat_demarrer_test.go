package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_Demarrer teste la méthode Demarrer()
func TestCombat_Demarrer(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	err := combat.Demarrer()

	// Assert
	assert.NoError(t, err, "Le démarrage ne devrait pas retourner d'erreur")
	assert.Equal(t, domain.EtatEnCours, combat.Etat(), "L'état devrait être EtatEnCours après démarrage")
	assert.Equal(t, 1, combat.TourActuel(), "Le tour actuel devrait être 1 après démarrage")
}
