package enums

import (
	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/permission"
)

var (
	// GrantTypeEnum represents permit grant targets
	GrantTypeEnum = enumset.New(nil)
	// GrantUser represents grant target user
	GrantUser = GrantTypeEnum.MustReg("user")
	// GrantKey represents grant target key
	GrantKey = GrantTypeEnum.MustReg("key")
	// GrantAPP represents grant target app
	GrantAPP = GrantTypeEnum.MustReg("app")
	// GrantService represents grant target service
	GrantService = GrantTypeEnum.MustReg("service")
)

var (
	// PermitTypeEnum represents permit objects
	PermitTypeEnum = enumset.New(nil)
	// PermitRaw represents permit object raw API
	PermitRaw = PermitTypeEnum.MustReg("raw")
	// PermitPoly represents permit object poly API
	PermitPoly = PermitTypeEnum.MustReg("poly")
	// PermitCKey represents permit object customer key
	PermitCKey = PermitTypeEnum.MustReg("ckey")
	// PermitNamespace represents permit object namespace
	PermitNamespace = PermitTypeEnum.MustReg("ns")
	// PermitService represents permit object service
	PermitService = PermitTypeEnum.MustReg("service")
)

// permission bit
var (
	PermitEnum    = permission.PermitEnum
	PermitRead    = permission.PermitRead
	PermitExecute = permission.PermitExecute
	PermitCreate  = permission.PermitCreate
	PermitUpdate  = permission.PermitUpdate
	PermitDelete  = permission.PermitDelete
	PermitGrant   = permission.PermitGrant
)
