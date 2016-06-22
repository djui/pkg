package os

import (
	"log"
	"os"
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
