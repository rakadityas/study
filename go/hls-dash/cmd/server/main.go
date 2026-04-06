package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type pageData struct {
	HLSURL  string
	DASHURL string
}

func main() {
	addr := envDefault("STREAM_ADDR", ":8085")
	hlsDir := envDefault("HLS_DIR", "./out/hls")
	dashDir := envDefault("DASH_DIR", "./out/dash")

	indexT := template.Must(template.New("index").Parse(indexHTML))

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = indexT.Execute(w, pageData{
			HLSURL:  "/player/hls",
			DASHURL: "/player/dash",
		})
	})

	mux.Handle("/hls/", http.StripPrefix("/hls/", http.FileServer(http.Dir(hlsDir))))
	mux.Handle("/dash/", http.StripPrefix("/dash/", http.FileServer(http.Dir(dashDir))))

	mux.HandleFunc("/player/hls", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(playerHLSHTML))
	})
	mux.HandleFunc("/player/dash", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(playerDASHHTML))
	})

	s := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("hls/dash server listening on %s", addr)
	log.Printf("hls dir:  %s", filepath.Clean(hlsDir))
	log.Printf("dash dir: %s", filepath.Clean(dashDir))
	log.Fatal(s.ListenAndServe())
}

func envDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

const indexHTML = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>HLS + DASH</title>
  </head>
  <body>
    <h1>HLS + DASH example</h1>
    <ul>
      <li><a href="/hls/master.m3u8">HLS master playlist</a></li>
      <li><a href="/dash/manifest.mpd">DASH manifest</a></li>
      <li><a href="{{.HLSURL}}">HLS player page</a> (requires hls.js via CDN)</li>
      <li><a href="{{.DASHURL}}">DASH player page</a> (requires dash.js via CDN)</li>
    </ul>
  </body>
</html>`

const playerHLSHTML = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>HLS Player</title>
    <script src="https://cdn.jsdelivr.net/npm/hls.js@latest"></script>
  </head>
  <body>
    <h1>HLS Player</h1>
    <video id="video" controls autoplay style="width: 90vw; max-width: 960px;"></video>
    <script>
      const video = document.getElementById('video');
      const src = '/hls/master.m3u8';
      if (video.canPlayType('application/vnd.apple.mpegurl')) {
        video.src = src;
      } else if (window.Hls) {
        const hls = new Hls();
        hls.loadSource(src);
        hls.attachMedia(video);
      } else {
        document.body.append('HLS not supported in this browser.');
      }
    </script>
  </body>
</html>`

const playerDASHHTML = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>DASH Player</title>
    <script src="https://cdn.dashjs.org/latest/dash.all.min.js"></script>
  </head>
  <body>
    <h1>DASH Player</h1>
    <video id="video" controls autoplay style="width: 90vw; max-width: 960px;"></video>
    <script>
      const url = '/dash/manifest.mpd';
      const player = dashjs.MediaPlayer().create();
      player.initialize(document.querySelector('#video'), url, true);
    </script>
  </body>
</html>`

