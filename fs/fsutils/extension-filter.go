package fsutils

import (
	"os"
	"strings"
)

// ExtensionFilter tests a files extension
type ExtensionFilter struct {
	extension string
}

// ExtensionFilter implementation of Filter interface. Assumes that the extension
// is everything after the first "."
func (f *ExtensionFilter) test(fi os.FileInfo) bool {
	ext := strings.Join(strings.Split(fi.Name(), ".")[1:], ".")

	if ext == f.extension {
		return true
	}

	return false
}

// NewExtensionFilter creates a new ExtensionFilter
func NewExtensionFilter(ext string) *ExtensionFilter {
	return &ExtensionFilter{
		extension: ext,
	}
}
