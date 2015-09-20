package password

import "golang.org/x/crypto/bcrypt"

var (
	// Cost is the default cost setting for hashing the password.
	// See http://godoc.org/golang.org/x/crypto/bcrypt#pkg-constants for more information about cost.
	Cost = bcrypt.DefaultCost
)

// P is the bcrypt hashed password byte slice named type.
// To create a new P from a raw unhashed string use password.Make(string).
type P string

// Make creates a new hashed password from string
func Make(pw string) (P, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), Cost)
	return P(hash), err
}

// String implements the Stringer interface and returns a hashed representation of the data
func (p P) String() string {
	return string(p)
}

// Compare tests a string pw against the Password and returns an error or nil
func (p P) Compare(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(p), []byte(pw))
}

// Match tests a string pw against the Password and returns a bool result
func (p P) Match(pw string) bool {
	return p.Compare(pw) == nil
}
