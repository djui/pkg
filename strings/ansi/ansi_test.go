package ansi

import (
	"bytes"
	"testing"
)

func TestToHTML(t *testing.T) {
	expected := []byte("")
	actual := ToHTML([]byte(""))

	if !bytes.Equal(expected, actual) {
		t.Fatalf("%s != %s", string(expected), string(actual))
	}
}

func TestRemoveEscapeSequences(t *testing.T) {
	expected := []byte("")
	actual := ToHTML([]byte(""))

	if !bytes.Equal(expected, actual) {
		t.Fatalf("%s != %s", string(expected), string(actual))
	}
}
