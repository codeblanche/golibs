package acl

import (
	"strings"

	"gopkg.in/mgo.v2/bson"
)

// Level constant type
type Level int

// Level constant values
const (
	LevelRead    Level = 04
	LevelWrite   Level = 02
	LevelExecute Level = 01
	LevelNone    Level = 00
)

// LevelFromRune converts an action rune into it's corresponding Level value
// r => 4, w => 2, x => 1
func LevelFromRune(action rune) Level {
	switch action {
	case Read:
		return LevelRead
	case Write:
		return LevelWrite
	case Execute:
		return LevelExecute
	}
	return 0
}

// RuneFromLevel converts a Level value into it's corresponding rune
// 4 => r, 2 => w, 1 => x
func RuneFromLevel(l Level) rune {
	switch l {
	case LevelRead:
		return Read
	case LevelWrite:
		return Write
	case LevelExecute:
		return Execute
	}
	return '-'
}

// LevelFromString converts an 9 character string of action runes to it's corresponding
// level value. Dashes can be used as a nil place holder for an action rune.
// For example:
// rwx------ means read, write, and execute for user only. Group and other receive no access.
// r--r--r-- means read access for all types with no write or execute access.
func LevelFromString(l string) Level {
	result := Level(0)
	shift := uint(6)
	for i, r := range []byte(l)[:9] {
		result = result | LevelFromRune(rune(r))<<shift

		if (i+1)%3 == 0 {
			shift = shift - 3
		}
	}
	return result
}

// LevelToString converts a Level to it's 9 character human readable form
func LevelToString(l Level) string {
	result := ""
	shift := uint(6)
	for i := 0; i < 3; i++ {
		result = result + string(RuneFromLevel(l>>shift&LevelRead))
		result = result + string(RuneFromLevel(l>>shift&LevelWrite))
		result = result + string(RuneFromLevel(l>>shift&LevelExecute))
		shift = shift - 3
	}
	return result
}

// String implements Stringer
func (l Level) String() string {
	return LevelToString(l)
}

// LevelFromExpression converts a permission express in the form of ugo+w or u+rw,g+r to it's corresponding
// level value. This function has the accidental genius ability to process a complex expression in the form
// of uwxgor which translates to ugo+r,u+wx though admitedly more complicated to read and understand.
func LevelFromExpression(e string) Level {
	u, g, o := LevelNone, LevelNone, LevelNone
	expressions := strings.Split(e, ",")
	for _, expression := range expressions {
		// Reverse the expression so level bytes (rwx) are process first followed by assignment bytes (ugo)
		bytes := reverse([]byte(expression))
		l := LevelNone
		for _, b := range bytes {
			switch b {
			case Read:
				l = l | LevelRead
			case Write:
				l = l | LevelWrite
			case Execute:
				l = l | LevelExecute
			case 'u':
				u = u | l
			case 'g':
				g = g | l
			case 'o':
				o = o | l
			}
		}
	}
	return (u << 6) | (g << 3) | o
}

// Reverse a byte slice
func reverse(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

// Need to create a new type for SetBSON to prevent recursion
type bsonLevel struct {
	UserCanRead     bool `bson:"user_can_read,omitempty"`
	UserCanWrite    bool `bson:"user_can_write,omitempty"`
	UserCanExecute  bool `bson:"user_can_execute,omitempty"`
	GroupCanRead    bool `bson:"group_can_read,omitempty"`
	GroupCanWrite   bool `bson:"group_can_write,omitempty"`
	GroupCanExecute bool `bson:"group_can_execute,omitempty"`
	OtherCanRead    bool `bson:"other_can_read,omitempty"`
	OtherCanWrite   bool `bson:"other_can_write,omitempty"`
	OtherCanExecute bool `bson:"other_can_execute,omitempty"`
}

// GetBSON implements bson.Getter
func (l Level) GetBSON() (interface{}, error) {
	t := bsonLevel{}
	t.UserCanRead = (l&0400 > 0)
	t.UserCanWrite = (l&0200 > 0)
	t.UserCanExecute = (l&0100 > 0)
	t.GroupCanRead = (l&0040 > 0)
	t.GroupCanWrite = (l&0020 > 0)
	t.GroupCanExecute = (l&0010 > 0)
	t.OtherCanRead = (l&0004 > 0)
	t.OtherCanWrite = (l&0002 > 0)
	t.OtherCanExecute = (l&0001 > 0)
	return t, nil
}

// SetBSON implements bson.Setter
func (l *Level) SetBSON(raw bson.Raw) (err error) {
	var t bsonLevel
	err = raw.Unmarshal(&t)
	if err != nil {
		return err
	}
	m := map[Level]bool{
		Level(0400): t.UserCanRead,
		Level(0200): t.UserCanWrite,
		Level(0100): t.UserCanExecute,
		Level(0040): t.GroupCanWrite,
		Level(0020): t.GroupCanRead,
		Level(0010): t.GroupCanExecute,
		Level(0004): t.OtherCanRead,
		Level(0002): t.OtherCanWrite,
		Level(0001): t.OtherCanExecute,
	}
	for p, b := range m {
		if !b {
			continue
		}
		*l = *l | p
	}
	return
}
