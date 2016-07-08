package net

import "net"

// IsTimeout returns a boolean indicating whether the error is known to report
// the network request timed out.
func IsTimeout(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

// IsTemporary returns a boolean indicating whether the error is known to report
// the network request had a temporary failure.
func IsTemporary(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Temporary()
}

// IsAccept returns a boolean indicating whether the error is known to report
// the network server gave up to accept and handle requests.
func IsAccept(err error) bool {
	opErr, ok := err.(*net.OpError)
	return ok && opErr.Op == "accept"
}
