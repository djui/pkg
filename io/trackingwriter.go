package io

import (
	"bufio"
	"io"
)

// TrackingWriter tracks the last byte written on every write so we can avoid
// printing a newline if one was already written or if there is no output at
// all.
type TrackingWriter struct {
	w    *bufio.Writer
	last byte
}

func NewTrackingWriter(w io.Writer) *TrackingWriter {
	return &TrackingWriter{
		w:    bufio.NewWriter(w),
		last: '\n',
	}
}

func (t *TrackingWriter) Write(p []byte) (n int, err error) {
	n, err = t.w.Write(p)
	if n > 0 {
		t.last = p[n-1]
	}
	return
}

func (t *TrackingWriter) Flush() {
	t.w.Flush()
}

func (t *TrackingWriter) NeedNL() bool {
	return t.last != '\n'
}
