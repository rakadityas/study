// consumer demonstrates consumer-group behavior: run several
// instances with the same KAFKA_GROUP_ID and the broker splits the
// topic's partitions between them, rebalancing automatically as
// instances join or leave. Run instances with different group IDs
// instead and each group gets its own full copy of every message --
// consumer groups are how Kafka supports both queue-like (shared)
// and pub/sub-like (broadcast) consumption on the same topic.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"protocols/kafka/internal/events"

	kafka "github.com/segmentio/kafka-go"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := getEnv("KAFKA_TOPIC", "orders.v1")
	groupID := getEnv("KAFKA_GROUP_ID", "order-processors")
	consumerName := getEnv("CONSUMER_NAME", "consumer-1")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID, // presence of GroupID is what enables group membership + rebalancing

		// CommitInterval == 0 means we commit offsets manually after
		// processing, giving at-least-once delivery: a crash between
		// FetchMessage and CommitMessages replays the message on
		// restart/rebalance instead of silently dropping it.
		CommitInterval: 0,

		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Printf("[%s] joining group %q on topic %q\n", consumerName, groupID, topic)

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				fmt.Printf("[%s] shutting down\n", consumerName)
				return
			}
			log.Fatalf("fetch message: %v", err)
		}

		var order events.OrderEvent
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("[%s] skipping malformed message: %v", consumerName, err)
			continue
		}

		fmt.Printf("[%s] partition=%d offset=%d key=%s order=%s item=%s qty=%d\n",
			consumerName, msg.Partition, msg.Offset, msg.Key, order.OrderID, order.Item, order.Quantity)

		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("[%s] commit failed: %v", consumerName, err)
		}
	}
}
