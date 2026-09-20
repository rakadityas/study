package memory

import (
	"fmt"
	"strings"
	"testing"
)

func TestAliasing(t *testing.T) {
	original, sub := Aliasing()

	if original[1] != 99 {
		t.Fatalf("original = %v, want the write through sub to be visible at index 1", original)
	}
	if sub[0] != 99 {
		t.Fatalf("sub = %v, want [99 3]", sub)
	}
	if cap(sub) != 4 {
		t.Fatalf("cap(sub) = %d, want 4 — capacity runs to the end of the backing array", cap(sub))
	}
}

func TestAppendCanClobber(t *testing.T) {
	if got := AppendCanClobber(); got[3] != 42 {
		t.Fatalf("got = %v, want index 3 clobbered to 42", got)
	}

	if got := SafeSubSlice(); got[3] != 4 {
		t.Fatalf("got = %v, want index 3 untouched by the capped sub-slice", got)
	}
}

func TestClone(t *testing.T) {
	in := []int{1, 2, 3}
	out := Clone(in)
	out[0] = 99

	if in[0] != 1 {
		t.Fatalf("Clone shared storage: in = %v", in)
	}
}

func TestGrowthPattern(t *testing.T) {
	caps := GrowthPattern(100)

	if len(caps) < 2 {
		t.Fatalf("expected several growth steps, got %v", caps)
	}
	for i := 1; i < len(caps); i++ {
		if caps[i] <= caps[i-1] {
			t.Fatalf("capacity should grow monotonically, got %v", caps)
		}
	}
	// growth is amortised: ~log(n) reallocations for n appends, not n
	if len(caps) > 20 {
		t.Fatalf("expected amortised growth, saw %d reallocations for 100 appends", len(caps))
	}
}

func TestBuildEquivalence(t *testing.T) {
	naive, prealloc := BuildNaive(50), BuildPrealloc(50)

	if len(naive) != len(prealloc) {
		t.Fatalf("lengths differ: %d vs %d", len(naive), len(prealloc))
	}
	if cap(prealloc) != 50 {
		t.Fatalf("cap(prealloc) = %d, want exactly 50", cap(prealloc))
	}
	for i := range naive {
		if naive[i] != prealloc[i] {
			t.Fatalf("contents differ at %d", i)
		}
	}
}

func TestConcat(t *testing.T) {
	parts := strings.Split("a b c d e", " ")
	want := "abcde"

	if got := ConcatNaive(parts); got != want {
		t.Fatalf("ConcatNaive = %q, want %q", got, want)
	}
	if got := ConcatBuilder(parts); got != want {
		t.Fatalf("ConcatBuilder = %q, want %q", got, want)
	}
}

func TestRuneVsByte(t *testing.T) {
	byteLen, runeLen := RuneVsByte("héllo")

	if byteLen != 6 {
		t.Fatalf("byteLen = %d, want 6 — é is two bytes in UTF-8", byteLen)
	}
	if runeLen != 5 {
		t.Fatalf("runeLen = %d, want 5", runeLen)
	}
}

func TestMapValueAddressability(t *testing.T) {
	byValue := map[string]Stat{}
	BumpByValue(byValue, "a")
	BumpByValue(byValue, "a")
	if byValue["a"].Count != 2 {
		t.Fatalf("BumpByValue = %d, want 2", byValue["a"].Count)
	}

	byPointer := map[string]*Stat{}
	BumpByPointer(byPointer, "a")
	BumpByPointer(byPointer, "a")
	if byPointer["a"].Count != 2 {
		t.Fatalf("BumpByPointer = %d, want 2", byPointer["a"].Count)
	}
}

// Benchmarks: run with
//
//	go test -bench=. -benchmem ./07_memory
//
// -benchmem is what turns these into an allocation story rather than just a timing one.
func BenchmarkBuildNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildNaive(1000)
	}
}

func BenchmarkBuildPrealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildPrealloc(1000)
	}
}

func BenchmarkConcat(b *testing.B) {
	parts := make([]string, 200)
	for i := range parts {
		parts[i] = fmt.Sprintf("part-%d-", i)
	}

	b.Run("naive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ConcatNaive(parts)
		}
	})
	b.Run("builder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ConcatBuilder(parts)
		}
	})
}
