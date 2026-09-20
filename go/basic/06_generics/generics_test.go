package generics

import (
	"errors"
	"strconv"
	"testing"
)

// Celsius exercises the ~ in the Number constraint: a named type over float64.
type Celsius float64

func TestSumAcrossTypes(t *testing.T) {
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Fatalf("Sum(ints) = %d, want 6", got)
	}
	if got := Sum([]float64{1.5, 2.5}); got != 4.0 {
		t.Fatalf("Sum(floats) = %v, want 4", got)
	}
	if got := Sum([]Celsius{10, 20}); got != 30 {
		t.Fatalf("Sum(Celsius) = %v, want 30 — the ~ in the constraint allows named types", got)
	}
}

func TestMax(t *testing.T) {
	if got, ok := Max([]int{3, 9, 2}); !ok || got != 9 {
		t.Fatalf("Max = (%d, %v), want (9, true)", got, ok)
	}
	if got, ok := Max([]string{"pear", "apple", "fig"}); !ok || got != "pear" {
		t.Fatalf("Max = (%q, %v), want (\"pear\", true)", got, ok)
	}
	if _, ok := Max([]int{}); ok {
		t.Fatal("Max of an empty slice should report ok == false")
	}
}

func TestMapFilterReduce(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}

	strs := Map(nums, strconv.Itoa)
	if len(strs) != 6 || strs[0] != "1" || strs[5] != "6" {
		t.Fatalf("Map = %v", strs)
	}

	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	if len(evens) != 3 || evens[0] != 2 {
		t.Fatalf("Filter = %v, want [2 4 6]", evens)
	}

	joined := Reduce(strs, "", func(acc, s string) string { return acc + s })
	if joined != "123456" {
		t.Fatalf("Reduce = %q, want %q", joined, "123456")
	}
}

func TestGroupByAndKeys(t *testing.T) {
	words := []string{"apple", "avocado", "banana", "blueberry", "cherry"}

	groups := GroupBy(words, func(w string) byte { return w[0] })
	if len(groups['a']) != 2 || len(groups['b']) != 2 || len(groups['c']) != 1 {
		t.Fatalf("GroupBy = %v", groups)
	}

	byLen := GroupBy(words, func(w string) int { return len(w) })
	keys := Keys(byLen)
	for i := 1; i < len(keys); i++ {
		if keys[i-1] > keys[i] {
			t.Fatalf("Keys = %v, want ascending order", keys)
		}
	}
}

func TestStack(t *testing.T) {
	var s Stack[string]

	if _, ok := s.Pop(); ok {
		t.Fatal("popping an empty stack should report ok == false")
	}

	s.Push("a")
	s.Push("b")
	if s.Len() != 2 {
		t.Fatalf("Len = %d, want 2", s.Len())
	}

	if v, ok := s.Pop(); !ok || v != "b" {
		t.Fatalf("Pop = (%q, %v), want (\"b\", true)", v, ok)
	}
	if v, ok := s.Pop(); !ok || v != "a" {
		t.Fatalf("Pop = (%q, %v), want (\"a\", true)", v, ok)
	}
	if s.Len() != 0 {
		t.Fatalf("Len = %d, want 0", s.Len())
	}
}

func TestResult(t *testing.T) {
	ok := Ok(42)
	if !ok.IsOk() {
		t.Fatal("Ok should be ok")
	}
	if v, err := ok.Unwrap(); v != 42 || err != nil {
		t.Fatalf("Unwrap = (%d, %v), want (42, nil)", v, err)
	}

	boom := errors.New("boom")
	bad := Err[int](boom)
	if bad.IsOk() {
		t.Fatal("Err should not be ok")
	}
	if _, err := bad.Unwrap(); !errors.Is(err, boom) {
		t.Fatalf("Unwrap err = %v, want boom", err)
	}
}
