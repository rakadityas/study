# SSE (Server-Sent Events) Example

## What this demonstrates

- One-way server → client streaming over plain HTTP.
- Automatic reconnect behavior on the client side (typical SSE clients do this).
- Typical use cases: live feeds, notifications, dashboards, progress updates.

SSE is a good fit when:

- You need server push, but not client-to-server messages over the same channel.
- You want to keep things HTTP-friendly (proxies, auth, infra).

## Run

Start the shared DB:

```bash
cd go
docker compose up -d db
```

Run the SSE server:

```bash
cd go/sse
DATABASE_URL='postgres://protocol:protocol@localhost:5432/protocols?sslmode=disable' \
SSE_ADDR=':8081' \
go run ./cmd/server
```

## Try it

Terminal A: subscribe (SSE stream):

```bash
curl -N http://localhost:8081/v1/stream/orders
```

Terminal B: create an order (triggers an SSE event):

```bash
curl -i http://localhost:8081/v1/orders \
  -H 'Content-Type: application/json' \
  -d '{"customer":"alice","amount_cents":2599}'
```

## Event format

You will see lines like:

```text
event: order.created
data: {"order":{"id":"...","customer":"alice","amount_cents":2599,"status":"created","created_at":"..."}}
```

Keep-alives are sent as comments:

```text
: ping
```

## HTTP endpoints

### Create order

Request:

```http
POST /v1/orders
Content-Type: application/json

{"customer":"alice","amount_cents":2599}
```

Response:

- `201 Created` with the created order JSON
- `400 Bad Request` with `{"error":"invalid_json"}` or `{"error":"invalid_payload"}`
- `500 Internal Server Error` with `{"error":"db_error"}`

