# Webhooks Example

## What this demonstrates

- Webhooks are HTTP callbacks used for asynchronous event delivery.
- Typical use cases: payments (Stripe), Git providers (GitHub), SaaS integrations, audit/event pipelines.
- Delivery concerns: signing, retries, idempotency, backpressure, failure handling.

This folder contains two services:

- `sender`: stores events + schedules deliveries + retries
- `receiver`: validates signatures and stores received payloads

## Run

Start the shared DB:

```bash
cd go
docker compose up -d db
```

Terminal A: run the receiver:

```bash
cd go/webhooks
DATABASE_URL='postgres://protocol:protocol@localhost:5432/protocols?sslmode=disable' \
WEBHOOK_RECEIVER_ADDR=':8083' \
WEBHOOK_SECRET='receiver-secret' \
go run ./cmd/receiver
```

Terminal B: run the sender:

```bash
cd go/webhooks
DATABASE_URL='postgres://protocol:protocol@localhost:5432/protocols?sslmode=disable' \
WEBHOOK_SENDER_ADDR=':8084' \
MAX_ATTEMPTS='5' \
go run ./cmd/sender
```

## Register an endpoint

Request:

```http
POST /v1/endpoints
Content-Type: application/json

{"url":"http://localhost:8083/webhook","secret":"receiver-secret"}
```

Response:

- `201 Created`

```json
{"id":"...","url":"http://localhost:8083/webhook","active":true}
```

## Create an event (triggers deliveries)

Request:

```http
POST /v1/events
Content-Type: application/json

{"event_type":"order.created","payload":{"order_id":"123","amount_cents":2599}}
```

Response:

```json
{"event_id":"...","delivery_ids":["..."],"endpoint_count":1}
```

## Signature scheme

Headers the sender includes:

- `X-Webhook-Timestamp`: unix seconds
- `X-Webhook-Signature`: `v1=<hex-hmac>`
- `X-Webhook-Event`: event type
- `X-Webhook-Delivery`: unique delivery id

Signing string:

```text
<timestamp> + "." + <raw_request_body>
```

HMAC:

```text
HMAC-SHA256(secret, signing_string)
```

Receiver returns:

- `200 OK` if signature valid and payload JSON is parseable
- `401 Unauthorized` if missing/invalid signature
- `500 Internal Server Error` if failure is simulated (see below)

## Simulate failures and watch retries

Trigger receiver failure:

```text
POST http://localhost:8083/webhook?fail=true
```

Or random failure percentage:

```bash
FAIL_PERCENT='30' go run ./cmd/receiver
```

Inspect deliveries:

```http
GET /v1/deliveries?limit=25
```

Inspect received webhook calls:

```http
GET http://localhost:8083/v1/received?limit=25
```

