package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_EstCibleValide teste la méthode EstCibleValide()
func TestCompetence_EstCibleValide(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("attack"), "Attaque", domain.CompetenceAttaque)
	teamID1 := domain.TeamID("team-1")
	teamID2 := domain.TeamID("team-2")

	lanceur := newTestUnite(domain.UnitID("u1"), "Guerrier", teamID1, 0, 0)
	ennemi := newTestUnite(domain.UnitID("u2"), "Goblin", teamID2, 1, 0)
	allie := newTestUnite(domain.UnitID("u3"), "Mage", teamID1, 2, 0)

	// Act
	estValideEnnemi := comp.EstCibleValide(lanceur, ennemi)
	estValideAllie := comp.EstCibleValide(lanceur, allie)

	// Assert
	assert.True(t, estValideEnnemi, "Un ennemi devrait être une cible valide pour une compétence CibleEnnemis")
	assert.False(t, estValideAllie, "Un allié ne devrait pas être une cible valide pour une compétence CibleEnnemis")
}
