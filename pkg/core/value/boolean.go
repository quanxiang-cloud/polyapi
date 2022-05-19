package value

// Boolean represents boolean value
type Boolean bool

// AddElement add a value as child
func (d *Boolean) AddElement(name string, val interface{}) error {
	return errNotSupport
}

// Set try set the value
func (d *Boolean) Set(val interface{}) error {
	switch v := val.(type) {
	case Boolean:
		*d = v
	case *Boolean:
		*d = *v
	case bool:
		*d = Boolean(v)
	case *bool:
		*d = Boolean(*v)
	default:
		return errNotSupport
	}
	return nil
}

// Type return the ValueType enum
func (d *Boolean) Type() Type {
	return TBoolean
}
