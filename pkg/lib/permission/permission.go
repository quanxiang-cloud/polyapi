package permission

// permit bit define
const (
	PermitBitRead PermitBit = 1 << iota
	PermitBitExecute
	PermitBitCreate
	PermitBitUpdate
	PermitBitDelete
	PermitBitGrant

	// PermitBitNone represents none permission
	PermitBitNone PermitBit = 0
	// PermitBitALL represents all permission
	PermitBitALL = PermitBitRead | PermitBitExecute | PermitBitCreate |
		PermitBitUpdate | PermitBitDelete | PermitBitGrant
)

// PermitBit represents permit privilege bitset
type PermitBit uint

// Uint convert permit as uint
func (p PermitBit) Uint() uint {
	return uint(p)
}

// String show name of the permit content
func (p PermitBit) String() string {
	s := permitBit[p]
	return s
}

// Contains check permit contans bitset
func (p PermitBit) Contains(permit PermitBit) bool {
	return p&permit == permit
}

// IsNone check if permit is none
func (p PermitBit) IsNone() bool {
	return p == PermitBitNone
}

// IsAll check if permit is all
func (p PermitBit) IsAll() bool {
	return p == PermitBitALL
}

// ToList convert permit bits to string list
func (p PermitBit) ToList() []string {
	var r []string
	for n := p; n != 0; n &= (n - 1) { // foreach '1' in permit
		b := (n ^ (n - 1) + 1) / 2 // the lowest '1' in n
		if s, ok := permitBit[b]; ok {
			r = append(r, s)
		} else {
			println("unregistered permit bit", b)
		}
	}
	return r
}

// ToPermitList convert access bits to
func ToPermitList(access uint) []string {
	return PermitBit(access).ToList()
}
