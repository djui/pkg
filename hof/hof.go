package hof

// Fixpoint represents a fixed-point combinator. It repeatedly call f with
// x until f(x) equals x."
func Fixpoint(f func(interface{}) interface{}, x interface{}) interface{} {
	return FixPointUntil(f, x, func(x, y interface{}) bool { return x == y })
}

// FixPointUntil represents a fixed-point combinator. It repeatedly call f with
// x until cond(x, f(x)) is true."
func FixPointUntil(f func(interface{}) interface{}, x interface{}, cond func(interface{}, interface{}) bool) interface{} {
	y := f(x)
	if cond(x, y) {
		return y
	}
	return FixPointUntil(f, y, cond)
}
