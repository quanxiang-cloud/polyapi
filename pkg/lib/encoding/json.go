package encoding

import (
	"encoding/json"
)

//-------------------------

// ToJSON encode a Go data to json string
func ToJSON(v interface{}, pretty bool) (string, error) {
	var b []byte
	var err error
	if pretty {
		b, err = json.MarshalIndent(v, "", "  ")
	} else {
		b, err = json.Marshal(v)
	}

	return unsafeByteString(b), err
}

// MustToJSON encode a Go data to json string.
// It panic if error.
func MustToJSON(v interface{}, pretty bool) string {
	s, err := ToJSON(v, pretty)
	if err != nil {
		panic(err)
	}
	return s
}

// FromJSON decode a json string to universal Go data
func FromJSON(s string) (interface{}, error) {
	var d interface{}
	err := json.Unmarshal(unsafeStringBytes(s), &d) // decode to *interface{}
	return d, err
}

// MustFromJSON decode a json string to universal Go data.
// It panic if error.
func MustFromJSON(s string) interface{} {
	d, err := FromJSON(s)
	if err != nil {
		panic(err)
	}
	return d
}
