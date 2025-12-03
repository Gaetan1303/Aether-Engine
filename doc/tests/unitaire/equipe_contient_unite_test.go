package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_ContientUnite teste la méthode ContientUnite()
func TestEquipe_ContientUnite(t *testing.T) {
	// Arrange
	teamID := domain.TeamID("team-1")
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(teamID, "Héros", "#0000FF", false, &joueurID)

	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	uniteID := domain.UnitID("u1")
	unite := domain.NewUnite(uniteID, "Archer", teamID, stats, newTestPosition(0, 0))

	equipe.AjouterMembre(unite)

	// Act
	contient := equipe.ContientUnite(uniteID)
	neContientPas := equipe.ContientUnite(domain.UnitID("u999"))

	// Assert
	assert.True(t, contient, "L'équipe devrait contenir l'unité ajoutée")
	assert.False(t, neContientPas, "L'équipe ne devrait pas contenir une unité non ajoutée")
}
