package apiprovider

import (
	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
)

func init() {
	enumset.FinishReg()
}

// doc type define
var (
	DocTypeEnum       = enumset.New(nil)
	DocTypeRaw        = DocTypeEnum.MustReg("raw")
	DocTypeSwag       = DocTypeEnum.MustReg("swag")
	DocTypeCurl       = DocTypeEnum.MustReg("curl")
	DocTypeJavascript = DocTypeEnum.MustReg("javascript")
	DocTypePython     = DocTypeEnum.MustReg("python")
)

// APIDoc is api doc structure
type APIDoc struct {
	URL    string      `json:"url"`
	Method string      `json:"method"`
	Desc   string      `json:"desc"`
	Input  interface{} `json:"input"`
	Output interface{} `json:"output"`
}
