package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_RegenererStatut teste la méthode RegenererStatut()
func TestUnite_RegenererStatut(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Moine", "team-1", 5, 5)
	mpInitial := unite.StatsActuelles().MP

	// Act
	unite.RegenererStatut()

	// Assert
	mpApres := unite.StatsActuelles().MP
	assert.GreaterOrEqual(t, mpApres, mpInitial, "Le MP devrait rester stable ou augmenter après régénération")
}
