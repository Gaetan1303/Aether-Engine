package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_PeutUtiliserCompetence teste la méthode PeutUtiliserCompetence()
func TestUnite_PeutUtiliserCompetence(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Enchanteur", "team-1", 5, 5)
	comp := newTestCompetence("heal", "Soin", domain.CompetenceSoin)
	unite.AjouterCompetence(comp)

	// Act
	peutUtiliser := unite.PeutUtiliserCompetence("heal")

	// Assert
	assert.True(t, peutUtiliser, "Devrait pouvoir utiliser la compétence avec les ressources disponibles")
}
