package patterns

import (
	"context"
	"errors"
	"sort"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	jobs := make([]Job, 20)
	for i := range jobs {
		jobs[i] = Job{ID: i, Input: i}
	}

	var peak, current int64
	results := WorkerPool(context.Background(), 4, jobs, func(_ context.Context, j Job) (int, error) {
		n := atomic.AddInt64(&current, 1)
		for {
			p := atomic.LoadInt64(&peak)
			if n <= p || atomic.CompareAndSwapInt64(&peak, p, n) {
				break
			}
		}
		defer atomic.AddInt64(&current, -1)

		time.Sleep(time.Millisecond)
		if j.Input == 7 {
			return 0, errors.New("job 7 failed")
		}
		return j.Input * 2, nil
	})

	if len(results) != len(jobs) {
		t.Fatalf("got %d results, want %d", len(results), len(jobs))
	}
	if peak > 4 {
		t.Fatalf("peak concurrency was %d, want at most 4 workers", peak)
	}

	sort.Slice(results, func(i, j int) bool { return results[i].ID < results[j].ID })
	for _, r := range results {
		if r.ID == 7 {
			if r.Err == nil {
				t.Fatal("job 7 should have reported an error")
			}
			continue
		}
		if r.Err != nil || r.Output != r.ID*2 {
			t.Fatalf("job %d = (%d, %v), want (%d, nil)", r.ID, r.Output, r.Err, r.ID*2)
		}
	}
}

func TestMapConcurrentPreservesOrder(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6, 7, 8}

	var peak, current int64
	out := MapConcurrent(context.Background(), 3, in, func(_ context.Context, v int) int {
		n := atomic.AddInt64(&current, 1)
		for {
			p := atomic.LoadInt64(&peak)
			if n <= p || atomic.CompareAndSwapInt64(&peak, p, n) {
				break
			}
		}
		defer atomic.AddInt64(&current, -1)

		time.Sleep(2 * time.Millisecond)
		return v * v
	})

	if peak > 3 {
		t.Fatalf("peak concurrency was %d, want at most 3", peak)
	}
	for i, v := range in {
		if out[i] != v*v {
			t.Fatalf("out = %v, want order preserved", out)
		}
	}
}

func TestErrGroupFirstErrorCancels(t *testing.T) {
	boom := errors.New("boom")
	g, ctx := NewErrGroup(context.Background())

	var cancelled atomic.Bool

	g.Go(func(context.Context) error { return boom })
	g.Go(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			cancelled.Store(true)
			return ctx.Err()
		case <-time.After(time.Second):
			return nil
		}
	})

	if err := g.Wait(); !errors.Is(err, boom) {
		t.Fatalf("Wait = %v, want boom", err)
	}
	if !cancelled.Load() {
		t.Fatal("the sibling task should have been cancelled by the first error")
	}
	if ctx.Err() == nil {
		t.Fatal("the group context should be cancelled after Wait")
	}
}

func TestErrGroupSuccess(t *testing.T) {
	g, _ := NewErrGroup(context.Background())

	var count atomic.Int64
	for i := 0; i < 5; i++ {
		g.Go(func(context.Context) error { count.Add(1); return nil })
	}

	if err := g.Wait(); err != nil {
		t.Fatalf("Wait = %v, want nil", err)
	}
	if count.Load() != 5 {
		t.Fatalf("ran %d tasks, want 5", count.Load())
	}
}

func TestTokenBucket(t *testing.T) {
	// drive the clock by hand so the test is deterministic
	now := time.Now()
	b := NewTokenBucket(3, 10) // 3 burst, 10/sec
	b.now = func() time.Time { return now }
	b.last = now // anchor the bucket to the fake clock too

	for i := 0; i < 3; i++ {
		if !b.Allow() {
			t.Fatalf("call %d should be allowed within the burst capacity", i)
		}
	}
	if b.Allow() {
		t.Fatal("the fourth call should be rate limited")
	}

	// 200ms at 10/sec refills 2 tokens
	now = now.Add(200 * time.Millisecond)
	if !b.Allow() || !b.Allow() {
		t.Fatal("two tokens should have refilled")
	}
	if b.Allow() {
		t.Fatal("only two tokens should have refilled")
	}
}

func TestGracefulShutdown(t *testing.T) {
	s := NewServer()
	s.Start(time.Millisecond)

	time.Sleep(20 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown = %v, want nil", err)
	}
	if s.Processed() == 0 {
		t.Fatal("the server should have processed at least one tick")
	}

	// the loop is stopped, so the count is now frozen
	before := s.Processed()
	time.Sleep(10 * time.Millisecond)
	if s.Processed() != before {
		t.Fatal("the loop kept running after shutdown")
	}
}
