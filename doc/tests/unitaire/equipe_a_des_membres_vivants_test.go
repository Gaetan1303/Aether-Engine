package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_ADesMembresVivants teste la méthode ADesMembresVivants()
func TestEquipe_ADesMembresVivants(t *testing.T) {
	// Arrange
	teamID := domain.TeamID("team-1")
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(teamID, "Héros", "#0000FF", false, &joueurID)

	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(domain.UnitID("u1"), "Guerrier", teamID, stats, newTestPosition(0, 0))

	equipe.AjouterMembre(unite)

	// Act
	aDesMembresVivants := equipe.ADesMembresVivants()

	// Assert
	assert.True(t, aDesMembresVivants, "L'équipe devrait avoir des membres vivants")
}
