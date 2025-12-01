package bases_donnees

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// Configuration de test
type TestConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisHost string
	RedisPort string
	RedisDB   int
}

// LoadTestConfig charge la configuration depuis les variables d'environnement
func LoadTestConfig() *TestConfig {
	return &TestConfig{
		DBHost:     getEnv("TEST_DB_HOST", "localhost"),
		DBPort:     getEnv("TEST_DB_PORT", "5432"),
		DBUser:     getEnv("TEST_DB_USER", "test"),
		DBPassword: getEnv("TEST_DB_PASSWORD", "test"),
		DBName:     getEnv("TEST_DB_NAME", "aether_test"),

		RedisHost: getEnv("TEST_REDIS_HOST", "localhost"),
		RedisPort: getEnv("TEST_REDIS_PORT", "6379"),
		RedisDB:   getEnvInt("TEST_REDIS_DB", 15),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		fmt.Sscanf(value, "%d", &intValue)
		return intValue
	}
	return defaultValue
}

// NewTestDB crée une connexion à la base de données PostgreSQL de test
func NewTestDB(t *testing.T) *pgxpool.Pool {
	config := LoadTestConfig()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	require.NoError(t, err, "Failed to parse database config")

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	require.NoError(t, err, "Failed to connect to test database")

	// Vérifier la connexion
	err = pool.Ping(context.Background())
	require.NoError(t, err, "Failed to ping test database")

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

// NewTestRedis crée un client Redis pour les tests
func NewTestRedis(t *testing.T) *redis.Client {
	config := LoadTestConfig()

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		DB:           config.RedisDB,
		Password:     "",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	ctx := context.Background()
	err := client.Ping(ctx).Err()
	require.NoError(t, err, "Failed to connect to test Redis")

	t.Cleanup(func() {
		client.Close()
	})

	return client
}

// CleanDatabase nettoie toutes les tables de la base de données de test
func CleanDatabase(t *testing.T, db *pgxpool.Pool) {
	ctx := context.Background()

	// Ordre de suppression pour respecter les contraintes de clés étrangères
	tables := []string{
		// Projections joueur
		"historique_niveau",
		"quetes_joueur",
		"competences_joueur",
		"items_inventaire",
		"sets_equipement",
		"inventaires",
		"joueurs",

		// Projections combat
		"historique_combat",
		"effets_statut",
		"actions_combat",
		"participants_combat",
		"instances_combat",

		// Projections monde
		"prix_marche",
		"transactions_economie",
		"ordres_economie",
		"etat_monde",
		"quetes",
		"competences",
		"items",

		// Event Store
		"snapshots",
		"evenements",
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		_, err := db.Exec(ctx, query)
		require.NoError(t, err, "Failed to truncate table %s", table)
	}
}

// CleanRedis vide complètement la base Redis de test
func CleanRedis(t *testing.T, redis *redis.Client) {
	ctx := context.Background()
	err := redis.FlushDB(ctx).Err()
	require.NoError(t, err, "Failed to flush Redis database")
}

// WaitForProjection attend qu'une projection soit créée jusqu'à une certaine séquence
func WaitForProjection(db *pgxpool.Pool, table string, minSequence int64, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE event_sequence >= $1", table)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for projection in %s at sequence %d", table, minSequence)
		case <-ticker.C:
			var count int64
			err := db.QueryRow(context.Background(), query, minSequence).Scan(&count)
			if err == nil && count > 0 {
				return nil
			}
		}
	}
}

// AssertEventStored vérifie qu'un événement a été stocké dans l'Event Store
func AssertEventStored(t *testing.T, db *pgxpool.Pool, eventID string) {
	ctx := context.Background()
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM evenements WHERE event_id = $1)"
	err := db.QueryRow(ctx, query, eventID).Scan(&exists)
	require.NoError(t, err, "Failed to check event existence")
	require.True(t, exists, "Event %s not found in Event Store", eventID)
}

// AssertRedisKeyExists vérifie qu'une clé existe dans Redis
func AssertRedisKeyExists(t *testing.T, redis *redis.Client, key string) {
	ctx := context.Background()
	exists, err := redis.Exists(ctx, key).Result()
	require.NoError(t, err, "Failed to check Redis key existence")
	require.Equal(t, int64(1), exists, "Redis key %s does not exist", key)
}

// AssertRedisKeyNotExists vérifie qu'une clé n'existe pas dans Redis
func AssertRedisKeyNotExists(t *testing.T, redis *redis.Client, key string) {
	ctx := context.Background()
	exists, err := redis.Exists(ctx, key).Result()
	require.NoError(t, err, "Failed to check Redis key existence")
	require.Equal(t, int64(0), exists, "Redis key %s should not exist", key)
}

// GetRedisKeyTTL récupère le TTL d'une clé Redis
func GetRedisKeyTTL(t *testing.T, redis *redis.Client, key string) time.Duration {
	ctx := context.Background()
	ttl, err := redis.TTL(ctx, key).Result()
	require.NoError(t, err, "Failed to get TTL for Redis key %s", key)
	return ttl
}

// AssertProjectionCount vérifie le nombre d'enregistrements dans une table de projection
func AssertProjectionCount(t *testing.T, db *pgxpool.Pool, table string, expected int64) {
	ctx := context.Background()
	var count int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	err := db.QueryRow(ctx, query).Scan(&count)
	require.NoError(t, err, "Failed to count rows in %s", table)
	require.Equal(t, expected, count, "Expected %d rows in %s, got %d", expected, table, count)
}

// CreateTestSchema crée toutes les tables nécessaires pour les tests
// Note: Dans un environnement réel, cela serait géré par des migrations
func CreateTestSchema(t *testing.T, db *pgxpool.Pool) {
	ctx := context.Background()

	// Event Store
	eventStoreSchema := `
	CREATE TABLE IF NOT EXISTS evenements (
		event_sequence BIGSERIAL PRIMARY KEY,
		event_id UUID NOT NULL UNIQUE,
		aggregate_type VARCHAR(50) NOT NULL,
		aggregate_id UUID NOT NULL,
		aggregate_version INT NOT NULL,
		event_type VARCHAR(100) NOT NULL,
		event_data JSONB NOT NULL,
		metadata JSONB,
		timestamp_utc TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		UNIQUE(aggregate_id, aggregate_version)
	);

	CREATE INDEX IF NOT EXISTS idx_evenements_aggregate ON evenements(aggregate_id, aggregate_version);
	CREATE INDEX IF NOT EXISTS idx_evenements_type ON evenements(event_type);
	CREATE INDEX IF NOT EXISTS idx_evenements_timestamp ON evenements(timestamp_utc);

	CREATE TABLE IF NOT EXISTS snapshots (
		aggregate_id UUID NOT NULL,
		aggregate_type VARCHAR(50) NOT NULL,
		aggregate_version INT NOT NULL,
		snapshot_data JSONB NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		PRIMARY KEY (aggregate_id, aggregate_version)
	);

	CREATE INDEX IF NOT EXISTS idx_snapshots_created ON snapshots(created_at);
	`

	// Projections simplifiées pour les tests
	projectionsSchema := `
	CREATE TABLE IF NOT EXISTS joueurs (
		joueur_id UUID PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		niveau INT NOT NULL DEFAULT 1,
		experience_actuelle BIGINT NOT NULL DEFAULT 0,
		hp_actuel INT NOT NULL,
		hp_max INT NOT NULL,
		mana_actuel INT NOT NULL,
		mana_max INT NOT NULL,
		"or" BIGINT NOT NULL DEFAULT 0,
		event_sequence BIGINT NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS instances_combat (
		combat_id UUID PRIMARY KEY,
		type_combat VARCHAR(20) NOT NULL,
		etat VARCHAR(20) NOT NULL,
		tour_actuel INT NOT NULL DEFAULT 1,
		participant_actif UUID,
		started_at TIMESTAMPTZ NOT NULL,
		ended_at TIMESTAMPTZ,
		event_sequence BIGINT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS participants_combat (
		combat_id UUID NOT NULL,
		joueur_id UUID NOT NULL,
		initiative INT NOT NULL,
		hp_actuel INT NOT NULL,
		hp_max INT NOT NULL,
		mana_actuel INT NOT NULL,
		mana_max INT NOT NULL,
		etat VARCHAR(20) NOT NULL DEFAULT 'ACTIF',
		PRIMARY KEY (combat_id, joueur_id)
	);

	CREATE TABLE IF NOT EXISTS inventaires (
		joueur_id UUID PRIMARY KEY,
		capacite_max INT NOT NULL DEFAULT 50,
		poids_actuel DECIMAL(10, 2) NOT NULL DEFAULT 0,
		nombre_items INT NOT NULL DEFAULT 0,
		event_sequence BIGINT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS items_inventaire (
		joueur_id UUID NOT NULL,
		slot INT NOT NULL,
		item_id VARCHAR(50) NOT NULL,
		quantite INT NOT NULL DEFAULT 1,
		equipe BOOLEAN NOT NULL DEFAULT FALSE,
		PRIMARY KEY (joueur_id, slot)
	);
	`

	_, err := db.Exec(ctx, eventStoreSchema)
	require.NoError(t, err, "Failed to create Event Store schema")

	_, err = db.Exec(ctx, projectionsSchema)
	require.NoError(t, err, "Failed to create projections schema")
}

// DropTestSchema supprime toutes les tables de test
func DropTestSchema(t *testing.T, db *pgxpool.Pool) {
	ctx := context.Background()

	dropSchema := `
	DROP TABLE IF EXISTS items_inventaire CASCADE;
	DROP TABLE IF EXISTS inventaires CASCADE;
	DROP TABLE IF EXISTS participants_combat CASCADE;
	DROP TABLE IF EXISTS instances_combat CASCADE;
	DROP TABLE IF EXISTS joueurs CASCADE;
	DROP TABLE IF EXISTS snapshots CASCADE;
	DROP TABLE IF EXISTS evenements CASCADE;
	`

	_, err := db.Exec(ctx, dropSchema)
	require.NoError(t, err, "Failed to drop test schema")
}
