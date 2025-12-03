package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCombat_ObtenirResultat teste la méthode ObtenirResultat()
func TestCombat_ObtenirResultat(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	resultat := combat.ObtenirResultat()

	// Assert
	assert.NotEmpty(t, resultat, "Le résultat ne devrait pas être vide")
	assert.Contains(t, []string{"CONTINUE", "VICTORY", "DEFEAT", "FLED"}, resultat, "Devrait être un résultat valide")
}
