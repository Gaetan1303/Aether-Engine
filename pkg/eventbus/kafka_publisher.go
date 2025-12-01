package eventbus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aether-engine/aether-engine/internal/combat/domain"
	"github.com/segmentio/kafka-go"
)

// KafkaEventPublisher implémente EventPublisher avec Kafka
type KafkaEventPublisher struct {
	writer *kafka.Writer
	topic  string
}

// NewKafkaEventPublisher crée une nouvelle instance de KafkaEventPublisher
func NewKafkaEventPublisher(brokers []string, topic string) *KafkaEventPublisher {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaEventPublisher{
		writer: writer,
		topic:  topic,
	}
}

// Publish publie un événement sur Kafka
func (p *KafkaEventPublisher) Publish(event domain.Evenement) error {
	// Sérialiser l'événement
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("erreur sérialisation événement: %w", err)
	}

	// Créer le message Kafka
	message := kafka.Message{
		Key:   []byte(event.AggregateID()),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(event.EventType())},
			{Key: "aggregate_id", Value: []byte(event.AggregateID())},
		},
	}

	// Publier le message
	ctx := context.Background()
	if err := p.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("erreur publication événement: %w", err)
	}

	return nil
}

// Close ferme le writer Kafka
func (p *KafkaEventPublisher) Close() error {
	return p.writer.Close()
}
