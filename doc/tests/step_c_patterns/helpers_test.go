package step_c_patterns_test

import (
	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// createTestCombat crée un combat de test avec deux équipes vides
func createTestCombat() *domain.Combat {
	// Créer une grille de test
	grille, err := shared.NewGrilleCombat(10, 10)
	if err != nil {
		panic(err)
	}

	// Créer des équipes vides (isIA = true pour les tests)
	team1, err := domain.NewEquipe(domain.TeamID("team1"), "Team 1", "blue", true, nil)
	if err != nil {
		panic(err)
	}
	team2, err := domain.NewEquipe(domain.TeamID("team2"), "Team 2", "red", true, nil)
	if err != nil {
		panic(err)
	}

	combat, err := domain.NewCombat("test_combat", []*domain.Equipe{team1, team2}, grille)
	if err != nil {
		panic(err)
	}
	return combat
}

// createTestUnit crée une unité de test avec les stats spécifiées
func createTestUnit(id string, speed int) *domain.Unite {
	return createTestUnitWithTeam(id, speed, "team1")
}

// createTestUnitWithTeam crée une unité avec une team spécifique
func createTestUnitWithTeam(id string, speed int, teamID string) *domain.Unite {
	// HP, MP, Stamina, ATK, DEF, MATK, MDEF, SPD, MOV
	stats, err := shared.NewStats(100, 50, 20, 15, 10, 15, 10, speed, 5, 80)
	if err != nil {
		panic(err)
	}
	position, err := shared.NewPosition(0, 0)
	if err != nil {
		panic(err)
	}
	return domain.NewUnite(
		domain.UnitID(id),
		"Test Unit",
		domain.TeamID(teamID),
		stats,
		position,
	)
}

// addUnitToCombat ajoute une unité au combat via son équipe
func addUnitToCombat(combat *domain.Combat, unit *domain.Unite) {
	teams := combat.Equipes()
	for _, team := range teams {
		if team.ID() == unit.TeamID() {
			team.AjouterMembre(unit)
			return
		}
	}
}

// createTestSkill crée une compétence de test
func createTestSkill(id string, cost int, skillType domain.TypeCompetence) *domain.Competence {
	zone := domain.ZoneEffet{} // Zone par défaut
	return domain.NewCompetence(
		domain.CompetenceID(id),
		"Test Skill",
		"Test skill description",
		skillType,
		2,    // Portée
		zone, // Zone d'effet
		cost, // CoutMP
		0,    // CoutStamina
		1,    // Cooldown
		50,   // DegatsBase
		1.0,  // Modificateur
		domain.CibleEnnemis,
	)
}

// createTestItem crée un objet de test
func createTestItem(id string, itemType string, effectValue int) *shared.Item {
	return &shared.Item{
		ID:          id,
		Name:        "Test Item",
		ItemType:    itemType,
		EffectVal:   effectValue,
		Range:       1,
		Description: "Test item",
	}
}
