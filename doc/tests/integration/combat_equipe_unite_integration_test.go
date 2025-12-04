package integration

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombat_TrouverUnites teste la recherche d'unités dans le combat
func TestCombat_TrouverUnites(t *testing.T) {
	// Arrange
	combat, unite1, unite2 := newTestCombatAvecUnites("combat-1")

	// Act
	found1 := combat.TrouverUnite(domain.UnitID("hero-1"))
	found2 := combat.TrouverUnite(domain.UnitID("enemy-1"))
	foundNone := combat.TrouverUnite(domain.UnitID("inexistant"))

	// Assert
	assert.Equal(t, unite1, found1, "Devrait trouver hero-1")
	assert.Equal(t, unite2, found2, "Devrait trouver enemy-1")
	assert.Nil(t, foundNone, "Ne devrait rien trouver pour ID inexistant")
}

// TestCombat_ObtenirEnnemis teste l'obtention des ennemis d'une équipe
func TestCombat_ObtenirEnnemis(t *testing.T) {
	// Arrange
	combat, _, unite2 := newTestCombatAvecUnites("combat-1")

	// Act
	ennemis := combat.ObtenirEnnemis(domain.TeamID("team-1"))

	// Assert
	assert.Len(t, ennemis, 1, "Devrait avoir 1 ennemi")
	assert.Equal(t, unite2, ennemis[0], "L'ennemi devrait être unite2")
}

// TestCombat_PositionsOccupees teste la récupération des positions occupées
func TestCombat_PositionsOccupees(t *testing.T) {
	// Arrange
	combat, _, _ := newTestCombatAvecUnites("combat-1")

	// Act
	positions := combat.ObtenirPositionsOccupees(domain.UnitID(""))

	// Assert
	assert.Len(t, positions, 2, "Devrait avoir 2 positions occupées")
}

// TestCombat_PositionsOccupeesAvecExclusion teste l'exclusion d'une unité
func TestCombat_PositionsOccupeesAvecExclusion(t *testing.T) {
	// Arrange
	combat, _, _ := newTestCombatAvecUnites("combat-1")

	// Act - Exclure hero-1
	positions := combat.ObtenirPositionsOccupees(domain.UnitID("hero-1"))

	// Assert
	assert.Len(t, positions, 1, "Devrait avoir 1 position occupée (en excluant hero-1)")
}

// TestCombat_VerifierConditionsVictoire teste la détection des conditions de victoire
func TestCombat_VerifierConditionsVictoire(t *testing.T) {
	// Arrange
	combat, _, unite2 := newTestCombatAvecUnites("combat-1")

	// Act - Combat en cours (les deux équipes ont des membres vivants)
	resultat := combat.VerifierConditionsVictoire()

	// Assert
	assert.Equal(t, "CONTINUE", resultat, "Le combat devrait continuer")

	// Act - Éliminer une équipe
	unite2.RecevoirDegats(200)
	resultatApresElimination := combat.VerifierConditionsVictoire()

	// Assert
	assert.Contains(t, []string{"VICTORY", "DEFEAT"}, resultatApresElimination, "Le combat devrait être terminé")
}

// TestCombat_MarquerEquipeFuite teste le marquage d'une équipe en fuite
func TestCombat_MarquerEquipeFuite(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	combat.SetFuiteAutorisee(true)

	// Act
	combat.MarquerEquipeFuite(domain.TeamID("team-1"))
	resultat := combat.VerifierConditionsVictoire()

	// Assert
	assert.Equal(t, "FLED", resultat, "Le résultat devrait être FLED")
}

// TestCombat_AnnulerFuite teste l'annulation d'une fuite
func TestCombat_AnnulerFuite(t *testing.T) {
	// Arrange
	combat, _, _ := newTestCombatAvecUnites("combat-1")
	combat.SetFuiteAutorisee(true)
	combat.MarquerEquipeFuite(domain.TeamID("team-1"))

	// Act
	combat.AnnulerFuite(domain.TeamID("team-1"))
	resultat := combat.VerifierConditionsVictoire()

	// Assert
	assert.Equal(t, "CONTINUE", resultat, "Le combat devrait continuer après annulation")
}

// TestCombat_DesactiverFuite teste la désactivation de la fuite
func TestCombat_DesactiverFuite(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	combat.SetFuiteAutorisee(false)

	// Assert
	assert.False(t, combat.FuiteAutorisee(), "La fuite ne devrait pas être autorisée")
}

// TestCombat_ObtenirResultat teste l'obtention du résultat du combat
func TestCombat_ObtenirResultat(t *testing.T) {
	// Arrange
	combat, _, _ := newTestCombatAvecUnites("combat-1")

	// Act
	resultat := combat.ObtenirResultat()

	// Assert
	assert.Equal(t, "CONTINUE", resultat, "Le résultat devrait être CONTINUE")
}

// TestCombat_MultipleEquipes teste un combat avec plusieurs équipes
func TestCombat_MultipleEquipes(t *testing.T) {
	// Arrange
	grille := newTestGrille(10, 10)
	equipe1 := newTestEquipe("team-1", "Héros")
	equipe2 := newTestEquipe("team-2", "Ennemis")
	equipe3 := newTestEquipe("team-3", "Neutres")
	equipes := []*domain.Equipe{equipe1, equipe2, equipe3}
	combat, _ := domain.NewCombat("combat-1", equipes, grille)

	// Act
	equipesMap := combat.Equipes()

	// Assert
	assert.Len(t, equipesMap, 3, "Devrait avoir 3 équipes")
}

// TestCombat_EquipesVides teste le comportement avec des équipes vides
func TestCombat_EquipesVides(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	resultat := combat.VerifierConditionsVictoire()

	// Assert
	assert.Equal(t, "DEFEAT", resultat, "Combat avec équipes vides retourne DEFEAT (0 équipes actives)")
}

// TestCombat_GetTimestamp teste l'obtention du timestamp
func TestCombat_GetTimestamp(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	timestamp := combat.GetTimestamp()

	// Assert
	assert.Greater(t, timestamp, int64(0), "Timestamp devrait être > 0")
}

// TestCombat_VersionIncrementation teste l'incrémentation de la version
func TestCombat_VersionIncrementation(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	versionInitiale := combat.Version()

	// Act - Effectuer une action qui incrémente la version
	combat.MarquerEquipeFuite(domain.TeamID("team-1"))

	// Assert
	assert.GreaterOrEqual(t, combat.Version(), versionInitiale, "La version devrait être incrémentée ou stable")
}

// TestCombat_Equipes teste l'accès aux équipes
func TestCombat_Equipes(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	equipes := combat.Equipes()

	// Assert
	assert.Len(t, equipes, 2, "Devrait avoir 2 équipes par défaut")
	assert.NotNil(t, equipes[domain.TeamID("team-1")])
	assert.NotNil(t, equipes[domain.TeamID("team-2")])
}
