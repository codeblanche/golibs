package fsutils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDir string
	file1   *os.File
	file2   *os.File
	file3   *os.File
	info1   os.FileInfo
	info2   os.FileInfo
	info3   os.FileInfo
)

type MockFilter struct {
	val    bool
	called bool
}

func (s *MockFilter) test(fi os.FileInfo) bool {
	s.called = true

	return s.val
}

type MockFilePath struct {
}

func (m *MockFilePath) Abs(path string) (string, error) {
	return "", errors.New("Mock Error")
}

func setup() {
	var err error
	testDir, err = ioutil.TempDir(os.TempDir(), "fsutils")

	file1, _ = os.Create(filepath.Join(testDir, "somefile1.abc"))
	file2, _ = os.Create(filepath.Join(testDir, "somefile2.123"))
	file3, _ = os.Create(filepath.Join(testDir, "somefile3.xyz"))

	info1, _ = file1.Stat()
	info2, _ = file2.Stat()
	info3, _ = file3.Stat()

	os.MkdirAll(filepath.Join(testDir, "src/a/b/c/"), 0777)

	for _, path := range []string{"src/", "scr/a/", "src/a/b/", "src/a/b/c/"} {
		for i := 0; i <= 3; i++ {
			filePath := filepath.Join(testDir, path, fmt.Sprintf("f%d.txt", i))

			ioutil.WriteFile(filePath, []byte(filePath), 0664)
		}
	}

	if err != nil {
		os.Exit(1)
	}
}

func teardown() {
	os.Remove(filepath.Join(testDir, "somefile1.abc"))
	os.Remove(filepath.Join(testDir, "somefile2.123"))
	os.Remove(filepath.Join(testDir, "somefile3.xyz"))
	os.Remove(testDir)
}

// pathExists checks whether the given path exists
func pathExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		return false
	}

	return os.IsNotExist(err)
}

func TestMain(m *testing.M) {
	// Set up
	setup()

	// Run tests
	result := m.Run()

	// Tear down
	teardown()

	// Exit
	os.Exit(result)
}

func TestExtensionFilter(t *testing.T) {
	assert := assert.New(t)

	f := NewExtensionFilter("123")

	assert.IsType(&ExtensionFilter{}, f, "Expected object to be of type ExtensionFilter")
	assert.False(f.test(info1), "Expected test method result to be false")
	assert.True(f.test(info2), "Expected test method result to be true")
}

func TestFilterListFilter(t *testing.T) {
	assert := assert.New(t)

	s1 := MockFilter{val: true}
	s2 := MockFilter{val: true}
	f := NewFilterListFilter([]Filter{&s1, &s2})
	info, _ := os.Stat(".")

	assert.IsType(&FilterListFilter{}, f, "Expected object to be of type FilterListFilter")

	assert.True(f.test(info), "Expected test method result to be true")

	s1.val = false

	assert.False(f.test(info), "Expected test method result to be false")
}

func TestNameFilter(t *testing.T) {
	assert := assert.New(t)

	f := NewNameFilter("somefile1")

	assert.IsType(&NameFilter{}, f, "Expected object to be of type ExtensionFilter")
	assert.True(f.test(info1), "Expected test method result to be true")
	assert.False(f.test(info2), "Expected test method result to be false")
}

func TestRegexpFilter(t *testing.T) {
	assert := assert.New(t)

	r := regexp.MustCompile("^somefile[0-9].(abc|123)$")
	f := NewRegexFilter(r)

	assert.IsType(&RegexFilter{}, f, "Expected object to be of type RegexpFilter")
	assert.True(f.test(info1), "Expected test method result to be true")
	assert.True(f.test(info2), "Expected test method result to be true")
	assert.False(f.test(info3), "Expected test method result to be false")
}

func TestSync(t *testing.T) {
	assert := assert.New(t)

	err := Sync(filepath.Join(testDir, "src/"), filepath.Join(testDir, "dest/"), true)

	assert.Nil(err, "Expected nil value for error result")

	err = Sync(filepath.Join(testDir, "srcs/"), filepath.Join(testDir, "dest/"), true)

	assert.Error(err, "Expected return value to be an error")

	err = Sync(filepath.Join(testDir, "src/"), "/dev/null", true)

	assert.Error(err, "Expected return value to be an error")
}

func TestResolveRel(t *testing.T) {
	assert.Equal(t, "a/b/c", resolveRel("a/b/c"), "Expected resolveRel to return target value")
}

func TestResolveAbs(t *testing.T) {
	assert := assert.New(t)

	_, _, err := resolveAbs("a/b/c/", "a/b/c/")

	assert.Nil(err, "Expected nil value for error result")
}
