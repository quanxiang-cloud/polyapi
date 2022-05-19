package expr

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/core/value"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/factory"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/jsonx"
)

const (
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

func first(s string, err error) string {
	return s
}

// Stringer define an interface with String()
type Stringer interface {
	String() string
	SetString(s string)
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
