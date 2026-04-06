# WebSocket (TCP) Example

## What this demonstrates

- Full-duplex, low-latency messaging over one TCP connection.
- Typical use cases: chat, collaborative editing, live trading UIs, multiplayer games, device control panels.
- A DB-backed message history endpoint alongside the live WebSocket feed.

## Run

Start the shared DB:

```bash
cd go
docker compose up -d db
```

Run the server:

```bash
cd go/websocket
DATABASE_URL='postgres://protocol:protocol@localhost:5432/protocols?sslmode=disable' \
WS_ADDR=':8082' \
go run ./cmd/server
```

Run a sample client (sends one message and prints whatever it receives for ~3 seconds):

```bash
cd go/websocket
WS_ROOM='general' WS_USER='alice' WS_MESSAGE='hello' go run ./cmd/client
```

## WebSocket endpoint

Connect:

```text
GET ws://localhost:8082/ws?room=general&user=alice
```

Client → server message payload:

```json
{"type":"chat.message","body":"hello"}
```

Server → clients broadcast payload:

```json
{"type":"chat.message","room":"general","sender":"alice","body":"hello","created_at":"2026-01-01T00:00:00Z"}
```

## HTTP endpoints

### List recent messages

Request:

```http
GET /v1/messages?room=general&limit=50
```

Responses:

- `200 OK`

```json
{"items":[{"type":"chat.message","room":"general","sender":"alice","body":"hello","created_at":"..."}]}
```

- `400 Bad Request`

```json
{"error":"room_required"}
```

```json
{"error":"invalid_limit"}
```

