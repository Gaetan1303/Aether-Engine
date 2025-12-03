package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_Nom teste la méthode Nom()
func TestCompetence_Nom(t *testing.T) {
	// Arrange
	expectedNom := "Frappe Puissante"
	comp := newTestCompetence(domain.CompetenceID("power-strike"), expectedNom, domain.CompetenceAttaque)

	// Act
	actualNom := comp.Nom()

	// Assert
	assert.Equal(t, expectedNom, actualNom, "Le nom retourné devrait correspondre au nom fourni")
}
