package docview

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/apiprovider"
	"github.com/quanxiang-cloud/polyapi/pkg/core/swagger"
)

// type exports
type (
	apiDocView = apiprovider.APIDoc
)

// func GetSampleIndex(titleFirst bool) int
var getSampleIndex = swagger.GetSampleIndex

// doc type exports
var (
	DocTypeEnum       = apiprovider.DocTypeEnum
	DocTypeRaw        = apiprovider.DocTypeRaw
	DocTypeSwag       = apiprovider.DocTypeSwag
	DocTypeCurl       = apiprovider.DocTypeCurl
	DocTypeJavascript = apiprovider.DocTypeJavascript
	DocTypePython     = apiprovider.DocTypePython
)
