package value

// String represents string value
type String string

// AddElement add a value as child
func (d *String) AddElement(name string, val interface{}) error {
	return errNotSupport
}

// Set try set the value
func (d *String) Set(val interface{}) error {
	switch v := val.(type) {
	case String:
		*d = v
	case *String:
		*d = *v
	case string:
		*d = String(v)
	case *string:
		*d = String(*v)
	default:
		return errNotSupport
	}
	return nil
}

// Type return the ValueType enum
func (d *String) Type() Type {
	return TString
}
