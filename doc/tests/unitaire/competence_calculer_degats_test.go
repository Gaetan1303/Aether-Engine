package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_CalculerDegats teste la méthode CalculerDegats()
func TestCompetence_CalculerDegats(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("attack"), "Attaque", domain.CompetenceAttaque)
	stats := &shared.Stats{ATK: 30, MATK: 20}

	// Act
	degats := comp.CalculerDegats(stats)

	// Assert
	// degatsBase (20) + (ATK (30) * modificateur (0.5)) = 20 + 15 = 35
	assert.Equal(t, 35, degats, "Les dégâts calculés devraient être 35 (20 + 30*0.5)")
}
