package value

import (
	"encoding/json"
)

var null = []byte("null")

// Null represents null value
type Null struct{}

// AddElement add a value as child
func (d *Null) AddElement(name string, val interface{}) error {
	return errNotSupport
}

// Set try set the value
func (d *Null) Set(val interface{}) error {
	return errNotSupport
}

// Type return the ValueType enum
func (d *Null) Type() Type {
	return TNull
}

// MarshalJSON adapt json Marshaler
func (d Null) MarshalJSON() ([]byte, error) {
	return null, nil
}

// UnmarshalJSON adapt json Unmarshaler
func (d *Null) UnmarshalJSON(data []byte) error {
	return nil
}

var _ json.Marshaler = (*Null)(nil)
var _ json.Unmarshaler = (*Null)(nil)
