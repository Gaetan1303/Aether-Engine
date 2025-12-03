package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipe_IsIA teste la méthode IsIA()
func TestEquipe_IsIA(t *testing.T) {
	// Arrange
	equipe, _ := domain.NewEquipe(domain.TeamID("team-1"), "IA Ennemis", "#FF0000", true, nil)

	// Act
	isIA := equipe.IsIA()

	// Assert
	assert.True(t, isIA, "L'équipe devrait être marquée comme IA")
}
