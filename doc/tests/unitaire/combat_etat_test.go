package unitaire

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_Etat teste la méthode Etat()
func TestCombat_Etat(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	etat := combat.Etat()

	// Assert
	assert.Equal(t, domain.EtatAttente, etat, "L'état initial devrait être EtatAttente")
}
