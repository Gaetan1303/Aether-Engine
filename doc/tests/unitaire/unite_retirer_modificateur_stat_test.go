package unitaire

import (
	"testing"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_RetirerModificateurStat teste la méthode RetirerModificateurStat()
func TestUnite_RetirerModificateurStat(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Titan", "team-1", 5, 5)
	modificateur := &shared.ModificateurStat{
		Stat:   "DEF",
		Valeur: 15,
	}
	unite.AppliquerModificateurStat(modificateur)
	defAvecBuff := unite.StatsActuelles().DEF

	// Act
	unite.RetirerModificateurStat(modificateur)

	// Assert
	defApres := unite.StatsActuelles().DEF
	assert.Equal(t, defAvecBuff-15, defApres, "DEF devrait diminuer de 15 après retrait")
}
