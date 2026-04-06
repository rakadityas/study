# HLS and DASH (Live Streaming) Example

## What this demonstrates

- HLS (HTTP Live Streaming) and MPEG-DASH are *adaptive bitrate streaming* formats.
- They work by serving:
  - a playlist/manifest (`.m3u8` for HLS, `.mpd` for DASH)
  - many small media segments (`.ts` for HLS in this example; `.m4s` for DASH)
- Clients fetch segments over plain HTTP and can switch representations/bitrates (not shown here, but the packaging format supports it).

Typical use cases:

- live or VOD video streaming
- large-scale distribution via CDN
- resilient playback under fluctuating network conditions

## Generate media assets

This repo uses a synthetic video source so you don't need any input files.

```bash
cd go
docker compose --profile media run --rm ffmpeg_packager
```

Output files are written under:

```text
go/hls-dash/out/hls/master.m3u8
go/hls-dash/out/dash/manifest.mpd
```

## Run the server

```bash
cd go/hls-dash
STREAM_ADDR=':8085' go run ./cmd/server
```

## URLs

- HLS playlist: `http://localhost:8085/hls/master.m3u8`
- DASH manifest: `http://localhost:8085/dash/manifest.mpd`

Player pages (use CDN JS libraries):

- `http://localhost:8085/player/hls`
- `http://localhost:8085/player/dash`

## Typical responses

Playlist/manifest requests return `200 OK` with text content.

Segment requests return `200 OK` with binary content.

Missing files return `404 Not Found`.

