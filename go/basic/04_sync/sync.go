// Package syncdemo covers the sync and sync/atomic primitives: when a mutex beats a channel,
// what RWMutex actually buys you, and why sync.Once and atomics exist.
//
// Channels are for transferring ownership; mutexes are for protecting state that several
// goroutines genuinely share. Reach for a mutex when the shared thing is a data structure
// and the critical section is short.
package syncdemo

import (
	"sync"
	"sync/atomic"
)

// ---------------------------------------------------------------------------
// Mutex: guard the data, not the code
// ---------------------------------------------------------------------------

// SafeCounter embeds the mutex directly above the fields it protects — the conventional way
// to document the association. Keep the zero value usable: a zero Mutex is unlocked, so
// SafeCounter needs no constructor.
type SafeCounter struct {
	mu sync.Mutex
	n  int
}

// Inc is safe from any number of goroutines. defer Unlock costs a few nanoseconds and buys
// correctness on every early return and panic path.
func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.n++
}

// Value reads under the same lock. Reading an int without the lock is a data race even
// though the read looks atomic — the race detector will flag it, and the compiler is free
// to reorder around it.
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

// ---------------------------------------------------------------------------
// RWMutex: many readers or one writer
// ---------------------------------------------------------------------------

// Cache is the classic read-heavy structure. RLock lets readers proceed concurrently, which
// only pays off when reads dominate and hold the lock long enough to matter — otherwise
// plain Mutex is faster because RWMutex has more bookkeeping.
type Cache struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewCache() *Cache { return &Cache{m: map[string]string{}} }

func (c *Cache) Get(k string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[k]
	return v, ok
}

func (c *Cache) Set(k, v string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[k] = v
}

// ---------------------------------------------------------------------------
// WaitGroup: wait for a known number of goroutines
// ---------------------------------------------------------------------------

// ParallelSum splits work across goroutines and waits for all of them. Add is called before
// starting the goroutine — calling it inside would race with Wait.
func ParallelSum(nums []int, workers int) int {
	if workers < 1 {
		workers = 1
	}

	partials := make([]int, workers)
	chunk := (len(nums) + workers - 1) / workers

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		start := w * chunk
		end := start + chunk
		if start > len(nums) {
			start = len(nums)
		}
		if end > len(nums) {
			end = len(nums)
		}

		wg.Add(1)
		go func(idx int, slice []int) {
			defer wg.Done()
			// each goroutine writes to its own slot, so no lock is needed —
			// disjoint writes to distinct elements of a slice are safe
			sum := 0
			for _, v := range slice {
				sum += v
			}
			partials[idx] = sum
		}(w, nums[start:end])
	}
	wg.Wait()

	total := 0
	for _, p := range partials {
		total += p
	}
	return total
}

// ---------------------------------------------------------------------------
// Once: exactly-once initialisation
// ---------------------------------------------------------------------------

// LazyConfig defers expensive setup until first use and guarantees it happens once even
// under concurrent access. Once.Do also blocks late callers until the first one finishes,
// so they never observe a half-built value.
type LazyConfig struct {
	once  sync.Once
	value string
	calls int32 // counted with an atomic so the test can assert without its own lock
}

func (l *LazyConfig) Get(build func() string) string {
	l.once.Do(func() {
		atomic.AddInt32(&l.calls, 1)
		l.value = build()
	})
	return l.value
}

func (l *LazyConfig) BuildCount() int32 { return atomic.LoadInt32(&l.calls) }

// ---------------------------------------------------------------------------
// Atomics: lock-free counters
// ---------------------------------------------------------------------------

// AtomicCounter uses the typed atomic wrappers (Go 1.19+), which are harder to misuse than
// the free functions because the value cannot be touched non-atomically by accident.
// Atomics suit single words; anything involving two fields still needs a mutex.
type AtomicCounter struct{ n atomic.Int64 }

func (a *AtomicCounter) Inc()         { a.n.Add(1) }
func (a *AtomicCounter) Value() int64 { return a.n.Load() }

// CompareAndSwap is the building block for lock-free algorithms: set to new only if the
// current value is still old, reporting whether it won the race.
func (a *AtomicCounter) CompareAndSwap(old, new int64) bool { return a.n.CompareAndSwap(old, new) }

// ---------------------------------------------------------------------------
// sync.Map: only for two specific shapes
// ---------------------------------------------------------------------------

// CountDistinct uses sync.Map, which is worth it only when keys are written once and read
// many times, or when goroutines touch disjoint key sets. For anything else, a plain map
// behind an RWMutex is faster and far easier to reason about.
func CountDistinct(keys []string) int {
	var m sync.Map
	var wg sync.WaitGroup

	for _, k := range keys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			m.LoadOrStore(key, true)
		}(k)
	}
	wg.Wait()

	n := 0
	m.Range(func(_, _ any) bool {
		n++
		return true // returning false stops the iteration
	})
	return n
}
