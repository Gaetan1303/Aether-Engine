package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_PuissanceTotale teste la méthode PuissanceTotale()
func TestEquipe_PuissanceTotale(t *testing.T) {
	// Arrange
	teamID := domain.TeamID("team-1")
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(teamID, "Héros", "#0000FF", false, &joueurID)

	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(domain.UnitID("u1"), "Guerrier", teamID, stats, newTestPosition(0, 0))

	equipe.AjouterMembre(unite)

	// Act
	puissanceTotale := equipe.PuissanceTotale()

	// Assert
	expectedPuissance := 100 + 50 + 80 + 30 + 20 + 10 + 15 + 12 + 5 // = 322
	assert.Equal(t, expectedPuissance, puissanceTotale, "La puissance totale devrait être la somme de toutes les stats")
}
