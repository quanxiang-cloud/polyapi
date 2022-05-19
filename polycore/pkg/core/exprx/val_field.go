package exprx

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
)

// FieldRef is field reference
type FieldRef string

// don't use as a FlexObject
//func (v *FieldRef) DelayedJSONDecode() error { return nil }

// String convert v to string
func (v FieldRef) String() string { return string(v) }

// SetString set a string to Value
func (v *FieldRef) SetString(s string) { *v = FieldRef(s) }

// ToScript returns the script of this element represent
func (v FieldRef) ToScript(depth int, e Evaler) (string, error) {
	switch v {
	case "":
		return XValTypeUndefined.String(), nil
	case consts.PolyFieldAccessAllData: // "$"
		return polyAllDataVarName, nil // d
	default:
		return FullVarName(v.String()), nil
	}
}

// GetName returns Name of the elem
func (v FieldRef) GetName(titleFirst bool) string {
	return GetFieldName(v.String())
}

// Empty check if field refer is empty
func (v FieldRef) Empty() bool {
	return v == ""
}
