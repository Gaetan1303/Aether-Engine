package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_Cibles teste la méthode Cibles()
func TestCompetence_Cibles(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("attack"), "Attaque", domain.CompetenceAttaque)

	// Act
	cibles := comp.Cibles()

	// Assert
	assert.Equal(t, domain.CibleEnnemis, cibles, "Le type de cibles devrait être CibleEnnemis")
}
