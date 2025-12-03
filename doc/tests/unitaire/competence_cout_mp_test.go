package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_CoutMP teste la méthode CoutMP()
func TestCompetence_CoutMP(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("spell"), "Sort", domain.CompetenceMagie)

	// Act
	coutMP := comp.CoutMP()

	// Assert
	assert.Equal(t, 10, coutMP, "Le coût MP devrait être 10")
}
