package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_RetirerMembre teste la méthode RetirerMembre()
func TestEquipe_RetirerMembre(t *testing.T) {
	// Arrange
	teamID := domain.TeamID("team-1")
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(teamID, "Héros", "#0000FF", false, &joueurID)

	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(domain.UnitID("u1"), "Guerrier", teamID, stats, newTestPosition(0, 0))

	equipe.AjouterMembre(unite)

	// Act
	err := equipe.RetirerMembre(unite.ID())

	// Assert
	assert.NoError(t, err, "Le retrait du membre ne devrait pas retourner d'erreur")
	assert.Equal(t, 0, equipe.NombreMembres(), "L'équipe ne devrait plus avoir de membre")
	assert.False(t, equipe.ContientUnite(unite.ID()), "L'équipe ne devrait plus contenir l'unité")
}
