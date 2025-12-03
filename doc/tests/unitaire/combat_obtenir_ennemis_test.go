package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_ObtenirEnnemis teste la méthode ObtenirEnnemis()
func TestCombat_ObtenirEnnemis(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	equipes := combat.Equipes()

	// Ajouter des unités aux deux équipes
	unite1 := newTestUnite(domain.UnitID("u1"), "Guerrier", domain.TeamID("team-1"), 0, 0)
	unite2 := newTestUnite(domain.UnitID("u2"), "Goblin", domain.TeamID("team-2"), 1, 0)
	unite3 := newTestUnite(domain.UnitID("u3"), "Orc", domain.TeamID("team-2"), 2, 0)

	equipes[domain.TeamID("team-1")].AjouterMembre(unite1)
	equipes[domain.TeamID("team-2")].AjouterMembre(unite2)
	equipes[domain.TeamID("team-2")].AjouterMembre(unite3)

	// Act
	ennemis := combat.ObtenirEnnemis(domain.TeamID("team-1"))

	// Assert
	assert.Equal(t, 2, len(ennemis), "L'équipe 1 devrait avoir 2 ennemis")
	assert.Equal(t, domain.UnitID("u2"), ennemis[0].ID(), "Le premier ennemi devrait être u2")
	assert.Equal(t, domain.UnitID("u3"), ennemis[1].ID(), "Le deuxième ennemi devrait être u3")
}
