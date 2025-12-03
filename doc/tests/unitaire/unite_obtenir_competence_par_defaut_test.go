package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnite_ObtenirCompetenceParDefaut teste la méthode ObtenirCompetenceParDefaut()
func TestUnite_ObtenirCompetenceParDefaut(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Recrue", "team-1", 5, 5)

	// Act
	competence := unite.ObtenirCompetenceParDefaut()

	// Assert
	assert.NotNil(t, competence, "La compétence par défaut ne devrait pas être nil")
	assert.Equal(t, "attaque-basique", string(competence.ID()), "Devrait être l'attaque basique")
	assert.Equal(t, "Attaque Basique", competence.Nom(), "Le nom devrait être Attaque Basique")
}
