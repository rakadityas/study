package main

import "testing"

func BenchmarkFib10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fib(10)
	}
}

func BenchmarkFib20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fib(20)
	}
}

// go test -bench=. -benchtime=100x

// go test -bench=. -benchtime=10s
