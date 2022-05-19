package encoding

import (
	"github.com/quanxiang-cloud/polyapi/pkg/lib/encoding/xml"
)

//-------------------------

// ToXML encode a Go data to xml string.
func ToXML(v interface{}, pretty bool) (string, error) {
	enc := xml.Encoder{}
	b, err := enc.Encode(v, nil, pretty)
	return unsafeByteString(b), err
}

// MustToXML encode a Go data to xml string.
// It panic if error.
func MustToXML(v interface{}, pretty bool) string {
	s, err := ToXML(v, pretty)
	if err != nil {
		panic(err)
	}
	return s
}

// FromXML decode a xml string to universal Go data.
func FromXML(s string) (interface{}, error) {
	dec := xml.Decoder{}
	d, err := dec.Decode(unsafeStringBytes(s))
	return d, err
}

// MustFromXML decode a xml string to universal Go data.
// It panic if error.
func MustFromXML(s string) interface{} {
	d, err := FromXML(s)
	if err != nil {
		panic(err)
	}
	return d
}
