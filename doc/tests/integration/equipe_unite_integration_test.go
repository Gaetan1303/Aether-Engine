package integration

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestEquipeUnite_AjouterPlusieursUnites teste l'ajout de plusieurs membres
func TestEquipeUnite_AjouterPlusieursUnites(t *testing.T) {
	// Arrange
	equipe := newTestEquipe("team-1", "Héros")
	unite1 := newTestUnite("u1", "Guerrier", "team-1", 0, 0)
	unite2 := newTestUnite("u2", "Mage", "team-1", 1, 0)
	unite3 := newTestUnite("u3", "Archer", "team-1", 2, 0)

	// Act
	equipe.AjouterMembre(unite1)
	equipe.AjouterMembre(unite2)
	equipe.AjouterMembre(unite3)

	// Assert
	assert.Len(t, equipe.Membres(), 3, "Devrait avoir 3 membres")
}

// TestEquipeUnite_RetirerMembre teste le retrait d'un membre
func TestEquipeUnite_RetirerMembre(t *testing.T) {
	// Arrange
	equipe := newTestEquipe("team-1", "Héros")
	unite1 := newTestUnite("u1", "Guerrier", "team-1", 0, 0)
	unite2 := newTestUnite("u2", "Mage", "team-1", 1, 0)

	equipe.AjouterMembre(unite1)
	equipe.AjouterMembre(unite2)

	// Act
	equipe.RetirerMembre(domain.UnitID("u1"))

	// Assert
	assert.Len(t, equipe.Membres(), 1, "Devrait avoir 1 membre restant")
	assert.False(t, equipe.ContientUnite(domain.UnitID("u1")), "Ne devrait plus contenir u1")
}

// TestEquipeUnite_MembresVivantsEtElimines teste la séparation entre vivants et éliminés
func TestEquipeUnite_MembresVivantsEtElimines(t *testing.T) {
	// Arrange
	equipe := newTestEquipe("team-1", "Héros")
	unite1 := newTestUnite("u1", "Guerrier", "team-1", 0, 0)
	unite2 := newTestUnite("u2", "Mage", "team-1", 1, 0)

	equipe.AjouterMembre(unite1)
	equipe.AjouterMembre(unite2)

	// Éliminer unite1
	unite1.RecevoirDegats(200)

	// Act
	vivants := equipe.MembresVivants()
	elimines := equipe.MembresElimines()

	// Assert
	assert.Len(t, vivants, 1, "Devrait avoir 1 membre vivant")
	assert.Len(t, elimines, 1, "Devrait avoir 1 membre éliminé")
	assert.True(t, unite1.EstEliminee(), "Unite1 devrait être éliminée")
	assert.False(t, unite2.EstEliminee(), "Unite2 devrait être vivante")
}

// TestEquipeUnite_TousElimines teste la détection de tous les membres éliminés
func TestEquipeUnite_TousElimines(t *testing.T) {
	// Arrange
	equipe := newTestEquipe("team-1", "Héros")
	unite1 := newTestUnite("u1", "Guerrier", "team-1", 0, 0)
	unite2 := newTestUnite("u2", "Mage", "team-1", 1, 0)

	equipe.AjouterMembre(unite1)
	equipe.AjouterMembre(unite2)

	// Act - Éliminer tous
	unite1.RecevoirDegats(200)
	unite2.RecevoirDegats(200)

	// Assert
	assert.False(t, equipe.ADesMembresVivants(), "Ne devrait plus avoir de membres vivants")
}

// TestEquipeUnite_StatsMoyennes teste le calcul des stats moyennes
func TestEquipeUnite_StatsMoyennes(t *testing.T) {
	// Arrange
	equipe := newTestEquipe("team-1", "Héros")
	unite1 := newTestUnite("u1", "Guerrier", "team-1", 0, 0)
	unite2 := newTestUnite("u2", "Mage", "team-1", 1, 0)

	equipe.AjouterMembre(unite1)
	equipe.AjouterMembre(unite2)

	// Act
	statsMoyennes := equipe.StatsMoyennes()

	// Assert
	assert.NotNil(t, statsMoyennes, "Stats moyennes ne devraient pas être nil")
	assert.Equal(t, 100, statsMoyennes.HP, "HP moyen devrait être 100")
	assert.Equal(t, 50, statsMoyennes.MP, "MP moyen devrait être 50")
}

// TestEquipeUnite_PuissanceTotale teste le calcul de la puissance totale
func TestEquipeUnite_PuissanceTotale(t *testing.T) {
	// Arrange
	equipe := newTestEquipe("team-1", "Héros")
	unite1 := newTestUnite("u1", "Guerrier", "team-1", 0, 0)
	unite2 := newTestUnite("u2", "Mage", "team-1", 1, 0)

	equipe.AjouterMembre(unite1)
	equipe.AjouterMembre(unite2)

	// Act
	puissance := equipe.PuissanceTotale()

	// Assert
	assert.Greater(t, puissance, 0, "La puissance totale devrait être > 0")
}

// TestEquipeUnite_ContientUnite teste la vérification de présence d'une unité
func TestEquipeUnite_ContientUnite(t *testing.T) {
	// Arrange
	equipe := newTestEquipe("team-1", "Héros")
	unite1 := newTestUnite("u1", "Guerrier", "team-1", 0, 0)
	equipe.AjouterMembre(unite1)

	// Act & Assert
	assert.True(t, equipe.ContientUnite(domain.UnitID("u1")), "Devrait contenir u1")
	assert.False(t, equipe.ContientUnite(domain.UnitID("u2")), "Ne devrait pas contenir u2")
}

// TestEquipeUnite_EstComplete teste la vérification d'équipe complète
func TestEquipeUnite_EstComplete(t *testing.T) {
	// Arrange
	equipe := newTestEquipe("team-1", "Héros")

	// Ajouter 4 membres (max généralement 4-6)
	for i := 1; i <= 4; i++ {
		id := "u" + string(rune('0'+i))
		unite := newTestUnite(id, "Membre", "team-1", i, 0)
		equipe.AjouterMembre(unite)
	}

	// Act
	estComplete := equipe.EstComplete()

	// Assert
	assert.True(t, estComplete, "L'équipe avec 4 membres devrait être considérée complète")
}

// TestEquipeUnite_EquipeEnnemie teste l'identification d'équipe ennemie
func TestEquipeUnite_EquipeEnnemie(t *testing.T) {
	// Arrange
	equipe1 := newTestEquipe("team-1", "Héros")
	equipe2 := newTestEquipe("team-2", "Ennemis")

	// Act & Assert
	assert.True(t, equipe1.EstEnnemie(equipe2), "team-2 devrait être ennemie de team-1")
	assert.False(t, equipe1.EstEnnemie(equipe1), "Une équipe n'est pas ennemie d'elle-même")
}

// TestEquipeUnite_RessusciterMembre teste la résurrection d'un membre
func TestEquipeUnite_RessusciterMembre(t *testing.T) {
	// Arrange
	equipe := newTestEquipe("team-1", "Héros")
	unite := newTestUnite("u1", "Prêtre", "team-1", 0, 0)
	equipe.AjouterMembre(unite)

	// Éliminer l'unité
	unite.RecevoirDegats(200)
	assert.True(t, unite.EstEliminee(), "Devrait être éliminée")

	// Act - Ressusciter
	unite.Ressusciter(50)

	// Assert
	assert.False(t, unite.EstEliminee(), "Ne devrait plus être éliminée")
	assert.Equal(t, 50, unite.HPActuels(), "Devrait avoir 50 HP")
	assert.True(t, equipe.ADesMembresVivants(), "L'équipe devrait avoir des membres vivants")
}
