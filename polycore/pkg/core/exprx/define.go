package exprx

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/core/value"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/factory"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/jsonx"
)

const (

	// DefaultBufLen is default buffer length
	DefaultBufLen = defaultBufLen
	defaultBufLen = 64

	polyAllDataVarName = consts.PolyAllDataVarName // all data of a poly api
)

// NamedType exports
type NamedType = factory.NamedType

// FlexJSONObject exports
type FlexJSONObject = jsonx.FlexJSONObject

// DelayedJSONDecoder define an object that need delay decode JSON
type DelayedJSONDecoder interface {
	DelayedJSONDecode() error //  delay unmarshal flex json object
}

// NamedValue define a value with name
type NamedValue interface {
	GetName(titleFirst bool) string // get name of this element
}

// ScriptElem define a script element
type ScriptElem interface {
	ToScript(depth int, e Evaler) (string, error) // element to script
}

func first(s string, err error) string {
	return s
}

// NamedScriptElem define a script element with name
type NamedScriptElem interface {
	ScriptElem
	NamedValue
}

// Stringer define an interface with String()
type Stringer interface {
	String() string
	SetString(s string)
}

// StringerWithError define an interface with SetStringWithError()
type StringerWithError interface {
	Validate() error
	SetStringWithError(s string) error
}

// Stringifier is an interface to convert to string only
type Stringifier interface {
	String() string
}

// GenSampler represents value that can generate sample JSON value
type GenSampler interface {
	GenSample(val interface{}, titleFirst bool) value.JSONValue
}

// CreateSampleDataor represents value that can generate sample JSON value
type CreateSampleDataor interface {
	CreateSampleData(val interface{}, titleFirst bool) value.JSONValue
}
