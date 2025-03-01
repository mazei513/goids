package main

import (
	"math/rand/v2"
	"testing"
)

func BenchmarkCalc(b *testing.B) {
	all := make([]Boid, 1000)
	out := make([]Vec, 1000)
	for i := range all {
		all[i] = Boid{rand.Float64(), rand.Float64(), rand.Float64(), rand.Float64()}
	}
	b.ResetTimer()
	for range b.N {
		Calc(all, all, out, 0)
	}
}
func BenchmarkCalc2(b *testing.B) {
	all := make([]Boid32, 1000)
	out := make([]Vec32, 1000)
	for i := range all {
		all[i] = Boid32{rand.Float32(), rand.Float32(), rand.Float32(), rand.Float32()}
	}
	b.ResetTimer()
	for range b.N {
		Calc2(all, all, out, 0)
	}
}
