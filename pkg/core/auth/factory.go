package auth

import (
	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/factory"
)

func init() {
	enumset.FinishReg()
}

// flexFactory regist some flex json objects
var flexFactory = factory.NewFlexObjFactory("flexObjCreator")

// GetFactory return the factory object
func GetFactory() *factory.FlexObjFactory { return flexFactory }

func init() {
	flexFactory.MustReg(authNone{})
	flexFactory.MustReg(authSystem{})
	flexFactory.MustReg(authSignature{
		Cmds: "sort query gonic asc|append begin GET\n/saas/\n|sha256 <SECRET_KEY>|base64 std encode",
	})
}
