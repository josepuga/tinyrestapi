package main

import (
	"errors"
	"os"
)

// fileExists check if file exists.
// This funcion works with negative !fileExists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true // other error (permissions, I/O, etc.)
	}
}
