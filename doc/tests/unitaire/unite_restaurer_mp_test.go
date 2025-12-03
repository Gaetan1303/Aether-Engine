package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_RestaurerMP teste la méthode RestaurerMP()
func TestUnite_RestaurerMP(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Alchimiste", "team-1", 5, 5)
	unite.ConsommerMP(20)
	mpAvant := unite.StatsActuelles().MP

	// Act
	unite.RestaurerMP(15)

	// Assert
	mpApres := unite.StatsActuelles().MP
	assert.Equal(t, mpAvant+15, mpApres, "MP devrait augmenter de 15")
}
