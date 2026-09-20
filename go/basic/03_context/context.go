// Package contextdemo covers context.Context: the standard way to carry cancellation,
// deadlines and request-scoped values across API boundaries.
//
// Rules that matter in review:
//   - context.Context is the first parameter, always named ctx, never stored in a struct.
//   - Every cancel func must be called, usually with defer, or the parent leaks the child.
//   - Cancellation is advisory: it closes ctx.Done(); your code still has to notice.
//   - Values are for request-scoped metadata (trace IDs, auth), never optional parameters.
package contextdemo

import (
	"context"
	"errors"
	"time"
)

// ---------------------------------------------------------------------------
// Cancellation
// ---------------------------------------------------------------------------

// Work simulates a job that polls ctx.Done() between units of work. Notice the shape: a
// select with a Done case is how a goroutine stays interruptible.
func Work(ctx context.Context, steps int, each time.Duration) (int, error) {
	done := 0
	for i := 0; i < steps; i++ {
		select {
		case <-ctx.Done():
			// ctx.Err() says why: Canceled or DeadlineExceeded
			return done, ctx.Err()
		case <-time.After(each):
			done++
		}
	}
	return done, nil
}

// ---------------------------------------------------------------------------
// Deadlines and timeouts
// ---------------------------------------------------------------------------

// WithBudget runs fn under a total time budget. A deadline propagates down the whole call
// tree, so every context-aware call below inherits the remaining time automatically.
func WithBudget(parent context.Context, budget time.Duration, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(parent, budget)
	defer cancel() // releases the timer immediately when fn returns early

	return fn(ctx)
}

// IsTimeout distinguishes a blown deadline from an explicit cancel. Compare with errors.Is,
// not ==, because the error may have been wrapped on the way up.
func IsTimeout(err error) bool { return errors.Is(err, context.DeadlineExceeded) }

// IsCanceled reports an explicit cancel.
func IsCanceled(err error) bool { return errors.Is(err, context.Canceled) }

// ---------------------------------------------------------------------------
// Values
// ---------------------------------------------------------------------------

// ctxKey is an unexported named type. Using a string literal as a key risks collisions with
// other packages writing to the same context; an unexported type cannot collide by design.
type ctxKey string

const requestIDKey ctxKey = "request-id"

// WithRequestID attaches a request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestID reads it back. Expose typed accessors rather than the key itself, so callers
// cannot store the wrong type under it.
func RequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}

// ---------------------------------------------------------------------------
// First-result-wins
// ---------------------------------------------------------------------------

// Race starts every fn with a derived context and returns the first successful result,
// cancelling the rest. This is the canonical use of cancellation: stop work nobody needs.
func Race(ctx context.Context, fns ...func(context.Context) (string, error)) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancels the losers as soon as we return

	type result struct {
		val string
		err error
	}
	// buffered to len(fns) so a losing goroutine can always send and exit rather than
	// blocking forever on a channel nobody reads — this is how you avoid a goroutine leak
	results := make(chan result, len(fns))

	for _, fn := range fns {
		go func(f func(context.Context) (string, error)) {
			v, err := f(ctx)
			results <- result{v, err}
		}(fn)
	}

	var lastErr error
	for range fns {
		select {
		case r := <-results:
			if r.err == nil {
				return r.val, nil
			}
			lastErr = r.err
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}

	return "", lastErr
}
