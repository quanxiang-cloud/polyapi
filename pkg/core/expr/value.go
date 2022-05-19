package expr

import (
	"encoding/json"
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/core/value"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/jsonx"
)

//------------------------------------------------------------------------------

// ValArray represents an array value, [...]
type ValArray []Value

// TypeName returns name of the type
func (v ValArray) TypeName() string { return ValTypeArray.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValArray) DelayedJSONDecode() error {
	for i := 0; i < len(*v); i++ {
		p := &(*v)[i]
		if err := p.DelayedJSONDecode(); err != nil {
			return err
		}
	}
	return nil
}

// GenSample generate a sample JSON value
func (v ValArray) GenSample(val interface{}, titleFirst bool) value.JSONValue {
	x := &value.Array{}
	for _, f := range v {
		d := f.Data.D
		if d == nil {
			d, _ = flexFactory.Create(f.Type.String())
		}
		if g, ok := d.(GenSampler); ok {
			x.AddElement("", g.GenSample(nil, titleFirst))
		}
	}
	return x
}

//------------------------------------------------------------------------------

// ValObject represents an object value, {...}
type ValObject ValArray

// TypeName returns name of the type
func (v ValObject) TypeName() string { return ValTypeObject.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValObject) DelayedJSONDecode() error {
	return (*ValArray)(v).DelayedJSONDecode()
}

// GenSample generate a sample JSON value
func (v ValObject) GenSample(val interface{}, titleFirst bool) value.JSONValue {
	x := &value.Object{}
	for _, f := range v {
		d := f.Data.D
		if d == nil {
			d, _ = flexFactory.Create(f.Type.String())
		}

		if g, ok := d.(GenSampler); ok {
			x.AddElement(f.GetName(titleFirst), g.GenSample(nil, titleFirst))
		}
	}
	return x
}

//------------------------------------------------------------------------------

// Value represents a value with given type
// It use Field value firstly.
type Value struct {
	// type of this value:
	// number|string|boolean|object|array|array_elem|array_string_elem|array_string|
	// undefined|null|mergeobj|filter|path|header|skey|action|signature|timestamp|
	// field|expr|exprcmp|exprsel|exprfunc|exprgroup|direct_expr
	Type  Enum   `json:"type"`
	Name  string `json:"name"` // new name of this value
	Title string `json:"title,omitempty"`
	Desc  string `json:"desc,omitempty"` // description of this value

	Appendix bool `json:"$appendix$,omitempty"` // NOTE: Appendix value, platform only
	Required bool `json:"required,omitempty"`   // required

	Field FieldRef       `json:"field,omitempty"` // field refer for this value, eg: "req1.data.x"
	Data  FlexJSONObject `json:"data,omitempty"`  // specific value for non-field content
}

// DelayedJSONDecode delay unmarshal flex json object
func (v *Value) DelayedJSONDecode() error {
	// don't decode Data if it contains a Field reference
	if v.Type.IsField() || v.Field != "" {
		if v.Field == "" { // from value
			v.Field = FieldRef(v.GetAsString())
			if v.Field == "" {
				return fmt.Errorf("missing field data of value '%s'", v.GetName(false))
			}
		}
		v.Data.D = nil
		return nil
	}

	if err := flexFactory.DelayedUnmarshalFlexJSONObject(v.Type.String(), &v.Data); err != nil {
		return err
	}

	// Delay decode for children
	if c, ok := v.Data.D.(DelayedJSONDecoder); ok {
		if err := c.DelayedJSONDecode(); err != nil {
			return err
		}
	}
	return nil
}

// GetName returns Name of the value.
// It returns field name if it is a filed value and Nanme not set.
func (v Value) GetName(titleFirst bool) string {
	if titleFirst { // use title as name first
		if v.Title != "" {
			return v.Title
		}
	}

	// name has set, direct return
	if v.Name != "" {
		return v.Name
	}

	if !v.Field.Empty() {
		return v.Field.GetName(titleFirst)
	}

	// use refered data name
	if p, ok := v.Data.D.(NamedValue); ok {
		return p.GetName(titleFirst)
	}

	// cannot access the correct name, _
	return consts.PolyDummyName
}

// GetAsString return Value as string if it contains that.
func (v Value) GetAsString() string {
	if v.Type.IsStringer() {
		if v.Field != "" {
			return v.Field.String()
		}
		switch d := v.Data.D.(type) {
		case Stringifier:
			return d.String()
		case jsonx.RawMessage:
			var s string
			_ = json.Unmarshal(d, &s)
			return s
		}
	}
	return ""
}

// DenyFieldRefer assert this value don't refer a field value
func (v Value) DenyFieldRefer() error {
	if v.Field != "" {
		return fmt.Errorf("value(%s) don't accept field(%s) refer", v.GetName(false), v.Field)
	}
	return nil
}

// Empty check if the value is empty
func (v Value) Empty() bool {
	return v.Field.Empty() && v.Data.Empty()
}

// CreateSampleData generate a sample JSON value
func (v Value) CreateSampleData(val interface{}, titleFirst bool) value.JSONValue {
	d := v.Data.D
	if d == nil {
		if c, err := flexFactory.Create(v.Type.String()); err == nil {
			d = c
		}
	}

	if g, ok := d.(GenSampler); ok {
		return g.GenSample(val, titleFirst)
	}

	return nil
}
