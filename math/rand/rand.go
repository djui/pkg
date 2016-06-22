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
