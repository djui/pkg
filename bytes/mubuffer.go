package bytes

import (
	"bytes"
	"sync"
)

// MuBuffer defines a thread-safe bytes buffer. The zero value for MuBuffer is
// an empty buffer ready to use.
type MuBuffer struct {
	mu sync.Mutex
	*bytes.Buffer
}

// Truncate discards all but the first n unread bytes from the buffer but
// continues to use the same allocated storage. It panics if n is negative or
// greater than the length of the buffer.
func (b *MuBuffer) Truncate(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Truncate(n)
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the buffer
// without another allocation. If n is negative, Grow will panic. If the buffer
// can't grow it will panic with ErrTooLarge.
func (b *MuBuffer) Grow(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Grow(n)
}

// Write appends the contents of p to the buffer, growing the buffer as needed.
// The return value n is the length of p; err is always nil. If the buffer
// becomes too large, Write will panic with ErrTooLarge.
func (b *MuBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Write(p)
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with ErrTooLarge.
func (b *MuBuffer) WriteString(s string) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.WriteString(s)
}

// WriteByte appends the byte c to the buffer, growing the buffer as needed. The
// returned error is always nil, but is included to match bufio.Writer's
// WriteByte. If the buffer becomes too large, WriteByte will panic with
// ErrTooLarge.
func (b *MuBuffer) WriteByte(c byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.WriteByte(c)
}

// NewBuffer creates and initializes a new MuBuffer using buf as its initial
// contents. It is intended to prepare a MuBuffer to read existing data. It can
// also be used to size the internal buffer for writing. To do that, buf should
// have the desired capacity but a length of zero.
//
// In most cases, new(MuBuffer) (or just declaring a MuBuffer variable) is
// sufficient to initialize a MuBuffer.
func NewBuffer(buf []byte) *MuBuffer {
	return &MuBuffer{Buffer: bytes.NewBuffer(buf)}
}

// NewBufferString creates and initializes a new MuBuffer using string s as its
// initial contents. It is intended to prepare a buffer to read an existing
// string.
//
// In most cases, new(MuBuffer) (or just declaring a MuBuffer variable) is
// sufficient to initialize a MuBuffer.
func NewBufferString(s string) *MuBuffer {
	return &MuBuffer{Buffer: bytes.NewBufferString(s)}
}
