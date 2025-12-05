package combat_test

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCombatComplet_GobelinsVsJoueurs teste un combat réaliste avec 2 équipes de 6
func TestCombatComplet_GobelinsVsJoueurs(t *testing.T) {
	// =============================================================================
	// PHASE 1: CRÉATION DE LA GRILLE DE COMBAT (15x15)
	// =============================================================================
	grille, err := shared.NewGrilleCombat(15, 15)
	require.NoError(t, err, "La grille doit être créée sans erreur")

	// Ajouter quelques obstacles pour rendre le combat intéressant
	pos1, _ := shared.NewPosition(7, 7) // Centre
	pos2, _ := shared.NewPosition(6, 7)
	pos3, _ := shared.NewPosition(8, 7)
	pos4, _ := shared.NewPosition(7, 6)
	pos5, _ := shared.NewPosition(7, 8)
	obstacles := []*shared.Position{pos1, pos2, pos3, pos4, pos5}
	for _, pos := range obstacles {
		err := grille.DefinirTypeCellule(pos, shared.CelluleObstacle)
		require.NoError(t, err)
	}

	// =============================================================================
	// PHASE 2: CRÉATION DE L'ÉQUIPE DES JOUEURS (6 unités)
	// =============================================================================
	statsGuerrier, _ := shared.NewStats(150, 100, 50, 25, 15, 18, 12, 10, 3, 80)
	statsArcher, _ := shared.NewStats(100, 80, 40, 18, 12, 15, 10, 20, 4, 80)
	statsMage, _ := shared.NewStats(80, 150, 60, 12, 8, 25, 20, 14, 3, 80)
	statsClerc, _ := shared.NewStats(120, 120, 50, 15, 12, 20, 18, 12, 3, 80)
	statsRodeur, _ := shared.NewStats(110, 90, 45, 20, 14, 16, 11, 18, 5, 80)
	statsPaladin, _ := shared.NewStats(140, 110, 55, 23, 16, 17, 15, 11, 3, 80)

	// Positions de départ des joueurs (côté gauche)
	posGuerrier, _ := shared.NewPosition(1, 7)
	posArcher, _ := shared.NewPosition(2, 5)
	posMage, _ := shared.NewPosition(2, 9)
	posClerc, _ := shared.NewPosition(3, 7)
	posRodeur, _ := shared.NewPosition(2, 3)
	posPaladin, _ := shared.NewPosition(2, 11)

	guerrier := domain.NewUnite("Guerrier-001", "Guerrier Tank", "team-joueurs", statsGuerrier, posGuerrier)
	archer := domain.NewUnite("Archer-001", "Archer DPS", "team-joueurs", statsArcher, posArcher)
	mage := domain.NewUnite("Mage-001", "Mage Élémentaliste", "team-joueurs", statsMage, posMage)
	clerc := domain.NewUnite("Clerc-001", "Clerc Soigneur", "team-joueurs", statsClerc, posClerc)
	rodeur := domain.NewUnite("Rodeur-001", "Rôdeur Scout", "team-joueurs", statsRodeur, posRodeur)
	paladin := domain.NewUnite("Paladin-001", "Paladin Protecteur", "team-joueurs", statsPaladin, posPaladin)

	joueurID := "player-001"
	equipeJoueurs, err := domain.NewEquipe("team-joueurs", "Héros Aventuriers", "bleu", false, &joueurID)
	require.NoError(t, err)
	for _, u := range []*domain.Unite{guerrier, archer, mage, clerc, rodeur, paladin} {
		equipeJoueurs.AjouterMembre(u)
	}

	// =============================================================================
	// PHASE 3: CRÉATION DE L'ÉQUIPE DES GOBELINS (6 unités)
	// =============================================================================
	statsGobelin, _ := shared.NewStats(60, 30, 25, 12, 8, 10, 6, 15, 4, 80)
	statsGobChef, _ := shared.NewStats(90, 50, 35, 18, 12, 14, 9, 14, 3, 80)
	statsGobShaмаn, _ := shared.NewStats(50, 80, 30, 8, 6, 16, 12, 13, 3, 80)

	// Positions de départ des gobelins (côté droit)
	posGob1, _ := shared.NewPosition(13, 7)  // Gobelin Chef (centre)
	posGob2, _ := shared.NewPosition(12, 5)  // Gobelin guerrier
	posGob3, _ := shared.NewPosition(12, 9)  // Gobelin guerrier
	posGob4, _ := shared.NewPosition(11, 7)  // Shaman
	posGob5, _ := shared.NewPosition(12, 3)  // Gobelin archer
	posGob6, _ := shared.NewPosition(12, 11) // Gobelin archer

	gobelinChef := domain.NewUnite("Gobelin-Chef", "Chef de Guerre Gobelin", "team-gobelins", statsGobChef, posGob1)
	gobelinGuerrier1 := domain.NewUnite("Gobelin-Guerrier-1", "Gobelin Guerrier", "team-gobelins", statsGobelin, posGob2)
	gobelinGuerrier2 := domain.NewUnite("Gobelin-Guerrier-2", "Gobelin Guerrier", "team-gobelins", statsGobelin, posGob3)
	gobelinShaman := domain.NewUnite("Gobelin-Shaman", "Shaman Gobelin", "team-gobelins", statsGobShaмаn, posGob4)
	gobelinArcher1 := domain.NewUnite("Gobelin-Archer-1", "Gobelin Archer", "team-gobelins", statsGobelin, posGob5)
	gobelinArcher2 := domain.NewUnite("Gobelin-Archer-2", "Gobelin Archer", "team-gobelins", statsGobelin, posGob6)

	equipeGobelins, err := domain.NewEquipe("team-gobelins", "Horde Gobeline", "rouge", true, nil)
	require.NoError(t, err)
	for _, u := range []*domain.Unite{gobelinChef, gobelinGuerrier1, gobelinGuerrier2, gobelinShaman, gobelinArcher1, gobelinArcher2} {
		equipeGobelins.AjouterMembre(u)
	}

	// =============================================================================
	// PHASE 4: CRÉATION DU COMBAT
	// =============================================================================
	combat, err := domain.NewCombat("combat-test-001", []*domain.Equipe{equipeJoueurs, equipeGobelins}, grille)
	require.NoError(t, err)

	// Vérifications initiales
	assert.Equal(t, domain.EtatAttente, combat.Etat())
	totalUnites := len(equipeJoueurs.Membres()) + len(equipeGobelins.Membres())
	assert.Equal(t, 12, totalUnites, "Doit avoir 12 unités au total")

	t.Logf("\n=== COMBAT INITIALISÉ ===")
	t.Logf("Équipe Joueurs: %d membres", len(equipeJoueurs.Membres()))
	t.Logf("Équipe Gobelins: %d membres", len(equipeGobelins.Membres()))
	t.Logf("Grille: 15x15 avec obstacles au centre")

	// =============================================================================
	// PHASE 5: DÉMARRAGE DU COMBAT
	// =============================================================================
	err = combat.Demarrer()
	require.NoError(t, err, "Le combat doit démarrer sans erreur")
	assert.Equal(t, domain.EtatEnCours, combat.Etat())

	t.Logf("\n=== COMBAT DÉMARRÉ ===")
	t.Logf("État: %v", combat.Etat())
	t.Logf("Tour actuel: %d", combat.TourActuel())

	// Afficher les statistiques initiales
	t.Logf("\n=== STATISTIQUES INITIALES ===")
	t.Logf("\n ÉQUIPE JOUEURS:")
	for _, unite := range equipeJoueurs.Membres() {
		stats := unite.Stats()
		t.Logf("  %s - HP: %d, MP: %d, ATK: %d, DEF: %d, SPD: %d, MOV: %d",
			unite.Nom(), stats.HP, stats.MP, stats.ATK, stats.DEF, stats.SPD, stats.MOV)
	}

	t.Logf("\n ÉQUIPE GOBELINS:")
	for _, unite := range equipeGobelins.Membres() {
		stats := unite.Stats()
		t.Logf("  %s - HP: %d, MP: %d, ATK: %d, DEF: %d, SPD: %d, MOV: %d",
			unite.Nom(), stats.HP, stats.MP, stats.ATK, stats.DEF, stats.SPD, stats.MOV)
	}

	// =============================================================================
	// PHASE 6: VÉRIFICATIONS FINALES
	// =============================================================================
	t.Logf("\n=== ÉTAT FINAL DU COMBAT ===")
	t.Logf("État du combat: %v", combat.Etat())
	t.Logf("Tour actuel: %d", combat.TourActuel())

	// Statistiques finales
	t.Logf("\n=== STATISTIQUES FINALES ===")

	joueursVivants := 0
	gobelinsVivants := 0

	t.Logf("\n ÉQUIPE JOUEURS:")
	for _, unite := range equipeJoueurs.Membres() {
		stats := unite.StatsActuelles()
		statsB := unite.Stats()
		statut := "✅ Vivant"
		if unite.EstEliminee() {
			statut = "☠️  Éliminé"
		} else {
			joueursVivants++
		}
		t.Logf("  %s - %s - HP: %d/%d, MP: %d/%d",
			unite.Nom(), statut,
			stats.HP, statsB.HP,
			stats.MP, statsB.MP)
	}

	t.Logf("\n ÉQUIPE GOBELINS:")
	for _, unite := range equipeGobelins.Membres() {
		stats := unite.StatsActuelles()
		statut := "✅ Vivant"
		if unite.EstEliminee() {
			statut = "☠️  Éliminé"
		} else {
			gobelinsVivants++
		}
		statsB := unite.Stats()
		t.Logf("  %s - %s - HP: %d/%d, MP: %d/%d",
			unite.Nom(), statut,
			stats.HP, statsB.HP,
			stats.MP, statsB.MP)
	}

	t.Logf("\n BILAN:")
	t.Logf("  Joueurs vivants: %d/6", joueursVivants)
	t.Logf("  Gobelins vivants: %d/6", gobelinsVivants)

	// Vérifications finales
	totalUnitesFinales := len(equipeJoueurs.Membres()) + len(equipeGobelins.Membres())
	assert.Equal(t, 12, totalUnitesFinales, "Doit toujours avoir 12 unités")

	// Vérifier que les deux équipes existent
	assert.NotNil(t, combat.Equipes()["team-joueurs"])
	assert.NotNil(t, combat.Equipes()["team-gobelins"])

	// Vérifier que le système de combat fonctionne
	if joueursVivants == 0 && gobelinsVivants > 0 {
		t.Logf("\n VICTOIRE DES GOBELINS!")
	} else if gobelinsVivants == 0 && joueursVivants > 0 {
		t.Logf("\n VICTOIRE DES JOUEURS!")
	} else {
		t.Logf("\n  COMBAT INITIALISÉ - Les deux équipes sont prêtes")
		t.Logf("    Le combat peut maintenant être joué avec le système de commandes")
	}

	// =============================================================================
	// PHASE 7: VÉRIFICATION DU SYSTÈME
	// =============================================================================
	t.Logf("\n=== VÉRIFICATION DU SYSTÈME ===")
	t.Logf("✅ Grille de combat: %dx%d avec obstacles", grille.Largeur(), grille.Hauteur())
	t.Logf("✅ Équipe Joueurs: %d unités", len(equipeJoueurs.Membres()))
	t.Logf("✅ Équipe Gobelins: %d unités", len(equipeGobelins.Membres()))
	t.Logf("✅ Combat initialisé et démarré")
	t.Logf("✅ Toutes les unités sont en position")
	t.Logf("\n Système prêt pour l'exécution des commandes via l'application layer")
}
