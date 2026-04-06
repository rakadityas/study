# gRPC Example (Unary RPC over HTTP/2)

## What this demonstrates

- Unary RPC calls (request/response) using gRPC over HTTP/2.
- A service boundary that is *method-based* (RPC) rather than resource-based (REST).
- Typical use cases: internal microservices, strongly-typed contracts, polyglot backends, low-latency service-to-service calls.

## Note about protobuf/codegen

Production gRPC typically uses Protocol Buffers (`.proto`) and generated code.

This example intentionally avoids `protoc` so it can run in a minimal environment: it registers a gRPC service manually and uses a JSON codec. The transport is still gRPC (HTTP/2 + gRPC framing), but the payload format is JSON.

## Run

Start the shared DB:

```bash
cd go
docker compose up -d db
```

Run the server:

```bash
cd go/grpc
DATABASE_URL='postgres://protocol:protocol@localhost:5432/protocols?sslmode=disable' \
GRPC_ADDR=':9090' \
go run ./cmd/server
```

Run the client:

```bash
cd go/grpc
GRPC_TARGET='localhost:9090' go run ./cmd/client
```

## Methods and payloads

### CreateOrder

Method:

```text
/orders.v1.OrderService/CreateOrder
```

Request JSON:

```json
{"customer":"alice","amount_cents":2599}
```

Response JSON (success):

```json
{"order":{"id":"...","customer":"alice","amount_cents":2599,"status":"created","created_at":"..."}}
```

Potential gRPC error codes:

- `InvalidArgument`: invalid JSON or invalid fields
- `Internal`: DB errors

### GetOrder

Request JSON:

```json
{"id":"<uuid>"}
```

Potential gRPC error codes:

- `NotFound`: order does not exist
- `InvalidArgument`: invalid id

### ListOrders

Request JSON:

```json
{"limit":5}
```

Response JSON:

```json
{"items":[{"id":"...","customer":"...","amount_cents":123,"status":"created","created_at":"..."}]}
```

