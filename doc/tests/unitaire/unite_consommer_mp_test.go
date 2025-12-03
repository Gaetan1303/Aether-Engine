package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_ConsommerMP teste la méthode ConsommerMP()
func TestUnite_ConsommerMP(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Arcaniste", "team-1", 5, 5)
	mpInitial := unite.StatsActuelles().MP

	// Act
	err := unite.ConsommerMP(10)

	// Assert
	assert.NoError(t, err, "La consommation de MP devrait réussir")
	mpApres := unite.StatsActuelles().MP
	assert.Equal(t, mpInitial-10, mpApres, "MP devrait diminuer de 10")
}
