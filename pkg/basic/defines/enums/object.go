package enums

import (
	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
)

// Enum exports
type Enum = enumset.Enum

// object enums
var (
	// ObjectTypeEnum represents objects kinds
	ObjectTypeEnum    = enumset.New(nil)
	ObjectRaw         = ObjectTypeEnum.MustReg("raw")       // raw api
	ObjectPoly        = ObjectTypeEnum.MustReg("poly")      // poly api
	ObjectKey         = ObjectTypeEnum.MustReg("apikey")    // platform api key
	ObjectCustomerKey = ObjectTypeEnum.MustReg("userkey")   // customer api key
	ObjectNamespace   = ObjectTypeEnum.MustReg("namespace") // namespace
	ObjectService     = ObjectTypeEnum.MustReg("service")   // service
)
