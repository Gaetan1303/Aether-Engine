package bases_donnees

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInsertEvent vérifie l'insertion d'un événement dans l'Event Store
func TestInsertEvent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()
	event := NewJoueurCreeEvent()

	// Insérer l'événement
	eventSequence, err := insertEvent(ctx, db, event)
	require.NoError(t, err)
	require.Greater(t, eventSequence, int64(0))

	// Vérifier l'insertion
	var stored Event
	query := `SELECT event_id, aggregate_type, aggregate_id, aggregate_version, 
	                 event_type, timestamp_utc 
	          FROM evenements WHERE event_id = $1`

	err = db.QueryRow(ctx, query, event.EventID).Scan(
		&stored.EventID,
		&stored.AggregateType,
		&stored.AggregateID,
		&stored.AggregateVersion,
		&stored.EventType,
		&stored.TimestampUTC,
	)

	require.NoError(t, err)
	assert.Equal(t, event.EventID, stored.EventID)
	assert.Equal(t, "JoueurCree", stored.EventType)
	assert.Equal(t, "Joueur", stored.AggregateType)
	assert.Equal(t, 1, stored.AggregateVersion)
}

// TestInsertMultipleEvents vérifie l'insertion de plusieurs événements séquentiels
func TestInsertMultipleEvents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()
	joueurID := uuid.New()

	// Créer 5 événements séquentiels
	events := []Event{
		NewEventBuilder().
			WithAggregateType("Joueur").
			WithAggregateID(joueurID).
			WithAggregateVersion(1).
			WithEventType("JoueurCree").
			WithEventData(map[string]interface{}{"username": "test"}).
			Build(),
		NewExperienceGagneeEvent(joueurID, 2, 100),
		NewExperienceGagneeEvent(joueurID, 3, 150),
		NewNiveauGagneEvent(joueurID, 4, 2),
		NewExperienceGagneeEvent(joueurID, 5, 200),
	}

	// Insérer tous les événements
	for _, event := range events {
		_, err := insertEvent(ctx, db, event)
		require.NoError(t, err)
	}

	// Vérifier l'ordre et la séquence
	var count int64
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM evenements WHERE aggregate_id = $1", joueurID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)

	// Vérifier l'ordre des versions
	rows, err := db.Query(ctx, `
		SELECT aggregate_version 
		FROM evenements 
		WHERE aggregate_id = $1 
		ORDER BY event_sequence`, joueurID)
	require.NoError(t, err)
	defer rows.Close()

	expectedVersion := 1
	for rows.Next() {
		var version int
		err := rows.Scan(&version)
		require.NoError(t, err)
		assert.Equal(t, expectedVersion, version)
		expectedVersion++
	}
}

// TestOptimisticConcurrency vérifie la détection des conflits de version
func TestOptimisticConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()
	joueurID := uuid.New()

	// Insérer le premier événement (version 1)
	event1 := NewEventBuilder().
		WithAggregateType("Joueur").
		WithAggregateID(joueurID).
		WithAggregateVersion(1).
		WithEventType("JoueurCree").
		WithEventData(map[string]interface{}{"username": "test"}).
		Build()

	_, err := insertEvent(ctx, db, event1)
	require.NoError(t, err)

	// Essayer d'insérer un autre événement avec la même version (devrait échouer)
	event2 := NewEventBuilder().
		WithAggregateType("Joueur").
		WithAggregateID(joueurID).
		WithAggregateVersion(1).
		WithEventType("ExperienceGagnee").
		WithEventData(map[string]interface{}{"xp": 100}).
		Build()

	_, err = insertEvent(ctx, db, event2)
	assert.Error(t, err)
	// Message peut être en français ou anglais selon la locale PostgreSQL
	assert.True(t, strings.Contains(err.Error(), "duplicate key value") ||
		strings.Contains(err.Error(), "clé dupliquée"))
}

// TestReconstructAggregate vérifie la reconstruction d'un agrégat depuis les événements
func TestReconstructAggregate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()
	joueurID := uuid.New()

	// Insérer plusieurs événements
	events := []Event{
		NewEventBuilder().
			WithAggregateType("Joueur").
			WithAggregateID(joueurID).
			WithAggregateVersion(1).
			WithEventType("JoueurCree").
			WithEventData(map[string]interface{}{
				"username": "test",
				"niveau":   1,
				"hp_max":   100,
			}).
			Build(),
		NewExperienceGagneeEvent(joueurID, 2, 100),
		NewNiveauGagneEvent(joueurID, 3, 2),
		NewExperienceGagneeEvent(joueurID, 4, 250),
	}

	for _, event := range events {
		_, err := insertEvent(ctx, db, event)
		require.NoError(t, err)
	}

	// Récupérer tous les événements pour reconstruire l'agrégat
	rows, err := db.Query(ctx, `
		SELECT event_type, event_data, aggregate_version 
		FROM evenements 
		WHERE aggregate_id = $1 
		ORDER BY aggregate_version`, joueurID)
	require.NoError(t, err)
	defer rows.Close()

	// Simuler la reconstruction (version simplifiée)
	niveau := 1
	xpTotal := 0

	for rows.Next() {
		var eventType string
		var eventData map[string]interface{}
		var version int

		err := rows.Scan(&eventType, &eventData, &version)
		require.NoError(t, err)

		switch eventType {
		case "ExperienceGagnee":
			// Dans un vrai système, on extrairait la valeur du JSON
			xpTotal += 100 // Simplifié
		case "NiveauGagne":
			niveau++
		}
	}

	assert.Equal(t, 2, niveau)
	assert.Greater(t, xpTotal, 0)
}

// TestSnapshot vérifie la création et l'utilisation de snapshots
func TestSnapshot(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()
	joueur := NewJoueurAggregate()

	// Insérer un snapshot
	snapshotData, err := joueur.ToJSON()
	require.NoError(t, err)

	query := `INSERT INTO snapshots (aggregate_id, aggregate_type, aggregate_version, snapshot_data) 
	          VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(ctx, query, joueur.ID, "Joueur", joueur.Version, snapshotData)
	require.NoError(t, err)

	// Récupérer le snapshot
	var retrievedData []byte
	var retrievedVersion int

	query = `SELECT snapshot_data, aggregate_version 
	         FROM snapshots 
	         WHERE aggregate_id = $1 
	         ORDER BY aggregate_version DESC 
	         LIMIT 1`

	err = db.QueryRow(ctx, query, joueur.ID).Scan(&retrievedData, &retrievedVersion)
	require.NoError(t, err)

	assert.Equal(t, joueur.Version, retrievedVersion)
	assert.NotEmpty(t, retrievedData)
}

// TestQueryEventsByTimeRange vérifie les requêtes temporelles
func TestQueryEventsByTimeRange(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()

	// Insérer des événements à différents moments
	now := time.Now()
	events := []Event{
		NewEventBuilder().
			WithAggregateType("Joueur").
			WithAggregateID(uuid.New()).
			WithAggregateVersion(1).
			WithEventType("JoueurCree").
			WithEventData(map[string]interface{}{}).
			WithTimestamp(now.Add(-2 * time.Hour)).
			Build(),
		NewEventBuilder().
			WithAggregateType("Joueur").
			WithAggregateID(uuid.New()).
			WithAggregateVersion(1).
			WithEventType("JoueurCree").
			WithEventData(map[string]interface{}{}).
			WithTimestamp(now.Add(-1 * time.Hour)).
			Build(),
		NewEventBuilder().
			WithAggregateType("Joueur").
			WithAggregateID(uuid.New()).
			WithAggregateVersion(1).
			WithEventType("JoueurCree").
			WithEventData(map[string]interface{}{}).
			WithTimestamp(now).
			Build(),
	}

	for _, event := range events {
		_, err := insertEventWithTimestamp(ctx, db, event)
		require.NoError(t, err)
	}

	// Requête pour les événements de la dernière heure
	var count int64
	query := `SELECT COUNT(*) FROM evenements WHERE timestamp_utc >= $1`

	err := db.QueryRow(ctx, query, now.Add(-90*time.Minute)).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// TestQueryEventsByType vérifie les requêtes par type d'événement
func TestQueryEventsByType(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()

	// Insérer différents types d'événements
	joueurID := uuid.New()
	_, err := insertEvent(ctx, db, NewJoueurCreeEvent())
	require.NoError(t, err)

	_, err = insertEvent(ctx, db, NewExperienceGagneeEvent(joueurID, 2, 100))
	require.NoError(t, err)

	_, err = insertEvent(ctx, db, NewExperienceGagneeEvent(joueurID, 3, 150))
	require.NoError(t, err)

	// Compter les événements de type ExperienceGagnee
	var count int64
	query := `SELECT COUNT(*) FROM evenements WHERE event_type = $1`

	err = db.QueryRow(ctx, query, "ExperienceGagnee").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// TestTransactionalInsert vérifie l'insertion transactionnelle d'événements
func TestTransactionalInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := NewTestDB(t)
	CreateTestSchema(t, db)
	defer DropTestSchema(t, db)

	ctx := context.Background()

	// Démarrer une transaction
	tx, err := db.Begin(ctx)
	require.NoError(t, err)

	joueurID := uuid.New()

	// Insérer plusieurs événements dans la transaction
	event1 := NewEventBuilder().
		WithAggregateType("Joueur").
		WithAggregateID(joueurID).
		WithAggregateVersion(1).
		WithEventType("JoueurCree").
		WithEventData(map[string]interface{}{}).
		Build()

	event2 := NewExperienceGagneeEvent(joueurID, 2, 100)

	_, err = insertEventTx(ctx, tx, event1)
	require.NoError(t, err)

	_, err = insertEventTx(ctx, tx, event2)
	require.NoError(t, err)

	// Rollback pour tester
	err = tx.Rollback(ctx)
	require.NoError(t, err)

	// Vérifier qu'aucun événement n'a été inséré
	var count int64
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM evenements WHERE aggregate_id = $1", joueurID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Réessayer avec commit
	tx, err = db.Begin(ctx)
	require.NoError(t, err)

	_, err = insertEventTx(ctx, tx, event1)
	require.NoError(t, err)

	_, err = insertEventTx(ctx, tx, event2)
	require.NoError(t, err)

	err = tx.Commit(ctx)
	require.NoError(t, err)

	// Vérifier que les événements ont été insérés
	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM evenements WHERE aggregate_id = $1", joueurID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

// --- Helpers ---

// insertEvent insère un événement dans l'Event Store
func insertEvent(ctx context.Context, db *pgxpool.Pool, event Event) (int64, error) {
	var eventSequence int64

	query := `INSERT INTO evenements 
	          (event_id, aggregate_type, aggregate_id, aggregate_version, 
	           event_type, event_data, metadata, timestamp_utc) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	          RETURNING event_sequence`

	err := db.QueryRow(ctx, query,
		event.EventID,
		event.AggregateType,
		event.AggregateID,
		event.AggregateVersion,
		event.EventType,
		event.EventData,
		event.Metadata,
		event.TimestampUTC,
	).Scan(&eventSequence)

	return eventSequence, err
}

// insertEventWithTimestamp insère un événement avec un timestamp spécifique
func insertEventWithTimestamp(ctx context.Context, db *pgxpool.Pool, event Event) (int64, error) {
	return insertEvent(ctx, db, event)
}

// insertEventTx insère un événement dans une transaction
func insertEventTx(ctx context.Context, tx pgx.Tx, event Event) (int64, error) {
	var eventSequence int64

	query := `INSERT INTO evenements 
	          (event_id, aggregate_type, aggregate_id, aggregate_version, 
	           event_type, event_data, metadata, timestamp_utc) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	          RETURNING event_sequence`

	err := tx.QueryRow(ctx, query,
		event.EventID,
		event.AggregateType,
		event.AggregateID,
		event.AggregateVersion,
		event.EventType,
		event.EventData,
		event.Metadata,
		event.TimestampUTC,
	).Scan(&eventSequence)

	return eventSequence, err
}
