package fs

import (
	"os"
)

// Exists checks if path exists
func Exists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		return false
	}

	return true
}

// IsFile checks whether path exists and is a file
func IsFile(path string) bool {
	s, err := os.Stat(path)

	if err != nil || s.IsDir() {
		return false
	}

	return true
}
