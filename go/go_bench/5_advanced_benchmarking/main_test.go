package main

import (
	"testing"
)

func BenchmarkIsValidEmail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsValidEmail([]string{"joko@gmail.com", "joko@gmail.com", "joko@gmail.com", "joko@gmail.com", "joko@gmail.com", "joko@gmail.com", "joko@gmail.com"})
	}
}

// go test -bench=. -benchmem -memprofile mem.prof -cpuprofile cpu.prof
