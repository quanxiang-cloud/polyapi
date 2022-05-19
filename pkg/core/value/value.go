package value

import (
	"errors"
)

var errNotSupport = errors.New("not support")

// Type enum
type Type int8

const (
	// TNull is type enum of null
	TNull Type = iota
	// TString is type enum of string
	TString
	// TNumber is type enum of number
	TNumber
	// TInteger is type enum of integer
	TInteger
	// TBoolean is type enum of boolean
	TBoolean
	// TObject is type enum of object
	TObject
	// TArray is type enum of array
	TArray
)

// Support check if a value is support
func Support(val interface{}) bool {
	switch val.(type) {
	case JSONValue:
	case String:
	case Number:
	case Boolean:
	case Object:
	case Array:
	case Null:
	case nil:
	default:
		return false
	}
	return true
}

// JSONValue is the abstract of JSON value
type JSONValue interface {
	AddElement(name string, val interface{}) error
	Set(val interface{}) error
	Type() Type
}

// NewValue generate a value
func (vt Type) NewValue(val interface{}) JSONValue {
	v := vt.newValue()
	if val != nil {
		v.Set(val)
	}
	return v
}

// String show name of the type
func (vt Type) String() string {
	switch vt {
	case TNull:
		return "null"
	case TString:
		return "string"
	case TNumber:
		return "number"
	case TInteger:
		return "integer"
	case TBoolean:
		return "boolean"
	case TObject:
		return "object"
	case TArray:
		return "array"
	}
	return "?"
}

func (vt Type) newValue() JSONValue {
	switch vt {
	case TNull:
		return &Null{}
	case TString:
		return new(String)
	case TNumber:
		return new(Number)
	case TInteger:
		return new(Integer)
	case TBoolean:
		return new(Boolean)
	case TObject:
		return new(Object)
	case TArray:
		return new(Array)
	}
	return nil
}
