package testing

import (
	"path/filepath"
	"runtime"
	"testing"
)

// AssertEqual compares an expected to an actual interface type, and if not
// equal, halts the test.
func AssertEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		_, fn, line, _ := runtime.Caller(1)
		t.Fatalf("%s:%d: %s != %s", filepath.Base(fn), line, expected, actual)
	}
}

// AssertEqualString compares an expected to an actual string, and if not equal,
// halts the test.
func AssertEqualString(t *testing.T, expected, actual string) {
	if expected != actual {
		_, fn, line, _ := runtime.Caller(1)
		t.Fatalf("%s:%d: %s != %s", filepath.Base(fn), line, expected, actual)
	}
}

// AssertNoError compares an actual error to nil, and if nil, halts the test.
func AssertNoError(t *testing.T, err error) {
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		t.Fatalf("%s:%d: nil != %v", filepath.Base(fn), line, err)
	}
}
