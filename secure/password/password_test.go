package password

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	PW = "i do not like green eggs and ham"
)

var (
	p  P
	p2 P
)

func TestMain(m *testing.M) {
	// Set up
	p, _ = Make(PW)
	p2 = P(p.String())

	// Run tests
	result := m.Run()

	// Tear down

	// Exit
	os.Exit(result)
}

func TestPassword(t *testing.T) {
	assert := assert.New(t)
	p, err := Make(PW)
	assert.NoError(err)
	assert.NotNil(p)
	assert.NotEmpty(p)
	assert.IsType(P(""), p)
}

func TestString(t *testing.T) {
	assert := assert.New(t)
	assert.NotEqual(PW, string(p))
}

func TestCompare(t *testing.T) {
	assert := assert.New(t)
	assert.NoError(p.Compare(PW))
	assert.NoError(p2.Compare(PW))
	assert.Error(p.Compare("I do not like them sam i am"))
	assert.Error(p2.Compare("I do not like them sam i am"))
}

func TestMatch(t *testing.T) {
	assert := assert.New(t)
	assert.True(p.Match(PW))
	assert.True(p2.Match(PW))
	assert.False(p.Match("I do not like them here or there"))
	assert.False(p2.Match("I do not like them here or there"))
}
