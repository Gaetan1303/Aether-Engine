package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_RecevoirDegats teste la méthode RecevoirDegats()
func TestUnite_RecevoirDegats(t *testing.T) {
	// Arrange
	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(domain.UnitID("u1"), "Guerrier", domain.TeamID("team-1"), stats, newTestPosition(0, 0))
	degats := 30

	// Act
	unite.RecevoirDegats(degats)

	// Assert
	assert.Equal(t, 70, unite.HPActuels(), "Les HP devraient être réduits de 30 (100 - 30 = 70)")
	assert.False(t, unite.EstEliminee(), "L'unité ne devrait pas être éliminée")
}
