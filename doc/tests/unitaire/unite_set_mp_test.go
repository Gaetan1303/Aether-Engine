package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_SetMP teste la méthode SetMP()
func TestUnite_SetMP(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Illusionniste", "team-1", 5, 5)

	// Act
	unite.SetMP(30)

	// Assert
	mpApres := unite.StatsActuelles().MP
	assert.Equal(t, 30, mpApres, "MP devrait être défini à 30")
}
