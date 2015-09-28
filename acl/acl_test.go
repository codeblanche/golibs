package acl

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set up

	// Run tests
	result := m.Run()

	// Tear down

	// Exit
	os.Exit(result)
}

func TestLoad(t *testing.T) {
	assert := assert.New(t)
	_ = assert
}

func TestMakeGroup(t *testing.T) {
	assert := assert.New(t)
	g := MakeGroup("g1")
	assert.IsType(Group(""), g)
	assert.Equal("g1", string(g))
}

func TestMakeGroups(t *testing.T) {
	assert := assert.New(t)
	g := MakeGroups("g1", "g2", "g3")
	assert.Len(g, 3)
	assert.IsType(Groups{}, g)
}

func TestNew(t *testing.T) {
	assert := assert.New(t)
	// ace := New("123", MakeGroups("g1", "g2", "g3"))
	// assert.IsType(&ACE{}, ace)
	_ = assert
}
