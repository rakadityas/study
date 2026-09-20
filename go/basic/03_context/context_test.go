package contextdemo

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWorkRespectsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	done, err := Work(ctx, 100, 5*time.Millisecond)
	if !IsCanceled(err) {
		t.Fatalf("Work err = %v, want context.Canceled", err)
	}
	if done == 0 || done == 100 {
		t.Fatalf("Work completed %d steps, want a partial count", done)
	}
}

func TestWorkRespectsDeadline(t *testing.T) {
	err := WithBudget(context.Background(), 30*time.Millisecond, func(ctx context.Context) error {
		_, err := Work(ctx, 100, 5*time.Millisecond)
		return err
	})

	if !IsTimeout(err) {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
	// a blown deadline is not an explicit cancel, even though both close Done()
	if IsCanceled(err) {
		t.Fatalf("err = %v should not classify as context.Canceled", err)
	}
}

func TestWorkCompletesWithinBudget(t *testing.T) {
	err := WithBudget(context.Background(), time.Second, func(ctx context.Context) error {
		done, err := Work(ctx, 3, time.Millisecond)
		if done != 3 {
			t.Fatalf("Work completed %d steps, want 3", done)
		}
		return err
	})
	if err != nil {
		t.Fatalf("WithBudget = %v, want nil", err)
	}
}

func TestRequestIDRoundTrip(t *testing.T) {
	ctx := WithRequestID(context.Background(), "abc-123")

	if id, ok := RequestID(ctx); !ok || id != "abc-123" {
		t.Fatalf("RequestID = (%q, %v), want (\"abc-123\", true)", id, ok)
	}

	if _, ok := RequestID(context.Background()); ok {
		t.Fatal("a bare context should carry no request ID")
	}
}

func TestRaceReturnsFirstSuccess(t *testing.T) {
	slow := func(ctx context.Context) (string, error) {
		select {
		case <-time.After(time.Second):
			return "slow", nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
	fast := func(ctx context.Context) (string, error) { return "fast", nil }

	got, err := Race(context.Background(), slow, fast)
	if err != nil || got != "fast" {
		t.Fatalf("Race = (%q, %v), want (\"fast\", nil)", got, err)
	}
}

func TestRaceAllFail(t *testing.T) {
	boom := errors.New("boom")
	fail := func(ctx context.Context) (string, error) { return "", boom }

	if _, err := Race(context.Background(), fail, fail); !errors.Is(err, boom) {
		t.Fatalf("Race err = %v, want boom", err)
	}
}
