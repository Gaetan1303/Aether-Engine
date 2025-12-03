package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_ConsommerStamina teste la méthode ConsommerStamina()
func TestUnite_ConsommerStamina(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Gladiateur", "team-1", 5, 5)
	staminaInitial := unite.StatsActuelles().Stamina

	// Act
	err := unite.ConsommerStamina(5)

	// Assert
	assert.NoError(t, err, "La consommation de Stamina devrait réussir")
	staminaApres := unite.StatsActuelles().Stamina
	assert.Equal(t, staminaInitial-5, staminaApres, "Stamina devrait diminuer de 5")
}
