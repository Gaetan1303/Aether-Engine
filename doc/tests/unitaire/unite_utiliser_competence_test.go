package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_UtiliserCompetence teste la méthode UtiliserCompetence()
func TestUnite_UtiliserCompetence(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Pyromancien", "team-1", 5, 5)
	comp := newTestCompetence("fireball", "Boule de Feu", domain.CompetenceMagie)
	unite.AjouterCompetence(comp)

	// Act
	err := unite.UtiliserCompetence("fireball")

	// Assert
	assert.NoError(t, err, "L'utilisation de la compétence devrait réussir")
	assert.True(t, comp.EstEnCooldown(), "La compétence devrait être en cooldown après utilisation")
}
