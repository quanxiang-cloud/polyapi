package expr

import (
	"bytes"

	"github.com/quanxiang-cloud/polyapi/pkg/core/value"
)

const sepElement = ',' // array element spliter
const escapeChar = '?' // ?, means ,
const defaultBufLen = 32

// ValArrayString enable parse a single string as string array
type ValArrayString []string

// TypeName returns name of the type
func (v ValArrayString) TypeName() string {
	return ValTypeArrayString.String()
}

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValArrayString) DelayedJSONDecode() error {
	return nil
}

// String encode ValArrayString as string
func (v ValArrayString) String() string {
	return v.encodeToString()
}

// SetString decode string as ValArrayString
func (v *ValArrayString) SetString(s string) {
	*v = ValArrayString(v.decode(unsafeStringBytes(s)))
}

// GenSample generate a sample JSON value
func (v ValArrayString) GenSample(val interface{}, titleFirst bool) value.JSONValue {
	var d value.Array
	for _, f := range v {
		g := ValString(f)
		d.AddElement("", g.GenSample(nil, titleFirst))
	}
	return &d
}

// , -> ?,
func (v ValArrayString) encode(s string, buf *bytes.Buffer) {
	size := len(s)
	for i := 0; i < size; i++ {
		switch c := s[i]; c {
		case sepElement:
			buf.WriteString(`?,`)
		default:
			buf.WriteByte(c)
		}
	}
}

func (v ValArrayString) encodeToString() string {
	var buf = bytes.NewBuffer(make([]byte, 0, defaultBufLen))
	buf.WriteByte('"')
	for i, e := range v {
		if i > 0 {
			buf.WriteByte(sepElement)
		}
		v.encode(e, buf)
	}
	buf.WriteByte('"')
	return buf.String()
}

// ?, -> ,
func (v ValArrayString) decode(b []byte) []string {
	//ignore ""
	if len(b) >= 2 && b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	size := len(b)
	var buf = bytes.NewBuffer(make([]byte, 0, defaultBufLen))
	var ret []string
	for i := 0; i < size; i++ {
		c := b[i]
		switch c {
		case escapeChar:
			if i < size-1 {
				if d := b[i+1]; d == sepElement {
					buf.WriteByte(d)
					i++
					continue
				}
			}
		case sepElement:
			ret = append(ret, buf.String())
			buf.Reset()
			continue
		}
		buf.WriteByte(c)
	}
	if buf.Len() > 0 {
		ret = append(ret, buf.String())
		buf.Reset()
	}
	return ret
}

// MarshalJSON encoding array as single string
func (v ValArrayString) MarshalJSON() ([]byte, error) {
	var buf = bytes.NewBuffer(make([]byte, 0, defaultBufLen))
	buf.WriteByte('[')
	for i, e := range v {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('"')
		buf.WriteString(e)
		buf.WriteByte('"')

	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

// UnmarshalJSON splist data as string array
func (v *ValArrayString) UnmarshalJSON(data []byte) error {
	*v = v.decode(data)
	return nil
}
