package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	_ = Add("pprof-test")
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add("pprof-test")
	}
}
