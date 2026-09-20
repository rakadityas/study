// Package memory covers what actually happens under slices, maps and strings, plus the
// allocation behaviour that shows up in benchmarks and profiles.
//
// The recurring theme: Go values are copied, but several of them can share one backing
// array. Knowing which is which is the difference between a fast program and a subtle bug.
package memory

import "strings"

// ---------------------------------------------------------------------------
// Slice internals
// ---------------------------------------------------------------------------

// A slice header is three words: pointer to a backing array, length, capacity. Copying a
// slice copies the header, not the data — so two slices can alias the same array.

// Aliasing demonstrates the trap. The sub-slice shares storage, so writing through it is
// visible in the original.
func Aliasing() (original, sub []int) {
	original = []int{1, 2, 3, 4, 5}
	sub = original[1:3] // len 2, cap 4 — capacity runs to the end of the backing array
	sub[0] = 99         // writes original[1]
	return original, sub
}

// AppendCanClobber shows the sharper edge of the same fact: appending to a sub-slice that
// still has spare capacity overwrites the next element of the parent.
func AppendCanClobber() []int {
	original := []int{1, 2, 3, 4, 5}
	sub := original[1:3] // cap is 4, so there is room to append in place
	_ = append(sub, 42)  // overwrites original[3]
	return original
}

// SafeSubSlice avoids that by using a full slice expression to cap the capacity. Now any
// append must allocate a fresh array instead of reaching into the parent.
func SafeSubSlice() []int {
	original := []int{1, 2, 3, 4, 5}
	sub := original[1:3:3] // low:high:max — capacity is pinned to 2
	_ = append(sub, 42)    // forced to allocate; original is untouched
	return original
}

// Clone copies the data so the result shares nothing with the input.
func Clone(in []int) []int {
	out := make([]int, len(in))
	copy(out, in)
	return out
}

// ---------------------------------------------------------------------------
// Growth and preallocation
// ---------------------------------------------------------------------------

// GrowthPattern records the capacity after each append, showing the amortised doubling
// (the exact factor is an implementation detail that tapers for large slices).
func GrowthPattern(n int) []int {
	var s []int
	caps := []int{}
	lastCap := -1

	for i := 0; i < n; i++ {
		s = append(s, i)
		if cap(s) != lastCap {
			lastCap = cap(s)
			caps = append(caps, lastCap)
		}
	}
	return caps
}

// BuildNaive reallocates repeatedly as the slice grows.
func BuildNaive(n int) []int {
	var out []int
	for i := 0; i < n; i++ {
		out = append(out, i)
	}
	return out
}

// BuildPrealloc does the same work with exactly one allocation. When the final size is
// known, make([]T, 0, n) is free speed — this is the single highest-yield Go optimisation.
func BuildPrealloc(n int) []int {
	out := make([]int, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, i)
	}
	return out
}

// ---------------------------------------------------------------------------
// Strings and bytes
// ---------------------------------------------------------------------------

// Strings are immutable, so every concatenation allocates a new one. In a loop that is
// quadratic: O(n²) bytes copied for n pieces.
func ConcatNaive(parts []string) string {
	out := ""
	for _, p := range parts {
		out += p
	}
	return out
}

// ConcatBuilder grows one buffer instead. strings.Builder also avoids the final copy from
// []byte to string, because it hands over the buffer it already owns.
func ConcatBuilder(parts []string) string {
	var b strings.Builder

	total := 0
	for _, p := range parts {
		total += len(p)
	}
	b.Grow(total) // one allocation for the whole result

	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}

// RuneVsByte shows that indexing a string yields bytes, while ranging yields runes.
// len() is in bytes — a fact that bites on any non-ASCII input.
func RuneVsByte(s string) (byteLen, runeLen int) {
	byteLen = len(s)
	for range s { // ranging decodes UTF-8
		runeLen++
	}
	return byteLen, runeLen
}

// ---------------------------------------------------------------------------
// Maps
// ---------------------------------------------------------------------------

// Map values are not addressable, so you cannot write m[k].Field = v on a struct value.
// Store pointers, or read-modify-write the whole value.
type Stat struct{ Count int }

// BumpByValue must write the whole struct back.
func BumpByValue(m map[string]Stat, k string) {
	s := m[k] // a copy
	s.Count++
	m[k] = s // put the copy back
}

// BumpByPointer mutates in place because the map holds a pointer.
func BumpByPointer(m map[string]*Stat, k string) {
	s, ok := m[k]
	if !ok {
		s = &Stat{}
		m[k] = s
	}
	s.Count++
}

// ---------------------------------------------------------------------------
// Escape analysis
// ---------------------------------------------------------------------------

// Go decides stack vs heap by itself: a value escapes to the heap when it outlives its
// frame. Returning a pointer to a local is safe here (unlike C) and simply forces the
// allocation. Check with: go build -gcflags='-m' ./07_memory
func EscapesToHeap() *Stat { return &Stat{Count: 1} }

// StaysOnStack keeps the value inside the frame, so the compiler can stack-allocate it.
func StaysOnStack() int {
	s := Stat{Count: 1}
	return s.Count
}
