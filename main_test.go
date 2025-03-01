package main

import (
	"fmt"
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
	fmt.Println(out[0])
}
