package unitaire

import (
	"testing"

	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUnite_RetirerStatut teste la méthode RetirerStatut()
func TestUnite_RetirerStatut(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Druide", "team-1", 5, 5)

	// Act - Retirer un statut (même s'il n'existe pas, ne devrait pas crasher)
	unite.RetirerStatut(shared.TypeStatutPoison)

	// Assert - Pas d'erreur, méthode void
	assert.NotNil(t, unite, "L'unité devrait toujours exister")
}
