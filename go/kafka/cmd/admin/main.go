// admin creates (or inspects) the topic used by the producer/consumer
// examples, so partition count and replication factor are explicit
// and reproducible instead of relying on broker auto-create defaults.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

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
	partitions, _ := strconv.Atoi(getEnv("KAFKA_PARTITIONS", "3"))
	replication, _ := strconv.Atoi(getEnv("KAFKA_REPLICATION", "1"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := kafka.DialContext(ctx, "tcp", broker)
	if err != nil {
		log.Fatalf("dial broker: %v", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("find controller: %v", err)
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatalf("dial controller: %v", err)
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replication,
	})
	if err != nil && err != kafka.TopicAlreadyExists {
		log.Fatalf("create topic: %v", err)
	}
	fmt.Printf("topic %q ready (partitions=%d, replication=%d)\n", topic, partitions, replication)

	partitionsInfo, err := conn.ReadPartitions(topic)
	if err != nil {
		log.Fatalf("read partitions: %v", err)
	}
	for _, p := range partitionsInfo {
		fmt.Printf("  partition %d leader=%d replicas=%v\n", p.ID, p.Leader.ID, p.Replicas)
	}
}
