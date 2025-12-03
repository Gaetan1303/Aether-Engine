package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_StatsMoyennes teste la méthode StatsMoyennes()
func TestEquipe_StatsMoyennes(t *testing.T) {
	// Arrange
	teamID := domain.TeamID("team-1")
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(teamID, "Héros", "#0000FF", false, &joueurID)

	stats1 := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	stats2 := &shared.Stats{HP: 80, MP: 100, Stamina: 60, ATK: 20, DEF: 10, MATK: 40, MDEF: 25, SPD: 14, MOV: 4}

	unite1 := domain.NewUnite(domain.UnitID("u1"), "Guerrier", teamID, stats1, newTestPosition(0, 0))
	unite2 := domain.NewUnite(domain.UnitID("u2"), "Mage", teamID, stats2, newTestPosition(1, 0))

	equipe.AjouterMembre(unite1)
	equipe.AjouterMembre(unite2)

	// Act
	statsMoyennes := equipe.StatsMoyennes()

	// Assert
	assert.NotNil(t, statsMoyennes, "Les stats moyennes ne devraient pas être nil")
	assert.Equal(t, 90, statsMoyennes.HP, "HP moyen devrait être (100+80)/2 = 90")
	assert.Equal(t, 75, statsMoyennes.MP, "MP moyen devrait être (50+100)/2 = 75")
	assert.Equal(t, 70, statsMoyennes.Stamina, "Stamina moyen devrait être (80+60)/2 = 70")
	assert.Equal(t, 25, statsMoyennes.ATK, "ATK moyen devrait être (30+20)/2 = 25")
}
