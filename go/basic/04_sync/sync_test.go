package syncdemo

import (
	"sync"
	"testing"
)

func TestSafeCounter(t *testing.T) {
	var c SafeCounter // zero value is ready to use
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()

	if got := c.Value(); got != 100 {
		t.Fatalf("SafeCounter = %d, want 100", got)
	}
}

func TestCacheConcurrentReadWrite(t *testing.T) {
	c := NewCache()
	c.Set("a", "1")

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); c.Get("a") }()
		go func() { defer wg.Done(); c.Set("b", "2") }()
	}
	wg.Wait()

	if v, ok := c.Get("a"); !ok || v != "1" {
		t.Fatalf("Get(\"a\") = (%q, %v), want (\"1\", true)", v, ok)
	}
	if v, ok := c.Get("b"); !ok || v != "2" {
		t.Fatalf("Get(\"b\") = (%q, %v), want (\"2\", true)", v, ok)
	}
}

func TestParallelSum(t *testing.T) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i + 1
	}
	want := 1000 * 1001 / 2

	for _, workers := range []int{1, 3, 8, 64} {
		if got := ParallelSum(nums, workers); got != want {
			t.Fatalf("ParallelSum(workers=%d) = %d, want %d", workers, got, want)
		}
	}
}

func TestOnceRunsExactlyOnce(t *testing.T) {
	var cfg LazyConfig
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if got := cfg.Get(func() string { return "built" }); got != "built" {
				t.Errorf("Get = %q, want %q", got, "built")
			}
		}()
	}
	wg.Wait()

	if got := cfg.BuildCount(); got != 1 {
		t.Fatalf("build ran %d times, want exactly 1", got)
	}
}

func TestAtomicCounter(t *testing.T) {
	var a AtomicCounter
	var wg sync.WaitGroup

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); a.Inc() }()
	}
	wg.Wait()

	if got := a.Value(); got != 200 {
		t.Fatalf("AtomicCounter = %d, want 200", got)
	}

	if !a.CompareAndSwap(200, 0) {
		t.Fatal("CompareAndSwap(200, 0) should have succeeded")
	}
	if a.CompareAndSwap(200, 1) {
		t.Fatal("CompareAndSwap should fail once the value has moved on")
	}
}

func TestCountDistinct(t *testing.T) {
	if got := CountDistinct([]string{"a", "b", "a", "c", "b", "a"}); got != 3 {
		t.Fatalf("CountDistinct = %d, want 3", got)
	}
}
