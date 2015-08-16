package slice

// S represents the slice type on which slice helper methods are attached.
// Slice may contains any value.
type S []interface{}

// Strings converts a slice of strings to a slice.S
func Strings(in []string) S {
	out := S{}
	for _, v := range in {
		out = append(out, v)
	}
	return out
}

// Contains checks a slice for the existence of the given value
func (s S) Contains(value interface{}) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}
