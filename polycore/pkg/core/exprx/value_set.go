package exprx

import (
	"bytes"
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/pkg/core/value"
)

// SwagConstValue is predef value
type SwagConstValue = XInputValue

// SwagConstValueSet is predef value set
type SwagConstValueSet = ValueSet

//------------------------------------------------------------------------------

// XInputValue is entended InputValue
type XInputValue expr.InputValue

func (v *XInputValue) base() *expr.InputValue { return (*expr.InputValue)(v) }

// Empty check if the value is empty
func (v *XInputValue) Empty() bool {
	return v.base().Empty()
}

// DelayedJSONDecode delay unmarshal flex json object
func (v *XInputValue) DelayedJSONDecode() error {
	return v.base().DelayedJSONDecode()
}

// GetName return name of input value
func (v XInputValue) GetName(titleFirst bool) string {
	return v.base().GetName(titleFirst)
}

// GetAsString return Value as string if it contains that.
func (v XInputValue) GetAsString() string {
	return v.base().GetAsString()
}

// ToScript returns the script of this element represent
func (v XInputValue) ToScript(depth int, e Evaler) (string, error) {
	if v.Type.IsField() || v.Field != "" { //field refer
		return FieldRef(v.Field).ToScript(depth, e)
	}

	if p, ok := v.Data.D.(ScriptElem); ok {
		return p.ToScript(depth, e)
	}
	return XValTypeUndefined.String(), nil
}

// DenyFieldRefer assert this value don't refer a field value
func (v XInputValue) DenyFieldRefer() error {
	return v.base().DenyFieldRefer()
}

// CreateSampleData generate a sample JSON value
func (v XInputValue) CreateSampleData(val interface{}, titleFirst bool) value.JSONValue {
	return v.base().CreateSampleData(val, titleFirst)
}

//------------------------------------------------------------------------------

// ValueSet represents a set of value for input
type ValueSet expr.ValueSet

func (v *ValueSet) base() *expr.ValueSet { return (*expr.ValueSet)(v) }

// AddKV add or update a named string value in this object
func (v *ValueSet) AddKV(key, value string, kind, in Enum) error {
	return v.base().AddKV(key, value, expr.Enum(kind), expr.Enum(in))
}

// AddExKV add or update a named string value in this object without type check
func (v *ValueSet) AddExKV(key, value string, kind, in Enum) error {
	return v.base().AddExKV(key, value, expr.Enum(kind), expr.Enum(in))
}

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValueSet) DelayedJSONDecode() error {
	return v.base().DelayedJSONDecode()
}

// RequestArgs represents a request parament of a request
type RequestArgs = expr.RequestArgs

// GetAction return the action predef value
func (v ValueSet) GetAction() *expr.InputValue {
	return v.base().GetAction()
}

// PrepareRequest solve the input parameters from predefined values
func (v ValueSet) PrepareRequest(args *RequestArgs) error {
	return v.base().PrepareRequest(args)
}

// ToScript returns the script of this element represent
func (v ValueSet) ToScript(depth int, e Evaler) (string, error) {
	var buf = bytes.NewBuffer(make([]byte, 0, defaultBufLen))
	lineHead := GenLinehead(depth)
	buf.WriteString(fmt.Sprintf("%s{\n", ""))
	for i := 0; i < len(v); i++ {
		p := (*XInputValue)(&v[i])
		if !p.In.IsBody() && !p.In.IsQuery() {
			continue
		}
		s, err := p.ToScript(depth+1, e)
		if err != nil {
			return "", err
		}
		if Enum(p.Type).IsArrayElem() {
			buf.WriteString(s)
		} else {
			buf.WriteString(fmt.Sprintf("%s  \"%s\": %s,\n", lineHead, p.GetName(false), s))
		}
	}
	buf.WriteString(fmt.Sprintf("%s}", lineHead))
	return buf.String(), nil
}
