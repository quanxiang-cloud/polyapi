package expr

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
)

// NewStringer create a string-like value
func NewStringer(v string, kind Enum) Stringer {
	if kind.IsStringer() {
		p, err := flexFactory.Create(kind.String())
		if err != nil {
			return nil
		}
		s := p.(Stringer)
		s.SetString(v)
		return s
	}

	return nil
}

// FullVarName get the full variable of poly API
func FullVarName(name string) string {
	if strings.HasPrefix(name, consts.PolyFieldAccessAllData) {
		return strings.Replace(name, consts.PolyFieldAccessAllData, polyAllDataVarName, 1)
	}
	return fmt.Sprintf("%s.%s", polyAllDataVarName, name)
}

// GetFieldName return the file name from field ref. eg: req1.data.x => x
func GetFieldName(f string) string {
	if f == consts.PolyFieldAccessAllData { // $
		return polyAllDataVarName // d
	}
	if i := strings.LastIndex(f, "."); i >= 0 { // req1.data.x => x
		return f[i+1:]
	}
	return f // return whole parts of f
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
