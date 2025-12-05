package domain_test

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
)

// TestStrategyPattern_PhysicalDamage teste la stratégie de dégâts physiques
func TestStrategyPattern_PhysicalDamage(t *testing.T) {
	// Arrange: Créer attaquant et défenseur
	statsAttacker, _ := shared.NewStats(100, 50, 50, 50, 10, 20, 10, 10, 5, 80)
	statsDefender, _ := shared.NewStats(100, 50, 50, 10, 30, 10, 20, 10, 5, 80)
	positionA, _ := shared.NewPosition(0, 0)
	positionD, _ := shared.NewPosition(1, 0)

	attacker := domain.NewUnite("attacker-1", "Guerrier", "team-1", statsAttacker, positionA)
	defender := domain.NewUnite("defender-1", "Chevalier", "team-2", statsDefender, positionD)

	// Créer une compétence physique
	competence := domain.NewCompetence(
		"slash",
		"Slash",
		"Attaque physique puissante",
		domain.CompetenceAttaque,
		1,
		domain.ZoneEffet{},
		0,
		0,
		1,
		20,  // 20 dégâts de base
		1.0, // Scaling 100% ATK
		domain.CibleEnnemis,
	)

	// Act: Calculer dégâts avec PhysicalDamageCalculator
	calculator := domain.NewPhysicalDamageCalculator()
	degats := calculator.Calculate(attacker, defender, competence)

	// Assert: (ATK 50 - DEF 30) + base 20 + scaling (1.0 * 50) = 20 + 20 + 50 = 90
	expected := 90
	assert.Equal(t, expected, degats, "Dégâts physiques incorrects")
	assert.Equal(t, "Physical", calculator.GetType())
}

// TestStrategyPattern_MagicalDamage teste la stratégie de dégâts magiques
func TestStrategyPattern_MagicalDamage(t *testing.T) {
	// Arrange
	statsAttacker, _ := shared.NewStats(100, 100, 50, 20, 10, 60, 10, 10, 5, 80)
	statsDefender, _ := shared.NewStats(100, 50, 50, 10, 30, 10, 40, 10, 5, 80)
	positionA, _ := shared.NewPosition(0, 0)
	positionD, _ := shared.NewPosition(1, 0)

	attacker := domain.NewUnite("attacker-1", "Mage", "team-1", statsAttacker, positionA)
	defender := domain.NewUnite("defender-1", "Tank", "team-2", statsDefender, positionD)

	// Créer une compétence magique
	competence := domain.NewCompetence(
		"fireball",
		"Fireball",
		"Boule de feu dévastatrice",
		domain.CompetenceMagie,
		3,
		domain.ZoneEffet{},
		30,
		0,
		3,
		30,  // 30 dégâts de base
		0.8, // Scaling 80% MATK
		domain.CibleEnnemis,
	)

	// Act: Calculer dégâts avec MagicalDamageCalculator
	calculator := domain.NewMagicalDamageCalculator()
	degats := calculator.Calculate(attacker, defender, competence)

	// Assert: (MATK 60 - MDEF 40) + base 30 + scaling (0.8 * 60) + pénétration (20% MDEF)
	// = 20 + 30 + 48 + 8 = 106
	expected := 106
	assert.Equal(t, expected, degats, "Dégâts magiques incorrects")
	assert.Equal(t, "Magical", calculator.GetType())
}

// TestStrategyPattern_FixedDamage teste la stratégie de dégâts fixes
func TestStrategyPattern_FixedDamage(t *testing.T) {
	// Arrange
	statsAttacker, _ := shared.NewStats(100, 50, 50, 50, 10, 20, 10, 10, 5, 80)
	statsDefender, _ := shared.NewStats(100, 50, 50, 10, 99, 10, 99, 10, 5) // DEF/MDEF énormes
	positionA, _ := shared.NewPosition(0, 0)
	positionD, _ := shared.NewPosition(1, 0)

	attacker := domain.NewUnite("attacker-1", "Assassin", "team-1", statsAttacker, positionA)
	defender := domain.NewUnite("defender-1", "Tank", "team-2", statsDefender, positionD)

	competence := attacker.ObtenirCompetenceParDefaut()

	// Act: Calculer dégâts fixes (ignore les stats)
	calculator := domain.NewFixedDamageCalculator(50)
	degats := calculator.Calculate(attacker, defender, competence)

	// Assert: Toujours 50, peu importe les stats
	assert.Equal(t, 50, degats, "Dégâts fixes incorrects")
	assert.Equal(t, "Fixed", calculator.GetType())
}

// TestStrategyPattern_HybridDamage teste la stratégie de dégâts hybrides
func TestStrategyPattern_HybridDamage(t *testing.T) {
	// Arrange
	statsAttacker, _ := shared.NewStats(100, 50, 50, 40, 10, 40, 10, 10, 5, 80)
	statsDefender, _ := shared.NewStats(100, 50, 50, 10, 20, 10, 20, 10, 5, 80)
	positionA, _ := shared.NewPosition(0, 0)
	positionD, _ := shared.NewPosition(1, 0)

	attacker := domain.NewUnite("attacker-1", "Paladin", "team-1", statsAttacker, positionA)
	defender := domain.NewUnite("defender-1", "Ennemi", "team-2", statsDefender, positionD)

	competence := domain.NewCompetence(
		"holy-strike",
		"Holy Strike",
		"Frappe sacrée (physique + magique)",
		domain.CompetenceAttaque,
		1,
		domain.ZoneEffet{},
		15,
		0,
		2,
		20,  // 20 dégâts de base
		0.6, // Scaling 60%
		domain.CibleEnnemis,
	)

	// Act: 60% physique, 40% magique
	calculator := domain.NewHybridDamageCalculator(0.6, 0.4)
	degats := calculator.Calculate(attacker, defender, competence)

	// Assert: Physique (40-20)*0.6=12 + Magique (40-20)*0.4=8 + base 20 + scaling 0.6*(40*0.6+40*0.4)=14.4
	// = 12 + 8 + 20 + 14 = 54
	assert.Greater(t, degats, 50, "Dégâts hybrides trop faibles")
	assert.Equal(t, "Hybrid", calculator.GetType())
}

// TestStrategyPattern_ProportionalDamage teste les dégâts proportionnels aux HP
func TestStrategyPattern_ProportionalDamage(t *testing.T) {
	// Arrange
	statsAttacker, _ := shared.NewStats(100, 50, 50, 40, 10, 40, 10, 10, 5, 80)
	statsDefender, _ := shared.NewStats(1000, 50, 50, 10, 20, 10, 20, 10, 5) // 1000 HP
	positionA, _ := shared.NewPosition(0, 0)
	positionD, _ := shared.NewPosition(1, 0)

	attacker := domain.NewUnite("attacker-1", "Exécuteur", "team-1", statsAttacker, positionA)
	defender := domain.NewUnite("defender-1", "Boss", "team-2", statsDefender, positionD)

	competence := attacker.ObtenirCompetenceParDefaut()

	// Act: 15% des HP max
	calculator := domain.NewProportionalDamageCalculator(0.15, false)
	degats := calculator.Calculate(attacker, defender, competence)

	// Assert: 15% de 1000 = 150
	assert.Equal(t, 150, degats, "Dégâts proportionnels incorrects")
	assert.Equal(t, "Proportional", calculator.GetType())
}

// TestStrategyPattern_SwitchingStrategies teste le changement dynamique de stratégie
func TestStrategyPattern_SwitchingStrategies(t *testing.T) {
	// Arrange: Créer un combat
	statsAttacker, _ := shared.NewStats(100, 50, 50, 50, 10, 50, 10, 10, 5, 80)
	statsDefender, _ := shared.NewStats(100, 50, 50, 10, 30, 10, 30, 10, 5, 80)
	positionA, _ := shared.NewPosition(0, 0)
	positionD, _ := shared.NewPosition(5, 5)

	team1ID := domain.TeamID("team-1")
	team2ID := domain.TeamID("team-2")

	attacker := domain.NewUnite("attacker-1", "Combattant", team1ID, statsAttacker, positionA)
	defender := domain.NewUnite("defender-1", "Défenseur", team2ID, statsDefender, positionD)

	dummyPlayer1 := "player-1"
	dummyPlayer2 := "player-2"
	equipe1, _ := domain.NewEquipe(team1ID, "Équipe A", "blue", false, &dummyPlayer1)
	equipe2, _ := domain.NewEquipe(team2ID, "Équipe B", "red", false, &dummyPlayer2)
	equipe1.AjouterMembre(attacker)
	equipe2.AjouterMembre(defender)

	grille, _ := shared.NewGrilleCombat(10, 10)
	combat, _ := domain.NewCombat("combat-1", []*domain.Equipe{equipe1, equipe2}, grille)

	competence := attacker.ObtenirCompetenceParDefaut()

	// Act & Assert 1: Mode physique par défaut
	combat.SetPhysicalDamageMode()
	calculator1 := combat.GetDamageCalculator()
	degats1 := calculator1.Calculate(attacker, defender, competence)
	assert.Equal(t, "Physical", calculator1.GetType())
	assert.Greater(t, degats1, 0)

	// Act & Assert 2: Changer vers mode magique
	combat.SetMagicalDamageMode()
	calculator2 := combat.GetDamageCalculator()
	degats2 := calculator2.Calculate(attacker, defender, competence)
	assert.Equal(t, "Magical", calculator2.GetType())
	// Les dégâts magiques devraient être différents des dégâts physiques
	assert.NotEqual(t, degats1, degats2)

	// Act & Assert 3: Changer vers mode hybride
	combat.SetHybridDamageMode(0.5, 0.5)
	calculator3 := combat.GetDamageCalculator()
	degats3 := calculator3.Calculate(attacker, defender, competence)
	assert.Equal(t, "Hybrid", calculator3.GetType())
	assert.Greater(t, degats3, 0, "Dégâts hybrides doivent être positifs")

	// Act & Assert 4: Injecter une stratégie personnalisée
	customCalculator := domain.NewFixedDamageCalculator(999)
	combat.SetDamageCalculator(customCalculator)
	calculator4 := combat.GetDamageCalculator()
	degats4 := calculator4.Calculate(attacker, defender, competence)
	assert.Equal(t, "Fixed", calculator4.GetType())
	assert.Equal(t, 999, degats4, "Stratégie personnalisée non appliquée")
}

// TestStrategyPattern_Factory teste la factory de calculators
func TestStrategyPattern_Factory(t *testing.T) {
	// Arrange
	factory := domain.NewDamageCalculatorFactory()

	// Créer différentes compétences
	physicalSkill := domain.NewCompetence("skill1", "Physical", "", domain.CompetenceAttaque, 1, domain.ZoneEffet{}, 0, 0, 1, 10, 1.0, domain.CibleEnnemis)
	magicalSkill := domain.NewCompetence("skill2", "Magical", "", domain.CompetenceMagie, 1, domain.ZoneEffet{}, 10, 0, 1, 20, 1.0, domain.CibleEnnemis)
	supportSkill := domain.NewCompetence("skill3", "Support", "", domain.CompetenceSoin, 1, domain.ZoneEffet{}, 10, 0, 1, 0, 0, domain.CibleAllies)

	// Act: Factory crée le bon calculator selon le type
	calc1 := factory.CreateCalculator(physicalSkill)
	calc2 := factory.CreateCalculator(magicalSkill)
	calc3 := factory.CreateCalculator(supportSkill)

	// Assert: Vérifier les types
	assert.Equal(t, "Physical", calc1.GetType())
	assert.Equal(t, "Magical", calc2.GetType())
	assert.Equal(t, "Fixed", calc3.GetType()) // Support n'inflige pas de dégâts

	// Test CreateHybridCalculator
	hybridCalc := factory.CreateHybridCalculator(0.6, 0.4)
	assert.Equal(t, "Hybrid", hybridCalc.GetType())

	// Test CreateProportionalCalculator
	propCalc := factory.CreateProportionalCalculator(0.2, true)
	assert.Equal(t, "Proportional", propCalc.GetType())
}

// TestStrategyPattern_MinimumDamage teste que les dégâts ne descendent jamais sous 1
func TestStrategyPattern_MinimumDamage(t *testing.T) {
	// Arrange: Défenseur avec DEF/MDEF énormes
	statsAttacker, _ := shared.NewStats(100, 50, 50, 10, 10, 10, 10, 10, 5)   // ATK/MATK faibles
	statsDefender, _ := shared.NewStats(100, 50, 50, 10, 999, 10, 999, 10, 5) // DEF/MDEF énormes
	positionA, _ := shared.NewPosition(0, 0)
	positionD, _ := shared.NewPosition(1, 0)

	attacker := domain.NewUnite("attacker-1", "Faible", "team-1", statsAttacker, positionA)
	defender := domain.NewUnite("defender-1", "Indestructible", "team-2", statsDefender, positionD)

	competence := attacker.ObtenirCompetenceParDefaut()

	// Act & Assert: Physical
	physicalCalc := domain.NewPhysicalDamageCalculator()
	physicalDmg := physicalCalc.Calculate(attacker, defender, competence)
	assert.GreaterOrEqual(t, physicalDmg, 1, "Dégâts physiques < 1")

	// Act & Assert: Magical
	magicalCalc := domain.NewMagicalDamageCalculator()
	magicalDmg := magicalCalc.Calculate(attacker, defender, competence)
	assert.GreaterOrEqual(t, magicalDmg, 1, "Dégâts magiques < 1")

	// Act & Assert: Hybrid
	hybridCalc := domain.NewHybridDamageCalculator(0.5, 0.5)
	hybridDmg := hybridCalc.Calculate(attacker, defender, competence)
	assert.GreaterOrEqual(t, hybridDmg, 1, "Dégâts hybrides < 1")
}
