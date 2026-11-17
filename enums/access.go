package enums

// Access is the access-right flag used in chat permissions.
type Access string

const (
	AccessRead  Access = "Read"
	AccessWrite Access = "Write"
)

// String implements fmt.Stringer.
func (a Access) String() string { return string(a) }
