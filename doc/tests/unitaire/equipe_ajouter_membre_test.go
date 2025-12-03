package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_AjouterMembre teste la méthode AjouterMembre()
func TestEquipe_AjouterMembre(t *testing.T) {
	// Arrange
	teamID := domain.TeamID("team-1")
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(teamID, "Héros", "#0000FF", false, &joueurID)

	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(domain.UnitID("u1"), "Guerrier", teamID, stats, newTestPosition(0, 0))

	// Act
	err := equipe.AjouterMembre(unite)

	// Assert
	assert.NoError(t, err, "L'ajout du membre ne devrait pas retourner d'erreur")
	assert.Equal(t, 1, equipe.NombreMembres(), "L'équipe devrait avoir 1 membre")
	assert.True(t, equipe.ContientUnite(unite.ID()), "L'équipe devrait contenir l'unité ajoutée")
}
