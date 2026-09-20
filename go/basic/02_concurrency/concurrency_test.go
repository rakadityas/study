package concurrency

import (
	"sort"
	"testing"
	"time"
)

func TestPipeline(t *testing.T) {
	got := []int{}
	for v := range Square(Generate(5)) {
		got = append(got, v)
	}

	want := []int{0, 1, 4, 9, 16}
	if len(got) != len(want) {
		t.Fatalf("pipeline produced %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("pipeline produced %v, want %v", got, want)
		}
	}
}

func TestFanIn(t *testing.T) {
	merged := FanIn(Generate(3), Generate(3))

	got := []int{}
	for v := range merged {
		got = append(got, v)
	}
	sort.Ints(got) // merge order is non-deterministic by design

	want := []int{0, 0, 1, 1, 2, 2}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("FanIn produced %v, want %v", got, want)
		}
	}
}

func TestFirstResponseAndTimeout(t *testing.T) {
	fast := make(chan string, 1)
	slow := make(chan string, 1)
	fast <- "fast"

	if got := FirstResponse(fast, slow, time.Second); got != "fast" {
		t.Fatalf("FirstResponse = %q, want %q", got, "fast")
	}

	quiet := make(chan string)
	if got := FirstResponse(quiet, quiet, 10*time.Millisecond); got != "" {
		t.Fatalf("FirstResponse on silence = %q, want empty (timeout)", got)
	}
}

func TestTryReceive(t *testing.T) {
	empty := make(chan int)
	if _, ok := TryReceive(empty); ok {
		t.Fatal("TryReceive on an empty channel should not succeed")
	}

	filled := make(chan int, 1)
	filled <- 7
	if v, ok := TryReceive(filled); !ok || v != 7 {
		t.Fatalf("TryReceive = (%d, %v), want (7, true)", v, ok)
	}

	// a closed channel is always ready, and reports ok == false
	closed := make(chan int)
	close(closed)
	if v, ok := TryReceive(closed); ok || v != 0 {
		t.Fatalf("TryReceive on closed = (%d, %v), want (0, false)", v, ok)
	}
}

func TestBroadcast(t *testing.T) {
	got := Broadcast(4)
	for i, v := range got {
		if v != i {
			t.Fatalf("Broadcast = %v, want index-matched values", got)
		}
	}
}
