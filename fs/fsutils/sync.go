package fsutils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Synchronise files/directories from a src path to a destination path.
// An optional number of filters may be added to filter which files should be
// synchronised
func Sync(src string, dest string, recurse bool, filters ...Filter) error {
	src, dest, err := resolveAbs(src, dest)

	if err == nil {
		err = copy(src, dest, NewFilterListFilter(filters...), recurse, 0)
	}

	return err
}

// Copy a file or directory from src path to dest path
func copy(src string, dest string, filter *FilterListFilter, recurse bool, depth int) error {
	info, err := os.Stat(src)

	if err == nil {
		if info.Mode().IsDir() && (recurse || depth < 1) {
			err = copyDir(src, dest, filter, recurse, depth)
		}

		if info.Mode().IsRegular() && filter.test(info) {
			err = copyFile(src, dest)
		}
	}

	return err
}

// Recursively iterate through dir contents
func copyDir(src string, dest string, filter *FilterListFilter, recurse bool, depth int) error {
	os.MkdirAll(dest, 0777)

	list, err := ioutil.ReadDir(src)

	if err == nil {
		for _, file := range list {
			err = copy(filepath.Join(src, file.Name()), filepath.Join(dest, file.Name()), filter, recurse, depth+1)

			if err != nil {
				break
			}
		}
	}

	return err
}

// Copy a file from src path to dest path
func copyFile(src string, dest string) error {
	var (
		srcf  *os.File
		destf *os.File
		err   error
	)

	// Open source file
	srcf, err = os.Open(src)

	if err == nil {
		defer srcf.Close()

		// Open destination file
		destf, err = os.Create(dest)

		if err == nil {
			defer destf.Close()

			// Copy data
			_, err = io.Copy(destf, srcf)
		}
	}

	return err
}

// Resolve absolute paths for both src and destination paths
func resolveAbs(src string, dest string) (string, string, error) {
	var err error

	src, err = filepath.Abs(src)

	if err == nil {
		dest, err = filepath.Abs(dest)
	}

	return src, dest, err
}

// Resolve path relative to application
func resolveRel(target string) string {
	cwd, _ := os.Getwd()

	path, err := filepath.Rel(cwd, target)

	if err != nil {
		return target
	}

	return path
}
