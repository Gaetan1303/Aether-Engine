package unitaire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEquipe_MembresElimines teste la méthode MembresElimines()
func TestEquipe_MembresElimines(t *testing.T) {
	// Arrange
	equipe := newTestEquipe("team-1", "Les Héros", "player-1")
	unite1 := newTestUnite("unite-1", "Guerrier", "team-1", 5, 5)
	unite2 := newTestUnite("unite-2", "Mage", "team-1", 6, 6)
	equipe.AjouterMembre(unite1)
	equipe.AjouterMembre(unite2)

	// Éliminer une unité
	unite1.RecevoirDegats(200)

	// Act
	elimines := equipe.MembresElimines()

	// Assert
	assert.Len(t, elimines, 1, "Devrait avoir 1 membre éliminé")
	assert.Equal(t, "unite-1", string(elimines[0].ID()), "L'unité éliminée devrait être unite-1")
}
