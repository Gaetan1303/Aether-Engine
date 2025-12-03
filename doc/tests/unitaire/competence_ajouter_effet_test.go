package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_AjouterEffet teste la méthode AjouterEffet()
func TestCompetence_AjouterEffet(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("spell"), "Sort", domain.CompetenceMagie)
	effet := domain.EffetCompetence{} // Effet vide pour le test

	// Act
	comp.AjouterEffet(effet)

	// Assert
	effets := comp.Effets()
	assert.Equal(t, 1, len(effets), "La compétence devrait avoir 1 effet après ajout")
}
