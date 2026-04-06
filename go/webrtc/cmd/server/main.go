package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type signalMsg struct {
	Type      string          `json:"type"`
	From      string          `json:"from,omitempty"`
	To        string          `json:"to,omitempty"`
	Room      string          `json:"room,omitempty"`
	SDP       json.RawMessage `json:"sdp,omitempty"`
	Candidate json.RawMessage `json:"candidate,omitempty"`
}

type hub struct {
	mu    sync.Mutex
	rooms map[string]map[*client]struct{}
}

type client struct {
	id   string
	room string
	conn *websocket.Conn
	send chan []byte
	hub  *hub
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func main() {
	addr := envDefault("WEBRTC_ADDR", ":8086")
	h := &hub{rooms: map[string]map[*client]struct{}{}}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(indexHTML))
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		room := strings.TrimSpace(r.URL.Query().Get("room"))
		if room == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		c := &client{
			id:   "c-" + randHex(6),
			room: room,
			conn: conn,
			send: make(chan []byte, 64),
			hub:  h,
		}
		h.add(c)
		go c.writeLoop()
		go c.readLoop()
	})

	s := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("webrtc signaling server listening on %s", addr)
	log.Fatal(s.ListenAndServe())
}

func (h *hub) add(c *client) {
	h.mu.Lock()
	if h.rooms[c.room] == nil {
		h.rooms[c.room] = map[*client]struct{}{}
	}
	h.rooms[c.room][c] = struct{}{}
	h.mu.Unlock()

	ready, _ := json.Marshal(signalMsg{Type: "ready", From: c.id, Room: c.room})
	c.send <- ready

	join, _ := json.Marshal(signalMsg{Type: "peer.joined", From: c.id, Room: c.room})
	h.broadcastExcept(c.room, c, join)
}

func (h *hub) remove(c *client) {
	h.mu.Lock()
	if h.rooms[c.room] != nil {
		delete(h.rooms[c.room], c)
		if len(h.rooms[c.room]) == 0 {
			delete(h.rooms, c.room)
		}
	}
	h.mu.Unlock()
}

func (h *hub) broadcastExcept(room string, except *client, b []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.rooms[room] {
		if c == except {
			continue
		}
		select {
		case c.send <- b:
		default:
		}
	}
}

func (c *client) readLoop() {
	defer func() {
		c.hub.remove(c)
		close(c.send)
		c.conn.Close()
		leave, _ := json.Marshal(signalMsg{Type: "peer.left", From: c.id, Room: c.room})
		c.hub.broadcastExcept(c.room, c, leave)
	}()

	c.conn.SetReadLimit(512 * 1024)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, b, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		var m signalMsg
		if err := json.Unmarshal(b, &m); err != nil {
			continue
		}
		m.From = c.id
		m.Room = c.room
		out, _ := json.Marshal(m)
		c.hub.broadcastExcept(c.room, c, out)
	}
}

func (c *client) writeLoop() {
	t := time.NewTicker(25 * time.Second)
	defer t.Stop()

	for {
		select {
		case b, ok := <-c.send:
			if !ok {
				return
			}
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
				return
			}
		case <-t.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func randHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
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
    <title>WebRTC DataChannel Demo</title>
    <style>
      body { font-family: ui-sans-serif, system-ui; }
      pre { background: #111; color: #eee; padding: 12px; height: 280px; overflow: auto; }
      input { width: 520px; }
    </style>
  </head>
  <body>
    <h1>WebRTC (UDP) demo with WebSocket signaling</h1>
    <div>
      <div>Open this page in two tabs.</div>
      <div>Tab A: add <code>?room=demo&amp;role=caller</code></div>
      <div>Tab B: add <code>?room=demo&amp;role=callee</code></div>
    </div>
    <div style="margin-top: 16px;">
      <div><b>Status:</b> <span id="status">init</span></div>
      <div><b>Role:</b> <span id="role"></span> <b>Room:</b> <span id="room"></span></div>
    </div>
    <div style="margin-top: 16px;">
      <input id="msg" placeholder="message" />
      <button id="send">Send</button>
    </div>
    <pre id="log"></pre>

    <script>
      const logEl = document.getElementById('log');
      const statusEl = document.getElementById('status');
      const roleEl = document.getElementById('role');
      const roomEl = document.getElementById('room');

      const params = new URLSearchParams(location.search);
      const room = params.get('room') || 'demo';
      const role = params.get('role') || 'caller';
      roleEl.textContent = role;
      roomEl.textContent = room;

      function log(msg) {
        logEl.textContent += msg + "\n";
        logEl.scrollTop = logEl.scrollHeight;
      }
      function setStatus(s) { statusEl.textContent = s; }

      const ws = new WebSocket('ws://' + location.host + '/ws?room=' + encodeURIComponent(room));
      let peerConn = null;
      let dataChannel = null;
      let selfId = null;

      function sendSignal(m) {
        ws.send(JSON.stringify(m));
      }

      function makePC() {
        const pc = new RTCPeerConnection({
          iceServers: [{ urls: ['stun:stun.l.google.com:19302'] }]
        });
        pc.onicecandidate = (e) => {
          if (e.candidate) {
            sendSignal({ type: 'candidate', candidate: e.candidate });
          }
        };
        pc.onconnectionstatechange = () => log('pc state=' + pc.connectionState);
        pc.ondatachannel = (e) => {
          dataChannel = e.channel;
          wireDC();
        };
        return pc;
      }

      function wireDC() {
        if (!dataChannel) return;
        dataChannel.onopen = () => { log('dc open'); setStatus('connected'); };
        dataChannel.onclose = () => { log('dc close'); setStatus('closed'); };
        dataChannel.onmessage = (e) => log('dc recv: ' + e.data);
      }

      ws.onopen = () => { setStatus('ws connected'); log('ws connected'); };
      ws.onmessage = async (e) => {
        const m = JSON.parse(e.data);
        if (m.type === 'ready') {
          selfId = m.from;
          log('ready self=' + selfId);
          if (role === 'caller') {
            peerConn = makePC();
            dataChannel = peerConn.createDataChannel('chat');
            wireDC();
          }
          return;
        }
        if (!peerConn) {
          peerConn = makePC();
        }
        if (m.type === 'offer') {
          await peerConn.setRemoteDescription(m.sdp);
          const answer = await peerConn.createAnswer();
          await peerConn.setLocalDescription(answer);
          sendSignal({ type: 'answer', sdp: peerConn.localDescription });
          log('sent answer');
        } else if (m.type === 'answer') {
          await peerConn.setRemoteDescription(m.sdp);
          log('set remote answer');
        } else if (m.type === 'candidate') {
          try {
            await peerConn.addIceCandidate(m.candidate);
          } catch (err) {
            log('addIceCandidate err=' + err);
          }
        } else if (m.type === 'peer.joined') {
          if (role === 'caller') {
            const offer = await peerConn.createOffer();
            await peerConn.setLocalDescription(offer);
            sendSignal({ type: 'offer', sdp: peerConn.localDescription });
            log('sent offer');
          }
        }
      };

      document.getElementById('send').onclick = () => {
        const v = document.getElementById('msg').value;
        if (!dataChannel || dataChannel.readyState !== 'open') {
          log('dc not open');
          return;
        }
        dataChannel.send(v);
        log('dc send: ' + v);
      };
    </script>
  </body>
</html>`
