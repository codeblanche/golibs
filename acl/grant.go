package acl

// Grant interface defines the methods needed for acl package to validate a user against
// the levels defined in the implementing structs.
type Grant interface {
	// CanUser tests the given user against the grant for the given level
	canUser(User, Level) bool
}

// MakeGrant creates a new grant model
func MakeGrant(uid, group string, level Level) G {
	return G{
		UserID: uid,
		Group:  MakeGroup(group),
		Level:  level,
	}
}

// G type defines the model for access grant
type G struct {
	UserID string
	Group  Group
	Level  Level
}

// canUser implements Grant interface
func (g G) canUser(u User, l Level) bool {
	ugo := Level(0)
	if g.UserID != "" && u.id == g.UserID {
		ugo = ugo | l<<6
	}
	if g.Group != "" && u.groups.Contains(g.Group) {
		ugo = ugo | l<<3
	}
	ugo = ugo | l
	return ugo&g.Level != 0
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
