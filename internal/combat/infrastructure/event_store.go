package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aether-engine/aether-engine/internal/combat/application"
	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresEventStore implémente EventStore avec PostgreSQL
type PostgresEventStore struct {
	pool *pgxpool.Pool
}

// NewPostgresEventStore crée une nouvelle instance de PostgresEventStore
func NewPostgresEventStore(pool *pgxpool.Pool) application.EventStore {
	return &PostgresEventStore{
		pool: pool,
	}
}

// AppendEvents ajoute des événements à un agrégat
func (s *PostgresEventStore) AppendEvents(aggregateID string, events []domain.Evenement, expectedVersion int) error {
	if len(events) == 0 {
		return nil
	}

	ctx := context.Background()

	// Commencer une transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erreur transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Vérifier la version actuelle (optimistic concurrency)
	var currentVersion int
	err = tx.QueryRow(ctx, `
		SELECT COALESCE(MAX(version), 0)
		FROM events
		WHERE aggregate_id = $1
	`, aggregateID).Scan(&currentVersion)

	if err != nil && err != pgx.ErrNoRows {
		return fmt.Errorf("erreur vérification version: %w", err)
	}

	// Vérifier le concurrency control
	if currentVersion != expectedVersion {
		return fmt.Errorf("conflit de version: attendu %d, actuel %d", expectedVersion, currentVersion)
	}

	// Insérer les événements
	for _, evt := range events {
		// Sérialiser le payload
		payload, err := json.Marshal(evt)
		if err != nil {
			return fmt.Errorf("erreur sérialisation événement: %w", err)
		}

		// Insérer l'événement
		_, err = tx.Exec(ctx, `
			INSERT INTO events (
				aggregate_id,
				version,
				event_type,
				payload,
				created_at
			) VALUES ($1, $2, $3, $4, $5)
		`, aggregateID, evt.AggregateVersion(), evt.EventType(), payload, evt.Timestamp())

		if err != nil {
			return fmt.Errorf("erreur insertion événement: %w", err)
		}
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erreur commit: %w", err)
	}

	return nil
}

// LoadEvents charge tous les événements d'un agrégat
func (s *PostgresEventStore) LoadEvents(aggregateID string) ([]domain.Evenement, error) {
	return s.LoadEventsFromVersion(aggregateID, 0)
}

// LoadEventsFromVersion charge les événements depuis une version
func (s *PostgresEventStore) LoadEventsFromVersion(aggregateID string, fromVersion int) ([]domain.Evenement, error) {
	ctx := context.Background()

	// Charger les événements depuis la version spécifiée
	rows, err := s.pool.Query(ctx, `
		SELECT event_type, payload, version, created_at
		FROM events
		WHERE aggregate_id = $1 AND version > $2
		ORDER BY version ASC
	`, aggregateID, fromVersion)

	if err != nil {
		return nil, fmt.Errorf("erreur chargement événements: %w", err)
	}
	defer rows.Close()

	events := make([]domain.Evenement, 0)

	for rows.Next() {
		var eventType string
		var payload []byte
		var version int
		var createdAt time.Time

		if err := rows.Scan(&eventType, &payload, &version, &createdAt); err != nil {
			return nil, fmt.Errorf("erreur scan événement: %w", err)
		}

		// Désérialiser l'événement
		evt, err := deserializeEvent(eventType, payload)
		if err != nil {
			return nil, fmt.Errorf("erreur désérialisation événement %s: %w", eventType, err)
		}

		evt.SetAggregateID(aggregateID)
		evt.SetAggregateVersion(version)
		evt.SetTimestamp(createdAt)

		events = append(events, evt)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("erreur itération événements: %w", rows.Err())
	}

	return events, nil
}

// SaveSnapshot sauvegarde un snapshot
func (s *PostgresEventStore) SaveSnapshot(aggregateID string, version int, data []byte) error {
	ctx := context.Background()

	_, err := s.pool.Exec(ctx, `
		INSERT INTO snapshots (aggregate_id, version, data, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (aggregate_id) DO UPDATE
		SET version = $2, data = $3, created_at = $4
	`, aggregateID, version, data, time.Now())

	if err != nil {
		return fmt.Errorf("erreur sauvegarde snapshot: %w", err)
	}

	return nil
}

// LoadSnapshot charge le dernier snapshot
func (s *PostgresEventStore) LoadSnapshot(aggregateID string) (version int, data []byte, err error) {
	ctx := context.Background()

	err = s.pool.QueryRow(ctx, `
		SELECT version, data
		FROM snapshots
		WHERE aggregate_id = $1
	`, aggregateID).Scan(&version, &data)

	if err == pgx.ErrNoRows {
		return 0, nil, nil // Pas de snapshot
	}

	if err != nil {
		return 0, nil, fmt.Errorf("erreur chargement snapshot: %w", err)
	}

	return version, data, nil
}

// deserializeEvent désérialise un événement depuis JSON
func deserializeEvent(eventType string, payload []byte) (domain.Evenement, error) {
	var evt domain.Evenement

	switch eventType {
	case "CombatDemarre":
		evt = &domain.CombatDemarreEvent{}
	case "TourDemarre":
		evt = &domain.TourDemarreEvent{}
	case "ActionExecutee":
		evt = &domain.ActionExecuteeEvent{}
	case "DegatsInfliges":
		evt = &domain.DegatsInfligesEvent{}
	case "SoinApplique":
		evt = &domain.SoinApliqueEvent{}
	case "StatutApplique":
		evt = &domain.StatutAppliqueEvent{}
	case "UniteEliminee":
		evt = &domain.UniteElimineeEvent{}
	case "UniteDeplacee":
		evt = &domain.UniteDeplaceeEvent{}
	case "CompetenceUtilisee":
		evt = &domain.CompetenceUtiliseeEvent{}
	case "CombatTermine":
		evt = &domain.CombatTermineEvent{}
	default:
		return nil, errors.New("type d'événement inconnu: " + eventType)
	}

	if err := json.Unmarshal(payload, evt); err != nil {
		return nil, err
	}

	return evt, nil
}

// InitSchema initialise le schéma de la base de données
func InitSchema(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Créer la table events
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			aggregate_id VARCHAR(255) NOT NULL,
			version INTEGER NOT NULL,
			event_type VARCHAR(100) NOT NULL,
			payload JSONB NOT NULL,
			created_at TIMESTAMP NOT NULL,
			UNIQUE (aggregate_id, version)
		)
	`)
	if err != nil {
		return fmt.Errorf("erreur création table events: %w", err)
	}

	// Index sur aggregate_id
	_, err = pool.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_events_aggregate_id
		ON events (aggregate_id)
	`)
	if err != nil {
		return fmt.Errorf("erreur création index aggregate_id: %w", err)
	}

	// Index sur created_at pour queries temporelles
	_, err = pool.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_events_created_at
		ON events (created_at)
	`)
	if err != nil {
		return fmt.Errorf("erreur création index created_at: %w", err)
	}

	// Créer la table snapshots
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS snapshots (
			aggregate_id VARCHAR(255) PRIMARY KEY,
			version INTEGER NOT NULL,
			data BYTEA NOT NULL,
			created_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("erreur création table snapshots: %w", err)
	}

	return nil
}
