package exprx

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
)

const (
	// DefaultEncoding is default encoding when empty
	DefaultEncoding = consts.DefaultEncoding
	// PolyEncoding is encoding by poly api
	PolyEncoding = consts.PolyEncoding
)

// enum of enum types
var (
	// EnumTypesEnum represents enum of enum types
	EnumTypesEnum = newEnumSet(nil)
	EnumNode      = EnumTypesEnum.MustReg(consts.EnumNode)
	EnumValue     = EnumTypesEnum.MustReg(consts.EnumValue)
	EnumOper      = EnumTypesEnum.MustReg(consts.EnumOper)
	EnumCond      = EnumTypesEnum.MustReg(consts.EnumCond)
	EnumCmp       = EnumTypesEnum.MustReg(consts.EnumCmp)
	EnumIn        = EnumTypesEnum.MustReg(consts.EnumIn)
	EnumAuth      = EnumTypesEnum.MustReg(consts.EnumAuth)
)

// basic types
var (
	// ValTypeEnum input value enum set, basic value accept in Input node
	ValTypeEnum = newEnumSet(nil)
	// ValTypeNumber represent value like 123
	ValTypeNumber = ValTypeEnum.MustReg("number")
	// ValTypeString represent value like "xyz"
	ValTypeString = ValTypeEnum.MustReg("string")

	// XValTypeAction represent value like action parameter
	ValTypeAction = ValTypeEnum.MustReg("action")
	// XValTypeTimestamp represent value like timestamp parameter
	ValTypeTimestamp = ValTypeEnum.MustReg("timestamp")

	// ValTypeBoolean represent value like true
	ValTypeBoolean = ValTypeEnum.MustReg("boolean")
	// ValTypeObject represent value like {...}
	ValTypeObject = ValTypeEnum.MustReg("object")
	// ValTypeArray represent value like [...]
	ValTypeArray = ValTypeEnum.MustReg("array")

	// ValTypeArrayString represent value like "foo,bar" => ["foo","bar"]
	ValTypeArrayString = ValTypeEnum.MustReg("array_string")
)

// extend value type
var (
	// XValTypeEnum extend value enum set, include ValEnum
	XValTypeEnum = newEnumSet(ValTypeEnum)
	// XValTypeUndefined represent value like undefined
	XValTypeUndefined = XValTypeEnum.MustReg("undefined") // undefined
	// XValTypeNull represent value like null
	XValTypeNull = XValTypeEnum.MustReg("null")
	// XValTypeMergeObj represent value like merged object, {a,b} + {c,d} => {a,b,c,d}
	XValTypeMergeObj = XValTypeEnum.MustReg("mergeobj")
	// XValTypeFilter represent value like use on Object|Array field
	XValTypeFilter = XValTypeEnum.MustReg("filter")

	// ValTypeArrayElem represent value like [...] as xxx.n
	//ValTypeArrayElem = ValTypeEnum.MustReg("array_elem")

	// ValTypeArrayStringElem represent value like [...] as xxx.n
	XValTypeArrayStringElem = XValTypeEnum.MustReg("array_string_elem")

	// XValTypeField represent value like req.data.userId
	XValTypeField = XValTypeEnum.MustReg("field")
)

// expression types
var (
	// ExprTypeEnum represents expression type, single expression, eg: const,field
	ExprTypeEnum = newEnumSet(XValTypeEnum)
	// ExprTypeExpr represents expr like (x + y) * 2
	ExprTypeExpr = ExprTypeEnum.MustReg("expr")
	// ExprTypeCmp represents compare expression, eg: a lt b
	ExprTypeCmp = ExprTypeEnum.MustReg("exprcmp")
	// ExprTypeSel represents select expression, eg: cond ? yesVal : noVal
	ExprTypeSel = ExprTypeEnum.MustReg("exprsel")
	// ExprTypeFunc represents // func(...), function call
	ExprTypeFunc = ExprTypeEnum.MustReg("exprfunc")
	// ExprTypeGroup represents group expression, eg: (a + b), (x and y)
	ExprTypeGroup = ExprTypeEnum.MustReg("exprgroup")
	// ExprTypeDirectExpr represents direct expression string, eg: "(req1.a+1)*2"
	ExprTypeDirectExpr = ExprTypeEnum.MustReg("direct_expr")
)

// operator enum set
var (
	// OpEnum represents operator enum set
	OpEnum = newEnumSet(nil)
	// OpAdd represents operator +
	OpAdd = OpEnum.MustRegWithContent("add", "+", "")
	// OpSub represents operator -
	OpSub = OpEnum.MustRegWithContent("sub", "-", "")
	// OpMul represents operator *
	OpMul = OpEnum.MustRegWithContent("mul", "*", "")
	// OpDiv represents operator /
	OpDiv = OpEnum.MustRegWithContent("div", "/", "")
)

// condition enum set
var (
	// CondEnum represents condition enum set
	CondEnum = newEnumSet(nil)
	// CondAnd represents logic operator &&
	CondAnd = CondEnum.MustRegWithContent("and", "&&", "") // &&
	// CondOr represents logic operator ||
	CondOr = CondEnum.MustRegWithContent("or", "||", "")
	// CondNot represents logic operator !
	CondNot = CondEnum.MustRegWithContent("not", "!", "")
)

// compare enum set
var (
	// CmpEnum represents compare enum set
	CmpEnum = newEnumSet(nil)
	CmpLT   = CmpEnum.MustRegWithContent("lt", "<", "")
	CmpGT   = CmpEnum.MustRegWithContent("gt", ">", "")
	CmpLE   = CmpEnum.MustRegWithContent("le", "<=", "")
	CmpGE   = CmpEnum.MustRegWithContent("ge", ">=", "")
	CmpEQ   = CmpEnum.MustRegWithContent("eq", "==", "")
	CmpNE   = CmpEnum.MustRegWithContent("ne", "!=", "")
)

var (
	// EncodingEnum represents encoding format
	EncodingEnum = newEnumSet(nil)
	// EncodingJSON represents encoding JSON
	EncodingJSON = EncodingEnum.MustReg(consts.EncodingJSON)
	// EncodingXML represents encoding XML
	EncodingXML = EncodingEnum.MustReg(consts.EncodingXML)
	// EncodingYAML represents encoding YAML
	EncodingYAML = EncodingEnum.MustReg(consts.EncodingYAML)
)

var (
	// SchemaEnum represents API scheme
	SchemaEnum = newEnumSet(nil)
	// SchemaHTTP represents API scheme http
	SchemaHTTP = SchemaEnum.MustReg(consts.SchemaHTTP)
	// SchemaHTTPS represents API scheme https
	SchemaHTTPS = SchemaEnum.MustReg(consts.SchemaHTTPS)
	// SchemaRPC represents API scheme rpc
	// SchemaRPC = SchemaEnum.MustReg(consts.SchemaRPC)
)

// http API methods
var (
	// MethodEnum represents http API methods
	MethodEnum    = newEnumSet(nil)
	MethodGet     = MethodEnum.MustReg(consts.MethodGet)
	MethodPost    = MethodEnum.MustReg(consts.MethodPost)
	MethodPut     = MethodEnum.MustReg(consts.MethodPut)
	MethodDelete  = MethodEnum.MustReg(consts.MethodDelete)
	MethodOPTIONS = MethodEnum.MustReg(consts.MethodOptions)
	MethodHEAD    = MethodEnum.MustReg(consts.MethodHead)
	MethodTRACE   = MethodEnum.MustReg(consts.MethodTrace)
	MethodCONNECT = MethodEnum.MustReg(consts.MethodConnect)
)

// parameter type
var (
	// ParaTypeEnum represents parameter type
	ParaTypeEnum     = newEnumSet(nil)
	ParaTypeHeader   = ParaTypeEnum.MustReg(consts.ParaInHeader)
	ParaTypePath     = ParaTypeEnum.MustReg(consts.ParaInPath)
	ParaTypeBody     = ParaTypeEnum.MustReg(consts.ParaInBody)
	ParaTypeQuery    = ParaTypeEnum.MustReg(consts.ParaInQuery)
	ParaTypeFormData = ParaTypeEnum.MustReg(consts.ParaInFormData)
	// ParaTypeHide represents hide parameter like skey
	ParaTypeHide = ParaTypeEnum.MustReg("hide")
)

//------------------------------------------------------------------------------

// Enum exports
type Enum enumset.Enum

// String show enum as string
func (e Enum) String() string {
	return string(e)
}

type enumSet struct{ *enumset.EnumSet }

func newEnumSet(exists *enumSet) *enumSet {
	var ex *enumset.EnumSet
	if exists != nil {
		ex = exists.EnumSet
	}
	return &enumSet{enumset.New(ex)}
}

func (es *enumSet) MustReg(val string) Enum {
	return Enum(es.EnumSet.MustReg(val))
}

func (es *enumSet) MustRegWithContent(val string, op, title string) Enum {
	return Enum(es.EnumSet.MustRegWithContent(val, op, title))
}

// ConvertOp eq->== and->&& add->+
func ConvertOp(e string) (string, error) {
	op, _, ok := enumset.FindContent(e,
		CondEnum.EnumSet, OpEnum.EnumSet, CmpEnum.EnumSet)
	if !ok {
		return e, errcode.ErrBuildInvalidOperator.FmtError(e)
	}
	return op, nil
}

//------------------------------------------------------------------------------

// IsHeader judege if is a header parameter
func (e Enum) IsHeader() bool {
	return e == ParaTypeHeader
}

// IsHide judege if is a hide parameter
func (e Enum) IsHide() bool {
	return e == ParaTypeHide
}

// IsHeaderAcceptable judege if is a header-acceptable parameter
func (e Enum) IsHeaderAcceptable() bool {
	switch e {
	case ValTypeString:
		return true
	}
	return false
}

// IsPath judege if is a path parameter
func (e Enum) IsPath() bool {
	return e == ParaTypePath
}

// IsQuery judege if is a query parameter
func (e Enum) IsQuery() bool {
	return e == ParaTypeQuery
}

// IsBody judege if is a body parameter
func (e Enum) IsBody() bool {
	return e == ParaTypeBody || e == "" // empty as body parameter
}

// IsArrayElem judge if a value type is array show as ary.1=xxx
func (e Enum) IsArrayElem() bool {
	return /*e == ValTypeArrayElem ||*/ e == XValTypeArrayStringElem
}

// IsAction judge if a value type is action
func (e Enum) IsAction() bool {
	return e == ValTypeAction
}

// IsTimestamp judge if a value type is timestamp
func (e Enum) IsTimestamp() bool {
	return e == ValTypeTimestamp
}

// IsField judge if a value type is field refer
func (e Enum) IsField() bool {
	return e == XValTypeField
}

// IsStringer judge if a value is able to convert to string
func (e Enum) IsStringer() bool {
	switch e {
	case ValTypeNumber, ValTypeString, ValTypeBoolean, ValTypeArrayString,
		ValTypeAction, ValTypeTimestamp, XValTypeField:
		return true
	}
	return false
}

// IsPredefineable judge if a value is predefineable
func (e Enum) IsPredefineable() bool {
	switch e {
	case ValTypeNumber, ValTypeString, ValTypeBoolean, ValTypeArrayString,
		XValTypeArrayStringElem, ValTypeAction, ValTypeTimestamp:
		return true
	}
	return false
}

// IsNullable judge if the data field is nullable
func (e Enum) IsNullable() bool {
	switch e {
	case ValTypeAction, ValTypeTimestamp:
		return true
	}
	return false
}

// SameIn compare if the parameter from the same input way
func (e Enum) SameIn(o Enum) bool {
	return e == o || (e.IsBody() && o.IsBody())
}

// EncodingToMIME change encoding like JSON to MIME like "application/json"
func (e Enum) EncodingToMIME() (string, error) {
	return consts.ToMIME(e.String())
}
