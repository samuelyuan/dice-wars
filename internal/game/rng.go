package game

import "math/rand/v2"

func rngFromSeed(seed uint64) *rand.Rand {
	return rand.New(rand.NewPCG(seed, 0))
}

// forEachShuffled visits indices 0..n-1 in random order. fn returns true to stop early.
func forEachShuffled(n int, rng *rand.Rand, fn func(i int) bool) {
	if n <= 0 {
		return
	}
	start := rng.IntN(n)
	for offset := 0; offset < n; offset++ {
		if fn((start + offset) % n) {
			return
		}
	}
}
