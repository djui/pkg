package testing

import "testing"

// AssertEqual compares an expected to an actual interface type, and if not
// equal, halts the test.
func AssertEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Fatalf("%s != %s", expected, actual)
	}
}

// AssertEqualString compares an expected to an actual string, and if not equal,
// halts the test.
func AssertEqualString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("%s != %s", expected, actual)
	}
}

// AssertNoError compares an actual error to nil, and if nil, halts the test.
func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("nil != %v", err)
	}
}
