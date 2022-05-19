package value

// Array represents array value
type Array []interface{}

// AddElement add a value as child
func (d *Array) AddElement(name string, val interface{}) error {
	if val != nil {
		*d = append(*d, val)
		return nil
	}
	return errNotSupport
}

// Set try set the value
func (d *Array) Set(val interface{}) error {
	switch v := val.(type) {
	case Array:
		*d = v
	case *Array:
		*d = *v
	default:
		return errNotSupport
	}
	return nil
}

// Type return the ValueType enum
func (d *Array) Type() Type {
	return TArray
}
