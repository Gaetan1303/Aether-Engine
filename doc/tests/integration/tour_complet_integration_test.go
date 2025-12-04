package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTourComplet_NouveauTourRestaureRessources teste la restauration des ressources
func TestTourComplet_NouveauTourRestaureRessources(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Mage", "team-1", 0, 0)
	mpInitial := unite.Stats().MP

	// Consommer des MP
	unite.ConsommerMP(20)

	// Act
	unite.NouveauTour()

	// Assert
	// NOTE: NouveauTour() régénère MP (10% = 5 MP) mais PAS les HP
	assert.Greater(t, unite.StatsActuelles().MP, mpInitial-20, "Les MP devraient avoir régénéré partiellement")
}

// TestTourComplet_CooldownsDecrementes teste la décrémentation des cooldowns
func TestTourComplet_CooldownsDecrementes(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Mage", "team-1", 0, 0)
	competence := newTestCompetenceAvecCouts("fireball", "Boule de Feu", 10, 5, 2)
	unite.AjouterCompetence(competence)

	// Utiliser la compétence
	unite.UtiliserCompetence("fireball")
	assert.True(t, competence.EstEnCooldown(), "Devrait être en cooldown")

	// Act
	unite.NouveauTour()

	// Assert - Cooldown devrait décrémenter (2 -> 1)
	assert.Equal(t, 1, competence.CooldownActuel(), "Cooldown devrait être à 1 après 1 tour")
}

// TestTourComplet_RegenerationMP teste la régénération de MP
func TestTourComplet_RegenerationMP(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Mage", "team-1", 0, 0)

	// Consommer des MP
	unite.ConsommerMP(30)
	mpApresConsommation := unite.StatsActuelles().MP

	// Act
	unite.NouveauTour()

	// Assert - Régénération MP = 10% de 50 = 5 MP
	mpApresRegeneration := unite.StatsActuelles().MP
	assert.Greater(t, mpApresRegeneration, mpApresConsommation, "MP devrait augmenter")
	assert.Equal(t, mpApresConsommation+5, mpApresRegeneration, "Devrait régénérer 5 MP (10%)")
}

// TestTourComplet_CycleComplet teste un cycle complet de tour
func TestTourComplet_CycleComplet(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Guerrier", "team-1", 0, 0)
	competence := newTestCompetenceAvecCouts("slash", "Entaille", 10, 5, 2)
	unite.AjouterCompetence(competence)

	// Tour 1: Utiliser compétence
	unite.UtiliserCompetence("slash")
	mpApresCast := unite.StatsActuelles().MP

	// Act - Tour 2: Nouveau tour
	unite.NouveauTour()

	// Assert
	assert.True(t, unite.PeutAgir(), "Devrait pouvoir agir après nouveau tour")
	assert.True(t, unite.PeutSeDeplacer(), "Devrait pouvoir se déplacer après nouveau tour")
	assert.Equal(t, 1, competence.CooldownActuel(), "Cooldown devrait être à 1")
	assert.Greater(t, unite.StatsActuelles().MP, mpApresCast, "MP devrait avoir régénéré")
}

// TestTourComplet_PlusieursUnites teste le cycle pour plusieurs unités
func TestTourComplet_PlusieursUnites(t *testing.T) {
	// Arrange
	unite1 := newTestUnite("u1", "Guerrier", "team-1", 0, 0)
	unite2 := newTestUnite("u2", "Mage", "team-1", 1, 0)

	// Consommer des ressources
	unite1.ConsommerMP(10)
	unite2.ConsommerMP(20)

	// Act - Nouveau tour pour les deux
	unite1.NouveauTour()
	unite2.NouveauTour()

	// Assert
	assert.True(t, unite1.PeutAgir(), "Unite1 devrait pouvoir agir")
	assert.True(t, unite2.PeutAgir(), "Unite2 devrait pouvoir agir")
}

// TestTourComplet_AvecCombat teste le cycle dans un contexte de combat
func TestTourComplet_AvecCombat(t *testing.T) {
	// Arrange
	combat, unite1, unite2 := newTestCombatAvecUnites("combat-1")

	// Consommer des ressources
	unite1.ConsommerMP(15)
	unite2.ConsommerMP(10)

	// Act - Nouveau tour
	unite1.NouveauTour()
	unite2.NouveauTour()

	// Assert
	assert.NotNil(t, combat, "Combat devrait exister")
	assert.True(t, unite1.PeutAgir(), "Unite1 devrait pouvoir agir")
	assert.True(t, unite2.PeutAgir(), "Unite2 devrait pouvoir agir")
	assert.Greater(t, unite1.StatsActuelles().MP, 35, "MP devrait avoir régénéré (50-15+5)")
}

// TestTourComplet_UniteElimineeNePasProceder teste qu'une unité morte ne procède pas
func TestTourComplet_UniteElimineeNePasProceder(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Guerrier", "team-1", 0, 0)
	unite.RecevoirDegats(200) // Éliminer

	// Act
	unite.NouveauTour()

	// Assert
	assert.False(t, unite.PeutAgir(), "Une unité éliminée ne peut pas agir")
	assert.False(t, unite.PeutSeDeplacer(), "Une unité éliminée ne peut pas se déplacer")
}

// TestTourComplet_ActionsEtDeplacementRestaurees teste la restauration des compteurs
func TestTourComplet_ActionsEtDeplacementRestaurees(t *testing.T) {
	// Arrange
	unite := newTestUnite("unite-1", "Scout", "team-1", 0, 0)

	// Consommer action et déplacement
	unite.SeDeplacer(newTestPosition(3, 0), 3)

	// Act
	unite.NouveauTour()

	// Assert
	assert.True(t, unite.PeutAgir(), "Actions devraient être restaurées")
	assert.True(t, unite.PeutSeDeplacer(), "Déplacement devrait être restauré")
}
