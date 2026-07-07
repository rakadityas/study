package events

import "time"

// OrderEvent is a small domain event used to demonstrate keyed
// ordering: all events for the same CustomerID are produced with
// that ID as the Kafka message key, so they always land on the same
// partition and are read back in the order they were written.
type OrderEvent struct {
	OrderID    string    `json:"order_id"`
	CustomerID string    `json:"customer_id"`
	Item       string    `json:"item"`
	Quantity   int       `json:"quantity"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
