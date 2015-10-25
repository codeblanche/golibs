package acl

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
)

func TestLevelToString(t *testing.T) {
	assert := assert.New(t)
	cases := map[Level]string{
		Level(0777): "rwxrwxrwx",
		Level(0770): "rwxrwx---",
		Level(0700): "rwx------",
		Level(0111): "--x--x--x",
		Level(0222): "-w--w--w-",
		Level(0444): "r--r--r--",
	}
	for l, expected := range cases {
		actual := LevelToString(l)
		assert.Equal(expected, actual, fmt.Sprintf("Expected %o to result in %s, got %s", l, expected, actual))
	}
}

func TestGetBSON(t *testing.T) {
	assert := assert.New(t)
	cases := map[Level]bsonLevel{
		Level(0777): bsonLevel{true, true, true, true, true, true, true, true, true},
	}
	for l, expected := range cases {
		actual, err := l.GetBSON()
		assert.NoError(err, fmt.Sprintf("Unexepected error getting bson for %o", l))
		assert.EqualValues(expected, actual, fmt.Sprintf("Expected %o to result in %+v, got %+v", l, expected, actual))
	}
}

func TestSetBSON(t *testing.T) {
	assert := assert.New(t)
	cases := map[Level]bsonLevel{
		Level(0777): bsonLevel{true, true, true, true, true, true, true, true, true},
	}
	for expected, bl := range cases {
		actual := Level(0)
		b, err := bson.Marshal(bl)
		assert.NoError(err, fmt.Sprintf("Unexepected error marshalling %+v", bl))
		actual.SetBSON(bson.Raw{Kind: 0x03, Data: b})
		fmt.Printf("%#4o == %#4o\n", expected, actual)
		assert.EqualValues(expected, actual, fmt.Sprintf("Expected %+v to result in %#4o, got %#4o", bl, expected, actual))
	}
	_ = assert
}
