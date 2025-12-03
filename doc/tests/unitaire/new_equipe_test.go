package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestNewEquipe teste la création d'une nouvelle équipe
func TestNewEquipe(t *testing.T) {
	// Arrange
	id := domain.TeamID("team-1")
	nom := "Héros"
	couleur := "#FF0000"
	isIA := false
	joueurID := "player-123"

	// Act
	equipe, err := domain.NewEquipe(id, nom, couleur, isIA, &joueurID)

	// Assert
	assert.NoError(t, err, "La création de l'équipe ne devrait pas retourner d'erreur")
	assert.NotNil(t, equipe, "L'équipe ne devrait pas être nil")
	assert.Equal(t, id, equipe.ID(), "L'ID devrait correspondre")
	assert.Equal(t, nom, equipe.Nom(), "Le nom devrait correspondre")
	assert.Equal(t, couleur, equipe.Couleur(), "La couleur devrait correspondre")
	assert.False(t, equipe.IsIA(), "L'équipe ne devrait pas être IA")
	assert.Equal(t, joueurID, *equipe.JoueurID(), "Le joueurID devrait correspondre")
}
