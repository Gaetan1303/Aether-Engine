package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_TeamID teste la méthode TeamID()
func TestUnite_TeamID(t *testing.T) {
	// Arrange
	expectedTeamID := domain.TeamID("team-heroes")
	stats := &shared.Stats{HP: 100, MP: 50, Stamina: 80, ATK: 30, DEF: 20, MATK: 10, MDEF: 15, SPD: 12, MOV: 5}
	unite := domain.NewUnite(domain.UnitID("u1"), "Héros", expectedTeamID, stats, newTestPosition(0, 0))

	// Act
	actualTeamID := unite.TeamID()

	// Assert
	assert.Equal(t, expectedTeamID, actualTeamID, "Le TeamID retourné devrait correspondre au TeamID fourni")
}
