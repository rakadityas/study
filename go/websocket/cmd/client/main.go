package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type chatIn struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

func main() {
	baseURL := envDefault("WS_BASE_URL", "ws://localhost:8082/ws")
	room := envDefault("WS_ROOM", "general")
	user := envDefault("WS_USER", "alice")
	msg := envDefault("WS_MESSAGE", "hello from client")

	u := fmt.Sprintf("%s?room=%s&user=%s", baseURL, room, user)
	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			fmt.Println(string(message))
		}
	}()

	b, _ := json.Marshal(chatIn{Type: "chat.message", Body: msg})
	if err := c.WriteMessage(websocket.TextMessage, b); err != nil {
		log.Fatal(err)
	}

	select {
	case <-done:
	case <-time.After(3 * time.Second):
	}
}

func envDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

