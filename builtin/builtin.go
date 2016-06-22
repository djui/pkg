package builtin

// IfString returns yes when cond is true, no otherwise.
func IfString(cond bool, yes, no string) string {
	if cond {
		return yes
	}
	return no
}

// IfInt returns yes when cond is true, no otherwise.
func IfInt(cond bool, yes, no int) int {
	if cond {
		return yes
	}
	return no
}

// UnlessString returns yes unless cond is true, no otherwise.
func UnlessString(cond bool, yes, no string) string {
	return IfString(!cond, yes, no)
}

// UnlessInt returns yes unless cond is true, no otherwise.
func UnlessInt(cond bool, yes, no int) int {
	return IfInt(!cond, yes, no)
}
