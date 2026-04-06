# WebRTC (UDP) Example

## What this demonstrates

- WebRTC creates peer-to-peer connections for real-time media and data.
- The actual transport is typically UDP (ICE + DTLS + SRTP/SRTCP for media, DTLS + SCTP for DataChannel).
- WebRTC does *not* define signaling; you must provide a signaling channel to exchange SDP offers/answers and ICE candidates.

This example provides:

- a Go signaling server using WebSocket (`/ws`)
- a browser page that establishes a WebRTC DataChannel between two tabs

## Run

```bash
cd go/webrtc
WEBRTC_ADDR=':8086' go run ./cmd/server
```

Open two tabs:

- `http://localhost:8086/?room=demo&role=caller`
- `http://localhost:8086/?room=demo&role=callee`

Send messages using the input box; they travel over the WebRTC DataChannel.

## Signaling protocol (WebSocket)

Connect:

```text
ws://localhost:8086/ws?room=demo
```

Messages are JSON:

- Offer

```json
{"type":"offer","sdp":{...}}
```

- Answer

```json
{"type":"answer","sdp":{...}}
```

- ICE candidate

```json
{"type":"candidate","candidate":{...}}
```

The server forwards signaling messages to other peers in the same room and injects:

- `ready` when connected
- `peer.joined` / `peer.left` when peers connect/disconnect

## Typical failure modes

- If peers are behind NATs/firewalls, direct connectivity may fail without STUN/TURN.
- If UDP is blocked, WebRTC may attempt TCP/TLS fallbacks depending on browser and network.

