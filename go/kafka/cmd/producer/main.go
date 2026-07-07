// producer demonstrates keyed writes: every event is published with
// its CustomerID as the message key. Kafka hashes the key to pick a
// partition, so one customer's events always go to the same
// partition and are therefore always read back in order -- ordering
// in Kafka is a per-partition guarantee, not a per-topic one.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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
	count, _ := strconv.Atoi(getEnv("PRODUCE_COUNT", "20"))

	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.Hash{}, // partition chosen from the message key
		// RequireAll waits for every in-sync replica to ack -- the
		// strongest durability setting kafka-go exposes, at the cost
		// of latency. Try kafka.RequireOne to feel the difference.
		RequiredAcks: kafka.RequireAll,
	}
	defer writer.Close()

	customers := []string{"cust-1", "cust-2", "cust-3"}
	items := []string{"widget", "gadget", "gizmo"}

	for i := 0; i < count; i++ {
		customer := customers[i%len(customers)]
		order := events.OrderEvent{
			OrderID:    fmt.Sprintf("order-%04d", i),
			CustomerID: customer,
			Item:       items[i%len(items)],
			Quantity:   (i % 5) + 1,
			Status:     "created",
			CreatedAt:  time.Now(),
		}

		payload, err := json.Marshal(order)
		if err != nil {
			log.Fatalf("marshal: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = writer.WriteMessages(ctx, kafka.Message{
			Key:   []byte(customer),
			Value: payload,
			Headers: []kafka.Header{
				{Key: "event-type", Value: []byte("order.created")},
			},
		})
		cancel()
		if err != nil {
			log.Fatalf("write message: %v", err)
		}

		fmt.Printf("produced %s key=%s\n", order.OrderID, customer)
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("done")
}
