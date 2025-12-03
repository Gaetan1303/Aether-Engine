package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_MembresVivants teste la méthode MembresVivants()
func TestEquipe_MembresVivants(t *testing.T) {
	// Arrange
	teamID := domain.TeamID("team-1")
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(teamID, "Héros", "#0000FF", false, &joueurID)

	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite1 := domain.NewUnite(domain.UnitID("u1"), "Guerrier", teamID, stats, newTestPosition(0, 0))
	unite2 := domain.NewUnite(domain.UnitID("u2"), "Mage", teamID, stats, newTestPosition(1, 0))

	equipe.AjouterMembre(unite1)
	equipe.AjouterMembre(unite2)

	// Éliminer une unité
	unite1.RecevoirDegats(100)

	// Act
	membresVivants := equipe.MembresVivants()

	// Assert
	assert.Equal(t, 1, len(membresVivants), "Il devrait y avoir 1 membre vivant")
	assert.Equal(t, domain.UnitID("u2"), membresVivants[0].ID(), "Le membre vivant devrait être u2")
}
