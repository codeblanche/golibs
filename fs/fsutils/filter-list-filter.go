package fsutils

import (
	"os"
)

// FilterListFilter is a list of filters but is also filter itself
type FilterListFilter struct {
	list []Filter
}

// FilterListFilter implementation of Filter interface
func (f *FilterListFilter) test(fi os.FileInfo) bool {
	state := true

	for _, filter := range f.list {
		state = state && filter.test(fi)

		if !state {
			break
		}
	}

	return state
}

// NewFilterListFilter creates a new FilterListFilter
func NewFilterListFilter(filter ...Filter) *FilterListFilter {
	return &FilterListFilter{
		list: filter,
	}
}
