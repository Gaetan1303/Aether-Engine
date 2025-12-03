package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_NombreMembres teste la méthode NombreMembres()
func TestEquipe_NombreMembres(t *testing.T) {
	// Arrange
	teamID := domain.TeamID("team-1")
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(teamID, "Héros", "#0000FF", false, &joueurID)

	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite1 := domain.NewUnite(domain.UnitID("u1"), "Guerrier", teamID, stats, newTestPosition(0, 0))
	unite2 := domain.NewUnite(domain.UnitID("u2"), "Mage", teamID, stats, newTestPosition(1, 0))

	equipe.AjouterMembre(unite1)
	equipe.AjouterMembre(unite2)

	// Act
	nombreMembres := equipe.NombreMembres()

	// Assert
	assert.Equal(t, 2, nombreMembres, "L'équipe devrait avoir 2 membres")
}
