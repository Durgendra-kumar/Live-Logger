// this is the log generator
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os/signal"
	"syscall"
	"time"

	internal "github.com/Durgendra-kumar/Live-Logger/internal"
	internalkafka "github.com/Durgendra-kumar/Live-Logger/internal/kafka"
	"github.com/google/uuid"
)

// Fake services — pretend these are real microservices in your system
var services = []string{
	"auth-service",
	"payment-service",
	"order-service",
	"inventory-service",
}

var levels = []internal.Level{
	internal.DEBUG,
	internal.INFO,
	internal.INFO, // INFO is most common in real systems
	internal.INFO,
	internal.WARN,
	internal.ERROR,
}

var messages = map[internal.Level][]string{
	internal.DEBUG: {"Cache miss for key user:42", "DB query took 12ms", "JWT decoded successfully"},
	internal.INFO:  {"User logged in", "Order #1234 placed", "Payment processed", "Item added to cart"},
	internal.WARN:  {"Retry attempt 2/3", "Response time > 500ms", "Memory usage at 80%"},
	internal.ERROR: {"DB connection failed", "Payment gateway timeout", "Unhandled exception in handler"},
}

func main() {
	producer := internalkafka.NewProducer("localhost:9092")
	defer producer.Close()

	// Listen for Ctrl+C to shut down gracefully
	//Creates a context that is automatically cancelled when OS signals arrive
	// context.Background() → base context (no timeout, no cancel)
	// signal.NotifyContext(...) -> wraps context and listens for OS signals: ctrl+c or contaner stop or kill command
	// When signal comes → ctx.Done() is triggered
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fmt.Println("Producer started — publishing logs every 800ms. Ctrl+C to stop.")
	fmt.Println()

	//A ticker sends a signal at a fixed interval
	ticker := time.NewTicker(800 * time.Millisecond)
	defer ticker.Stop()

	for {
		select { // Waits for multiple events at same time
		case <-ctx.Done():
			fmt.Println("\nShutting down producer...")
			return
		case <-ticker.C:
			event := randomEvent()
			if err := producer.Publish(ctx, event); err != nil {
				log.Printf("publish error: %v", err)
				continue
			}
			fmt.Printf("published: [%s] %s — %s\n", event.Level, event.Service, event.Message)
		}
	}
}

func randomEvent() internal.LogEvent {
	svc := services[rand.Intn(len(services))]
	lvl := levels[rand.Intn(len(levels))]
	msgs := messages[lvl]
	msg := msgs[rand.Intn(len(msgs))]

	return internal.LogEvent{
		ID:        uuid.NewString(),
		Timestamp: time.Now(),
		Service:   svc,
		Level:     lvl,
		Message:   msg,
	}
}
