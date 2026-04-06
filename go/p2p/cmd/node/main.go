package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type msg struct {
	Type       string   `json:"type"`
	NodeID     string   `json:"node_id,omitempty"`
	ListenAddr string   `json:"listen_addr,omitempty"`
	Peers      []string `json:"peers,omitempty"`
	ID         string   `json:"id,omitempty"`
	From       string   `json:"from,omitempty"`
	Body       string   `json:"body,omitempty"`
	TS         string   `json:"ts,omitempty"`
}

type node struct {
	id         string
	listenAddr string

	mu       sync.Mutex
	peers    map[string]net.Conn
	knownAdd map[string]struct{}
	seen     map[string]struct{}
}

func main() {
	var listen string
	var peersCSV string
	var name string
	var say string
	flag.StringVar(&listen, "listen", ":7001", "listen address")
	flag.StringVar(&peersCSV, "peers", "", "comma-separated peer addresses to dial")
	flag.StringVar(&name, "name", "", "node id (defaults to random)")
	flag.StringVar(&say, "say", "", "broadcast one gossip message on startup")
	flag.Parse()

	if name == "" {
		name = "node-" + randHex(4)
	}

	n := &node{
		id:         name,
		listenAddr: listen,
		peers:      map[string]net.Conn{},
		knownAdd:   map[string]struct{}{listen: {}},
		seen:       map[string]struct{}{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go n.listenLoop(ctx)

	for _, p := range splitCSV(peersCSV) {
		go n.dialLoop(ctx, p)
	}

	time.Sleep(500 * time.Millisecond)
	if say != "" {
		n.broadcast(msg{
			Type: "gossip",
			ID:   randHex(12),
			From: n.id,
			Body: say,
			TS:   time.Now().UTC().Format(time.RFC3339Nano),
		})
	}

	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			n.printState()
		}
	}
}

func (n *node) listenLoop(ctx context.Context) {
	ln, err := net.Listen("tcp", n.listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s listening on %s", n.id, n.listenAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go n.handleConn(ctx, conn, conn.RemoteAddr().String())
	}
}

func (n *node) dialLoop(ctx context.Context, addr string) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return
	}
	backoff := 500 * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			time.Sleep(backoff)
			if backoff < 5*time.Second {
				backoff *= 2
			}
			continue
		}
		n.handleConn(ctx, conn, addr)
		time.Sleep(backoff)
	}
}

func (n *node) handleConn(ctx context.Context, conn net.Conn, peerKey string) {
	defer conn.Close()

	n.mu.Lock()
	n.peers[peerKey] = conn
	n.mu.Unlock()

	enc := json.NewEncoder(conn)
	_ = enc.Encode(msg{Type: "hello", NodeID: n.id, ListenAddr: n.listenAddr})
	_ = enc.Encode(msg{Type: "peers", Peers: n.snapshotKnown()})

	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		var m msg
		if err := json.Unmarshal(sc.Bytes(), &m); err != nil {
			continue
		}
		n.onMsg(ctx, peerKey, m)
	}

	n.mu.Lock()
	delete(n.peers, peerKey)
	n.mu.Unlock()
}

func (n *node) onMsg(ctx context.Context, fromPeer string, m msg) {
	switch m.Type {
	case "hello":
		if m.ListenAddr != "" {
			n.addKnown(m.ListenAddr)
		}
	case "peers":
		for _, a := range m.Peers {
			n.addKnown(a)
		}
	case "gossip":
		if m.ID == "" {
			return
		}
		if n.markSeen(m.ID) {
			return
		}
		fmt.Printf("gossip from=%s id=%s body=%q ts=%s\n", m.From, m.ID, m.Body, m.TS)
		n.forward(fromPeer, m)
	default:
	}

	select {
	case <-ctx.Done():
	default:
	}
}

func (n *node) broadcast(m msg) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, c := range n.peers {
		_ = json.NewEncoder(c).Encode(m)
	}
}

func (n *node) forward(exceptPeer string, m msg) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for k, c := range n.peers {
		if k == exceptPeer {
			continue
		}
		_ = json.NewEncoder(c).Encode(m)
	}
}

func (n *node) addKnown(addr string) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return
	}
	n.mu.Lock()
	if _, ok := n.knownAdd[addr]; ok {
		n.mu.Unlock()
		return
	}
	n.knownAdd[addr] = struct{}{}
	n.mu.Unlock()
	if addr != n.listenAddr {
		go n.dialLoop(context.Background(), addr)
	}
}

func (n *node) snapshotKnown() []string {
	n.mu.Lock()
	defer n.mu.Unlock()
	out := make([]string, 0, len(n.knownAdd))
	for a := range n.knownAdd {
		out = append(out, a)
	}
	return out
}

func (n *node) markSeen(id string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	if _, ok := n.seen[id]; ok {
		return true
	}
	n.seen[id] = struct{}{}
	return false
}

func (n *node) printState() {
	n.mu.Lock()
	peerCount := len(n.peers)
	knownCount := len(n.knownAdd)
	n.mu.Unlock()
	log.Printf("%s peers=%d known_addrs=%d", n.id, peerCount, knownCount)
}

func randHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
