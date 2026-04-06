package main

import (
	"context"
	"database/sql"
	"encoding/json"
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

func main() {
	addr := envDefault("WEBHOOK_RECEIVER_ADDR", ":8083")
	secret := envDefault("WEBHOOK_SECRET", "receiver-secret")
	failPercent := parseIntDefault(envDefault("FAIL_PERCENT", "0"), 0)

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
	mux.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Query().Get("fail") == "true" || (failPercent > 0 && time.Now().UnixNano()%100 < int64(failPercent)) {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "simulated_failure"})
			return
		}

		raw, err := io.ReadAll(io.LimitReader(r.Body, 1024*1024))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_body"})
			return
		}
		ts := strings.TrimSpace(r.Header.Get("X-Webhook-Timestamp"))
		sig := strings.TrimSpace(r.Header.Get("X-Webhook-Signature"))
		if ts == "" || sig == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "missing_signature"})
			return
		}
		if !strings.HasPrefix(sig, "v1=") {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid_signature_format"})
			return
		}
		if !signing.VerifyV1(secret, ts, raw, strings.TrimPrefix(sig, "v1=")) {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid_signature"})
			return
		}

		var body any
		if err := json.Unmarshal(raw, &body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_json"})
			return
		}

		headers := map[string]string{}
		for k, vv := range r.Header {
			if len(vv) > 0 {
				headers[k] = vv[0]
			}
		}
		hb, _ := json.Marshal(headers)
		bb, _ := json.Marshal(body)

		_, _ = db.ExecContext(r.Context(), `
			INSERT INTO webhook_received(headers, body)
			VALUES ($1::jsonb, $2::jsonb)
		`, string(hb), string(bb))

		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	})
	mux.HandleFunc("/v1/received", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		limit := parseIntDefault(r.URL.Query().Get("limit"), 25)
		if limit < 1 || limit > 200 {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_limit"})
			return
		}
		items, err := listReceived(r.Context(), db, limit)
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
	log.Printf("webhook receiver listening on %s", addr)
	log.Fatal(s.ListenAndServe())
}

type receivedItem struct {
	ID         string `json:"id"`
	ReceivedAt string `json:"received_at"`
	Body       any    `json:"body"`
}

func listReceived(ctx context.Context, db *sql.DB, limit int) ([]receivedItem, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id::text, received_at, body
		FROM webhook_received
		ORDER BY received_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []receivedItem
	for rows.Next() {
		var it receivedItem
		var receivedAt time.Time
		var bodyRaw []byte
		if err := rows.Scan(&it.ID, &receivedAt, &bodyRaw); err != nil {
			return nil, err
		}
		it.ReceivedAt = receivedAt.UTC().Format(time.RFC3339Nano)
		_ = json.Unmarshal(bodyRaw, &it.Body)
		out = append(out, it)
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
