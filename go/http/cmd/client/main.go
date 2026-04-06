package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	baseURL := envDefault("HTTP_BASE_URL", "http://localhost:8080")

	body, _ := json.Marshal(map[string]any{
		"customer":     "alice",
		"amount_cents": 2599,
	})
	res, err := http.Post(baseURL+"/v1/orders", "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)
	fmt.Printf("POST /v1/orders -> %d\n%s\n\n", res.StatusCode, string(b))

	var created map[string]any
	_ = json.Unmarshal(b, &created)
	id, _ := created["id"].(string)
	if id == "" {
		return
	}

	res2, err := http.Get(baseURL + "/v1/orders/" + id)
	if err != nil {
		panic(err)
	}
	defer res2.Body.Close()
	b2, _ := io.ReadAll(res2.Body)
	fmt.Printf("GET /v1/orders/%s -> %d\n%s\n", id, res2.StatusCode, string(b2))
}

func envDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

