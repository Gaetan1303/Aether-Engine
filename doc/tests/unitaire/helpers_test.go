package unitaire

import (
	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// newTestPosition crée une position pour les tests (panic si erreur)
func newTestPosition(x, y int) *shared.Position {
	pos, err := shared.NewPosition(x, y)
	if err != nil {
		panic("erreur création position de test: " + err.Error())
	}
	return pos
}

// newTestUnite crée une unité de test avec des paramètres simplifiés
func newTestUnite(id domain.UnitID, nom string, teamID domain.TeamID, x, y int) *domain.Unite {
	stats := &shared.Stats{
		HP:      100,
		MP:      50,
		Stamina: 80,
		ATK:     30,
		DEF:     20,
		MATK:    10,
		MDEF:    15,
		SPD:     12,
		MOV:     5,
	}
	return domain.NewUnite(id, nom, teamID, stats, newTestPosition(x, y))
}

// newTestEquipe crée une équipe de test simplifiée
func newTestEquipe(id domain.TeamID, nom string, joueurID string) *domain.Equipe {
	equipe, err := domain.NewEquipe(id, nom, "#0000FF", false, &joueurID)
	if err != nil {
		panic("erreur création équipe de test: " + err.Error())
	}
	return equipe
}

// newTestStats crée des stats de test
func newTestStats() *shared.Stats {
	return &shared.Stats{
		HP:      100,
		MP:      50,
		Stamina: 80,
		ATK:     30,
		DEF:     20,
		MATK:    10,
		MDEF:    15,
		SPD:     12,
		MOV:     5,
	}
}

// newTestCompetence crée une compétence de test simplifiée
func newTestCompetence(id domain.CompetenceID, nom string, typeComp domain.TypeCompetence) *domain.Competence {
	zone := domain.ZoneEffet{} // ZoneSingle par défaut
	return domain.NewCompetence(
		id,
		nom,
		"Description de test",
		typeComp,
		5,    // portée
		zone, // zone
		10,   // coutMP
		5,    // coutStamina
		2,    // cooldown
		20,   // degatsBase
		0.5,  // modificateur
		domain.CibleEnnemis,
	)
}

// newTestGrille crée une grille de combat de test
func newTestGrille(largeur, hauteur int) *shared.GrilleCombat {
	grille, err := shared.NewGrilleCombat(largeur, hauteur)
	if err != nil {
		panic("erreur création grille de test: " + err.Error())
	}
	return grille
}

// newTestCombat crée un combat de test avec 2 équipes
func newTestCombat(id string) *domain.Combat {
	joueur1 := "player-1"
	joueur2 := "player-2"

	equipe1, _ := domain.NewEquipe(domain.TeamID("team-1"), "Héros", "#0000FF", false, &joueur1)
	equipe2, _ := domain.NewEquipe(domain.TeamID("team-2"), "Ennemis", "#FF0000", true, &joueur2)

	grille := newTestGrille(10, 10)

	combat, err := domain.NewCombat(id, []*domain.Equipe{equipe1, equipe2}, grille)
	if err != nil {
		panic("erreur création combat de test: " + err.Error())
	}
	return combat
}
