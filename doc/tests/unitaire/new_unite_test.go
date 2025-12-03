package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestNewUnite teste la création d'une nouvelle unité
func TestNewUnite(t *testing.T) {
	// Arrange
	id := domain.UnitID("unit-1")
	nom := "Guerrier"
	teamID := domain.TeamID("team-1")
	stats := &shared.Stats{
		HP:      100,
		MP:      50,
		Stamina: 80,
		ATK:     30,
		DEF:     20,
		MATK:    10,
		MDEF:    15,
		SPD:     12,
		MOV:     5,
	}
	position := newTestPosition(0, 0)

	// Act
	unite := domain.NewUnite(id, nom, teamID, stats, position)

	// Assert
	assert.NotNil(t, unite, "L'unité ne devrait pas être nil")
	assert.Equal(t, id, unite.ID(), "L'ID devrait correspondre")
	assert.Equal(t, nom, unite.Nom(), "Le nom devrait correspondre")
	assert.Equal(t, teamID, unite.TeamID(), "Le TeamID devrait correspondre")
	assert.Equal(t, position, unite.Position(), "La position devrait correspondre")
	assert.NotNil(t, unite.Stats(), "Les stats ne devraient pas être nil")
	assert.Equal(t, 100, unite.Stats().HP, "Les HP devraient être 100")
}
