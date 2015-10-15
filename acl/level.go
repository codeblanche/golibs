package acl

// Level constant type
type Level int

// Level constant values
const (
	LevelRead    Level = 04
	LevelWrite   Level = 02
	LevelExecute Level = 01
)

// LevelFromRune converts an action rune into it's corresponding Level value
// r => 4, w => 2, x => 1
func LevelFromRune(action rune) Level {
	switch action {
	case 'r':
		return LevelRead
	case 'w':
		return LevelWrite
	case 'x':
		return LevelExecute
	}
	return 0
}

// RuneFromLevel converts a Level value into it's corresponding rune
// 4 => r, 2 => w, 1 => x
func RuneFromLevel(l Level) rune {
	switch l {
	case LevelRead:
		return 'r'
	case LevelWrite:
		return 'w'
	case LevelExecute:
		return 'x'
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
