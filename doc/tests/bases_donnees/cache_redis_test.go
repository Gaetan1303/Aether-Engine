package bases_donnees

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRedisCombatState vérifie le stockage de l'état d'un combat dans Redis
func TestRedisCombatState(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewTestRedis(t)
	CleanRedis(t, client)

	ctx := context.Background()
	combatID := uuid.New()
	key := fmt.Sprintf("combat:%s:state", combatID)

	// Stocker l'état du combat
	combatState := map[string]interface{}{
		"combat_id":   combatID.String(),
		"etat":        "EN_COURS",
		"tour_actuel": 1,
		"phase":       "ACTION",
	}

	err := client.HSet(ctx, key, combatState).Err()
	require.NoError(t, err)

	// Définir un TTL de 1 heure
	err = client.Expire(ctx, key, time.Hour).Err()
	require.NoError(t, err)

	// Récupérer l'état
	retrieved, err := client.HGetAll(ctx, key).Result()
	require.NoError(t, err)

	assert.Equal(t, "EN_COURS", retrieved["etat"])
	assert.Equal(t, "1", retrieved["tour_actuel"])
	assert.Equal(t, "ACTION", retrieved["phase"])

	// Vérifier le TTL
	ttl, err := client.TTL(ctx, key).Result()
	require.NoError(t, err)
	assert.Greater(t, ttl.Seconds(), 3500.0)
	assert.LessOrEqual(t, ttl.Seconds(), 3600.0)
}

// TestRedisCombatParticipants vérifie le stockage des participants
func TestRedisCombatParticipants(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewTestRedis(t)
	CleanRedis(t, client)

	ctx := context.Background()
	combatID := uuid.New()
	joueur1ID := uuid.New()
	joueur2ID := uuid.New()

	keyParticipants := fmt.Sprintf("combat:%s:participants", combatID)

	// Ajouter les participants
	err := client.SAdd(ctx, keyParticipants, joueur1ID.String(), joueur2ID.String()).Err()
	require.NoError(t, err)

	// Récupérer les participants
	participants, err := client.SMembers(ctx, keyParticipants).Result()
	require.NoError(t, err)
	assert.Len(t, participants, 2)
	assert.Contains(t, participants, joueur1ID.String())
	assert.Contains(t, participants, joueur2ID.String())

	// Vérifier qu'un joueur est participant
	isMember, err := client.SIsMember(ctx, keyParticipants, joueur1ID.String()).Result()
	require.NoError(t, err)
	assert.True(t, isMember)
}

// TestRedisTurnOrder vérifie l'ordre des tours
func TestRedisTurnOrder(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewTestRedis(t)
	CleanRedis(t, client)

	ctx := context.Background()
	combatID := uuid.New()
	keyTurnOrder := fmt.Sprintf("combat:%s:turn_order", combatID)

	// Ajouter l'ordre des tours avec les scores d'initiative
	joueur1ID := uuid.New()
	joueur2ID := uuid.New()
	joueur3ID := uuid.New()

	err := client.ZAdd(ctx, keyTurnOrder,
		redis.Z{Score: 15, Member: joueur1ID.String()},
		redis.Z{Score: 12, Member: joueur2ID.String()},
		redis.Z{Score: 10, Member: joueur3ID.String()},
	).Err()
	require.NoError(t, err)

	// Récupérer l'ordre (du plus grand au plus petit score)
	turnOrder, err := client.ZRevRange(ctx, keyTurnOrder, 0, -1).Result()
	require.NoError(t, err)

	assert.Len(t, turnOrder, 3)
	assert.Equal(t, joueur1ID.String(), turnOrder[0])
	assert.Equal(t, joueur2ID.String(), turnOrder[1])
	assert.Equal(t, joueur3ID.String(), turnOrder[2])

	// Récupérer le joueur avec la plus haute initiative
	topPlayer, err := client.ZRevRange(ctx, keyTurnOrder, 0, 0).Result()
	require.NoError(t, err)
	assert.Equal(t, joueur1ID.String(), topPlayer[0])
}

// TestRedisActionQueue vérifie la file d'actions
func TestRedisActionQueue(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewTestRedis(t)
	CleanRedis(t, client)

	ctx := context.Background()
	combatID := uuid.New()
	keyActions := fmt.Sprintf("combat:%s:actions", combatID)

	// Ajouter des actions à la file
	actions := []string{
		`{"type":"ATTAQUE","acteur":"player1","cible":"player2"}`,
		`{"type":"COMPETENCE","acteur":"player2","cible":"player1","skill":"FIREBALL"}`,
		`{"type":"SOIN","acteur":"player1","cible":"player1"}`,
	}

	for _, action := range actions {
		err := client.RPush(ctx, keyActions, action).Err()
		require.NoError(t, err)
	}

	// Récupérer la taille de la file
	length, err := client.LLen(ctx, keyActions).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(3), length)

	// Consommer les actions (FIFO)
	action1, err := client.LPop(ctx, keyActions).Result()
	require.NoError(t, err)
	assert.Contains(t, action1, "ATTAQUE")

	// Vérifier la nouvelle taille
	length, err = client.LLen(ctx, keyActions).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(2), length)
}

// TestRedisDistributedLock vérifie les verrous distribués
func TestRedisDistributedLock(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewTestRedis(t)
	CleanRedis(t, client)

	ctx := context.Background()
	combatID := uuid.New()
	lockKey := fmt.Sprintf("combat:%s:action", combatID)

	session1 := "session-1"
	session2 := "session-2"
	lockTTL := 5 * time.Second

	// Session 1 acquiert le verrou
	acquired, err := client.SetNX(ctx, lockKey, session1, lockTTL).Result()
	require.NoError(t, err)
	assert.True(t, acquired)

	// Session 2 ne peut pas acquérir le verrou
	acquired, err = client.SetNX(ctx, lockKey, session2, lockTTL).Result()
	require.NoError(t, err)
	assert.False(t, acquired)

	// Vérifier que le verrou appartient à session1
	owner, err := client.Get(ctx, lockKey).Result()
	require.NoError(t, err)
	assert.Equal(t, session1, owner)

	// Session 1 libère le verrou (avec script Lua pour sécurité)
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	result, err := client.Eval(ctx, script, []string{lockKey}, session1).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(1), result)

	// Maintenant session 2 peut acquérir
	acquired, err = client.SetNX(ctx, lockKey, session2, lockTTL).Result()
	require.NoError(t, err)
	assert.True(t, acquired)
}

// TestRedisLockExpiration vérifie l'expiration automatique des verrous
func TestRedisLockExpiration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewTestRedis(t)
	CleanRedis(t, client)

	ctx := context.Background()
	lockKey := "test:lock"
	lockTTL := 1 * time.Second

	// Acquérir le verrou avec un TTL court
	acquired, err := client.SetNX(ctx, lockKey, "session1", lockTTL).Result()
	require.NoError(t, err)
	assert.True(t, acquired)

	// Attendre l'expiration
	time.Sleep(1100 * time.Millisecond)

	// Le verrou devrait avoir expiré
	exists, err := client.Exists(ctx, lockKey).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists)

	// Un autre client peut acquérir
	acquired, err = client.SetNX(ctx, lockKey, "session2", lockTTL).Result()
	require.NoError(t, err)
	assert.True(t, acquired)
}

// TestRedisPubSub vérifie le système de notifications
func TestRedisPubSub(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewTestRedis(t)
	CleanRedis(t, client)

	ctx := context.Background()
	combatID := uuid.New()
	channel := fmt.Sprintf("combat:%s:notifications", combatID)

	// S'abonner au canal
	pubsub := client.Subscribe(ctx, channel)
	defer pubsub.Close()

	// Attendre la confirmation d'abonnement
	_, err := pubsub.Receive(ctx)
	require.NoError(t, err)

	// Publier un message
	message := `{"type":"DEGATS","joueur_id":"player1","degats":25}`
	err = client.Publish(ctx, channel, message).Err()
	require.NoError(t, err)

	// Recevoir le message (avec timeout)
	msgChan := pubsub.Channel()
	select {
	case msg := <-msgChan:
		assert.Equal(t, channel, msg.Channel)
		assert.Contains(t, msg.Payload, "DEGATS")
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
}

// TestRedisLeaderboard vérifie les classements
func TestRedisLeaderboard(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewTestRedis(t)
	CleanRedis(t, client)

	ctx := context.Background()
	leaderboardKey := "leaderboard:pvp"

	// Ajouter des scores
	players := []redis.Z{
		{Score: 1500, Member: "player1"},
		{Score: 2000, Member: "player2"},
		{Score: 1200, Member: "player3"},
		{Score: 1800, Member: "player4"},
		{Score: 2200, Member: "player5"},
	}

	err := client.ZAdd(ctx, leaderboardKey, players...).Err()
	require.NoError(t, err)

	// Récupérer le top 3
	top3, err := client.ZRevRangeWithScores(ctx, leaderboardKey, 0, 2).Result()
	require.NoError(t, err)

	assert.Len(t, top3, 3)
	assert.Equal(t, "player5", top3[0].Member)
	assert.Equal(t, 2200.0, top3[0].Score)
	assert.Equal(t, "player2", top3[1].Member)
	assert.Equal(t, "player4", top3[2].Member)

	// Récupérer le rang d'un joueur
	rank, err := client.ZRevRank(ctx, leaderboardKey, "player1").Result()
	require.NoError(t, err)
	assert.Equal(t, int64(3), rank) // 4ème position (0-indexed)

	// Incrémenter le score d'un joueur
	newScore, err := client.ZIncrBy(ctx, leaderboardKey, 500, "player1").Result()
	require.NoError(t, err)
	assert.Equal(t, 2000.0, newScore)
}

// TestRedisInventoryCache vérifie le cache d'inventaire
func TestRedisInventoryCache(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewTestRedis(t)
	CleanRedis(t, client)

	ctx := context.Background()
	joueurID := uuid.New()
	keyInventory := fmt.Sprintf("inventory:%s", joueurID)

	// Stocker des items
	items := map[string]interface{}{
		"EPEE_FER":   3,
		"POTION_VIE": 10,
		"BOUCLIER":   1,
	}

	err := client.HSet(ctx, keyInventory, items).Err()
	require.NoError(t, err)

	// Définir un TTL
	err = client.Expire(ctx, keyInventory, 30*time.Minute).Err()
	require.NoError(t, err)

	// Récupérer un item spécifique
	quantity, err := client.HGet(ctx, keyInventory, "POTION_VIE").Result()
	require.NoError(t, err)
	assert.Equal(t, "10", quantity)

	// Incrémenter la quantité d'un item
	newQty, err := client.HIncrBy(ctx, keyInventory, "POTION_VIE", -2).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(8), newQty)

	// Vérifier tous les items
	allItems, err := client.HGetAll(ctx, keyInventory).Result()
	require.NoError(t, err)
	assert.Len(t, allItems, 3)
	assert.Equal(t, "8", allItems["POTION_VIE"])
}

// TestRedisPlayerStats vérifie le cache des statistiques joueur
func TestRedisPlayerStats(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewTestRedis(t)
	CleanRedis(t, client)

	ctx := context.Background()
	joueurID := uuid.New()
	keyStats := fmt.Sprintf("player:%s:stats", joueurID)

	// Stocker les stats
	stats := map[string]interface{}{
		"hp_actuel":   "100",
		"hp_max":      "100",
		"mana_actuel": "50",
		"mana_max":    "50",
		"niveau":      "5",
		"xp":          "2500",
	}

	err := client.HSet(ctx, keyStats, stats).Err()
	require.NoError(t, err)

	// Réduire HP
	newHP, err := client.HIncrBy(ctx, keyStats, "hp_actuel", -25).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(75), newHP)

	// Récupérer plusieurs champs
	values, err := client.HMGet(ctx, keyStats, "hp_actuel", "hp_max", "niveau").Result()
	require.NoError(t, err)
	assert.Equal(t, "75", values[0])
	assert.Equal(t, "100", values[1])
	assert.Equal(t, "5", values[2])
}

// TestRedisExpiration vérifie l'expiration automatique des clés
func TestRedisExpiration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewTestRedis(t)
	CleanRedis(t, client)

	ctx := context.Background()
	key := "test:expiration"

	// Créer une clé avec un TTL de 1 seconde
	err := client.Set(ctx, key, "value", 1*time.Second).Err()
	require.NoError(t, err)

	// Vérifier que la clé existe
	exists, err := client.Exists(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(1), exists)

	// Attendre l'expiration
	time.Sleep(1100 * time.Millisecond)

	// La clé devrait avoir disparu
	exists, err = client.Exists(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists)
}

// TestRedisPipeline vérifie l'utilisation des pipelines pour les opérations batch
func TestRedisPipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client := NewTestRedis(t)
	CleanRedis(t, client)

	ctx := context.Background()

	// Créer un pipeline
	pipe := client.Pipeline()

	// Ajouter plusieurs opérations
	combatID := uuid.New()
	pipe.HSet(ctx, fmt.Sprintf("combat:%s:state", combatID), "etat", "EN_COURS")
	pipe.HSet(ctx, fmt.Sprintf("combat:%s:state", combatID), "tour", 1)
	pipe.SAdd(ctx, fmt.Sprintf("combat:%s:participants", combatID), "p1", "p2")
	pipe.Expire(ctx, fmt.Sprintf("combat:%s:state", combatID), time.Hour)

	// Exécuter le pipeline
	_, err := pipe.Exec(ctx)
	require.NoError(t, err)

	// Vérifier les résultats
	etat, err := client.HGet(ctx, fmt.Sprintf("combat:%s:state", combatID), "etat").Result()
	require.NoError(t, err)
	assert.Equal(t, "EN_COURS", etat)

	members, err := client.SMembers(ctx, fmt.Sprintf("combat:%s:participants", combatID)).Result()
	require.NoError(t, err)
	assert.Len(t, members, 2)
}
