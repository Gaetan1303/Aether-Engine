package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_Competences teste la méthode Competences()
func TestUnite_Competences(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Mage", "team-1", 5, 5)
	comp := newTestCompetence("fireball", "Boule de Feu", domain.CompetenceMagie)

	// Act
	unite.AjouterCompetence(comp)
	competences := unite.Competences()

	// Assert
	assert.NotNil(t, competences, "La liste de compétences ne devrait pas être nil")
	assert.Len(t, competences, 1, "Devrait avoir 1 compétence")
	assert.Equal(t, "fireball", string(competences[0].ID()), "Devrait contenir la compétence ajoutée")
}
