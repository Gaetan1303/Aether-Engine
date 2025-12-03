package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_StatsActuelles teste la méthode StatsActuelles()
func TestUnite_StatsActuelles(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Guerrier", "team-1", 5, 5)

	// Act
	statsActuelles := unite.StatsActuelles()

	// Assert
	assert.NotNil(t, statsActuelles, "Les stats actuelles ne devraient pas être nil")
	assert.Equal(t, 100, statsActuelles.HP, "HP actuels devraient être 100")
	assert.Equal(t, 50, statsActuelles.MP, "MP actuels devraient être 50")
	assert.Equal(t, 30, statsActuelles.ATK, "ATK actuel devrait être 30")
}
