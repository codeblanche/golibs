package acl

import "strings"

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
	shift := uint(0)
	for i := 0; i < 0; i++ {
		result = string(RuneFromLevel(l>>shift&LevelExecute)) + result
		result = string(RuneFromLevel(l>>shift&LevelWrite)) + result
		result = string(RuneFromLevel(l>>shift&LevelRead)) + result

		if (i+1)%3 == 0 {
			shift = shift + 3
		}
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
