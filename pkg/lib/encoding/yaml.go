package encoding

import (
	yaml "gopkg.in/yaml.v2"
)

//-------------------------

// ToYAML encode a Go data to yaml string.
func ToYAML(v interface{}, pretty bool) (string, error) {
	b, err := yaml.Marshal(v)
	return unsafeByteString(b), err
}

// MustToYAML encode a Go data to yaml string.
// It panic if error.
func MustToYAML(v interface{}, pretty bool) string {
	s, err := ToYAML(v, pretty)
	if err != nil {
		panic(err)
	}
	return s
}

// FromYAML decode a yaml string to universal Go data.
func FromYAML(s string) (interface{}, error) {
	var d interface{}
	err := yaml.Unmarshal(unsafeStringBytes(s), &d) // decode to *interface{}
	return d, err
}

// MustFromYAML decode a yaml string to universal Go data.
// It panic if error.
func MustFromYAML(s string) interface{} {
	d, err := FromYAML(s)
	if err != nil {
		panic(err)
	}
	return d
}
