package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_Portee teste la méthode Portee()
func TestCompetence_Portee(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("arrow"), "Flèche", domain.CompetenceAttaque)

	// Act
	portee := comp.Portee()

	// Assert
	assert.Equal(t, 5, portee, "La portée devrait être 5")
}
