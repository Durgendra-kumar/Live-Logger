package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	internal "github.com/Durgendra-kumar/Live-Logger/internal"
	"github.com/segmentio/kafka-go"
)

// Consumer wraps a kafka-go Reader.
// A Reader is the thing that polls Kafka for new messages.
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a Kafka consumer belonging to a named group.
// GroupID is key: two consumers with the same GroupID share the work.
// Two consumers with different GroupIDs each get ALL messages independently.
func NewConsumer(brokerAddr, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{brokerAddr},
			Topic:          topic,       //  topic to read form
			GroupID:        groupID,     // consumer group ID
			MinBytes:       1,           // fetch as soon as 1 byte is available
			MaxBytes:       1 << 20,     // 1MB max per fetch
			CommitInterval: time.Second, // auto-commit offsets every second
		}),
	}
}

// ReadOne blocks until the next message arrives, then returns the decoded event.
// This is the core loop — call it in a for{} to stream continuously.
func (c *Consumer) ReadOne(ctx context.Context) (internal.LogEvent, error) {
	//ctx is just passed as a control signal
	//ReadMessage() is the one that -> talks to Kafka fetches the message returns it
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return internal.LogEvent{}, fmt.Errorf("read message: %w", err)
	}

	var event internal.LogEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return internal.LogEvent{}, fmt.Errorf("unmarshal log event: %w", err)
	}

	return event, nil
}

// Close shuts down the consumer cleanly.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
