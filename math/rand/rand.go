package rand

import (
	"math/rand"
	"time"
)

// Duration returns, a pseudo-random duration in [0,d) from the default Source.
// It panics if d <= 0.
func Duration(d time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(d)))
}

// PermInt permutates a given slice of integers and returns a new slice.
func PermInt(x []int) []int {
	y := make([]int, len(x))
	perm := rand.Perm(len(x))
	for i, v := range perm {
		y[v] = x[i]
	}
	return y
}

// PermIntInplace permutates a given slice of integers in-place.
func PermIntInplace(x []int) {
	for i := range x {
		j := rand.Intn(i + 1)
		x[i], x[j] = x[j], x[i]
	}
}

// PermString permutates a given slice of integers and returns a new slice.
func PermString(x []string) []string {
	y := make([]string, len(x))
	perm := rand.Perm(len(x))
	for i, v := range perm {
		y[v] = x[i]
	}
	return y
}

// PermStringInplace permutates a given slice of integers in-place.
func PermStringInplace(x []string) {
	for i := range x {
		j := rand.Intn(i + 1)
		x[i], x[j] = x[j], x[i]
	}
}
