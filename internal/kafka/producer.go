package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	internal "github.com/Durgendra-kumar/Live-Logger/internal"
	"github.com/segmentio/kafka-go"
)

const topic = "app.logs"

// Producer wraps a kafka-go Writer.
// A Writer is the thing that sends messages to Kafka.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a connected Kafka producer.
func NewProducer(brokerAddr string) *Producer {

	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokerAddr), // tells where kafka is running
			Topic:        topic,                 // tell about topic
			Balancer:     &kafka.LeastBytes{},   // inseter in partition which as less Byte
			WriteTimeout: 5 * time.Second,
		},
	}
}

// Publish sends one LogEvent to Kafka.
// The event is JSON-encoded — Kafka only cares about bytes.
func (p *Producer) Publish(ctx context.Context, event internal.LogEvent) error {
	data, err := json.Marshal(event) // convet struct to josn
	if err != nil {
		return fmt.Errorf("marshal log event %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.Service), // same service → same partition → ordered logs
		Value: data,
	})
}

// Close flushes and closes the writer. Always call this on shutdown.
func (p *Producer) Close() error {
	return p.writer.Close()
}
