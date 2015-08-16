package fsutils

import (
	"os"
)

// Filter interface for filesystem utils
type Filter interface {
	test(fi os.FileInfo) bool
}
