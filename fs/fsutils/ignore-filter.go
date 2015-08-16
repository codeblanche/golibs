package fsutils

import "os"

// NameFilter tests a files name
type IgnoreFilter struct {
	name string
}

// IgnoreFilter implementation of Filter interface.
// Ignores files matched by name (including extension)
func (f *IgnoreFilter) test(fi os.FileInfo) bool {
	if fi.Name() != f.name {
		return true
	}

	return false
}

// NewNameFilter creates a new NameFilter
func NewIgnoreFilter(name string) *IgnoreFilter {
	return &IgnoreFilter{
		name: name,
	}
}
