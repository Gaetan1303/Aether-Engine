package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCompetence_Clone teste la méthode Clone()
func TestCompetence_Clone(t *testing.T) {
	// Arrange
	comp := newTestCompetence(domain.CompetenceID("skill"), "Compétence", domain.CompetenceAttaque)
	comp.ActiverCooldown()

	// Act
	clone := comp.Clone()

	// Assert
	assert.NotNil(t, clone, "Le clone ne devrait pas être nil")
	assert.Equal(t, comp.ID(), clone.ID(), "Le clone devrait avoir le même ID")
	assert.Equal(t, comp.Nom(), clone.Nom(), "Le clone devrait avoir le même nom")
	assert.Equal(t, comp.CooldownActuel(), clone.CooldownActuel(), "Le clone devrait avoir le même cooldown actuel")
	
	// Vérifier que c'est une copie indépendante
	clone.DecrémenterCooldown()
	assert.NotEqual(t, comp.CooldownActuel(), clone.CooldownActuel(), "Modifier le clone ne devrait pas affecter l'original")
}
