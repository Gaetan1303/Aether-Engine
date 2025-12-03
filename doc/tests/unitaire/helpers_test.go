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
