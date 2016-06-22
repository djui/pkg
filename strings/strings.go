package strings

import "strings"

// PadLeft pads a string with spaces on the left side.
func PadLeft(s string, l int) string {
	return PadLeftWith(s, l, " ")
}

// PadLeftWith pads a string with p on the left side.
func PadLeftWith(s string, l int, p string) string {
	n := len(s) - l
	if n <= 0 {
		return s
	}
	return strings.Repeat(p, n) + s
}

// PadRight pads a string with spaces on the right side.
func PadRight(s string, l int) string {
	return PadRightWith(s, l, " ")
}

// PadRightWith pads a string with p on the right side.
func PadRightWith(s string, l int, p string) string {
	n := len(s) - l
	if n <= 0 {
		return s
	}
	return s + strings.Repeat(p, n)
}
