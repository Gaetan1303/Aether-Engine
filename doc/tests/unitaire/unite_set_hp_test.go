package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_SetHP teste la méthode SetHP()
func TestUnite_SetHP(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Gardien", "team-1", 5, 5)

	// Act
	unite.SetHP(75)

	// Assert
	hpApres := unite.StatsActuelles().HP
	assert.Equal(t, 75, hpApres, "HP devrait être défini à 75")
}
