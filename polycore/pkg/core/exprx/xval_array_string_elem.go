package exprx

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/core/value"
)

// ValArrayStringElem enable parse a single string as string array that output as ary.1=xxx
type ValArrayStringElem struct {
	Name  string         `json:"name"`
	Array ValArrayString `json:"array"`
}

// TypeName returns name of the type
func (v ValArrayStringElem) TypeName() string { return XValTypeArrayStringElem.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValArrayStringElem) DelayedJSONDecode() error {
	if v.Name == "" {
		return errors.New("value array_string_elem missing name")
	}
	return nil
}

// ToScript returns the script of this element represent
func (v ValArrayStringElem) ToScript(depth int, e Evaler) (string, error) {
	var buf = bytes.NewBuffer(make([]byte, 0, defaultBufLen))

	lineHead := GenLinehead(depth)
	for i := 0; i < len(v.Array); i++ {
		str := v.Array[i]
		buf.WriteString(fmt.Sprintf("%s\"%s.%d\": \"%s\",\n", lineHead, v.Name, i+1, str))
	}
	return buf.String(), nil
}

// GetName returns Name of the value.
func (v ValArrayStringElem) GetName(titleFirst bool) string {
	return v.Name
}

// GenSample generate a sample JSON value
func (v ValArrayStringElem) GenSample(val interface{}, titleFirst bool) value.JSONValue {
	var d value.Array
	for _, f := range v.Array {
		g := ValString(f)
		d.AddElement("", g.GenSample(nil, titleFirst))
	}
	return &d
}
