package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	internalkafka "github.com/Durgendra-kumar/Live-Logger/internal/kafka"
	"github.com/Durgendra-kumar/Live-Logger/internal/logger"
)

func main() {
	consumer := internalkafka.NewConsumer("localhost:9092", "log-viewer")
	defer consumer.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fmt.Println("Consumer started — watching for logs. Ctrl+C to stop.")
	fmt.Println()

	// This loop is the entire consumer — it blocks on ReadOne() until
	// a message arrives, then prints it, then waits for the next one.
	// This is identical to how your Kafka notification worker will work.
	for {
		event, err := consumer.ReadOne(ctx)
		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("\nShutting down consumer...")
				return
			}
			log.Printf("read error: %v", err)
			continue
		}

		logger.Print(event)
	}
}
