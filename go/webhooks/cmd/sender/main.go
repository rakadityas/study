package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"protocols/webhooks/internal/signing"
)

type registerEndpointRequest struct {
	URL    string `json:"url"`
	Secret string `json:"secret"`
}

type registerEndpointResponse struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Active bool   `json:"active"`
}

type createEventRequest struct {
	EventType string         `json:"event_type"`
	Payload   map[string]any `json:"payload"`
}

type createEventResponse struct {
	EventID      string   `json:"event_id"`
	DeliveryIDs  []string `json:"delivery_ids"`
	EndpointCount int     `json:"endpoint_count"`
}

func main() {
	addr := envDefault("WEBHOOK_SENDER_ADDR", ":8084")
	maxAttempts := parseIntDefault(envDefault("MAX_ATTEMPTS", "5"), 5)

	db, err := openDB(context.Background(), envDefault("DATABASE_URL", "postgres://protocol:protocol@localhost:5432/protocols?sslmode=disable"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go worker(ctx, db, maxAttempts)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	})
	mux.HandleFunc("/v1/endpoints", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req registerEndpointRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_json"})
			return
		}
		req.URL = strings.TrimSpace(req.URL)
		req.Secret = strings.TrimSpace(req.Secret)
		if req.URL == "" || req.Secret == "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_payload"})
			return
		}
		var id string
		if err := db.QueryRowContext(r.Context(), `
			INSERT INTO webhook_endpoints(url, secret, active)
			VALUES ($1, $2, true)
			RETURNING id::text
		`, req.URL, req.Secret).Scan(&id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db_error"})
			return
		}
		writeJSON(w, http.StatusCreated, registerEndpointResponse{ID: id, URL: req.URL, Active: true})
	})
	mux.HandleFunc("/v1/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req createEventRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_json"})
			return
		}
		req.EventType = strings.TrimSpace(req.EventType)
		if req.EventType == "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_payload"})
			return
		}
		eventID, deliveryIDs, endpointCount, err := createEvent(r.Context(), db, req.EventType, req.Payload)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db_error"})
			return
		}
		writeJSON(w, http.StatusCreated, createEventResponse{EventID: eventID, DeliveryIDs: deliveryIDs, EndpointCount: endpointCount})
	})
	mux.HandleFunc("/v1/deliveries", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		limit := parseIntDefault(r.URL.Query().Get("limit"), 25)
		if limit < 1 || limit > 200 {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_limit"})
			return
		}
		items, err := listDeliveries(r.Context(), db, limit)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db_error"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
	})

	s := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("webhook sender listening on %s", addr)
	log.Fatal(s.ListenAndServe())
}

func createEvent(ctx context.Context, db *sql.DB, eventType string, payload map[string]any) (eventID string, deliveryIDs []string, endpointCount int, err error) {
	payloadBytes, _ := json.Marshal(payload)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", nil, 0, err
	}
	defer func() { _ = tx.Rollback() }()

	if err := tx.QueryRowContext(ctx, `
		INSERT INTO webhook_events(event_type, payload)
		VALUES ($1, $2::jsonb)
		RETURNING id::text
	`, eventType, string(payloadBytes)).Scan(&eventID); err != nil {
		return "", nil, 0, err
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT id::text
		FROM webhook_endpoints
		WHERE active = true
	`)
	if err != nil {
		return "", nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var endpointID string
		if err := rows.Scan(&endpointID); err != nil {
			return "", nil, 0, err
		}
		var deliveryID string
		if err := tx.QueryRowContext(ctx, `
			INSERT INTO webhook_deliveries(endpoint_id, event_id, attempt, status, next_attempt_at)
			VALUES ($1::uuid, $2::uuid, 0, 'pending', now())
			RETURNING id::text
		`, endpointID, eventID).Scan(&deliveryID); err != nil {
			return "", nil, 0, err
		}
		deliveryIDs = append(deliveryIDs, deliveryID)
		endpointCount++
	}
	if err := rows.Err(); err != nil {
		return "", nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return "", nil, 0, err
	}
	return eventID, deliveryIDs, endpointCount, nil
}

type deliveryRow struct {
	DeliveryID string
	Attempt    int
	URL        string
	Secret     string
	EventID    string
	EventType  string
	Payload    []byte
}

func worker(ctx context.Context, db *sql.DB, maxAttempts int) {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			for i := 0; i < 10; i++ {
				row, err := claimOne(ctx, db)
				if errors.Is(err, sql.ErrNoRows) {
					break
				}
				if err != nil {
					break
				}
				_ = deliverOnce(ctx, db, row, maxAttempts)
			}
		}
	}
}

func claimOne(ctx context.Context, db *sql.DB) (deliveryRow, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return deliveryRow{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var row deliveryRow
	err = tx.QueryRowContext(ctx, `
		SELECT d.id::text, d.attempt, e.url, e.secret, ev.id::text, ev.event_type, ev.payload
		FROM webhook_deliveries d
		JOIN webhook_endpoints e ON e.id = d.endpoint_id
		JOIN webhook_events ev ON ev.id = d.event_id
		WHERE d.next_attempt_at <= now()
		  AND d.status IN ('pending', 'retrying')
		  AND e.active = true
		ORDER BY d.next_attempt_at ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`).Scan(&row.DeliveryID, &row.Attempt, &row.URL, &row.Secret, &row.EventID, &row.EventType, &row.Payload)
	if err != nil {
		return deliveryRow{}, err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE webhook_deliveries
		SET status = 'sending', updated_at = now()
		WHERE id = $1::uuid
	`, row.DeliveryID)
	if err != nil {
		return deliveryRow{}, err
	}

	if err := tx.Commit(); err != nil {
		return deliveryRow{}, err
	}
	return row, nil
}

func deliverOnce(ctx context.Context, db *sql.DB, row deliveryRow, maxAttempts int) error {
	eventBody, _ := json.Marshal(map[string]any{
		"id":         row.EventID,
		"type":       row.EventType,
		"payload":    json.RawMessage(row.Payload),
		"created_at": time.Now().UTC().Format(time.RFC3339Nano),
	})

	ts := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	sig := signing.SignV1(row.Secret, ts, eventBody)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, row.URL, bytes.NewReader(eventBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Timestamp", ts)
	req.Header.Set("X-Webhook-Signature", "v1="+sig)
	req.Header.Set("X-Webhook-Event", row.EventType)
	req.Header.Set("X-Webhook-Delivery", row.DeliveryID)

	client := &http.Client{Timeout: 3 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return markAttempt(ctx, db, row.DeliveryID, row.Attempt, 0, err.Error(), maxAttempts)
	}
	defer res.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(res.Body, 1024*1024))

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		_, err := db.ExecContext(ctx, `
			UPDATE webhook_deliveries
			SET status = 'delivered',
			    last_status_code = $2,
			    last_error = NULL,
			    updated_at = now()
			WHERE id = $1::uuid
		`, row.DeliveryID, res.StatusCode)
		return err
	}
	return markAttempt(ctx, db, row.DeliveryID, row.Attempt, res.StatusCode, "non_2xx", maxAttempts)
}

func markAttempt(ctx context.Context, db *sql.DB, deliveryID string, prevAttempt int, statusCode int, errMsg string, maxAttempts int) error {
	nextAttempt := prevAttempt + 1
	if nextAttempt >= maxAttempts {
		_, err := db.ExecContext(ctx, `
			UPDATE webhook_deliveries
			SET attempt = $2,
			    status = 'failed',
			    last_status_code = NULLIF($3, 0),
			    last_error = $4,
			    updated_at = now()
			WHERE id = $1::uuid
		`, deliveryID, nextAttempt, statusCode, errMsg)
		return err
	}
	backoff := time.Duration(1<<nextAttempt) * time.Second
	_, err := db.ExecContext(ctx, `
		UPDATE webhook_deliveries
		SET attempt = $2,
		    status = 'retrying',
		    last_status_code = NULLIF($3, 0),
		    last_error = $4,
		    next_attempt_at = now() + $5::interval,
		    updated_at = now()
		WHERE id = $1::uuid
	`, deliveryID, nextAttempt, statusCode, errMsg, fmtInterval(backoff))
	return err
}

func fmtInterval(d time.Duration) string {
	return strconv.FormatFloat(d.Seconds(), 'f', 0, 64) + " seconds"
}

type deliveryView struct {
	ID            string `json:"id"`
	EndpointID    string `json:"endpoint_id"`
	EventID       string `json:"event_id"`
	Attempt       int    `json:"attempt"`
	Status        string `json:"status"`
	LastError     string `json:"last_error"`
	LastStatus    int    `json:"last_status_code"`
	NextAttemptAt string `json:"next_attempt_at"`
	UpdatedAt     string `json:"updated_at"`
}

func listDeliveries(ctx context.Context, db *sql.DB, limit int) ([]deliveryView, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id::text, endpoint_id::text, event_id::text, attempt, status,
		       COALESCE(last_error, ''), COALESCE(last_status_code, 0),
		       next_attempt_at, updated_at
		FROM webhook_deliveries
		ORDER BY updated_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []deliveryView
	for rows.Next() {
		var v deliveryView
		var next time.Time
		var updated time.Time
		if err := rows.Scan(&v.ID, &v.EndpointID, &v.EventID, &v.Attempt, &v.Status, &v.LastError, &v.LastStatus, &next, &updated); err != nil {
			return nil, err
		}
		v.NextAttemptAt = next.UTC().Format(time.RFC3339Nano)
		v.UpdatedAt = updated.UTC().Format(time.RFC3339Nano)
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
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
