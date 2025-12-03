package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_ObtenirMembre teste la méthode ObtenirMembre()
func TestEquipe_ObtenirMembre(t *testing.T) {
	// Arrange
	teamID := domain.TeamID("team-1")
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(teamID, "Héros", "#0000FF", false, &joueurID)

	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	uniteID := domain.UnitID("u1")
	unite := domain.NewUnite(uniteID, "Mage", teamID, stats, newTestPosition(0, 0))

	equipe.AjouterMembre(unite)

	// Act
	membreRecupere := equipe.ObtenirMembre(uniteID)

	// Assert
	assert.NotNil(t, membreRecupere, "Le membre récupéré ne devrait pas être nil")
	assert.Equal(t, uniteID, membreRecupere.ID(), "L'ID du membre récupéré devrait correspondre")
	assert.Equal(t, "Mage", membreRecupere.Nom(), "Le nom du membre récupéré devrait correspondre")
}
