package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_TrouverUnite teste la méthode TrouverUnite()
func TestCombat_TrouverUnite(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	equipes := combat.Equipes()

	// Ajouter une unité à l'équipe 1
	unite := newTestUnite(domain.UnitID("u1"), "Guerrier", domain.TeamID("team-1"), 0, 0)
	equipes[domain.TeamID("team-1")].AjouterMembre(unite)

	// Act
	uniteRecuperee := combat.TrouverUnite(domain.UnitID("u1"))
	uniteInexistante := combat.TrouverUnite(domain.UnitID("u999"))

	// Assert
	assert.NotNil(t, uniteRecuperee, "L'unité devrait être trouvée")
	assert.Equal(t, domain.UnitID("u1"), uniteRecuperee.ID(), "L'ID de l'unité devrait correspondre")
	assert.Nil(t, uniteInexistante, "Une unité inexistante devrait retourner nil")
}
