// Package concurrency covers goroutines and channels: the mechanics, the directional types,
// the select statement, and the closing/draining rules that decide whether a program
// deadlocks or leaks.
//
// The guiding maxim: "Don't communicate by sharing memory; share memory by communicating."
// Channels move ownership of a value from one goroutine to another, which is what makes the
// handoff race-free without a mutex.
package concurrency

import (
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Buffered vs unbuffered
// ---------------------------------------------------------------------------

// Unbuffered channels are a *rendezvous*: the send blocks until a receiver is ready, so the
// two goroutines synchronise. A buffered channel decouples them up to its capacity, which
// makes the send non-blocking while there is room — useful for smoothing bursts, not for
// "making things faster".

// Generate returns a channel that emits n values and is then closed. Returning a
// receive-only channel (<-chan) makes it a compile error for the caller to send or close.
func Generate(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out) // the producer closes — never the consumer
		for i := 0; i < n; i++ {
			out <- i
		}
	}()
	return out
}

// ---------------------------------------------------------------------------
// Pipelines: each stage reads from one channel and writes to the next
// ---------------------------------------------------------------------------

// Square is a pipeline stage. It owns its output channel and closes it when its input dries
// up, which propagates the shutdown signal downstream without any extra coordination.
func Square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in { // ranging over a channel ends when it is closed and drained
			out <- v * v
		}
	}()
	return out
}

// ---------------------------------------------------------------------------
// Fan-out / fan-in
// ---------------------------------------------------------------------------

// FanIn merges several channels into one. The WaitGroup counts the forwarding goroutines so
// the output can be closed exactly once, after the last input is exhausted.
func FanIn(inputs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(in) // pre-Go 1.22 the loop variable had to be passed explicitly; still clearest

	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ---------------------------------------------------------------------------
// select
// ---------------------------------------------------------------------------

// FirstResponse returns whichever channel produces a value first, or "" after the timeout.
// select blocks until one case is ready; if several are ready it picks uniformly at random,
// which prevents one busy channel from starving another.
func FirstResponse(a, b <-chan string, timeout time.Duration) string {
	timer := time.NewTimer(timeout)
	defer timer.Stop() // always stop a timer you abandon, or it holds memory until it fires

	select {
	case v := <-a:
		return v
	case v := <-b:
		return v
	case <-timer.C:
		return ""
	}
}

// TryReceive shows the non-blocking form. A default case makes select give up immediately
// instead of waiting, which is how you poll a channel without blocking.
func TryReceive(c <-chan int) (int, bool) {
	select {
	case v, ok := <-c:
		return v, ok
	default:
		return 0, false
	}
}

// ---------------------------------------------------------------------------
// Closed-channel semantics
// ---------------------------------------------------------------------------

// Receiving from a closed channel yields the zero value immediately with ok == false, which
// is how a closed channel broadcasts to every waiting goroutine at once. Sending on a closed
// channel panics, and so does closing one twice — hence "the producer closes".

// Broadcast uses a closed channel as a one-shot signal to every listener.
func Broadcast(listeners int) []int {
	done := make(chan struct{}) // struct{} carries no data: this is a pure signal
	results := make([]int, listeners)
	var wg sync.WaitGroup

	for i := 0; i < listeners; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			<-done // every listener unblocks on the single close below
			results[idx] = idx
		}(i)
	}

	close(done)
	wg.Wait()
	return results
}
