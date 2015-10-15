package acl

// User type defines the model for an access controlled user/entity
type User struct {
	id     string
	groups Groups
}

// NewUser creates a new User
func NewUser(id string, groups []string) *User {
	return &User{
		id:     id,
		groups: MakeGroups(groups...),
	}
}

// Can tests a user agains a action and grant
func (u User) Can(action rune, g Grant) bool {
	return u.groups.Contains(adminGroup) || g.canUser(u, LevelFromRune(action))
}
