# HTTP (REST/JSON) Example

## What this demonstrates

- Request/response over HTTP with JSON payloads.
- Stateless handlers (each request contains all the information needed).
- Typical use cases: CRUD APIs, public web APIs, internal microservices when streaming/bidirectional is not required.

## Run

Start the shared DB:

```bash
cd go
docker compose up -d db
```

Run the server:

```bash
cd go/http
DATABASE_URL='postgres://protocol:protocol@localhost:5432/protocols?sslmode=disable' \
HTTP_ADDR=':8080' \
go run ./cmd/server
```

Run the client:

```bash
cd go/http
HTTP_BASE_URL='http://localhost:8080' go run ./cmd/client
```

## API

### Create order

Request:

```http
POST /v1/orders
Content-Type: application/json

{"customer":"alice","amount_cents":2599}
```

Responses:

- `201 Created`

```json
{"id":"...","customer":"alice","amount_cents":2599,"status":"created","created_at":"..."}
```

- `400 Bad Request`

```json
{"error":"invalid_json"}
```

```json
{"error":"invalid_payload"}
```

- `500 Internal Server Error`

```json
{"error":"db_error"}
```

### Get order

Request:

```http
GET /v1/orders/{id}
```

Responses:

- `200 OK` (same JSON shape as create)
- `404 Not Found`

```json
{"error":"not_found"}
```

### List orders

Request:

```http
GET /v1/orders?limit=25
```

Response:

```json
{"items":[{"id":"...","customer":"...","amount_cents":123,"status":"created","created_at":"..."}]}
```

