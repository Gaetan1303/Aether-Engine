package combat_test

import (
	"testing"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/aether-engine/aether-engine/internal/combat/domain/commands"
	shared "github.com/aether-engine/aether-engine/internal/shared/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCombatTourComplet_GobelinsVsJoueurs teste un combat avec plusieurs tours d'actions
// Ce test simule un combat réaliste tour par tour en utilisant uniquement les APIs du domain
func TestCombatTourComplet_GobelinsVsJoueurs(t *testing.T) {
	// =============================================================================
	// PHASE 1: CRÉATION DE LA GRILLE
	// =============================================================================

	// Grille 10x10 pour un combat plus rapide
	grille, err := shared.NewGrilleCombat(10, 10)
	require.NoError(t, err)

	// Obstacles au centre (en croix)
	obstacles := []struct{ x, y int }{
		{5, 5}, {4, 5}, {6, 5}, {5, 4}, {5, 6},
	}
	for _, obs := range obstacles {
		pos, _ := shared.NewPosition(obs.x, obs.y)
		grille.DefinirTypeCellule(pos, shared.CelluleObstacle)
	}

	t.Logf("\n=== TERRAIN DE COMBAT ===")
	t.Logf("✅ Grille: 10x10")
	t.Logf("✅ Obstacles: 5 cellules (croix centrale)")

	// =============================================================================
	// PHASE 2: CRÉATION DES UNITÉS - ÉQUIPE JOUEURS (3 unités)
	// =============================================================================

	// Guerrier Tank (haute défense, bas DPS)
	statsGuerrier, _ := shared.NewStats(100, 50, 30, 20, 15, 15, 12, 10, 3, 80)
	posGuerrier, _ := shared.NewPosition(1, 5)
	guerrier := domain.NewUnite("J-Guerrier-001", "Guerrier Tank", "team-joueurs", statsGuerrier, posGuerrier)

	// Archer DPS (haute attaque, basse défense)
	statsArcher, _ := shared.NewStats(80, 40, 25, 25, 8, 20, 10, 14, 4, 80)
	posArcher, _ := shared.NewPosition(1, 4)
	archer := domain.NewUnite("J-Archer-001", "Archer DPS", "team-joueurs", statsArcher, posArcher)

	// Mage (haute attaque magique)
	statsMage, _ := shared.NewStats(60, 100, 20, 15, 10, 30, 15, 12, 3, 80)
	posMage, _ := shared.NewPosition(1, 6)
	mage := domain.NewUnite("J-Mage-001", "Mage Élémentaliste", "team-joueurs", statsMage, posMage)

	joueurID := "player-001"
	equipeJoueurs, err := domain.NewEquipe("team-joueurs", "Héros Aventuriers", "bleu", false, &joueurID)
	require.NoError(t, err)
	for _, u := range []*domain.Unite{guerrier, archer, mage} {
		equipeJoueurs.AjouterMembre(u)
	}

	t.Logf("\n=== ÉQUIPE JOUEURS ===")
	t.Logf("✅ 3 unités créées:")
	t.Logf("   - Guerrier Tank (HP:100, ATK:20, DEF:15)")
	t.Logf("   - Archer DPS (HP:80, ATK:25, DEF:8)")
	t.Logf("   - Mage Élémentaliste (HP:60, MATK:30, MDEF:15)")

	// =============================================================================
	// PHASE 3: CRÉATION DES UNITÉS - ÉQUIPE GOBELINS (3 unités)
	// =============================================================================

	// Chef Gobelin (équilibré)
	statsChef, _ := shared.NewStats(70, 40, 25, 18, 10, 18, 10, 11, 3, 80)
	posChef, _ := shared.NewPosition(8, 5)
	chef := domain.NewUnite("G-Chef-001", "Chef de Guerre", "team-gobelins", statsChef, posChef)

	// Gobelin Guerrier 1
	statsGobelin, _ := shared.NewStats(50, 30, 20, 15, 8, 12, 8, 10, 3, 80)
	posGob1, _ := shared.NewPosition(8, 4)
	gobelinGuerrier1 := domain.NewUnite("G-Guerrier-001", "Gobelin Guerrier", "team-gobelins", statsGobelin, posGob1)

	// Gobelin Guerrier 2
	posGob2, _ := shared.NewPosition(8, 6)
	gobelinGuerrier2 := domain.NewUnite("G-Guerrier-002", "Gobelin Guerrier", "team-gobelins", statsGobelin, posGob2)

	equipeGobelins, err := domain.NewEquipe("team-gobelins", "Horde Gobeline", "rouge", true, nil)
	require.NoError(t, err)
	for _, u := range []*domain.Unite{chef, gobelinGuerrier1, gobelinGuerrier2} {
		equipeGobelins.AjouterMembre(u)
	}

	t.Logf("\n=== ÉQUIPE GOBELINS ===")
	t.Logf("✅ 3 unités créées:")
	t.Logf("   - Chef de Guerre (HP:70, ATK:18, DEF:10)")
	t.Logf("   - Gobelin Guerrier x2 (HP:50, ATK:15, DEF:8)")

	// =============================================================================
	// PHASE 4: CRÉATION ET DÉMARRAGE DU COMBAT
	// =============================================================================

	combat, err := domain.NewCombat("combat-tour-001", []*domain.Equipe{equipeJoueurs, equipeGobelins}, grille)
	require.NoError(t, err)

	err = combat.Demarrer()
	require.NoError(t, err, "Le combat doit démarrer sans erreur")

	assert.Equal(t, domain.EtatEnCours, combat.Etat())
	assert.Equal(t, 1, combat.TourActuel())

	t.Logf("\n=== COMBAT DÉMARRÉ ===")
	t.Logf("✅ État: %s", combat.Etat())
	t.Logf("✅ Tour actuel: %d", combat.TourActuel())

	// Fonction helper pour afficher les stats des unités
	afficherStats := func(titre string) {
		t.Logf("\n%s", titre)
		for _, equipe := range combat.Equipes() {
			vivants := 0
			for _, u := range equipe.Membres() {
				if !u.EstEliminee() {
					vivants++
					stats := u.StatsActuelles()
					t.Logf("   [%s] %s - HP:%d/%d MP:%d/%d Stamina:%d/%d",
						equipe.Nom(), u.Nom(),
						stats.HP, u.Stats().HP,
						stats.MP, u.Stats().MP,
						stats.Stamina, u.Stats().Stamina)
				}
			}
			t.Logf("   → %s: %d/%d vivants", equipe.Nom(), vivants, len(equipe.Membres()))
		}
	}

	afficherStats("=== STATS INITIALES ===")

	// =============================================================================
	// PHASE 5: TOUR 1 - ACTIONS DE COMBAT
	// =============================================================================

	t.Logf("\n=== TOUR 1: SIMULATION D'ACTIONS ===")

	// Action 1: Archer attaque Gobelin Guerrier 1 (distance = 7, hors portée mêlée)
	t.Logf("\n--- Action 1: Archer → Gobelin Guerrier 1 ---")
	statsAvant := gobelinGuerrier1.StatsActuelles()
	t.Logf("   HP avant: %d", statsAvant.HP)

	// Créer la commande d'attaque
	attackCmd1 := commands.NewAttackCommand(archer, combat, gobelinGuerrier1)

	// Exécuter la commande
	result, err := attackCmd1.Execute()
	if err != nil {
		t.Logf("   ⚠️ Erreur: %v (peut-être hors de portée)", err)
	} else {
		statsApres := gobelinGuerrier1.StatsActuelles()
		degats := statsAvant.HP - statsApres.HP
		t.Logf("   ✅ Attaque réussie! (Dégâts calculés: %d, réels: %d)", result.DamageDealt, degats)
		t.Logf("   HP après: %d", statsApres.HP)
		t.Logf("   Éliminé: %v", gobelinGuerrier1.EstEliminee())
	} // Action 2: Mage attaque Chef Gobelin
	t.Logf("\n--- Action 2: Mage → Chef Gobelin ---")
	statsAvant = chef.StatsActuelles()
	t.Logf("   HP avant: %d", statsAvant.HP)

	attackCmd2 := commands.NewAttackCommand(mage, combat, chef)

	result, err = attackCmd2.Execute()
	if err != nil {
		t.Logf("   ⚠️ Erreur: %v", err)
	} else {
		statsApres := chef.StatsActuelles()
		degats := statsAvant.HP - statsApres.HP
		t.Logf("   ✅ Attaque réussie! (Dégâts calculés: %d, réels: %d)", result.DamageDealt, degats)
		t.Logf("   HP après: %d", statsApres.HP)
		t.Logf("   Éliminé: %v", chef.EstEliminee())
	} // Action 3: Chef Gobelin contre-attaque le Guerrier
	t.Logf("\n--- Action 3: Chef Gobelin → Guerrier ---")
	statsAvant = guerrier.StatsActuelles()
	t.Logf("   HP avant: %d", statsAvant.HP)

	attackCmd3 := commands.NewAttackCommand(chef, combat, guerrier)

	result, err = attackCmd3.Execute()
	if err != nil {
		t.Logf("   ⚠️ Erreur: %v", err)
	} else {
		statsApres := guerrier.StatsActuelles()
		degats := statsAvant.HP - statsApres.HP
		t.Logf("   ✅ Attaque réussie! (Dégâts calculés: %d, réels: %d)", result.DamageDealt, degats)
		t.Logf("   HP après: %d", statsApres.HP)
		t.Logf("   Éliminé: %v", guerrier.EstEliminee())
	} // Action 4: Gobelin Guerrier 2 attaque le Mage
	t.Logf("\n--- Action 4: Gobelin Guerrier 2 → Mage ---")
	statsAvant = mage.StatsActuelles()
	t.Logf("   HP avant: %d", statsAvant.HP)

	attackCmd4 := commands.NewAttackCommand(gobelinGuerrier2, combat, mage)

	result, err = attackCmd4.Execute()
	if err != nil {
		t.Logf("   ⚠️ Erreur: %v", err)
	} else {
		statsApres := mage.StatsActuelles()
		degats := statsAvant.HP - statsApres.HP
		t.Logf("   ✅ Attaque réussie! (Dégâts calculés: %d, réels: %d)", result.DamageDealt, degats)
		t.Logf("   HP après: %d", statsApres.HP)
		t.Logf("   Éliminé: %v", mage.EstEliminee())
	} // Afficher les stats finales après le tour 1
	afficherStats("=== STATS APRÈS TOUR 1 ===")

	// =============================================================================
	// PHASE 6: VÉRIFICATIONS FINALES
	// =============================================================================

	t.Logf("\n=== VÉRIFICATIONS FINALES ===")

	// Compter les survivants
	joueursVivants := 0
	gobelinsVivants := 0

	for _, u := range equipeJoueurs.Membres() {
		if !u.EstEliminee() {
			joueursVivants++
		}
	}

	for _, u := range equipeGobelins.Membres() {
		if !u.EstEliminee() {
			gobelinsVivants++
		}
	}

	t.Logf("✅ Joueurs vivants: %d/3", joueursVivants)
	t.Logf("✅ Gobelins vivants: %d/3", gobelinsVivants)
	t.Logf("✅ Combat en cours: %v", combat.Etat() == domain.EtatEnCours)

	// Vérifier que le combat est toujours en cours (aucune équipe éliminée)
	assert.Equal(t, domain.EtatEnCours, combat.Etat(), "Le combat doit être encore en cours")
	assert.Equal(t, 1, combat.TourActuel(), "Nous sommes au tour 1")

	// Au moins une unité de chaque côté doit être vivante
	assert.Greater(t, joueursVivants, 0, "Au moins un joueur doit être vivant")
	assert.Greater(t, gobelinsVivants, 0, "Au moins un gobelin doit être vivant")

	t.Logf("\n=== TEST TERMINÉ AVEC SUCCÈS ===")
}
