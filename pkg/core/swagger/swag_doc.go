package swagger

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/polysign"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
)

// SwagVersion is current swagger version
const SwagVersion = "2.0"

var errDefineNotFound = errors.New("define not found")

var inputBodyPolyReservedValue = expr.Value{
	Appendix: true, //NOTE: appendix value, required only in platform request
	Name:     polysign.XPolyBodyHideArgs,
	Type:     expr.ValTypeObject,
	Title:    "隐藏参数",
	Desc:     "polyapi reserved hide args like path args in raw api.",
	Data: expr.FlexJSONObject{
		D: &expr.ValObject{},
	},
}

var inputBodyPolySignatureValue = expr.Value{
	Appendix: true, //NOTE: appendix value, required only in platform request
	Name:     polysign.XBodyPolySignSignature,
	Type:     expr.ValTypeString,
	Title:    "参数签名",
	Desc: `required if Access-Token doesn't use.
HmacSHA256 signature of input body: sort query gonic asc|sha256 <SECRET_KEY>|base64 std encode`,
	Data: expr.FlexJSONObject{
		D: expr.NewStringer("EJML8aQ3BkbciPwMYHlffv2BagW0kdoI3L_qOedQylw", expr.ValTypeString),
	},
}

var inputSignatureDefine = []expr.ValueDefine{
	expr.ValueDefine{
		InputValue: expr.InputValue{
			Appendix: true, //NOTE: appendix value, required only in platform request
			Name:     polysign.XHeaderPolySignKeyID,
			In:       expr.ParaTypeHeader,
			Type:     expr.ValTypeString,
			Title:    "签名密钥序号",
			Desc:     "access_key_id dispatched by poly api server",
			Data: expr.FlexJSONObject{
				D: expr.NewStringer("KeiIY8098435rty", expr.ValTypeString),
			},
			Required: false,
		},
		Mock: "KeiIY8098435rty",
	},
	expr.ValueDefine{
		InputValue: expr.InputValue{
			Appendix: true, //NOTE: appendix value, required only in platform request
			Name:     polysign.XHeaderPolySignTimestamp,
			In:       expr.ParaTypeHeader,
			Type:     expr.ValTypeString,
			Title:    "签名时间戳",
			Desc:     "timestamp format ISO8601: 2006-01-02T15:04:05-0700",
			Data: expr.FlexJSONObject{
				D: expr.NewStringer("2020-12-31T12:34:56+0800", expr.ValTypeString),
			},
			Required: false,
		},
		Mock: "2020-12-31T12:34:56+0800",
	},
	expr.ValueDefine{
		InputValue: expr.InputValue{
			Appendix: true, //NOTE: appendix value, required only in platform request
			Name:     polysign.XHeaderPolySignVersion,
			In:       expr.ParaTypeHeader,
			Type:     expr.ValTypeString,
			Title:    "签名版本",
			Desc:     fmt.Sprintf("%q only current", polysign.XHeaderPolySignVersionVal),
			Data: expr.FlexJSONObject{
				D: expr.NewStringer(polysign.XHeaderPolySignVersionVal, expr.ValTypeString),
			},
			Required: false,
		},
		Mock: polysign.XHeaderPolySignVersionVal,
	},
	expr.ValueDefine{
		InputValue: expr.InputValue{
			Appendix: true, //NOTE: appendix value, required only in platform request
			Name:     polysign.XHeaderPolySignMethod,
			In:       expr.ParaTypeHeader,
			Type:     expr.ValTypeString,
			Title:    "签名方法",
			Desc:     fmt.Sprintf("%q only current", polysign.XHeaderPolySignMethodVal),
			Data: expr.FlexJSONObject{
				D: expr.NewStringer(polysign.XHeaderPolySignMethodVal, expr.ValTypeString),
			},
			Required: false,
		},

		Mock: polysign.XHeaderPolySignMethodVal,
	},

	expr.ValueDefine{
		InputValue: expr.InputValue{
			Appendix: true, //NOTE: appendix value, required only in platform request
			Name:     "Access-Token",
			In:       expr.ParaTypeHeader,
			Type:     expr.ValTypeString,
			Title:    "登录授权码",
			Desc:     `Access-Token from oauth2 if use token access mode`,
			Required: false,
		},
		Mock: "H3K56789lHIUkjfkslds",
	},
	// BUG: https://home.yunify.com/distributor.action?serviceName=clogin
	// response 500 with "Content-Type" header
	//
	// expr.ValueDefine{
	// 	InputValue: expr.InputValue{
	// 		Name:  "Content-Type",
	// 		In:    expr.ParaTypeHeader,
	// 		Type:  expr.ValTypeString,
	// 		Title: "数据格式",
	// 		Desc:  `application/json`,
	// 		Data: expr.FlexJSONObject{
	// 			D: expr.NewStringer("application/json", expr.ValTypeString),
	// 		},
	// 		Required: true,
	// 	},
	// 	Mock: "application/json",
	// },
}

// SwagDoc is the top swagger structure
type SwagDoc struct {
	Consts expr.SwagConstValueSet `json:"x-consts"` // extended consts define
	//Auth                adaptor.APIAuth      `json:"x-auth"`
	Defines             json.RawMessage `json:"definitions,omitempty"`
	SecurityDefinitions json.RawMessage `json:"securityDefinitions,omitempty"`

	Host         string                 `json:"host"` // host
	Version      string                 `json:"swagger"`
	Info         SwagInfo               `json:"info"`
	Tags         []SwagTag              `json:"tags,omitempty"`
	Schemes      []string               `json:"schemes"`
	BasePath     string                 `json:"basePath"`
	EncodingsIn  []string               `json:"consumes,omitempty"`
	EncodingsOut []string               `json:"produces,omitempty"`
	Paths        map[string]SwagMethods `json:"paths"` // path -> methods
}

// SwagMethods is the method map
type SwagMethods map[string]*SwagAPI // method -> api

// SwagAPI is the api specific
type SwagAPI struct {
	Consts       expr.SwagConstValueSet `json:"x-consts"` // extended consts define
	Name         string                 `json:"operationId"`
	Tags         []string               `json:"tags,omitempty"`
	Parameters   json.RawMessage        `json:"parameters"`
	Responses    json.RawMessage        `json:"responses"`
	EncodingsIn  []string               `json:"consumes"`
	EncodingsOut []string               `json:"produces"`
	Summary      string                 `json:"summary"`
	Desc         string                 `json:"description"`
	Deprecated   bool                   `json:"deprecated,omitempty"`
	Security     json.RawMessage        `json:"security,omitempty"`
}

// SwagTag is the swag tag
type SwagTag struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

// SwagContact is the contact info of this swag
type SwagContact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

// SwagInfo is the infomation of swagger
type SwagInfo struct {
	Title   string      `json:"title"`
	Version string      `json:"version"`
	Desc    string      `json:"description"`
	Contact SwagContact `json:"contact"`
}
