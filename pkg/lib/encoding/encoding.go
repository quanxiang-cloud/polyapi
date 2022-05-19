package encoding

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
)

// MustFromEncoding parse specific encoding string to universal Go data.
// It panic if error occur.
func MustFromEncoding(encoding string, source string) interface{} {
	d, err := FromEncoding(encoding, source)
	if err != nil {
		panic(err)
	}
	return d
}

// MustToEncoding encode a universal Go data into specific encoding string.
// It panic if error occur.
func MustToEncoding(encoding string, d interface{}, pretty bool) string {
	s, err := ToEncoding(encoding, d, pretty)
	if err != nil {
		panic(err)
	}
	return s
}

// MustConvertEncoding converts source text from encoding to dstEncoding.
// It panic if error occur.
func MustConvertEncoding(encoding string, source string, dstEncoding string, pretty bool) string {
	s, err := ConvertEncoding(encoding, source, dstEncoding, pretty)
	if err != nil {
		panic(err)
	}
	return s
}

//------------------------------------------------------------------------------

// FromEncoding parse specific encoding string to universal Go data.
func FromEncoding(encoding string, source string) (interface{}, error) {
	switch enc := encoding; enc {
	case "", consts.EncodingJSON:
		return FromJSON(source)
	case consts.EncodingXML:
		return FromXML(source)
	case consts.EncodingYAML:
		return FromYAML(source)
	default:
		return nil, fmt.Errorf("encoding: unsupported encoding %s", enc)
	}
}

// ToEncoding encode a universal Go data into specific encoding string.
func ToEncoding(encoding string, d interface{}, pretty bool) (string, error) {
	switch enc := encoding; enc {
	case "", consts.EncodingJSON:
		return ToJSON(d, pretty)
	case consts.EncodingXML:
		return ToXML(d, pretty)
	case consts.EncodingYAML:
		return ToYAML(d, pretty)
	default:
		return "", fmt.Errorf("encoding: unsupported encoding %s", enc)
	}
}

// ConvertEncoding converts source text from encoding to dstEncoding.
func ConvertEncoding(encoding string, source string, dstEncoding string, pretty bool) (string, error) {
	if encoding == dstEncoding {
		return source, nil
	}
	obj, err := FromEncoding(encoding, source)
	if err != nil {
		return "", err
	}
	dst, err := ToEncoding(dstEncoding, obj, pretty)
	if err != nil {
		return "", err
	}
	return dst, err
}

//------------------------------------------------------------------------------

// unsafeByteString convert []byte to string without copy
// the origin []byte **MUST NOT** accessed after that
func unsafeByteString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// unsafeStringBytes return GoString's buffer slice
// ** NEVER modify returned []byte **
func unsafeStringBytes(s string) []byte {
	var bh reflect.SliceHeader
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return *(*[]byte)(unsafe.Pointer(&bh))
}
