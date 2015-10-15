package acl

// Group ...
type Group string

// Groups ...
type Groups []Group

// MakeGroup creates a new Group
func MakeGroup(name string) Group {
	return Group(name)
}

// MakeGroups creates a new Groups
func MakeGroups(n ...string) Groups {
	g := Groups{}
	for _, name := range n {
		g = append(g, MakeGroup(name))
	}
	return g
}

func (g Group) String() string {
	return string(g)
}

// Contains checks if groups list contains a given group
func (g Groups) Contains(group Group) bool {
	for _, v := range g {
		if v == group {
			return true
		}
	}
	return false
}
