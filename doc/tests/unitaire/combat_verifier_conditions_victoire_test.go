package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_VerifierConditionsVictoire teste la méthode VerifierConditionsVictoire()
func TestCombat_VerifierConditionsVictoire(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	equipes := combat.Equipes()
	
	// Ajouter des unités aux deux équipes
	unite1 := newTestUnite(domain.UnitID("u1"), "Guerrier", domain.TeamID("team-1"), 0, 0)
	unite2 := newTestUnite(domain.UnitID("u2"), "Goblin", domain.TeamID("team-2"), 1, 0)
	
	equipes[domain.TeamID("team-1")].AjouterMembre(unite1)
	equipes[domain.TeamID("team-2")].AjouterMembre(unite2)

	// Act
	resultat := combat.VerifierConditionsVictoire()

	// Assert
	assert.Equal(t, "CONTINUE", resultat, "Le combat devrait continuer quand les deux équipes ont des unités vivantes")
}
