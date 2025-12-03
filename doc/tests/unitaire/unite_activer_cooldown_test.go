package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_ActiverCooldown teste la méthode ActiverCooldown()
func TestUnite_ActiverCooldown(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Ranger", "team-1", 5, 5)
	comp := newTestCompetence("arrow", "Flèche", domain.CompetenceAttaque)
	unite.AjouterCompetence(comp)

	// Act
	unite.ActiverCooldown("arrow", 3)

	// Assert
	competence := unite.ObtenirCompetence("arrow")
	assert.Equal(t, 3, competence.CooldownActuel(), "Le cooldown devrait être défini à 3")
}
