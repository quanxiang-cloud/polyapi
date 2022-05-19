package consts

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	// Debug debug mode
	Debug = false
)

// some predefined var name in JS script
const (
	// PolyAPIInputVarName is function input object from global var
	PolyAPIInputVarName    = "__input"
	PolyRequestTempVar     = "_t" // temp var to call request
	PolyDummyName          = "_"  // DummyName
	PolyAllDataVarName     = "d"  // PolyAllDataVarName is a global object in JS script
	PolyFieldAccessAllData = "$"  // access root of all data field (d)
	PolyIfNodeResultVar    = "y"  // ifNode.y

	ScriptValueUndefined = "undefined" // js.undefined

	APIRequestPath = "/api/v1/polyapi/request/%s"
)

// name of operator system
const (
	SystemName           = "system"
	SystemTitle          = "系统"
	PathArgNamespacePath = "namespacePath" // api/v1/xxx/:fullNamespace
)

const (
	// SampleIndexNormal is normal index in sample api doc
	SampleIndexNormal = 0
	// SampleIndexTitle is tilte first index in sample api doc
	SampleIndexTitle = 1
)

const (
	// MaxNameLength is max name length in db
	MaxNameLength = 64
)

// encodings
const (
	// EncodingJSON is josn encoding
	EncodingJSON = "json"
	// EncodingXML is xml encoding
	EncodingXML = "xml"
	// EncodingYAML is yaml encoding
	EncodingYAML = "yaml"
	// DefaultEncoding is default encoding if missing
	DefaultEncoding = EncodingJSON
	// PolyEncoding is poly api using encoding
	PolyEncoding = EncodingJSON

	// PolyMethod is poly api default method
	PolyDefaultMethod = "POST"
)

const (
	// SchemaHTTP is http scheme
	SchemaHTTP = "http"
	// SchemaHTTPS is https scheme
	SchemaHTTPS = "https"
	// SchemaRPC is rpc scheme
	SchemaRPC = "[rpc]"
)

// where parameter from
const (
	ParaInHeader   = "header"
	ParaInPath     = "path"
	ParaInBody     = "body"
	ParaInQuery    = "query"
	ParaInFormData = "formData"
)

// Predef JS functions
const (
	PDFromJSON   = "pdFromJson"
	PDToJSON     = "pdToJson"
	PDToJSONP    = "pdToJsonP"
	PDMergeObjs  = "pdMergeObjs"
	PDFiltObject = "pdFiltObject"
	PDSelect     = "pdSelect"
	PDToJsobj    = "pdToJsobj"

	PDHttpRequest        = "pdHttpRequest"
	PDNewHTTPHeader      = "pdNewHttpHeader"
	PDUpdateReferPath    = "pdUpdateReferPath"
	PDAddHTTPHeader      = "pdAddHttpHeader"
	PDFromXML            = "pdFromXml"
	PDToXML              = "pdToXml"
	PDFromYAML           = "pdFromYaml"
	PDToYAML             = "pdToYaml"
	PDQueryUser          = "pdQueryUser"
	PDCreateNS           = "pdCreateNS"
	PDAppendAuth         = "pdAppendAuth"
	PDQueryGrantedAPIKey = "pdQueryGrantedAPIKey"

	PDtimestamp = "timestamp"
	PDformat    = "format"
	PDsel       = "sel"
)

// enum types
const (
	EnumNode  = "node"
	EnumValue = "value"
	EnumOper  = "oper"
	EnumCond  = "cond"
	EnumCmp   = "cmp"
	EnumIn    = "in"
	EnumAuth  = "auth"
)

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
)

// http methods
const (
	MethodGet     = http.MethodGet
	MethodHead    = http.MethodHead
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodPatch   = http.MethodPatch
	MethodDelete  = http.MethodDelete
	MethodConnect = http.MethodConnect
	MethodOptions = http.MethodOptions
	MethodTrace   = http.MethodTrace
)

// FromMIME converts MIME to encoding string
func FromMIME(mime string) (string, error) {
	switch mime {
	case MIMEJSON: // "application/xml"
		return EncodingJSON, nil
	case MIMEXML, MIMEXML2: // "text/xml", "application/xml"
		return EncodingXML, nil
	case MIMEYAML: //"application/x-yaml"
		return EncodingYAML, nil
	}
	return "", fmt.Errorf("unsupported MIME: %s", mime)
}

// ToMIME converts encoding string to MIME
func ToMIME(e string) (string, error) {
	switch e {
	case EncodingJSON, "":
		return MIMEJSON, nil // "application/json"
	case EncodingXML:
		return MIMEXML, nil // "application/xml"
	case EncodingYAML:
		return MIMEYAML, nil // "application/x-yaml"
	}
	return "", fmt.Errorf("unsupported encoding: %s", e)
}

// IsDefaultEncoding check if encoding is default
func IsDefaultEncoding(encoding string) bool {
	enc := strings.ToLower(encoding)
	return enc == DefaultEncoding
}

// IsDefaultSchema check if schema is default
func IsDefaultSchema(schema string) bool {
	switch schema {
	case SchemaHTTP, SchemaHTTPS:
		return true
	}
	return false
}

// ValidSchema check if schema is valid
func ValidSchema(schema string) (string, error) {
	s := strings.ToLower(schema)
	switch s {
	case SchemaHTTP, SchemaHTTPS:
		return s, nil
	}
	return "", fmt.Errorf("unsupported schema: %s", schema)
}
