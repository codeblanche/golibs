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
// r = 4, w = 2, x = 1
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
