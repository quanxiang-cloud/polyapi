package value

// Number represents number value
type Number float64

// AddElement add a value as child
func (d *Number) AddElement(name string, val interface{}) error {
	return errNotSupport
}

// Set try set the value
func (d *Number) Set(val interface{}) error {
	switch v := val.(type) {
	case Number:
		*d = v
	case *Number:
		*d = *v
	case float64:
		*d = Number(v)
	case *float64:
		*d = Number(*v)
	case float32:
		*d = Number(v)
	case *float32:
		*d = Number(*v)
	case int:
		*d = Number(v)
	case *int:
		*d = Number(*v)
	case int8:
		*d = Number(v)
	case *int8:
		*d = Number(*v)
	case int16:
		*d = Number(v)
	case *int16:
		*d = Number(*v)
	case int32:
		*d = Number(v)
	case *int32:
		*d = Number(*v)
	case int64:
		*d = Number(v)
	case *int64:
		*d = Number(*v)
	case uint:
		*d = Number(v)
	case *uint:
		*d = Number(*v)
	case uint8:
		*d = Number(v)
	case *uint8:
		*d = Number(*v)
	case uint16:
		*d = Number(v)
	case *uint16:
		*d = Number(*v)
	case uint32:
		*d = Number(v)
	case *uint32:
		*d = Number(*v)
	case uint64:
		*d = Number(v)
	case *uint64:
		*d = Number(*v)
	default:
		return errNotSupport
	}
	return nil
}

// Type return the ValueType enum
func (d *Number) Type() Type {
	return TNumber
}
