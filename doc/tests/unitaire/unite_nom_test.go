package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_Nom teste la méthode Nom()
func TestUnite_Nom(t *testing.T) {
	// Arrange
	expectedNom := "Chevalier Noir"
	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(domain.UnitID("u1"), expectedNom, domain.TeamID("team-1"), stats, newTestPosition(0, 0))

	// Act
	actualNom := unite.Nom()

	// Assert
	assert.Equal(t, expectedNom, actualNom, "Le nom retourné devrait correspondre au nom fourni")
}
