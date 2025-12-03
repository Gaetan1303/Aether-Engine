package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_EstEnnemie teste la méthode EstEnnemie()
func TestEquipe_EstEnnemie(t *testing.T) {
	// Arrange
	joueurID1 := "player-1"
	joueurID2 := "player-2"
	equipe1, _ := domain.NewEquipe(domain.TeamID("team-1"), "Héros", "#0000FF", false, &joueurID1)
	equipe2, _ := domain.NewEquipe(domain.TeamID("team-2"), "Ennemis", "#FF0000", false, &joueurID2)
	equipe3, _ := domain.NewEquipe(domain.TeamID("team-1"), "Héros Bis", "#00FF00", false, &joueurID1)

	// Act
	estEnnemieDeEquipe2 := equipe1.EstEnnemie(equipe2)
	estEnnemieDeEquipe3 := equipe1.EstEnnemie(equipe3)

	// Assert
	assert.True(t, estEnnemieDeEquipe2, "equipe1 devrait être ennemie de equipe2")
	assert.False(t, estEnnemieDeEquipe3, "equipe1 ne devrait pas être ennemie d'une équipe avec le même ID")
}
