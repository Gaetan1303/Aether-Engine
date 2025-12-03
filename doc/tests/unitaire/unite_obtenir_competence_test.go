package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_ObtenirCompetence teste la méthode ObtenirCompetence()
func TestUnite_ObtenirCompetence(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Invocateur", "team-1", 5, 5)
	comp := newTestCompetence("summon", "Invocation", domain.CompetenceMagie)
	unite.AjouterCompetence(comp)

	// Act
	competence := unite.ObtenirCompetence("summon")

	// Assert
	assert.NotNil(t, competence, "La compétence devrait être trouvée")
	assert.Equal(t, "summon", string(competence.ID()), "ID devrait correspondre")
}
