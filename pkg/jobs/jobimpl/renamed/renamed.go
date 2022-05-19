package renamed

// map of name change
var (
	Raw       = RenameMap{}
	Poly      = RenameMap{}
	Namespace = RenameMap{}
	Service   = RenameMap{}
)

//------------------------------------------------------------------------------

// RenameMap map name that has changed
type RenameMap map[string]string

// Add add a name change to map
func (m RenameMap) Add(oldPath, newPath string) {
	m[oldPath] = newPath // replace mode
}

// Query find if name has change
func (m RenameMap) Query(oldPath string) (string, bool) {
	s, ok := m[oldPath]
	return s, ok
}

// Empty check if map is empty
func (m RenameMap) Empty() bool {
	return len(m) == 0
}
