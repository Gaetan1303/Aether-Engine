package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

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

// createPlayerProjection crée une projection joueur
func createPlayerProjection(ctx context.Context, db *pgxpool.Pool, joueur JoueurAggregate) error {
	query := `INSERT INTO joueurs 
	          (joueur_id, username, niveau, experience_actuelle, 
	           hp_actuel, hp_max, mana_actuel, mana_max, or, event_sequence) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1)
	          ON CONFLICT (joueur_id) DO NOTHING`

	_, err := db.Exec(ctx, query,
		joueur.ID,
		joueur.Username,
		joueur.Niveau,
		joueur.ExperienceActuelle,
		joueur.HPActuel,
		joueur.HPMax,
		joueur.ManaActuel,
		joueur.ManaMax,
		joueur.Or,
	)

	return err
}

// createCombatProjection crée une projection combat
func createCombatProjection(ctx context.Context, db *pgxpool.Pool,
	combatID uuid.UUID, j1, j2 JoueurAggregate, eventSeq int64) error {

	// Insérer l'instance de combat
	query := `INSERT INTO instances_combat 
	          (combat_id, type_combat, etat, tour_actuel, participant_actif, started_at, event_sequence) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.Exec(ctx, query,
		combatID,
		"PVP",
		"EN_COURS",
		1,
		j1.ID,
		time.Now(),
		eventSeq,
	)

	if err != nil {
		return err
	}

	// Insérer les participants
	queryParticipant := `INSERT INTO participants_combat 
	                     (combat_id, joueur_id, initiative, hp_actuel, hp_max, 
	                      mana_actuel, mana_max, etat) 
	                     VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = db.Exec(ctx, queryParticipant,
		combatID, j1.ID, 15, j1.HPActuel, j1.HPMax, j1.ManaActuel, j1.ManaMax, "ACTIF")
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, queryParticipant,
		combatID, j2.ID, 12, j2.HPActuel, j2.HPMax, j2.ManaActuel, j2.ManaMax, "ACTIF")

	return err
}

// createCombatCache crée le cache Redis pour un combat
func createCombatCache(ctx context.Context, client *redis.Client,
	combatID uuid.UUID, j1, j2 JoueurAggregate) error {

	// État du combat
	keyState := fmt.Sprintf("combat:%s:state", combatID)
	err := client.HSet(ctx, keyState, map[string]interface{}{
		"combat_id":   combatID.String(),
		"etat":        "EN_COURS",
		"tour_actuel": 1,
		"phase":       "ACTION",
	}).Err()

	if err != nil {
		return err
	}

	err = client.Expire(ctx, keyState, time.Hour).Err()
	if err != nil {
		return err
	}

	// Participants
	keyParticipants := fmt.Sprintf("combat:%s:participants", combatID)
	err = client.SAdd(ctx, keyParticipants, j1.ID.String(), j2.ID.String()).Err()
	if err != nil {
		return err
	}

	// Ordre des tours (Sorted Set avec initiative comme score)
	keyTurnOrder := fmt.Sprintf("combat:%s:turn_order", combatID)
	err = client.ZAdd(ctx, keyTurnOrder,
		redis.Z{Score: 15, Member: j1.ID.String()},
		redis.Z{Score: 12, Member: j2.ID.String()},
	).Err()

	if err != nil {
		return err
	}

	// Cache des participants individuels
	keyP1 := fmt.Sprintf("combat:%s:participant:%s", combatID, j1.ID)
	err = client.HSet(ctx, keyP1, map[string]interface{}{
		"hp_actuel":   j1.HPActuel,
		"hp_max":      j1.HPMax,
		"mana_actuel": j1.ManaActuel,
		"mana_max":    j1.ManaMax,
		"initiative":  15,
	}).Err()

	if err != nil {
		return err
	}

	keyP2 := fmt.Sprintf("combat:%s:participant:%s", combatID, j2.ID)
	err = client.HSet(ctx, keyP2, map[string]interface{}{
		"hp_actuel":   j2.HPActuel,
		"hp_max":      j2.HPMax,
		"mana_actuel": j2.ManaActuel,
		"mana_max":    j2.ManaMax,
		"initiative":  12,
	}).Err()

	return err
}

// updateCombatTurn met à jour le tour actuel du combat
func updateCombatTurn(ctx context.Context, db *pgxpool.Pool, combatID uuid.UUID, tour int, eventSeq int64) error {
	query := `UPDATE instances_combat 
	          SET tour_actuel = $1, event_sequence = $2 
	          WHERE combat_id = $3`

	_, err := db.Exec(ctx, query, tour, eventSeq, combatID)
	return err
}

// updateParticipantHP met à jour les HP d'un participant
func updateParticipantHP(ctx context.Context, db *pgxpool.Pool, combatID, joueurID uuid.UUID, hp int) error {
	query := `UPDATE participants_combat 
	          SET hp_actuel = $1 
	          WHERE combat_id = $2 AND joueur_id = $3`

	_, err := db.Exec(ctx, query, hp, combatID, joueurID)
	return err
}

// createEffetsStatutTable crée la table effets_statut si elle n'existe pas
func createEffetsStatutTable(t *testing.T, db *pgxpool.Pool) {
	ctx := context.Background()

	query := `CREATE TABLE IF NOT EXISTS effets_statut (
		combat_id UUID NOT NULL,
		joueur_id UUID NOT NULL,
		type_effet VARCHAR(50) NOT NULL,
		puissance INT NOT NULL,
		tours_restants INT NOT NULL,
		event_sequence BIGINT NOT NULL,
		PRIMARY KEY (combat_id, joueur_id, type_effet)
	)`

	_, err := db.Exec(ctx, query)
	require.NoError(t, err)
}

// createEffetStatutProjection crée une projection d'effet de statut
func createEffetStatutProjection(ctx context.Context, db *pgxpool.Pool,
	combatID, joueurID uuid.UUID, typeEffet string, puissance, tours int, eventSeq int64) error {

	query := `INSERT INTO effets_statut 
	          (combat_id, joueur_id, type_effet, puissance, tours_restants, event_sequence) 
	          VALUES ($1, $2, $3, $4, $5, $6)
	          ON CONFLICT (combat_id, joueur_id, type_effet) 
	          DO UPDATE SET tours_restants = $5, event_sequence = $6`

	_, err := db.Exec(ctx, query,
		combatID,
		joueurID,
		typeEffet,
		puissance,
		tours,
		eventSeq,
	)

	return err
}

// terminateCombat marque un combat comme terminé
func terminateCombat(ctx context.Context, db *pgxpool.Pool, combatID uuid.UUID, eventSeq int64) error {
	query := `UPDATE instances_combat 
	          SET etat = 'TERMINE', ended_at = $1, event_sequence = $2 
	          WHERE combat_id = $3`

	_, err := db.Exec(ctx, query, time.Now(), eventSeq, combatID)
	return err
}
