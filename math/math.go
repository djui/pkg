package math

import "math"

// LogN calculates the Log_N of a float64.
func LogN(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}
