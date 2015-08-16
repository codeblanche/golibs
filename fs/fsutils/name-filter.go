package fsutils

import "os"

// NameFilter tests a files name
type NameFilter struct {
	name string
}

// NameFilter implementation of Filter interface.
// Allows files matched by name (including extension)
func (f *NameFilter) test(fi os.FileInfo) bool {
	if fi.Name() == f.name {
		return true
	}

	return false
}

// NewNameFilter creates a new NameFilter
func NewNameFilter(name string) *NameFilter {
	return &NameFilter{
		name: name,
	}
}
