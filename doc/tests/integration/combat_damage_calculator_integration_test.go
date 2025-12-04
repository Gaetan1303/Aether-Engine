package integration

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestCombatDamageCalculator_Physical teste le calculateur de dégâts physiques
func TestCombatDamageCalculator_Physical(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	combat.SetPhysicalDamageMode()

	attaquant := newTestUnite("u1", "Guerrier", "team-1", 0, 0)
	cible := newTestUnite("u2", "Mage", "team-2", 1, 0)
	competence := newTestCompetence("attack", "Attaque", domain.CompetenceAttaque)

	// Act
	calculator := combat.GetDamageCalculator()
	degats := calculator.CalculerDegats(attaquant, cible, competence)

	// Assert
	assert.Greater(t, degats, 0, "Les dégâts physiques devraient être > 0")
}

// TestCombatDamageCalculator_ChangerModePhysique teste le changement de mode
func TestCombatDamageCalculator_ChangerModePhysique(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")

	// Act
	combat.SetPhysicalDamageMode()
	calculator := combat.GetDamageCalculator()

	// Assert
	assert.NotNil(t, calculator, "Le calculateur ne devrait pas être nil")
}

// TestCombatDamageCalculator_Magical teste les dégâts magiques
func TestCombatDamageCalculator_Magical(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	combat.SetMagicalDamageMode()

	attaquant := newTestUnite("u1", "Mage", "team-1", 0, 0)
	cible := newTestUnite("u2", "Guerrier", "team-2", 1, 0)
	competence := newTestCompetence("fireball", "Boule de Feu", domain.CompetenceMagie)

	// Act
	calculator := combat.GetDamageCalculator()
	degats := calculator.CalculerDegats(attaquant, cible, competence)

	// Assert
	assert.Greater(t, degats, 0, "Les dégâts magiques devraient être > 0")
}

// TestCombatDamageCalculator_Hybrid teste les dégâts hybrides
func TestCombatDamageCalculator_Hybrid(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	combat.SetHybridDamageMode(0.5, 0.5) // 50% physique, 50% magique

	attaquant := newTestUnite("u1", "Paladin", "team-1", 0, 0)
	cible := newTestUnite("u2", "Ennemi", "team-2", 1, 0)
	competence := newTestCompetence("smite", "Châtiment", domain.CompetenceAttaque)

	// Act
	calculator := combat.GetDamageCalculator()
	degats := calculator.CalculerDegats(attaquant, cible, competence)

	// Assert
	assert.Greater(t, degats, 0, "Les dégâts hybrides devraient être > 0")
}

// TestCombatDamageCalculator_HybridRatios teste différents ratios hybrides
func TestCombatDamageCalculator_HybridRatios(t *testing.T) {
	// Arrange
	combat1 := newTestCombat("combat-1")
	combat2 := newTestCombat("combat-2")

	combat1.SetHybridDamageMode(0.8, 0.2) // 80% physique
	combat2.SetHybridDamageMode(0.2, 0.8) // 80% magique

	attaquant := newTestUnite("u1", "Hybrid", "team-1", 0, 0)
	cible := newTestUnite("u2", "Tank", "team-2", 1, 0)
	competence := newTestCompetence("slash", "Entaille", domain.CompetenceAttaque)

	// Act
	degats1 := combat1.GetDamageCalculator().CalculerDegats(attaquant, cible, competence)
	degats2 := combat2.GetDamageCalculator().CalculerDegats(attaquant, cible, competence)

	// Assert
	assert.Greater(t, degats1, 0, "Dégâts hybrides (80/20) devraient être > 0")
	assert.Greater(t, degats2, 0, "Dégâts hybrides (20/80) devraient être > 0")
}

// TestCombatDamageCalculator_ComparePhysicalVsMagical compare physique vs magique
func TestCombatDamageCalculator_ComparePhysicalVsMagical(t *testing.T) {
	// Arrange
	combatPhys := newTestCombat("combat-phys")
	combatMag := newTestCombat("combat-mag")

	combatPhys.SetPhysicalDamageMode()
	combatMag.SetMagicalDamageMode()

	attaquant := newTestUnite("u1", "Attaquant", "team-1", 0, 0)
	cible := newTestUnite("u2", "Cible", "team-2", 1, 0)
	compPhys := newTestCompetence("slash", "Entaille", domain.CompetenceAttaque)
	compMag := newTestCompetence("fireball", "Boule de Feu", domain.CompetenceMagie)

	// Act
	degatsPhys := combatPhys.GetDamageCalculator().CalculerDegats(attaquant, cible, compPhys)
	degatsMag := combatMag.GetDamageCalculator().CalculerDegats(attaquant, cible, compMag)

	// Assert
	assert.Greater(t, degatsPhys, 0, "Dégâts physiques > 0")
	assert.Greater(t, degatsMag, 0, "Dégâts magiques > 0")
}

// TestCombatDamageCalculator_SetCustomCalculator teste l'injection d'un calculateur custom
func TestCombatDamageCalculator_SetCustomCalculator(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	// Créer un calculateur fixe et l'injecter
	fixedCalc := domain.NewFixedDamageCalculator(42)
	combat.SetDamageCalculator(fixedCalc)

	attaquant := newTestUnite("u1", "Attaquant", "team-1", 0, 0)
	cible := newTestUnite("u2", "Cible", "team-2", 1, 0)
	competence := newTestCompetence("attack", "Attaque", domain.CompetenceAttaque)

	// Act
	calculator := combat.GetDamageCalculator()
	degats := calculator.CalculerDegats(attaquant, cible, competence)

	// Assert
	assert.Equal(t, 42, degats, "Les dégâts fixes devraient être 42")
}

// TestCombatDamageCalculator_MinimumDamage teste la règle de dégâts minimum
func TestCombatDamageCalculator_MinimumDamage(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	combat.SetPhysicalDamageMode()

	// Attaquant très faible vs cible très résistante
	attaquantFaible := newTestUnite("u1", "Faible", "team-1", 0, 0)
	cibleTank := newTestUnite("u2", "Tank", "team-2", 1, 0)
	competence := newTestCompetence("weak-attack", "Attaque Faible", domain.CompetenceAttaque)

	// Act
	calculator := combat.GetDamageCalculator()
	degats := calculator.CalculerDegats(attaquantFaible, cibleTank, competence)

	// Assert
	assert.GreaterOrEqual(t, degats, 1, "Les dégâts minimum devraient être >= 1")
}

// TestCombatDamageCalculator_AvecModificateurCompetence teste les modificateurs
func TestCombatDamageCalculator_AvecModificateurCompetence(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	combat.SetPhysicalDamageMode()

	attaquant := newTestUnite("u1", "Guerrier", "team-1", 0, 0)
	cible := newTestUnite("u2", "Ennemi", "team-2", 1, 0)

	compNormale := newTestCompetenceAvecCouts("slash", "Entaille", 10, 5, 2)
	compPuissante := newTestCompetenceAvecCouts("power-slash", "Entaille Puissante", 20, 10, 4)

	// Act
	calculator := combat.GetDamageCalculator()
	degatsNormaux := calculator.CalculerDegats(attaquant, cible, compNormale)
	degatsPuissants := calculator.CalculerDegats(attaquant, cible, compPuissante)

	// Assert
	assert.Greater(t, degatsNormaux, 0, "Les dégâts normaux devraient être > 0")
	assert.Greater(t, degatsPuissants, 0, "Les dégâts puissants devraient être > 0")
}

// TestCombatDamageCalculator_InfligerDegatsIntegration teste l'infliction complète de dégâts
func TestCombatDamageCalculator_InfligerDegatsIntegration(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-1")
	combat.SetPhysicalDamageMode()

	attaquant := newTestUnite("u1", "Guerrier", "team-1", 0, 0)
	cible := newTestUnite("u2", "Ennemi", "team-2", 1, 0)
	hpInitial := cible.HPActuels()
	competence := newTestCompetence("attack", "Attaque", domain.CompetenceAttaque)

	// Act
	calculator := combat.GetDamageCalculator()
	degats := calculator.CalculerDegats(attaquant, cible, competence)
	cible.RecevoirDegats(degats)

	// Assert
	assert.Less(t, cible.HPActuels(), hpInitial, "Les HP devraient avoir diminué")
	assert.Equal(t, hpInitial-degats, cible.HPActuels(), "Les HP devraient correspondre")
}

// TestCombatDamageCalculator_CalculerPourDifferentsTypes teste différents types de compétences
func TestCombatDamageCalculator_CalculerPourDifferentsTypes(t *testing.T) {
	// Arrange
	combat := newTestCombat("combat-physical")
	combat.SetPhysicalDamageMode()

	attaquant := newTestUnite("u1", "Polyvalent", "team-1", 0, 0)
	cible := newTestUnite("u2", "Cible", "team-2", 1, 0)

	compPhys := newTestCompetence("slash", "Entaille", domain.CompetenceAttaque)
	compMag := newTestCompetence("fireball", "Boule de Feu", domain.CompetenceMagie)

	// Act
	degatsPhys := combat.GetDamageCalculator().CalculerDegats(attaquant, cible, compPhys)
	degatsMag := combat.GetDamageCalculator().CalculerDegats(attaquant, cible, compMag)

	// Assert
	assert.Greater(t, degatsPhys, 0, "Dégâts physiques devraient être calculés")
	assert.Greater(t, degatsMag, 0, "Dégâts magiques devraient être calculés")
}
