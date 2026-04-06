package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func main() {
	addr := envDefault("HTTP_ADDR", ":8080")
	db, err := openDB(context.Background(), envDefault("DATABASE_URL", "postgres://protocol:protocol@localhost:5432/protocols?sslmode=disable"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	})
	mux.HandleFunc("/v1/orders", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		switch r.Method {
		case http.MethodPost:
		case http.MethodGet:
			limit := parseIntDefault(r.URL.Query().Get("limit"), 25)
			if limit < 1 || limit > 200 {
				writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_limit"})
				return
			}
			orders, err := listOrders(ctx, db, limit)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db_error"})
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"items": orders})
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

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

		o, err := createOrder(ctx, db, req.Customer, req.AmountCents)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db_error"})
			return
		}

		writeJSON(w, http.StatusCreated, o)
	})
	mux.HandleFunc("/v1/orders/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		ctx := r.Context()
		id := strings.TrimPrefix(r.URL.Path, "/v1/orders/")
		id = strings.TrimSpace(id)
		if id == "" || strings.Contains(id, "/") {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_id"})
			return
		}
		o, err := getOrder(ctx, db, id)
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not_found"})
			return
		}
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db_error"})
			return
		}
		writeJSON(w, http.StatusOK, o)
	})

	s := &http.Server{
		Addr:              addr,
		Handler:           requestLog(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("http listening on %s", addr)
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

func createOrder(ctx context.Context, db *sql.DB, customer string, amountCents int) (order, error) {
	var o order
	row := db.QueryRowContext(ctx, `
		INSERT INTO orders(customer, amount_cents, status)
		VALUES ($1, $2, 'created')
		RETURNING id::text, customer, amount_cents, status, created_at
	`, customer, amountCents)
	if err := row.Scan(&o.ID, &o.Customer, &o.AmountCents, &o.Status, &o.CreatedAt); err != nil {
		return order{}, err
	}
	payload, _ := json.Marshal(map[string]any{
		"order_id":      o.ID,
		"customer":      o.Customer,
		"amount_cents":  o.AmountCents,
		"status":        o.Status,
		"occurred_at":   time.Now().UTC().Format(time.RFC3339Nano),
		"protocol_demo": "http",
	})
	_, _ = db.ExecContext(ctx, `
		INSERT INTO order_events(order_id, event_type, payload)
		VALUES ($1, $2, $3::jsonb)
	`, o.ID, "order.created", string(payload))
	return o, nil
}

func getOrder(ctx context.Context, db *sql.DB, id string) (order, error) {
	var o order
	row := db.QueryRowContext(ctx, `
		SELECT id::text, customer, amount_cents, status, created_at
		FROM orders
		WHERE id = $1::uuid
	`, id)
	if err := row.Scan(&o.ID, &o.Customer, &o.AmountCents, &o.Status, &o.CreatedAt); err != nil {
		return order{}, err
	}
	return o, nil
}

func listOrders(ctx context.Context, db *sql.DB, limit int) ([]order, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id::text, customer, amount_cents, status, created_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []order
	for rows.Next() {
		var o order
		if err := rows.Scan(&o.ID, &o.Customer, &o.AmountCents, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

func envDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, fmt.Sprintf("%s", time.Since(start)))
	})
}
