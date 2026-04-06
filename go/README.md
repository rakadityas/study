# Go Protocol Examples

This folder contains runnable, minimal examples for common protocols and when you would use them.

## Quickstart

Bring up the shared DB (and optional tools):

```bash
cd go
docker compose up -d db pgadmin
```

Generate synthetic HLS + DASH assets (optional):

```bash
cd go
docker compose --profile media run --rm ffmpeg_packager
```

Connection string used by the DB-backed examples:

```text
DATABASE_URL=postgres://protocol:protocol@localhost:5432/protocols?sslmode=disable
```

## Examples

- [http](./http): REST/JSON over HTTP with Postgres persistence.
- [grpc](./grpc): gRPC request/response RPCs with Postgres persistence.
- [sse](./sse): server pushes events to clients over HTTP using Server-Sent Events.
- [websocket](./websocket): bidirectional low-latency messaging over a single TCP connection.
- [webhooks](./webhooks): async HTTP callbacks with signatures + retry worker.
- [hls-dash](./hls-dash): live-like streaming using playlists/manifests + segment files over HTTP.
- [webrtc](./webrtc): UDP-based peer-to-peer media/data via WebRTC (with WebSocket signaling).
- [p2p](./p2p): simple peer discovery + gossip over a custom TCP protocol.

