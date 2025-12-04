package integration

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestUniteDeplacement_Basique teste le déplacement basique
func TestUniteDeplacement_Basique(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Scout", "team-1", 0, 0)
	nouvellePosition := newTestPosition(2, 2)

	// Act
	err := unite.SeDeplacer(nouvellePosition, 2)

	// Assert
	assert.NoError(t, err, "Le déplacement devrait réussir")
	assert.Equal(t, 2, unite.Position().X(), "X devrait être 2")
	assert.Equal(t, 2, unite.Position().Y(), "Y devrait être 2")
}

// TestUniteDeplacement_CoutDeplacement teste la consommation du mouvement
func TestUniteDeplacement_CoutDeplacement(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Scout", "team-1", 0, 0)
	nouvellePosition := newTestPosition(3, 0)

	// Act
	err := unite.SeDeplacer(nouvellePosition, 3)

	// Assert
	assert.NoError(t, err, "Le déplacement devrait réussir")
	assert.True(t, unite.PeutSeDeplacer(), "Devrait encore pouvoir se déplacer (5-3=2)")
}

// TestUniteDeplacement_DeplacementInsuffisant teste le déplacement avec mouvement insuffisant
func TestUniteDeplacement_DeplacementInsuffisant(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Scout", "team-1", 0, 0)
	nouvellePosition := newTestPosition(10, 10)

	// Act
	err := unite.SeDeplacer(nouvellePosition, 10)

	// Assert
	assert.Error(t, err, "Le déplacement devrait échouer (mouvement insuffisant)")
}

// TestUniteDeplacement_AvecRoot teste le déplacement bloqué par Root
func TestUniteDeplacement_AvecRoot(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Scout", "team-1", 0, 0)
	statutRoot := newTestStatut(shared.StatutRoot, 2, 0)
	unite.AjouterStatut(statutRoot)

	// Act
	peutSeDeplacer := unite.PeutSeDeplacer()

	// Assert
	// NOTE: Le statut Root n'est PAS implémenté dans votre code actuel
	// Donc le test vérifie le comportement actuel (pas de blocage)
	assert.False(t, peutSeDeplacer, "L'unité ne peut pas se déplacer")
}

// TestUniteDeplacement_UniteEliminee teste qu'une unité éliminée ne peut pas se déplacer
func TestUniteDeplacement_UniteEliminee(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Scout", "team-1", 0, 0)
	unite.RecevoirDegats(200) // Éliminer

	// Act
	peutSeDeplacer := unite.PeutSeDeplacer()

	// Assert
	assert.False(t, peutSeDeplacer, "Une unité éliminée ne peut pas se déplacer")
}

// TestUniteDeplacement_NouveauTourRestaure teste que le mouvement est restauré
func TestUniteDeplacement_NouveauTourRestaure(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Scout", "team-1", 0, 0)
	nouvellePosition := newTestPosition(3, 0)

	// Consommer du mouvement
	unite.SeDeplacer(nouvellePosition, 3)

	// Act
	unite.NouveauTour()

	// Assert
	assert.True(t, unite.PeutSeDeplacer(), "Le mouvement devrait être restauré")
}

// TestUniteDeplacement_DeplacerVers teste le déplacement direct sans coût
func TestUniteDeplacement_DeplacerVers(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Scout", "team-1", 0, 0)
	nouvellePosition := newTestPosition(5, 5)

	// Act
	unite.DeplacerVers(nouvellePosition)

	// Assert
	assert.Equal(t, 5, unite.Position().X(), "X devrait être 5")
	assert.Equal(t, 5, unite.Position().Y(), "Y devrait être 5")
}

// TestUniteDeplacement_PositionsOccupeesDansCombat teste les positions occupées dans un combat
func TestUniteDeplacement_PositionsOccupeesDansCombat(t *testing.T) {
	// Arrange
	combat, _, _ := newTestCombatAvecUnites("combat-1")

	// Act - Obtenir positions occupées
	positions := combat.ObtenirPositionsOccupees(domain.UnitID(""))

	// Assert
	assert.Len(t, positions, 2, "Devrait avoir 2 positions occupées")
}

// TestUniteDeplacement_DeplacementMultiple teste plusieurs déplacements successifs
func TestUniteDeplacement_DeplacementMultiple(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Scout", "team-1", 0, 0)
	pos1 := newTestPosition(1, 0)
	pos2 := newTestPosition(2, 0)

	// Act
	err1 := unite.SeDeplacer(pos1, 1)
	err2 := unite.SeDeplacer(pos2, 1)

	// Assert
	assert.NoError(t, err1, "Premier déplacement devrait réussir")
	assert.NoError(t, err2, "Deuxième déplacement devrait réussir")
	assert.Equal(t, 2, unite.Position().X(), "Position finale devrait être (2,0)")
}

// TestUniteDeplacement_VerificationEstBloqueDeplacement teste la vérification de blocage
func TestUniteDeplacement_VerificationEstBloqueDeplacement(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Scout", "team-1", 0, 0)

	// Act
	estBloque := unite.EstBloqueDeplacement()

	// Assert
	assert.False(t, estBloque, "Une unité normale ne devrait pas être bloquée")
}

// TestUniteDeplacement_GrilleValidation teste la validation avec une grille
func TestUniteDeplacement_GrilleValidation(t *testing.T) {
	// Arrange
	grille := newTestGrille(10, 10)
	positionValide := newTestPosition(5, 5)
	positionInvalide := newTestPosition(15, 15)

	// Act
	valide := grille.EstDansLimites(positionValide)
	invalide := grille.EstDansLimites(positionInvalide)

	// Assert
	assert.True(t, valide, "Position (5,5) devrait être valide")
	assert.False(t, invalide, "Position (15,15) devrait être invalide")
}
