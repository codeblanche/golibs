package acl

var (
	adminGroup = MakeGroup("admin")
)

// RWX acl contants
const (
	Read    = 'r'
	Write   = 'w'
	Execute = 'x'
)
