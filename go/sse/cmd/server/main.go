package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type order struct {
	ID          string    `json:"id"`
	Customer    string    `json:"customer"`
	AmountCents int       `json:"amount_cents"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type createOrderRequest struct {
	Customer    string `json:"customer"`
	AmountCents int    `json:"amount_cents"`
}

type sseEvent struct {
	Name string
	Data []byte
}

type hub struct {
	mu      sync.Mutex
	clients map[chan sseEvent]struct{}
}

func newHub() *hub {
	return &hub{clients: map[chan sseEvent]struct{}{}}
}

func (h *hub) subscribe() chan sseEvent {
	ch := make(chan sseEvent, 16)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *hub) unsubscribe(ch chan sseEvent) {
	h.mu.Lock()
	delete(h.clients, ch)
	h.mu.Unlock()
	close(ch)
}

func (h *hub) publish(evt sseEvent) {
	h.mu.Lock()
	for ch := range h.clients {
		select {
		case ch <- evt:
		default:
		}
	}
	h.mu.Unlock()
}

func main() {
	addr := envDefault("SSE_ADDR", ":8081")
	db, err := openDB(context.Background(), envDefault("DATABASE_URL", "postgres://protocol:protocol@localhost:5432/protocols?sslmode=disable"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	h := newHub()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	})
	mux.HandleFunc("/v1/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		ctx := r.Context()
		var req createOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_json"})
			return
		}
		req.Customer = strings.TrimSpace(req.Customer)
		if req.Customer == "" || req.AmountCents <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_payload"})
			return
		}
		o, payload, err := createOrder(ctx, db, req.Customer, req.AmountCents)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db_error"})
			return
		}
		h.publish(sseEvent{Name: "order.created", Data: payload})
		writeJSON(w, http.StatusCreated, o)
	})
	mux.HandleFunc("/v1/stream/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		flusher, ok := w.(http.Flusher)
		if !ok {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "stream_not_supported"})
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch := h.subscribe()
		defer h.unsubscribe(ch)

		fmt.Fprintf(w, "event: ready\ndata: %s\n\n", `{"ok":true}`)
		flusher.Flush()

		keepAlive := time.NewTicker(15 * time.Second)
		defer keepAlive.Stop()

		for {
			select {
			case <-r.Context().Done():
				return
			case <-keepAlive.C:
				fmt.Fprint(w, ": ping\n\n")
				flusher.Flush()
			case evt := <-ch:
				fmt.Fprintf(w, "event: %s\ndata: %s\n\n", evt.Name, evt.Data)
				flusher.Flush()
			}
		}
	})

	s := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("sse listening on %s", addr)
	log.Fatal(s.ListenAndServe())
}

func openDB(ctx context.Context, url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

func createOrder(ctx context.Context, db *sql.DB, customer string, amountCents int) (order, []byte, error) {
	var o order
	row := db.QueryRowContext(ctx, `
		INSERT INTO orders(customer, amount_cents, status)
		VALUES ($1, $2, 'created')
		RETURNING id::text, customer, amount_cents, status, created_at
	`, customer, amountCents)
	if err := row.Scan(&o.ID, &o.Customer, &o.AmountCents, &o.Status, &o.CreatedAt); err != nil {
		return order{}, nil, err
	}
	payload, _ := json.Marshal(map[string]any{
		"order": map[string]any{
			"id":           o.ID,
			"customer":     o.Customer,
			"amount_cents": o.AmountCents,
			"status":       o.Status,
			"created_at":   o.CreatedAt.UTC().Format(time.RFC3339Nano),
		},
	})
	_, _ = db.ExecContext(ctx, `
		INSERT INTO order_events(order_id, event_type, payload)
		VALUES ($1::uuid, $2, $3::jsonb)
	`, o.ID, "order.created", string(payload))
	return o, payload, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func envDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
