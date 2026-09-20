// Package patterns collects the concurrency patterns that show up in real services:
// worker pools, bounded parallelism, graceful shutdown, and the leak-avoidance rules that
// keep them honest.
//
// Two rules underpin all of it:
//  1. Never start a goroutine without knowing how it will stop.
//  2. Every send needs a guaranteed receiver, or a select with a Done case to bail out.
package patterns

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// Worker pool
// ---------------------------------------------------------------------------

// Job is a unit of work with a stable ID so results can be matched back to inputs.
type Job struct {
	ID    int
	Input int
}

// JobResult pairs the job ID with its outcome.
type JobResult struct {
	ID     int
	Output int
	Err    error
}

// WorkerPool runs fn over jobs with a fixed number of workers. Fixed workers, not one
// goroutine per job: that is what bounds memory and downstream pressure when the input is
// large or unbounded.
//
// The shape to memorise: a jobs channel, N workers ranging over it, a WaitGroup, and a
// closer goroutine that closes results once every worker has finished.
func WorkerPool(ctx context.Context, workers int, jobs []Job, fn func(context.Context, Job) (int, error)) []JobResult {
	if workers < 1 {
		workers = 1
	}

	jobCh := make(chan Job)
	resultCh := make(chan JobResult, len(jobs))

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				out, err := fn(ctx, job)
				resultCh <- JobResult{ID: job.ID, Output: out, Err: err}
			}
		}()
	}

	// feed the workers, bailing out early if the caller cancels
	go func() {
		defer close(jobCh) // closing unblocks every worker's range loop
		for _, j := range jobs {
			select {
			case jobCh <- j:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	results := make([]JobResult, 0, len(jobs))
	for r := range resultCh {
		results = append(results, r)
	}
	return results
}

// ---------------------------------------------------------------------------
// Bounded parallelism with a semaphore
// ---------------------------------------------------------------------------

// Semaphore is a counting semaphore built from a buffered channel — the idiomatic Go way to
// cap concurrency without a worker pool. Acquire blocks while the buffer is full.
type Semaphore chan struct{}

func NewSemaphore(n int) Semaphore { return make(Semaphore, n) }

// Acquire takes a slot, honouring cancellation so a caller is never stuck behind a full
// semaphore after its context has expired.
func (s Semaphore) Acquire(ctx context.Context) error {
	select {
	case s <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s Semaphore) Release() { <-s }

// MapConcurrent applies fn to every input with at most `limit` running at once, preserving
// input order in the output because each goroutine writes to its own index.
func MapConcurrent[T, U any](ctx context.Context, limit int, in []T, fn func(context.Context, T) U) []U {
	out := make([]U, len(in))
	sem := NewSemaphore(limit)

	var wg sync.WaitGroup
	for i, v := range in {
		if err := sem.Acquire(ctx); err != nil {
			break
		}

		wg.Add(1)
		go func(idx int, val T) {
			defer wg.Done()
			defer sem.Release()
			out[idx] = fn(ctx, val)
		}(i, v)
	}
	wg.Wait()

	return out
}

// ---------------------------------------------------------------------------
// First error wins
// ---------------------------------------------------------------------------

// ErrGroup is a minimal stand-in for golang.org/x/sync/errgroup: run tasks concurrently,
// cancel the rest on the first failure, and return that error.
type ErrGroup struct {
	wg     sync.WaitGroup
	once   sync.Once
	err    error
	cancel context.CancelFunc
	ctx    context.Context
}

func NewErrGroup(parent context.Context) (*ErrGroup, context.Context) {
	ctx, cancel := context.WithCancel(parent)
	return &ErrGroup{cancel: cancel, ctx: ctx}, ctx
}

// Go runs fn. The first non-nil error is recorded (sync.Once makes that race-free) and
// cancels the group's context so the remaining tasks can wind down.
func (g *ErrGroup) Go(fn func(context.Context) error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := fn(g.ctx); err != nil {
			g.once.Do(func() {
				g.err = err
				g.cancel()
			})
		}
	}()
}

// Wait blocks for every task and returns the first error, if any.
func (g *ErrGroup) Wait() error {
	g.wg.Wait()
	g.cancel() // release the context even on the success path
	return g.err
}

// ---------------------------------------------------------------------------
// Rate limiting
// ---------------------------------------------------------------------------

// ErrRateLimited is returned when a token is unavailable.
var ErrRateLimited = errors.New("rate limited")

// TokenBucket allows bursts up to its capacity while holding the long-run rate steady.
// Tokens are computed lazily from elapsed time rather than by a background goroutine, so an
// idle bucket costs nothing.
type TokenBucket struct {
	mu       sync.Mutex
	tokens   float64
	capacity float64
	refill   float64 // tokens per second
	last     time.Time
	now      func() time.Time // injectable for deterministic tests
}

func NewTokenBucket(capacity int, perSecond float64) *TokenBucket {
	return &TokenBucket{
		tokens:   float64(capacity),
		capacity: float64(capacity),
		refill:   perSecond,
		last:     time.Now(),
		now:      time.Now,
	}
}

// Allow takes one token, reporting whether the call may proceed.
func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now

	b.tokens += elapsed * b.refill
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// ---------------------------------------------------------------------------
// Graceful shutdown
// ---------------------------------------------------------------------------

// Server models the shutdown handshake every long-running component needs: a quit signal, a
// WaitGroup for in-flight work, and a Shutdown that waits with a deadline instead of
// hanging forever.
type Server struct {
	quit      chan struct{}
	wg        sync.WaitGroup
	processed int
	mu        sync.Mutex
}

func NewServer() *Server { return &Server{quit: make(chan struct{})} }

// Start runs a loop that ticks until told to quit.
func (s *Server) Start(interval time.Duration) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-s.quit:
				return
			case <-ticker.C:
				s.mu.Lock()
				s.processed++
				s.mu.Unlock()
			}
		}
	}()
}

// Shutdown signals the loop and waits, up to the context's deadline.
func (s *Server) Shutdown(ctx context.Context) error {
	close(s.quit)

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err() // in-flight work outlived the grace period
	}
}

func (s *Server) Processed() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.processed
}
