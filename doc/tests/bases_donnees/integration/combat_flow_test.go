package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCombatCompleteFlow teste un scénario de combat complet de bout en bout
func TestCombatCompleteFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup
	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	redisClient := NewTestRedis(t)
	CleanRedis(t, redisClient)

	ctx := context.Background()

	// Créer deux joueurs
	joueur1 := NewJoueurAggregate()
	joueur2 := NewJoueurAggregate()

	// Insérer les événements de création des joueurs
	eventJoueur1 := createJoueurEvent(joueur1)
	eventJoueur2 := createJoueurEvent(joueur2)

	_, err := insertEvent(ctx, db, eventJoueur1)
	require.NoError(t, err)
	_, err = insertEvent(ctx, db, eventJoueur2)
	require.NoError(t, err)

	// Créer les projections joueurs
	err = createPlayerProjection(ctx, db, joueur1)
	require.NoError(t, err)
	err = createPlayerProjection(ctx, db, joueur2)
	require.NoError(t, err)

	t.Run("1. Démarrage du combat", func(t *testing.T) {
		// Créer un événement CombatDemarre
		combatID := uuid.New()
		eventCombatDemarre := createCombatDemarreEvent(combatID, joueur1, joueur2)
		eventSeq, err := insertEvent(ctx, db, eventCombatDemarre)
		require.NoError(t, err)

		// Créer la projection combat
		err = createCombatProjection(ctx, db, combatID, joueur1, joueur2, eventSeq)
		require.NoError(t, err)

		// Créer le cache Redis pour le combat
		err = createCombatCache(ctx, redisClient, combatID, joueur1, joueur2)
		require.NoError(t, err)

		// Vérifier la projection
		var combat struct {
			CombatID   uuid.UUID
			TypeCombat string
			Etat       string
			TourActuel int
		}

		query := `SELECT combat_id, type_combat, etat, tour_actuel 
		          FROM instances_combat WHERE combat_id = $1`

		err = db.QueryRow(ctx, query, combatID).Scan(
			&combat.CombatID,
			&combat.TypeCombat,
			&combat.Etat,
			&combat.TourActuel,
		)

		require.NoError(t, err)
		assert.Equal(t, combatID, combat.CombatID)
		assert.Equal(t, "PVP", combat.TypeCombat)
		assert.Equal(t, "EN_COURS", combat.Etat)
		assert.Equal(t, 1, combat.TourActuel)

		// Vérifier le cache Redis
		combatState, err := redisClient.HGetAll(ctx, fmt.Sprintf("combat:%s:state", combatID)).Result()
		require.NoError(t, err)
		assert.Equal(t, "EN_COURS", combatState["etat"])
		assert.Equal(t, "1", combatState["tour_actuel"])

		// Vérifier les participants dans Redis
		participants, err := redisClient.SMembers(ctx, fmt.Sprintf("combat:%s:participants", combatID)).Result()
		require.NoError(t, err)
		assert.Len(t, participants, 2)

		// Vérifier l'ordre des tours
		turnOrder, err := redisClient.ZRevRange(ctx, fmt.Sprintf("combat:%s:turn_order", combatID), 0, -1).Result()
		require.NoError(t, err)
		assert.Len(t, turnOrder, 2)
		assert.Equal(t, joueur1.ID.String(), turnOrder[0]) // Plus haute initiative

		// Stocker combatID pour les tests suivants
		t.Cleanup(func() {
			// Nettoyage du cache à la fin du test
			redisClient.Del(ctx, fmt.Sprintf("combat:%s:*", combatID))
		})
	})

	t.Run("2. Exécution d'une action de combat", func(t *testing.T) {
		combatID := getCombatID(t, db)

		// Acquérir un verrou distribué pour l'action
		lockKey := fmt.Sprintf("combat:%s:action", combatID)
		sessionID := "session-test-1"
		acquired, err := redisClient.SetNX(ctx, lockKey, sessionID, 5*time.Second).Result()
		require.NoError(t, err)
		assert.True(t, acquired)

		// Créer un événement ActionCombatExecutee
		eventAction := createActionCombatEvent(combatID, joueur1.ID, joueur2.ID, 2)
		eventSeq, err := insertEvent(ctx, db, eventAction)
		require.NoError(t, err)

		// Mettre à jour la projection
		err = updateCombatTurn(ctx, db, combatID, 2, eventSeq)
		require.NoError(t, err)

		// Mettre à jour le cache Redis
		err = redisClient.HSet(ctx, fmt.Sprintf("combat:%s:state", combatID), "tour_actuel", 2).Err()
		require.NoError(t, err)

		// Publier une notification
		notification := map[string]interface{}{
			"type":      "ACTION_EXECUTEE",
			"combat_id": combatID.String(),
			"acteur_id": joueur1.ID.String(),
			"cible_id":  joueur2.ID.String(),
			"tour":      2,
		}
		notifJSON, _ := json.Marshal(notification)
		err = redisClient.Publish(ctx, fmt.Sprintf("combat:%s:notifications", combatID), notifJSON).Err()
		require.NoError(t, err)

		// Libérer le verrou
		releaseLockScript := `
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`
		result, err := redisClient.Eval(ctx, releaseLockScript, []string{lockKey}, sessionID).Result()
		require.NoError(t, err)
		assert.Equal(t, int64(1), result)

		// Vérifier la mise à jour du tour
		var tourActuel int
		err = db.QueryRow(ctx, "SELECT tour_actuel FROM instances_combat WHERE combat_id = $1", combatID).Scan(&tourActuel)
		require.NoError(t, err)
		assert.Equal(t, 2, tourActuel)
	})

	t.Run("3. Application de dégâts", func(t *testing.T) {
		combatID := getCombatID(t, db)

		// Récupérer les HP actuels de joueur2
		var hpAvant int
		query := `SELECT hp_actuel FROM participants_combat WHERE combat_id = $1 AND joueur_id = $2`
		err := db.QueryRow(ctx, query, combatID, joueur2.ID).Scan(&hpAvant)
		require.NoError(t, err)

		// Créer un événement DegatsInfliges
		degats := 25
		eventDegats := createDegatsEvent(combatID, joueur2.ID, degats, 3, hpAvant)
		_, err = insertEvent(ctx, db, eventDegats)
		require.NoError(t, err)

		// Mettre à jour la projection
		hpApres := hpAvant - degats
		err = updateParticipantHP(ctx, db, combatID, joueur2.ID, hpApres)
		require.NoError(t, err)

		// Mettre à jour le cache Redis
		keyParticipant := fmt.Sprintf("combat:%s:participant:%s", combatID, joueur2.ID)
		err = redisClient.HSet(ctx, keyParticipant, "hp_actuel", hpApres).Err()
		require.NoError(t, err)

		// Vérifier la mise à jour des HP
		var hpActuel int
		err = db.QueryRow(ctx, query, combatID, joueur2.ID).Scan(&hpActuel)
		require.NoError(t, err)
		assert.Equal(t, hpApres, hpActuel)
		assert.Less(t, hpActuel, hpAvant)

		// Vérifier le cache Redis
		hpRedis, err := redisClient.HGet(ctx, keyParticipant, "hp_actuel").Result()
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%d", hpApres), hpRedis)

		t.Logf("Dégâts appliqués: %d HP (avant: %d, après: %d)", degats, hpAvant, hpActuel)
	})

	t.Run("4. Application d'un effet de statut", func(t *testing.T) {
		combatID := getCombatID(t, db)

		// Créer la table effets_statut si elle n'existe pas
		createEffetsStatutTable(t, db)

		// Créer un événement EffetStatutApplique (Poison)
		eventEffet := createEffetStatutEvent(combatID, joueur2.ID, 4)
		eventSeq, err := insertEvent(ctx, db, eventEffet)
		require.NoError(t, err)

		// Créer la projection de l'effet
		err = createEffetStatutProjection(ctx, db, combatID, joueur2.ID, "POISON", 5, 3, eventSeq)
		require.NoError(t, err)

		// Ajouter l'effet au cache Redis
		keyEffets := fmt.Sprintf("combat:%s:participant:%s:effets", combatID, joueur2.ID)
		effetData := map[string]interface{}{
			"type":           "POISON",
			"puissance":      5,
			"tours_restants": 3,
		}
		effetJSON, _ := json.Marshal(effetData)
		err = redisClient.SAdd(ctx, keyEffets, effetJSON).Err()
		require.NoError(t, err)

		// Vérifier l'effet dans la projection
		var typeEffet string
		var toursRestants int
		query := `SELECT type_effet, tours_restants FROM effets_statut 
		          WHERE combat_id = $1 AND joueur_id = $2`
		err = db.QueryRow(ctx, query, combatID, joueur2.ID).Scan(&typeEffet, &toursRestants)
		require.NoError(t, err)
		assert.Equal(t, "POISON", typeEffet)
		assert.Equal(t, 3, toursRestants)

		t.Logf("Effet POISON appliqué: %d dégâts par tour pendant %d tours", 5, 3)
	})

	t.Run("5. Plusieurs tours avec dégâts de poison", func(t *testing.T) {
		combatID := getCombatID(t, db)

		// Simuler 3 tours avec dégâts de poison
		for tour := 1; tour <= 3; tour++ {
			// Récupérer les HP actuels
			var hpAvant int
			query := `SELECT hp_actuel FROM participants_combat WHERE combat_id = $1 AND joueur_id = $2`
			err := db.QueryRow(ctx, query, combatID, joueur2.ID).Scan(&hpAvant)
			require.NoError(t, err)

			// Appliquer les dégâts de poison
			degatsPoison := 5
			hpApres := hpAvant - degatsPoison

			// Créer un événement pour les dégâts de poison
			eventVersion := 4 + tour
			eventDegats := createDegatsEvent(combatID, joueur2.ID, degatsPoison, eventVersion, hpAvant)
			_, err = insertEvent(ctx, db, eventDegats)
			require.NoError(t, err)

			// Mettre à jour les HP
			err = updateParticipantHP(ctx, db, combatID, joueur2.ID, hpApres)
			require.NoError(t, err)

			// Décrémenter les tours restants de l'effet
			query = `UPDATE effets_statut SET tours_restants = tours_restants - 1 
			         WHERE combat_id = $1 AND joueur_id = $2 AND type_effet = 'POISON'`
			_, err = db.Exec(ctx, query, combatID, joueur2.ID)
			require.NoError(t, err)

			t.Logf("Tour %d: Dégâts de poison appliqués (%d HP), HP restants: %d", tour, degatsPoison, hpApres)
		}

		// Vérifier que l'effet a expiré
		var toursRestants int
		query := `SELECT tours_restants FROM effets_statut 
		          WHERE combat_id = $1 AND joueur_id = $2 AND type_effet = 'POISON'`
		err := db.QueryRow(ctx, query, combatID, joueur2.ID).Scan(&toursRestants)
		require.NoError(t, err)
		assert.Equal(t, 0, toursRestants)

		// Supprimer l'effet expiré
		_, err = db.Exec(ctx, `DELETE FROM effets_statut 
		                       WHERE combat_id = $1 AND joueur_id = $2 AND tours_restants <= 0`,
			combatID, joueur2.ID)
		require.NoError(t, err)
	})

	t.Run("6. Utilisation d'une potion de soin", func(t *testing.T) {
		combatID := getCombatID(t, db)

		// Récupérer les HP actuels et max
		var hpAvant, hpMax int
		query := `SELECT hp_actuel, hp_max FROM participants_combat 
		          WHERE combat_id = $1 AND joueur_id = $2`
		err := db.QueryRow(ctx, query, combatID, joueur2.ID).Scan(&hpAvant, &hpMax)
		require.NoError(t, err)

		// Appliquer les soins (50 HP)
		soins := 50
		hpApres := hpAvant + soins
		if hpApres > hpMax {
			hpApres = hpMax
		}

		// Créer un événement SoinsRecus
		eventSoins := createSoinsEvent(combatID, joueur2.ID, soins, 8, hpAvant, hpApres)
		_, err = insertEvent(ctx, db, eventSoins)
		require.NoError(t, err)

		// Mettre à jour les HP
		err = updateParticipantHP(ctx, db, combatID, joueur2.ID, hpApres)
		require.NoError(t, err)

		// Vérifier les HP après soins
		var hpFinal int
		err = db.QueryRow(ctx, `SELECT hp_actuel FROM participants_combat 
		                        WHERE combat_id = $1 AND joueur_id = $2`,
			combatID, joueur2.ID).Scan(&hpFinal)
		require.NoError(t, err)
		assert.Greater(t, hpFinal, hpAvant)
		assert.LessOrEqual(t, hpFinal, hpMax)

		t.Logf("Soins appliqués: %d HP (avant: %d, après: %d, max: %d)", soins, hpAvant, hpFinal, hpMax)
	})

	t.Run("7. Fin du combat", func(t *testing.T) {
		combatID := getCombatID(t, db)

		// Joueur1 gagne le combat
		eventTermine := createCombatTermineEvent(combatID, joueur1.ID, 9)
		eventSeq, err := insertEvent(ctx, db, eventTermine)
		require.NoError(t, err)

		// Mettre à jour la projection
		err = terminateCombat(ctx, db, combatID, eventSeq)
		require.NoError(t, err)

		// Distribuer les récompenses (XP, Or)
		xpRecompense := 100
		orRecompense := 50

		eventXP := createExperienceEvent(joueur1.ID, 10, xpRecompense)
		_, err = insertEvent(ctx, db, eventXP)
		require.NoError(t, err)

		// Mettre à jour le joueur (simplifié)
		_, err = db.Exec(ctx, `UPDATE joueurs SET experience_actuelle = experience_actuelle + $1 
		                       WHERE joueur_id = $2`, xpRecompense, joueur1.ID)
		require.NoError(t, err)

		_, err = db.Exec(ctx, `UPDATE joueurs SET or = or + $1 WHERE joueur_id = $2`,
			orRecompense, joueur1.ID)
		require.NoError(t, err)

		// Vérifier l'état du combat
		var etat string
		err = db.QueryRow(ctx, "SELECT etat FROM instances_combat WHERE combat_id = $1", combatID).Scan(&etat)
		require.NoError(t, err)
		assert.Equal(t, "TERMINE", etat)

		// Nettoyer le cache Redis
		pattern := fmt.Sprintf("combat:%s:*", combatID)
		iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			err := redisClient.Del(ctx, iter.Val()).Err()
			require.NoError(t, err)
		}
		require.NoError(t, iter.Err())

		// Vérifier que le cache est nettoyé
		AssertRedisKeyNotExists(t, redisClient, fmt.Sprintf("combat:%s:state", combatID))

		t.Logf("Combat terminé. Vainqueur: %s (récompenses: %d XP, %d Or)", joueur1.ID, xpRecompense, orRecompense)
	})

	t.Run("8. Vérification finale de la cohérence", func(t *testing.T) {
		// Vérifier que tous les événements ont été stockés
		var countEvents int64
		err := db.QueryRow(ctx, "SELECT COUNT(*) FROM evenements").Scan(&countEvents)
		require.NoError(t, err)
		assert.Greater(t, countEvents, int64(10)) // Au moins 10 événements

		// Vérifier les projections joueurs
		var xpJoueur1 int64
		err = db.QueryRow(ctx, "SELECT experience_actuelle FROM joueurs WHERE joueur_id = $1", joueur1.ID).Scan(&xpJoueur1)
		require.NoError(t, err)
		assert.Greater(t, xpJoueur1, int64(0))

		// Vérifier qu'il n'y a plus de combats en cours
		var countCombatsEnCours int64
		err = db.QueryRow(ctx, "SELECT COUNT(*) FROM instances_combat WHERE etat = 'EN_COURS'").Scan(&countCombatsEnCours)
		require.NoError(t, err)
		assert.Equal(t, int64(0), countCombatsEnCours)

		t.Log("Vérification finale: Toutes les projections sont cohérentes avec l'Event Store")
	})
}

// --- Helpers pour le test d'intégration ---

func getCombatID(t *testing.T, db *pgxpool.Pool) uuid.UUID {
	var combatID uuid.UUID
	err := db.QueryRow(context.Background(), "SELECT combat_id FROM instances_combat LIMIT 1").Scan(&combatID)
	require.NoError(t, err)
	return combatID
}

func createJoueurEvent(joueur JoueurAggregate) Event {
	return NewEventBuilder().
		WithAggregateType("Joueur").
		WithAggregateID(joueur.ID).
		WithAggregateVersion(1).
		WithEventType("JoueurCree").
		WithEventData(map[string]interface{}{
			"username": joueur.Username,
			"niveau":   joueur.Niveau,
			"hp_max":   joueur.HPMax,
		}).
		Build()
}

func createCombatDemarreEvent(combatID uuid.UUID, j1, j2 JoueurAggregate) Event {
	return NewEventBuilder().
		WithAggregateType("Combat").
		WithAggregateID(combatID).
		WithAggregateVersion(1).
		WithEventType("CombatDemarre").
		WithEventData(map[string]interface{}{
			"type_combat": "PVP",
			"participants": []map[string]interface{}{
				{"joueur_id": j1.ID.String(), "initiative": 15, "hp_max": j1.HPMax},
				{"joueur_id": j2.ID.String(), "initiative": 12, "hp_max": j2.HPMax},
			},
		}).
		Build()
}

func createActionCombatEvent(combatID, acteurID, cibleID uuid.UUID, version int) Event {
	return NewEventBuilder().
		WithAggregateType("Combat").
		WithAggregateID(combatID).
		WithAggregateVersion(version).
		WithEventType("ActionCombatExecutee").
		WithEventData(map[string]interface{}{
			"acteur_id": acteurID.String(),
			"cible_id":  cibleID.String(),
			"tour":      version,
		}).
		Build()
}

func createDegatsEvent(combatID, joueurID uuid.UUID, degats, version, hpAvant int) Event {
	return NewEventBuilder().
		WithAggregateType("Combat").
		WithAggregateID(combatID).
		WithAggregateVersion(version).
		WithEventType("DegatsInfliges").
		WithEventData(map[string]interface{}{
			"joueur_id": joueurID.String(),
			"degats":    degats,
			"hp_avant":  hpAvant,
			"hp_apres":  hpAvant - degats,
		}).
		Build()
}

func createEffetStatutEvent(combatID, joueurID uuid.UUID, version int) Event {
	return NewEventBuilder().
		WithAggregateType("Combat").
		WithAggregateID(combatID).
		WithAggregateVersion(version).
		WithEventType("EffetStatutApplique").
		WithEventData(map[string]interface{}{
			"joueur_id":   joueurID.String(),
			"type_effet":  "POISON",
			"puissance":   5,
			"duree_tours": 3,
		}).
		Build()
}

func createSoinsEvent(combatID, joueurID uuid.UUID, soins, version, hpAvant, hpApres int) Event {
	return NewEventBuilder().
		WithAggregateType("Combat").
		WithAggregateID(combatID).
		WithAggregateVersion(version).
		WithEventType("SoinsRecus").
		WithEventData(map[string]interface{}{
			"joueur_id": joueurID.String(),
			"soins":     soins,
			"hp_avant":  hpAvant,
			"hp_apres":  hpApres,
		}).
		Build()
}

func createCombatTermineEvent(combatID, vainqueurID uuid.UUID, version int) Event {
	return NewEventBuilder().
		WithAggregateType("Combat").
		WithAggregateID(combatID).
		WithAggregateVersion(version).
		WithEventType("CombatTermine").
		WithEventData(map[string]interface{}{
			"vainqueur_id": vainqueurID.String(),
		}).
		Build()
}

func createExperienceEvent(joueurID uuid.UUID, version, xp int) Event {
	return NewEventBuilder().
		WithAggregateType("Joueur").
		WithAggregateID(joueurID).
		WithAggregateVersion(version).
		WithEventType("ExperienceGagnee").
		WithEventData(map[string]interface{}{
			"experience_gagnee": xp,
		}).
		Build()
}

// Suite des helpers dans le prochain message...
