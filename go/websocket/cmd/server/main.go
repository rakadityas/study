package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

type chatIn struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

type chatOut struct {
	Type      string `json:"type"`
	Room      string `json:"room"`
	Sender    string `json:"sender"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type hub struct {
	mu    sync.Mutex
	rooms map[string]map[*client]struct{}
}

type client struct {
	room  string
	user  string
	conn  *websocket.Conn
	send  chan []byte
	hub   *hub
	db    *sql.DB
	close chan struct{}
}

func newHub() *hub {
	return &hub{rooms: map[string]map[*client]struct{}{}}
}

func (h *hub) add(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[c.room] == nil {
		h.rooms[c.room] = map[*client]struct{}{}
	}
	h.rooms[c.room][c] = struct{}{}
}

func (h *hub) remove(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[c.room] != nil {
		delete(h.rooms[c.room], c)
		if len(h.rooms[c.room]) == 0 {
			delete(h.rooms, c.room)
		}
	}
}

func (h *hub) broadcast(room string, msg []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.rooms[room] {
		select {
		case c.send <- msg:
		default:
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	addr := envDefault("WS_ADDR", ":8082")
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
	mux.HandleFunc("/v1/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		room := strings.TrimSpace(r.URL.Query().Get("room"))
		if room == "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "room_required"})
			return
		}
		limit := parseIntDefault(r.URL.Query().Get("limit"), 50)
		if limit < 1 || limit > 500 {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_limit"})
			return
		}
		items, err := listMessages(r.Context(), db, room, limit)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "db_error"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		room := strings.TrimSpace(r.URL.Query().Get("room"))
		user := strings.TrimSpace(r.URL.Query().Get("user"))
		if room == "" || user == "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "room_and_user_required"})
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		c := &client{
			room:  room,
			user:  user,
			conn:  conn,
			send:  make(chan []byte, 16),
			hub:   h,
			db:    db,
			close: make(chan struct{}),
		}
		h.add(c)
		go c.writePump()
		go c.readPump()
	})

	s := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("websocket listening on %s", addr)
	log.Fatal(s.ListenAndServe())
}

func (c *client) readPump() {
	defer func() {
		c.hub.remove(c)
		close(c.close)
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(64 * 1024)
	_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var in chatIn
		if err := c.conn.ReadJSON(&in); err != nil {
			return
		}
		if strings.TrimSpace(in.Type) != "chat.message" {
			continue
		}
		body := strings.TrimSpace(in.Body)
		if body == "" {
			continue
		}
		now := time.Now().UTC()
		_, _ = c.db.ExecContext(context.Background(), `
			INSERT INTO ws_messages(room, sender, body)
			VALUES ($1, $2, $3)
		`, c.room, c.user, body)
		out := chatOut{
			Type:      "chat.message",
			Room:      c.room,
			Sender:    c.user,
			Body:      body,
			CreatedAt: now.Format(time.RFC3339Nano),
		}
		b, _ := json.Marshal(out)
		c.hub.broadcast(c.room, b)
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.close:
			return
		case msg := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
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

func listMessages(ctx context.Context, db *sql.DB, room string, limit int) ([]chatOut, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT room, sender, body, created_at
		FROM ws_messages
		WHERE room = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, room, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []chatOut
	for rows.Next() {
		var m chatOut
		var created time.Time
		if err := rows.Scan(&m.Room, &m.Sender, &m.Body, &created); err != nil {
			return nil, err
		}
		m.Type = "chat.message"
		m.CreatedAt = created.UTC().Format(time.RFC3339Nano)
		out = append(out, m)
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
