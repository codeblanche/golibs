package fsutils

import (
	"os"
	"regexp"
)

// RegexFilter tests a file name.ext combination against a regular expression
type RegexFilter struct {
	expression *regexp.Regexp
}

// RegexFilter implementation of Filter interface
func (f *RegexFilter) test(fi os.FileInfo) bool {
	return f.expression.Match([]byte(fi.Name()))
}

// NewRegexFilter creates a new RegexFilter
func NewRegexFilter(expression *regexp.Regexp) *RegexFilter {
	return &RegexFilter{
		expression: expression,
	}
}
