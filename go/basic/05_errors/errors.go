// Package errorsdemo covers Go's error model: errors are values, not control flow.
//
// The three tools that cover almost every case:
//   - fmt.Errorf with %w to add context while keeping the original reachable
//   - errors.Is to test identity against a sentinel, through any depth of wrapping
//   - errors.As to recover a concrete type so you can read its fields
//
// Never compare with == across a package boundary, and never match on err.Error() text.
package errorsdemo

import (
	"errors"
	"fmt"
)

// ---------------------------------------------------------------------------
// Sentinel errors
// ---------------------------------------------------------------------------

// Sentinels are exported package-level values callers can test for. Keep them few: each one
// is part of your public API forever.
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

// ---------------------------------------------------------------------------
// Custom error types
// ---------------------------------------------------------------------------

// ValidationError carries structured detail a sentinel cannot. Use a type when the caller
// needs the fields; use a sentinel when it only needs to know which category it is.
type ValidationError struct {
	Field string
	Value any
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid value %v for field %q", e.Value, e.Field)
}

// QueryError wraps a lower-level cause. Implementing Unwrap is what puts it on the chain
// that errors.Is and errors.As walk.
type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string { return fmt.Sprintf("query %q: %v", e.Query, e.Err) }
func (e *QueryError) Unwrap() error { return e.Err }

// ---------------------------------------------------------------------------
// Wrapping with context
// ---------------------------------------------------------------------------

// store is a stand-in for a repository layer.
type store map[string]string

func (s store) lookup(key string) (string, error) {
	v, ok := s[key]
	if !ok {
		// wrap the sentinel with the detail that makes the log line useful
		return "", fmt.Errorf("lookup key %q: %w", key, ErrNotFound)
	}
	return v, nil
}

// FetchUser adds a second layer of context. Each layer says what *it* was doing — not what
// went wrong, which the innermost error already said.
func FetchUser(s store, id string) (string, error) {
	name, err := s.lookup(id)
	if err != nil {
		return "", &QueryError{Query: "users." + id, Err: err}
	}
	return name, nil
}

// NewStore builds the demo store.
func NewStore(entries map[string]string) store { return store(entries) }

// ---------------------------------------------------------------------------
// Validation returning a typed error
// ---------------------------------------------------------------------------

// ValidateAge reports a *ValidationError so the caller can read Field and Value back out.
func ValidateAge(age int) error {
	if age < 0 || age > 150 {
		return fmt.Errorf("validate: %w", &ValidationError{Field: "age", Value: age})
	}
	return nil
}

// ---------------------------------------------------------------------------
// Joining multiple errors
// ---------------------------------------------------------------------------

// ValidateAll collects every failure instead of stopping at the first. errors.Join (Go 1.20+)
// returns an error whose chain contains all of them, so errors.Is matches any member.
func ValidateAll(name string, age int) error {
	var errs []error

	if name == "" {
		errs = append(errs, &ValidationError{Field: "name", Value: name})
	}
	if err := ValidateAge(age); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...) // returns nil when errs is empty
}

// ---------------------------------------------------------------------------
// Inspecting
// ---------------------------------------------------------------------------

// Classify shows the two inspection verbs side by side.
func Classify(err error) string {
	if err == nil {
		return "ok"
	}

	// Is: identity check against a sentinel, anywhere in the chain
	if errors.Is(err, ErrNotFound) {
		return "not-found"
	}
	if errors.Is(err, ErrUnauthorized) {
		return "unauthorized"
	}

	// As: find a concrete type in the chain and bind it so its fields are readable
	var ve *ValidationError
	if errors.As(err, &ve) {
		return "validation:" + ve.Field
	}

	return "unknown"
}
