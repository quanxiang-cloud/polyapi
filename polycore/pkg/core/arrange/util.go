package arrange

import (
	"reflect"
	"regexp"
	"unsafe"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
)

var regexNodeNameRule = regexp.MustCompile(`(?sm:^[a-zA-Z_][\w]*$)`)

const maxNodeNameLen = 32

// ValidateNodeName verify node name format
func ValidateNodeName(name string) error {
	if name == "" {
		return errcode.ErrBuildNodeName.FmtError(name)
	}
	if len(name) > maxNodeNameLen {
		return errcode.ErrNameTooLong.FmtError(name, maxNodeNameLen)
	}

	if !regexNodeNameRule.MatchString(name) {
		return errcode.ErrBuildNodeName.FmtError(name)
	}
	return nil
}

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

// convert filed name to a valid variable name
func toValidVarName(fieldName string) string {
	if fieldName == "" {
		return "_"
	}
	var s []rune = []rune(fieldName)
	var t = make([]rune, len(s))
	for i, v := range s {
		if isValidVarCh(v, i) {
			t[i] = v
		} else {
			t[i] = '_'
		}
	}
	return string(t)
}

// check if filed name is a valid variable name
func isValidVarName(fieldName string) bool {
	if fieldName == "" {
		return false
	}
	for i, v := range fieldName {
		if !isValidVarCh(v, i) {
			return false
		}
	}
	return true
}

// check if ch is a valid var rune
func isValidVarCh(ch rune, idx int) bool {
	if ch >= '0' && ch <= '9' {
		return idx > 0
	}
	if ch == '_' {
		return true
	}
	if ch >= 'a' && ch <= 'z' ||
		ch >= 'A' && ch <= 'Z' {
		return true
	}
	return false
}

// \n => \\n
// \r => \\r
var expToJsString = regexp.MustCompile(`(?sm:\r|\n)`)

// toJsString convert a Go string to JS string
func toJsString(s string) string {
	// \n => \\n
	// \r => \\r
	ss := expToJsString.ReplaceAllStringFunc(s, func(src string) string {
		switch src {
		case "\n":
			return "\\n"
		case "\r":
			return "\\r"
		}
		return ""
	})

	return ss
}
