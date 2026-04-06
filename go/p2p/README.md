# P2P (Peer-to-Peer) Example

## What this demonstrates

- A simple custom peer-to-peer protocol over TCP.
- Peer discovery by exchanging known peer addresses.
- Gossip/broadcast with loop prevention (seen-message IDs).

Typical P2P use cases:

- decentralized messaging
- content distribution
- node discovery and gossip in distributed systems

This example is intentionally small and does not implement NAT traversal, encryption, or robust membership protocols.

## Run

Terminal A:

```bash
cd go/p2p
go run ./cmd/node -listen :7001 -name a
```

Terminal B:

```bash
cd go/p2p
go run ./cmd/node -listen :7002 -name b -peers localhost:7001
```

Terminal C:

```bash
cd go/p2p
go run ./cmd/node -listen :7003 -name c -peers localhost:7002
```

Send a single gossip message:

```bash
cd go/p2p
go run ./cmd/node -listen :7010 -name sender -peers localhost:7001 -say 'hello p2p'
```

## Message shapes

All messages are JSON objects written as newline-delimited frames.

Handshake:

```json
{"type":"hello","node_id":"a","listen_addr":":7001"}
```

Peer list:

```json
{"type":"peers","peers":[":7001",":7002"]}
```

Gossip:

```json
{"type":"gossip","id":"<random>","from":"a","body":"hello","ts":"2026-01-01T00:00:00Z"}
```

