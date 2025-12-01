package bases_donnees

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCombatProjectionCreation vérifie la création d'une projection combat
func TestCombatProjectionCreation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()

	// Insérer un événement CombatDemarre
	event := NewCombatDemarreEvent()
	eventSequence, err := insertEvent(ctx, db, event)
	require.NoError(t, err)

	// Appliquer le handler de projection
	err = handleCombatDemarre(ctx, db, event, eventSequence)
	require.NoError(t, err)

	// Vérifier la projection instances_combat
	var combat struct {
		CombatID      uuid.UUID
		TypeCombat    string
		Etat          string
		TourActuel    int
		EventSequence int64
	}

	query := `SELECT combat_id, type_combat, etat, tour_actuel, event_sequence 
	          FROM instances_combat WHERE combat_id = $1`

	err = db.QueryRow(ctx, query, event.AggregateID).Scan(
		&combat.CombatID,
		&combat.TypeCombat,
		&combat.Etat,
		&combat.TourActuel,
		&combat.EventSequence,
	)

	require.NoError(t, err)
	assert.Equal(t, event.AggregateID, combat.CombatID)
	assert.Equal(t, "PVP", combat.TypeCombat)
	assert.Equal(t, "EN_COURS", combat.Etat)
	assert.Equal(t, 1, combat.TourActuel)
	assert.Equal(t, eventSequence, combat.EventSequence)

	// Vérifier la projection participants_combat
	var countParticipants int64
	query = `SELECT COUNT(*) FROM participants_combat WHERE combat_id = $1`
	err = db.QueryRow(ctx, query, event.AggregateID).Scan(&countParticipants)
	require.NoError(t, err)
	assert.Equal(t, int64(2), countParticipants)
}

// TestCombatActionProjection vérifie la projection d'une action de combat
func TestCombatActionProjection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()

	// Setup: Créer un combat
	eventDemarre := NewCombatDemarreEvent()
	seq1, err := insertEvent(ctx, db, eventDemarre)
	require.NoError(t, err)

	err = handleCombatDemarre(ctx, db, eventDemarre, seq1)
	require.NoError(t, err)

	// Extraire les IDs des participants
	participants := eventDemarre.EventData["participants"].([]map[string]interface{})
	acteurID, _ := uuid.Parse(participants[0]["joueur_id"].(string))
	cibleID, _ := uuid.Parse(participants[1]["joueur_id"].(string))

	// Insérer un événement ActionCombatExecutee
	eventAction := NewActionCombatExecuteeEvent(eventDemarre.AggregateID, 2, acteurID, cibleID)
	seq2, err := insertEvent(ctx, db, eventAction)
	require.NoError(t, err)

	err = handleActionCombatExecutee(ctx, db, eventAction, seq2)
	require.NoError(t, err)

	// Vérifier que le tour a été incrémenté
	var tourActuel int
	query := `SELECT tour_actuel FROM instances_combat WHERE combat_id = $1`
	err = db.QueryRow(ctx, query, eventDemarre.AggregateID).Scan(&tourActuel)
	require.NoError(t, err)
	assert.Equal(t, 2, tourActuel)
}

// TestCombatDegatsProjection vérifie la mise à jour des HP après des dégâts
func TestCombatDegatsProjection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()

	// Setup: Créer un combat
	eventDemarre := NewCombatDemarreEvent()
	seq1, err := insertEvent(ctx, db, eventDemarre)
	require.NoError(t, err)

	err = handleCombatDemarre(ctx, db, eventDemarre, seq1)
	require.NoError(t, err)

	// Extraire un participant
	participants := eventDemarre.EventData["participants"].([]map[string]interface{})
	joueurID, _ := uuid.Parse(participants[0]["joueur_id"].(string))
	hpMax := int(participants[0]["hp_max"].(int))

	// Insérer un événement DegatsInfliges
	eventDegats := NewDegatsInfligesEvent(eventDemarre.AggregateID, 2, joueurID, 25)
	seq2, err := insertEvent(ctx, db, eventDegats)
	require.NoError(t, err)

	err = handleDegatsInfliges(ctx, db, eventDegats, seq2)
	require.NoError(t, err)

	// Vérifier la mise à jour des HP
	var hpActuel int
	query := `SELECT hp_actuel FROM participants_combat 
	          WHERE combat_id = $1 AND joueur_id = $2`

	err = db.QueryRow(ctx, query, eventDemarre.AggregateID, joueurID).Scan(&hpActuel)
	require.NoError(t, err)
	assert.Equal(t, hpMax-25, hpActuel)
}

// TestCombatEffetStatutProjection vérifie l'application d'un effet de statut
func TestCombatEffetStatutProjection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()

	// Setup: Créer un combat
	eventDemarre := NewCombatDemarreEvent()
	seq1, err := insertEvent(ctx, db, eventDemarre)
	require.NoError(t, err)

	err = handleCombatDemarre(ctx, db, eventDemarre, seq1)
	require.NoError(t, err)

	// Extraire un participant
	participants := eventDemarre.EventData["participants"].([]map[string]interface{})
	joueurID, _ := uuid.Parse(participants[0]["joueur_id"].(string))

	// Insérer un événement EffetStatutApplique
	eventEffet := NewEffetStatutAppliqueEvent(eventDemarre.AggregateID, 2, joueurID)
	seq2, err := insertEvent(ctx, db, eventEffet)
	require.NoError(t, err)

	// Note: Le handler devrait créer une table effets_statut
	// Pour ce test, on simule la création
	query := `CREATE TABLE IF NOT EXISTS effets_statut (
		combat_id UUID NOT NULL,
		joueur_id UUID NOT NULL,
		type_effet VARCHAR(50) NOT NULL,
		puissance INT NOT NULL,
		tours_restants INT NOT NULL,
		event_sequence BIGINT NOT NULL
	)`
	_, err = db.Exec(ctx, query)
	require.NoError(t, err)

	err = handleEffetStatutApplique(ctx, db, eventEffet, seq2)
	require.NoError(t, err)

	// Vérifier l'effet de statut
	var typeEffet string
	var toursRestants int

	query = `SELECT type_effet, tours_restants FROM effets_statut 
	         WHERE combat_id = $1 AND joueur_id = $2`

	err = db.QueryRow(ctx, query, eventDemarre.AggregateID, joueurID).Scan(&typeEffet, &toursRestants)
	require.NoError(t, err)
	assert.Equal(t, "POISON", typeEffet)
	assert.Equal(t, 3, toursRestants)
}

// TestCombatTermineProjection vérifie la fin d'un combat
func TestCombatTermineProjection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()

	// Setup: Créer un combat
	eventDemarre := NewCombatDemarreEvent()
	seq1, err := insertEvent(ctx, db, eventDemarre)
	require.NoError(t, err)

	err = handleCombatDemarre(ctx, db, eventDemarre, seq1)
	require.NoError(t, err)

	// Extraire le vainqueur
	participants := eventDemarre.EventData["participants"].([]map[string]interface{})
	vainqueurID, _ := uuid.Parse(participants[0]["joueur_id"].(string))

	// Insérer un événement CombatTermine
	eventTermine := NewCombatTermineEvent(eventDemarre.AggregateID, 2, vainqueurID)
	seq2, err := insertEvent(ctx, db, eventTermine)
	require.NoError(t, err)

	err = handleCombatTermine(ctx, db, eventTermine, seq2)
	require.NoError(t, err)

	// Vérifier l'état du combat
	var etat string
	query := `SELECT etat FROM instances_combat WHERE combat_id = $1`
	err = db.QueryRow(ctx, query, eventDemarre.AggregateID).Scan(&etat)
	require.NoError(t, err)
	assert.Equal(t, "TERMINE", etat)
}

// TestProjectionIdempotence vérifie l'idempotence des handlers
func TestProjectionIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()

	// Insérer un événement
	event := NewCombatDemarreEvent()
	eventSequence, err := insertEvent(ctx, db, event)
	require.NoError(t, err)

	// Appliquer le handler deux fois
	err = handleCombatDemarre(ctx, db, event, eventSequence)
	require.NoError(t, err)

	err = handleCombatDemarre(ctx, db, event, eventSequence)
	require.NoError(t, err)

	// Vérifier qu'il n'y a qu'un seul enregistrement
	var count int64
	query := `SELECT COUNT(*) FROM instances_combat WHERE combat_id = $1`
	err = db.QueryRow(ctx, query, event.AggregateID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

// --- Handlers de projection simplifiés pour les tests ---

// handleCombatDemarre crée les projections pour un combat démarré
func handleCombatDemarre(ctx context.Context, db *pgxpool.Pool, event Event, eventSequence int64) error {
	// Insérer dans instances_combat
	query := `INSERT INTO instances_combat 
	          (combat_id, type_combat, etat, tour_actuel, started_at, event_sequence) 
	          VALUES ($1, $2, $3, $4, $5, $6)
	          ON CONFLICT (combat_id) DO NOTHING`

	_, err := db.Exec(ctx, query,
		event.AggregateID,
		event.EventData["type_combat"],
		"EN_COURS",
		1,
		event.TimestampUTC,
		eventSequence,
	)

	if err != nil {
		return err
	}

	// Insérer les participants
	participants := event.EventData["participants"].([]map[string]interface{})
	for _, p := range participants {
		joueurID, _ := uuid.Parse(p["joueur_id"].(string))

		query := `INSERT INTO participants_combat 
		          (combat_id, joueur_id, initiative, hp_actuel, hp_max, mana_actuel, mana_max, etat) 
		          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		          ON CONFLICT (combat_id, joueur_id) DO NOTHING`

		_, err := db.Exec(ctx, query,
			event.AggregateID,
			joueurID,
			int(p["initiative"].(int)),
			int(p["hp_actuel"].(int)),
			int(p["hp_max"].(int)),
			int(p["mana_actuel"].(int)),
			int(p["mana_max"].(int)),
			"ACTIF",
		)

		if err != nil {
			return err
		}
	}

	return nil
}

// handleActionCombatExecutee met à jour le tour actuel
func handleActionCombatExecutee(ctx context.Context, db *pgxpool.Pool, event Event, eventSequence int64) error {
	tour := event.EventData["tour"].(int)

	query := `UPDATE instances_combat 
	          SET tour_actuel = $1, event_sequence = $2 
	          WHERE combat_id = $3 AND event_sequence < $2`

	_, err := db.Exec(ctx, query, tour, eventSequence, event.AggregateID)
	return err
}

// handleDegatsInfliges met à jour les HP d'un participant
func handleDegatsInfliges(ctx context.Context, db *pgxpool.Pool, event Event, eventSequence int64) error {
	joueurID, _ := uuid.Parse(event.EventData["joueur_id"].(string))
	hpApres := int(event.EventData["hp_apres"].(int))

	query := `UPDATE participants_combat 
	          SET hp_actuel = $1 
	          WHERE combat_id = $2 AND joueur_id = $3`

	_, err := db.Exec(ctx, query, hpApres, event.AggregateID, joueurID)
	return err
}

// handleEffetStatutApplique crée un effet de statut
func handleEffetStatutApplique(ctx context.Context, db *pgxpool.Pool, event Event, eventSequence int64) error {
	joueurID, _ := uuid.Parse(event.EventData["joueur_id"].(string))
	typeEffet := event.EventData["type_effet"].(string)
	puissance := int(event.EventData["puissance"].(int))
	dureeTours := int(event.EventData["duree_tours"].(int))

	query := `INSERT INTO effets_statut 
	          (combat_id, joueur_id, type_effet, puissance, tours_restants, event_sequence) 
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(ctx, query,
		event.AggregateID,
		joueurID,
		typeEffet,
		puissance,
		dureeTours,
		eventSequence,
	)

	return err
}

// handleCombatTermine marque le combat comme terminé
func handleCombatTermine(ctx context.Context, db *pgxpool.Pool, event Event, eventSequence int64) error {
	query := `UPDATE instances_combat 
	          SET etat = $1, ended_at = $2, event_sequence = $3 
	          WHERE combat_id = $4`

	_, err := db.Exec(ctx, query,
		"TERMINE",
		event.TimestampUTC,
		eventSequence,
		event.AggregateID,
	)

	return err
}
