package password

import (
	"crypto/rand"
	"io"

	"github.com/codeblanche/golibs/logr"

	"golang.org/x/crypto/bcrypt"
)

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

// Generate a new password of given length
func Generate(length int) string {
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")
	n := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, r); err != nil {
			logr.Errorf("Unable to read from rand.Reader with error: %s", err.Error())
		}
		for _, c := range r {
			if c >= maxrb {
				continue
			}
			n[i] = chars[c%clen]
			i++
			if i == length {
				return string(n)
			}
		}
	}
}
