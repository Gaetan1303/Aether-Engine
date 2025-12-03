package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_SkillEstPret teste la méthode SkillEstPret()
func TestUnite_SkillEstPret(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Assassin", "team-1", 5, 5)
	comp := newTestCompetence("backstab", "Coup de Poignard", domain.CompetenceAttaque)
	unite.AjouterCompetence(comp)

	// Act
	estPret := unite.SkillEstPret("backstab")

	// Assert
	assert.True(t, estPret, "La compétence devrait être prête (pas en cooldown)")
}
