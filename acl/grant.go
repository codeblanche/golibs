package acl

// Grant interface defines the methods needed for acl package to validate a user against
// the levels defined in the implementing structs.
type Grant interface {
	// CanUser tests the given user against the grant for the given level
	canUser(User, Level) bool
	// Marshal the Grant into a byte array
	Marshal(interface{}) ([]byte, error)
	// Unmarshal the Grant from a byte array
	Unmarshal([]byte, interface{}) error
}

// NewGrant creates a new grant model
func NewGrant(uid, group, level string) *G {
	return &G{
		uid:   uid,
		group: MakeGroup(group),
		level: LevelFromString(level),
	}
}

// G type defines the model for access grant
type G struct {
	uid   string
	group Group
	level Level
}

// canUser implements Grant interface
func (g *G) canUser(u User, l Level) bool {
	ugo := Level(0)
	if u.id == g.uid {
		ugo = ugo | l<<6
	}
	if u.groups.Contains(g.group) {
		ugo = ugo | l<<3
	}
	ugo = ugo | l
	return ugo&g.level != 0
}

// GL type is a list of G that also implements the Grant interface
type GL []G

// canUser implements Grant interface
func (g GL) canUser(u User, l Level) bool {
	for _, grant := range g {
		if grant.canUser(u, l) {
			return true
		}
	}
	return false
}
