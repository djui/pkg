package os

import (
	"log"
	"os"
)

const (
	// PrivateFileMode grants owner to read/write a file.
	PrivateFileMode = 0600
	// PrivateDirMode grants owner to make/remove files inside the directory.
	PrivateDirMode = 0700
)

// Exists checks if a given filepath exists. Note that it's usually better to
// just use os.Open/1 when interacting with the file to avoid race-conditions.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// LogCloseError is a helper function which logs close errors. Useful in defers.
// Optionally accepts a format string. e.g, defer LogCloseError(fd.Close,
// "Failed to close fd")
func LogCloseError(doClose func() error, fmtString ...string) {
	if err := doClose(); err != nil {
		if len(fmtString) == 1 {
			log.Printf("%v: %v", fmtString[0], err)
		} else {
			log.Printf("Error closing object: %T %v", err, err)
		}
	}
}
