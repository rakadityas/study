package main

import (
	"testing"
)

func BenchmarkAppend(b *testing.B) {
	var ints []int
	for i := 0; i < b.N; i++ {
		ints = append(ints, i)
	}
}

func BenchmarkPreallocAssign(b *testing.B) {
	ints := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		ints[i] = i
	}
}
