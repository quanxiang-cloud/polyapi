package value

// Integer represents integer value
type Integer int64

// AddElement add a value as child
func (d *Integer) AddElement(name string, val interface{}) error {
	return errNotSupport
}

// Set try set the value
func (d *Integer) Set(val interface{}) error {
	switch v := val.(type) {
	case Integer:
		*d = v
	case *Integer:
		*d = *v
	case float64:
		*d = Integer(v)
	case *float64:
		*d = Integer(*v)
	case float32:
		*d = Integer(v)
	case *float32:
		*d = Integer(*v)
	case int:
		*d = Integer(v)
	case *int:
		*d = Integer(*v)
	case int8:
		*d = Integer(v)
	case *int8:
		*d = Integer(*v)
	case int16:
		*d = Integer(v)
	case *int16:
		*d = Integer(*v)
	case int32:
		*d = Integer(v)
	case *int32:
		*d = Integer(*v)
	case int64:
		*d = Integer(v)
	case *int64:
		*d = Integer(*v)
	case uint:
		*d = Integer(v)
	case *uint:
		*d = Integer(*v)
	case uint8:
		*d = Integer(v)
	case *uint8:
		*d = Integer(*v)
	case uint16:
		*d = Integer(v)
	case *uint16:
		*d = Integer(*v)
	case uint32:
		*d = Integer(v)
	case *uint32:
		*d = Integer(*v)
	case uint64:
		*d = Integer(v)
	case *uint64:
		*d = Integer(*v)
	default:
		return errNotSupport
	}
	return nil
}

// Type return the ValueType enum
func (d *Integer) Type() Type {
	return TInteger
}
