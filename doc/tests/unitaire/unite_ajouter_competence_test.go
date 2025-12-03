package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_AjouterCompetence teste la méthode AjouterCompetence()
func TestUnite_AjouterCompetence(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Sorcier", "team-1", 5, 5)
	comp := newTestCompetence("lightning", "Éclair", domain.CompetenceMagie)

	// Act
	err := unite.AjouterCompetence(comp)

	// Assert
	assert.NoError(t, err, "L'ajout de compétence ne devrait pas échouer")
	assert.Len(t, unite.Competences(), 1, "Devrait avoir 1 compétence")
}
