package internal

import "time"

type Level string

const (
	DEBUG Level = "DEBUG"
	INFO  Level = "INFO"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
)

// LogEvent is what every producer publishes and every consumer reads.
// This is your Kafka message payload — must be serializable to JSON.
type LogEvent struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"` // e.g. "auth-service", "payment-service"
	Level     Level     `json:"level"`
	Message   string    `json:"message"`
}
