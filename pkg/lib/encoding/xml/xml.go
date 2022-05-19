package xml

import (
	"reflect"
)

const (
	// ValueNumber ValueNumber
	ValueNumber = "number"
	// ValueString ValueString
	ValueString = "string"
	// ValueBool ValueBool
	ValueBool = "boolean"
	// ValueMap ValueMap
	ValueMap = "map"
	// ValueArray ValueArray
	ValueArray = "array"
)

func typeName(k reflect.Kind) string {
	switch k {
	case reflect.Int, reflect.Uint, reflect.Float32, reflect.Float64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ValueNumber
	case reflect.Bool:
		return ValueBool
	case reflect.String:
		return ValueString
	case reflect.Struct, reflect.Map:
		return ValueMap
	case reflect.Array, reflect.Slice:
		return ValueArray
	}
	return ""
}

func elemName(k reflect.Kind) string {
	switch k {
	case reflect.Int, reflect.Uint, reflect.Float32, reflect.Float64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ValueNumber
	case reflect.Bool:
		return ValueBool
	case reflect.String:
		return ValueString
	default:
		return ValueMap
	}
}
