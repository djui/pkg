package math

// MaxInt returns the larger of x or y.
func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// MaxInt64 returns the larger of x or y.
func MaxInt64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
