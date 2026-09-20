// Package generics covers type parameters (Go 1.18+): constraints, inference, generic
// containers, and — just as important — when not to reach for them.
//
// The test for a good generic: would you otherwise write the same function two or three
// times for different types? If instead you would write *different* logic, an interface or
// plain duplication is clearer.
package generics

import (
	"cmp"
	"sort"
)

// ---------------------------------------------------------------------------
// Constraints
// ---------------------------------------------------------------------------

// Number is a union constraint. The ~ prefix means "any type whose underlying type is this",
// so a `type Celsius float64` still satisfies it.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Sum works for any numeric type. The type argument is inferred from the call site, so
// callers write Sum(xs) rather than Sum[int](xs).
func Sum[T Number](values []T) T {
	var total T // the zero value of the type parameter
	for _, v := range values {
		total += v
	}
	return total
}

// Max uses the standard library's cmp.Ordered rather than hand-rolling a constraint.
// Prefer stdlib constraints (cmp.Ordered, comparable) when they fit.
func Max[T cmp.Ordered](values []T) (T, bool) {
	var zero T
	if len(values) == 0 {
		return zero, false
	}

	best := values[0]
	for _, v := range values[1:] {
		if v > best {
			best = v
		}
	}
	return best, true
}

// ---------------------------------------------------------------------------
// Generic slice helpers
// ---------------------------------------------------------------------------

// Map converts a slice element-wise. Two type parameters, both inferred — T from the input
// slice, U from the function's return type.
func Map[T, U any](in []T, fn func(T) U) []U {
	out := make([]U, 0, len(in)) // preallocate: the length is known
	for _, v := range in {
		out = append(out, fn(v))
	}
	return out
}

// Filter keeps elements matching the predicate.
func Filter[T any](in []T, keep func(T) bool) []T {
	out := []T{}
	for _, v := range in {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

// Reduce folds a slice into a single value.
func Reduce[T, A any](in []T, initial A, fn func(A, T) A) A {
	acc := initial
	for _, v := range in {
		acc = fn(acc, v)
	}
	return acc
}

// GroupBy buckets elements by a derived key. K is constrained to comparable because it has
// to be usable as a map key — that is the only thing `comparable` guarantees.
func GroupBy[T any, K comparable](in []T, key func(T) K) map[K][]T {
	out := map[K][]T{}
	for _, v := range in {
		k := key(v)
		out[k] = append(out[k], v)
	}
	return out
}

// Keys returns a map's keys in sorted order, which makes iteration deterministic —
// Go randomises map iteration on purpose, so never rely on its order.
func Keys[K cmp.Ordered, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// ---------------------------------------------------------------------------
// Generic containers
// ---------------------------------------------------------------------------

// Stack is a type-safe LIFO. Before generics this was either []any with assertions at every
// pop, or one hand-written stack per element type.
type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}

	v := s.items[len(s.items)-1]
	// clear the vacated slot so the popped value can be garbage collected —
	// shrinking the slice alone leaves the backing array holding a live reference
	s.items[len(s.items)-1] = zero
	s.items = s.items[:len(s.items)-1]
	return v, true
}

func (s *Stack[T]) Len() int { return len(s.items) }

// Result pairs a value with an error, a shape that only becomes expressible with generics.
// Go idiom still prefers multiple return values; this is for places where a result has to
// travel through a channel or a slice.
type Result[T any] struct {
	Value T
	Err   error
}

func Ok[T any](v T) Result[T]          { return Result[T]{Value: v} }
func Err[T any](err error) Result[T]   { return Result[T]{Err: err} }
func (r Result[T]) IsOk() bool         { return r.Err == nil }
func (r Result[T]) Unwrap() (T, error) { return r.Value, r.Err }
