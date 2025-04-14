package main

import (
	"fmt"
	"testing"
)

func benchmarkFibCore(num int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fib(num)
	}
}

func BenchmarkFib(b *testing.B) {
	nums := []int{10, 20, 30, 40}
	for _, num := range nums {
		b.Run(fmt.Sprintf("BenchmarkFib-%d", num), func(b *testing.B) {
			benchmarkFibCore(num, b)
		})
	}
}
