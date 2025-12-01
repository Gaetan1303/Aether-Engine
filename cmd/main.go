package main

import (
	"context"
	"log"
	"os"

	"github.com/aether-engine/aether-engine/api/handlers"
	"github.com/aether-engine/aether-engine/internal/combat/application"
	"github.com/aether-engine/aether-engine/internal/combat/infrastructure"
	"github.com/aether-engine/aether-engine/pkg/eventbus"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Configuration depuis les variables d'environnement
	dbURL := getEnv("DATABASE_URL", "postgres://test:test@localhost:5432/aether_test?sslmode=disable")
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("KAFKA_TOPIC", "combat-events")
	port := getEnv("PORT", "8080")

	// Connexion à PostgreSQL
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Erreur connexion PostgreSQL: %v", err)
	}
	defer pool.Close()

	// Initialiser le schéma
	if err := infrastructure.InitSchema(pool); err != nil {
		log.Fatalf("Erreur initialisation schéma: %v", err)
	}

	log.Println("Connexion PostgreSQL établie")

	// Créer l'Event Store
	eventStore := infrastructure.NewPostgresEventStore(pool)

	// Créer l'Event Publisher (Kafka)
	publisher := eventbus.NewKafkaEventPublisher([]string{kafkaBrokers}, kafkaTopic)
	defer publisher.Close()

	log.Println("Event Publisher Kafka créé")

	// Créer le Combat Engine
	combatEngine := application.NewCombatEngine(eventStore, publisher)

	// Créer le router Gin
	router := gin.Default()

	// Health check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Créer les handlers
	combatHandler := handlers.NewCombatHandler(combatEngine)

	// Enregistrer les routes
	api := router.Group("/api/v1")
	combatHandler.RegisterRoutes(api)

	// Démarrer le serveur
	log.Printf("Serveur Fabric démarré sur le port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Erreur démarrage serveur: %v", err)
	}
}

// getEnv récupère une variable d'environnement avec une valeur par défaut
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
