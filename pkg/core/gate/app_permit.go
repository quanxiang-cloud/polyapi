package gate

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/business/app"
	"github.com/quanxiang-cloud/polyapi/pkg/config"

	"github.com/gin-gonic/gin"
)

type appPermit struct{}

func createAppPermit(cfg *config.Config) (*appPermit, error) {
	n := &appPermit{}
	return n, nil
}

func (p *appPermit) Handle(c *gin.Context, apiType apiType) (err error) {
	namespacePath := GetNamespacePath(c)
	if err := app.ValidateAppAccessPermit(namespacePath, c.Request.Header, apiType.IsWriter()); err != nil {
		return errcode.ErrInternal.FmtError("app-auth: " + err.Error())
	}
	return nil
}
