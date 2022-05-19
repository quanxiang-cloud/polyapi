package expr

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/polysign"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/timestamp"
	"github.com/quanxiang-cloud/polyapi/pkg/core/value"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/encoding"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/jsonx"
)

// SwagConstValue is predef value
type SwagConstValue = InputValue

// SwagConstValueSet is predef value set
type SwagConstValueSet = ValueSet

//------------------------------------------------------------------------------

// InputValue represent a const value input
type InputValue struct {
	//Value
	// type of this value:
	// number|string|boolean|object|array|arrayelem|null|undefined|field|
	// mergeobj|filter|expr|exprcmp|exprsel|exprgroup
	Type  Enum   `json:"type"`
	Name  string `json:"name,omitempty"` // new name of this value
	Title string `json:"title,omitempty"`
	Desc  string `json:"desc,omitempty"` // description of this value

	Appendix bool `json:"$appendix$,omitempty"` // NOTE: Appendix value, platform only
	Required bool `json:"required,omitempty"`   // required

	Field FieldRef       `json:"field,omitempty"` // field refer for this value, eg: "req1.data.x"
	Data  FlexJSONObject `json:"data,omitempty"`  // specific value for non-field content

	In Enum `json:"in"` // header|path|body
}

// Empty check if the value is empty
func (v *InputValue) Empty() bool {
	return v.Field.Empty() && v.Data.Empty()
}

// DelayedJSONDecode delay unmarshal flex json object
func (v *InputValue) DelayedJSONDecode() error {
	// don't decode Data if it contains a Field reference
	if v.Type.IsField() || v.Field != "" {
		if v.Field == "" { // from value
			v.Field = FieldRef(v.GetAsString())
			if v.Field == "" {
				return fmt.Errorf("missing field data of input value '%s'", v.GetName(false))
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

// GetName return name of input value
func (v InputValue) GetName(titleFirst bool) string {
	if titleFirst {
		if v.Title != "" {
			return v.Title
		}
	}

	// name has set, direct return
	if v.Name != "" {
		return v.Name
	}

	// cannot access the correct name, _
	return consts.PolyDummyName
}

// GetAsString return Value as string if it contains that.
func (v InputValue) GetAsString() string {
	if v.Type.IsStringer() {
		if v.Field != "" {
			return v.Field.String()
		}
		//BUG: don't use Stringer here, ValString is not Stringer(*ValString is).
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
func (v InputValue) DenyFieldRefer() error {
	if v.Type.IsField() || v.Field != "" {
		return fmt.Errorf("value(%s) dont't accept field(%s) refer", v.GetName(false), v.Field)
	}
	return nil
}

// CreateSampleData generate a sample JSON value
func (v InputValue) CreateSampleData(val interface{}, titleFirst bool) value.JSONValue {
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

//------------------------------------------------------------------------------

// ValueSet represents a set of value for input
type ValueSet []InputValue

// AddKV add or update a named string value in this object
func (v *ValueSet) AddKV(key, value string, kind, in Enum) error {
	return v.addKV(key, value, kind, in, true)
}

// AddExKV add or update a named string value in this object without type check
func (v *ValueSet) AddExKV(key, value string, kind, in Enum) error {
	return v.addKV(key, value, kind, in, false)
}

// AddKV add or update a named string value in this object
func (v *ValueSet) addKV(key, value string, kind, in Enum, typeCheck bool) error {
	if key == "" && value == "" { // ignore empty input
		return nil
	}
	if key == "" || (!kind.IsNullable() && value == "") { // error check
		return fmt.Errorf("missing key(%s) or value(%s)", key, value)
	}
	if typeCheck && !kind.IsPredefineable() {
		return fmt.Errorf("unsupported type %s for AddKV(%s,%s)", kind, key, value)
	}

	s := *v
	for i := 0; i < len(s); i++ {
		p := &s[i]
		if p.Name == key && p.In.SameIn(in) {
			v.setValue(p, value, kind)
			return nil
		}
	}

	n := InputValue{
		Name: key,
		Type: kind,
		In:   in,
	}
	v.setValue(&n, value, kind)
	*v = append(*v, n)
	return nil
}

func (v *ValueSet) setValue(p *InputValue, value string, kind Enum) {
	p.Type = kind
	if kind.IsField() {
		p.Field = FieldRef(value)
		p.Data.D = nil
	} else {
		p.Field = ""
		p.Data.D = NewStringer(value, kind)
	}
}

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValueSet) DelayedJSONDecode() error {
	for i := 0; i < len(*v); i++ {
		p := &(*v)[i]
		if err := p.DelayedJSONDecode(); err != nil {
			return err
		}
	}
	return nil
}

//------------------------------------------------------------------------------

// RequestArgs represents a request parament of a request
type RequestArgs struct {
	Name         string      // Name of the request node
	URL          string      // URL of a request
	EncodingPoly string      // Content-Type, of poly call
	EncodingIn   string      // Content-Type of request api, json/xml/yaml
	Body         []byte      // request body
	Header       http.Header // request header
	Action       string      // action of this request

	parsedBody map[string]interface{} // universal body data
}

func (args *RequestArgs) preparePathArgs() error {
	hideArg, ok := args.parsedBody[polysign.XPolyBodyHideArgs]
	if !ok {
		return nil
	}
	delete(args.parsedBody, polysign.XPolyBodyHideArgs) //remove polysign.XPolyBodyHideArgs args from body

	arg, ok := hideArg.(map[string]interface{})
	if !ok {
		return fmt.Errorf(" body.%s is not object", polysign.XPolyBodyHideArgs)
	}

	r, err := replacePathArgsUniversal(args.URL, args.Name, arg)
	if err != nil {
		return err
	}
	args.URL = r
	return nil
}

// GetAction return the action predef value
func (v ValueSet) GetAction() *SwagConstValue {
	for i := 0; i < len(v); i++ {
		p := &v[i]
		if p.Type.IsAction() {
			return p
		}
	}
	return nil
}

// PrepareRequest solve the input parameters from predefined values
func (v ValueSet) PrepareRequest(args *RequestArgs) error {
	d, err := encoding.FromEncoding(args.EncodingPoly, unsafeByteString(args.Body))
	if err != nil {
		return err
	}
	if p, ok := d.(map[string]interface{}); ok {
		args.parsedBody = p
	} else {
		return fmt.Errorf("PrepareRequest for %s error: body not object", args.Name)
	}

	if err := v.prepareArgs(args); err != nil {
		return err
	}
	if err := args.preparePathArgs(); err != nil { // NOTE: must after prepareArgs
		return err
	}

	body, err := encoding.ToEncoding(args.EncodingIn, v.getBodyRoot(args), true)
	if err != nil {
		return err
	}
	args.Body = unsafeStringBytes(body)

	return nil
}

// NOTE: query body with child $root$ firstly
func (v ValueSet) getBodyRoot(args *RequestArgs) interface{} {
	if root, ok := args.parsedBody[polysign.XPolyCustomerBodyRoot]; ok {
		return root
	}
	return args.parsedBody
}

func (v *ValueSet) prepareArgs(args *RequestArgs) error {
	for i := 0; i < len(*v); i++ {
		p := &(*v)[i]
		switch {
		case p.Type.IsAction():
			if args.Action != "" {
				p.Data.D = ValAction(args.Action)
			}

		case p.Type.IsTimestamp():
			p.Data.D = ValTimestamp(timestamp.Timestamp(p.GetAsString()))

		}

		switch {
		case p.In.IsBody() || p.In.IsQuery():
			args.parsedBody[p.Name] = p.Data.D
		case p.In.IsHeader():
			if v, ok := p.Data.D.(Stringifier); ok {
				args.Header.Add(p.Name, v.String())
			}
		}
	}

	return nil
}

//------------------------------------------------------------------------------
