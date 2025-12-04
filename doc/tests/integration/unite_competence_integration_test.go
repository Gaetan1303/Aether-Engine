package integration

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/stretchr/testify/assert"
)

// TestUniteUtiliserCompetence_FluxComplet teste le flux complet d'utilisation d'une compétence
func TestUniteUtiliserCompetence_FluxComplet(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Mage", "team-1", 0, 0)
	competence := newTestCompetence("fireball", "Boule de Feu", domain.CompetenceMagie)
	unite.AjouterCompetence(competence)

	mpInitial := unite.Stats().MP
	staminaInitial := unite.Stats().Stamina

	// Act
	err := unite.UtiliserCompetence(domain.CompetenceID("fireball"))

	// Assert
	assert.NoError(t, err, "L'utilisation de la compétence devrait réussir")
	assert.True(t, competence.EstEnCooldown(), "La compétence devrait être en cooldown")
	assert.Equal(t, mpInitial-10, unite.StatsActuelles().MP, "Les MP devraient être consommés")
	assert.Equal(t, staminaInitial-5, unite.StatsActuelles().Stamina, "La stamina devrait être consommée")
}

// TestUniteUtiliserCompetence_MPInsuffisant teste l'échec quand MP insuffisant
func TestUniteUtiliserCompetence_MPInsuffisant(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Mage", "team-1", 0, 0)
	competence := newTestCompetenceAvecCouts("meteor", "Météore", 100, 5, 2)
	unite.AjouterCompetence(competence)

	// Act
	err := unite.UtiliserCompetence(domain.CompetenceID("meteor"))

	// Assert
	assert.Error(t, err, "L'utilisation devrait échouer (MP insuffisant)")
}

// TestUniteUtiliserCompetence_EnCooldown teste qu'on ne peut pas utiliser une compétence en cooldown
func TestUniteUtiliserCompetence_EnCooldown(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Mage", "team-1", 0, 0)
	competence := newTestCompetence("fireball", "Boule de Feu", domain.CompetenceMagie)
	unite.AjouterCompetence(competence)

	// Utiliser une première fois
	unite.UtiliserCompetence(domain.CompetenceID("fireball"))

	// Act - Essayer d'utiliser à nouveau
	err := unite.UtiliserCompetence(domain.CompetenceID("fireball"))

	// Assert
	assert.Error(t, err, "Ne devrait pas pouvoir utiliser une compétence en cooldown")
}

// TestUniteCompetence_CooldownDecremente teste que les cooldowns décrément avec NouveauTour()
func TestUniteCompetence_CooldownDecremente(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Mage", "team-1", 0, 0)
	competence := newTestCompetence("fireball", "Boule de Feu", domain.CompetenceMagie)
	unite.AjouterCompetence(competence)

	// Utiliser la compétence
	unite.UtiliserCompetence(domain.CompetenceID("fireball"))
	assert.True(t, competence.EstEnCooldown(), "Devrait être en cooldown")

	// Act - Nouveau tour
	unite.NouveauTour()

	// Assert - Le cooldown devrait décrémenter (2 -> 1)
	assert.True(t, competence.EstEnCooldown(), "Devrait encore être en cooldown après 1 tour")
	assert.Equal(t, 1, competence.CooldownActuel(), "Cooldown devrait être à 1")

	// Act - Encore un tour
	unite.NouveauTour()

	// Assert - Le cooldown devrait être terminé (1 -> 0)
	assert.False(t, competence.EstEnCooldown(), "Ne devrait plus être en cooldown après 2 tours")
}

// TestUniteCompetence_ObtenirCompetenceParDefaut teste l'attaque basique
func TestUniteCompetence_ObtenirCompetenceParDefaut(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Guerrier", "team-1", 0, 0)

	// Act
	competence := unite.ObtenirCompetenceParDefaut()

	// Assert
	assert.NotNil(t, competence, "Devrait avoir une compétence par défaut")
	assert.Equal(t, "attaque-basique", string(competence.ID()), "L'ID devrait être 'attaque-basique'")
}

// TestUniteCompetence_AjouterPlusieursCompetences teste l'ajout de plusieurs compétences
func TestUniteCompetence_AjouterPlusieursCompetences(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Mage", "team-1", 0, 0)
	comp1 := newTestCompetence("fireball", "Boule de Feu", domain.CompetenceMagie)
	comp2 := newTestCompetence("heal", "Soin", domain.CompetenceSoin)
	comp3 := newTestCompetence("shield", "Bouclier", domain.CompetenceBuff)

	// Act
	unite.AjouterCompetence(comp1)
	unite.AjouterCompetence(comp2)
	unite.AjouterCompetence(comp3)

	// Assert
	competences := unite.Competences()
	assert.Len(t, competences, 3, "Devrait avoir 3 compétences")
	assert.NotNil(t, unite.ObtenirCompetence(domain.CompetenceID("fireball")))
	assert.NotNil(t, unite.ObtenirCompetence(domain.CompetenceID("heal")))
	assert.NotNil(t, unite.ObtenirCompetence(domain.CompetenceID("shield")))
}
