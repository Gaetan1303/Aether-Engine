package integration

import (
	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
)

// Helpers pour créer des objets de test avec le code existant

func newTestPosition(x, y int) *shared.Position {
	pos, _ := shared.NewPosition(x, y)
	return pos
}

func newTestStats(hp, mp, stamina, atk, def, matk, mdef, spd, mov int) *shared.Stats {
	stats, _ := shared.NewStats(hp, mp, stamina, atk, def, matk, mdef, spd, mov)
	return stats
}

func newTestUnite(id, nom, teamID string, x, y int) *domain.Unite {
	stats := newTestStats(100, 50, 80, 30, 20, 10, 15, 12, 5)
	position := newTestPosition(x, y)
	return domain.NewUnite(domain.UnitID(id), nom, domain.TeamID(teamID), stats, position)
}

func newTestEquipe(id, nom string) *domain.Equipe {
	joueurID := "player-1"
	equipe, _ := domain.NewEquipe(domain.TeamID(id), nom, "#FF0000", false, &joueurID)
	return equipe
}

func newTestEquipeIA(id, nom string) *domain.Equipe {
	equipe, _ := domain.NewEquipe(domain.TeamID(id), nom, "#0000FF", true, nil)
	return equipe
}

func newTestCompetence(id, nom string, typeComp domain.TypeCompetence) *domain.Competence {
	return domain.NewCompetence(
		domain.CompetenceID(id),
		nom,
		"Description",
		typeComp,
		5, // portée
		domain.ZoneEffet{},
		10,  // coutMP
		5,   // coutStamina
		2,   // cooldown
		20,  // degatsBase
		1.0, // modificateur
		domain.CibleEnnemis,
	)
}

func newTestCompetenceAvecCouts(id, nom string, coutMP, coutStamina, cooldown int) *domain.Competence {
	return domain.NewCompetence(
		domain.CompetenceID(id),
		nom,
		"Description",
		domain.CompetenceAttaque,
		5, // portée
		domain.ZoneEffet{},
		coutMP,
		coutStamina,
		cooldown,
		20,  // degatsBase
		1.0, // modificateur
		domain.CibleEnnemis,
	)
}

func newTestGrille(largeur, hauteur int) *shared.GrilleCombat {
	grille, _ := shared.NewGrilleCombat(largeur, hauteur)
	return grille
}

func newTestCombat(id string) *domain.Combat {
	grille := newTestGrille(10, 10)

	// Créer les équipes
	equipe1 := newTestEquipe("team-1", "Héros")
	equipe2 := newTestEquipe("team-2", "Ennemis")
	equipes := []*domain.Equipe{equipe1, equipe2}

	combat, _ := domain.NewCombat(id, equipes, grille)
	return combat
}

func newTestCombatAvecUnites(id string) (*domain.Combat, *domain.Unite, *domain.Unite) {
	combat := newTestCombat(id)

	// Créer deux unités
	unite1 := newTestUnite("hero-1", "Héros", "team-1", 0, 0)
	unite2 := newTestUnite("enemy-1", "Ennemi", "team-2", 5, 5)

	// Ajouter aux équipes
	equipes := combat.Equipes()
	equipes[domain.TeamID("team-1")].AjouterMembre(unite1)
	equipes[domain.TeamID("team-2")].AjouterMembre(unite2)

	return combat, unite1, unite2
}

func newTestStatut(typeStatut shared.TypeStatut, duree, puissance int) *shared.Statut {
	return shared.NewStatut(typeStatut, duree, puissance)
}
