package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_AppliquerBuff teste la méthode AppliquerBuff()
func TestUnite_AppliquerBuff(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Prêtre", "team-1", 5, 5)
	defInitial := unite.StatsActuelles().DEF

	// Act
	unite.AppliquerBuff("DEF", 20, 3)

	// Assert
	defApres := unite.StatsActuelles().DEF
	assert.Equal(t, defInitial+20, defApres, "DEF devrait augmenter de 20 après le buff")
}
