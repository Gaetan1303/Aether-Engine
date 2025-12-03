package unitaire

import (
	"testing"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_AppliquerModificateurStat teste la méthode AppliquerModificateurStat()
func TestUnite_AppliquerModificateurStat(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Champion", "team-1", 5, 5)
	atkInitial := unite.StatsActuelles().ATK
	modificateur := &shared.ModificateurStat{
		Stat:   "ATK",
		Valeur: 10,
	}

	// Act
	unite.AppliquerModificateurStat(modificateur)

	// Assert
	atkApres := unite.StatsActuelles().ATK
	assert.Equal(t, atkInitial+10, atkApres, "ATK devrait augmenter de 10")
}
