package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_Type teste la méthode Type()
func TestCompetence_Type(t *testing.T) {
	// Arrange
	expectedType := domain.CompetenceMagie
	comp := newTestCompetence(domain.CompetenceID("fireball"), "Boule de Feu", expectedType)

	// Act
	actualType := comp.Type()

	// Assert
	assert.Equal(t, expectedType, actualType, "Le type retourné devrait correspondre au type fourni")
}
