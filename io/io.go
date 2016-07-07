package io

import (
	"io"
	"io/ioutil"
)

// Discard discards an reader's content.
func Discard(r io.ReadCloser) {
	io.CopyN(ioutil.Discard, r, 2<<10)
}

// DiscardAndClose discards an reader's content and closes it.
func DiscardAndClose(r io.ReadCloser) {
	Discard(r)
	r.Close()
}
