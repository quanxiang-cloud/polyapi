package value

// Object represents object value
type Object map[string]interface{}

func notSupportErr(name string, val interface{}) error {
	return errNotSupport
	//return fmt.Errorf("%s: name=%q type=%t", errNotSupport.Error(), name, val)
}

// AddElement add a value as child
func (d *Object) AddElement(name string, val interface{}) error {
	if name == "" || !Support(val) { // BUG: dont add none-value data
		return notSupportErr(name, val)
	}

	if name != "" && val != nil {
		(*d)[name] = val
		return nil
	}
	return errNotSupport
}

// Set try set the value
func (d *Object) Set(val interface{}) error {
	switch v := val.(type) {
	case Object:
		*d = v
	case *Object:
		*d = *v
	default:
		return errNotSupport
	}
	return nil
}

// Type return the ValueType enum
func (d *Object) Type() Type {
	return TObject
}
