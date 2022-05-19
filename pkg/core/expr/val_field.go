package expr

// FieldRef is field reference
type FieldRef string

// don't use as a FlexObject
//func (v *FieldRef) DelayedJSONDecode() error { return nil }

// String convert v to string
func (v FieldRef) String() string { return string(v) }

// SetString set a string to Value
func (v *FieldRef) SetString(s string) { *v = FieldRef(s) }

// GetName returns Name of the elem
func (v FieldRef) GetName(titleFirst bool) string {
	return GetFieldName(v.String())
}

// Empty check if field refer is empty
func (v FieldRef) Empty() bool {
	return v == ""
}
