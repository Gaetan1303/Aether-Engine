package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_Description teste la méthode Description()
func TestCompetence_Description(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("heal"), "Soin", domain.CompetenceSoin)

	// Act
	description := comp.Description()

	// Assert
	assert.NotEmpty(t, description, "La description ne devrait pas être vide")
	assert.Equal(t, "Description de test", description, "La description devrait correspondre")
}
