package logr

import "strings"

//go:generate stringer -type=Type

type (
	// Type of log item
	Type int
)

// Available log types
const (
	None Type = iota
	P    Type = 1 << iota
	E
	W
	I
	D
	S

	Critical = P | E
	Monitor  = Critical | W
	Verbose  = Monitor | I | S
	All      = Verbose | D
)

var (
	types = "PEWIDS"
)

// Itot converts and integer to a log type/level
func Itot(i int) Type {
	return Type(i)
}

// Atot converts a string of runes to a log type/level
func Atot(s string) Type {
	t := Type(0)
	for _, r := range s {
		t |= Rtot(r)
	}
	return t
}

// Rtot converts a rune to a log type/level
func Rtot(r rune) Type {
	i := strings.IndexRune(types, r)
	if i == -1 {
		return Type(0)
	}
	return Type(i)
}
