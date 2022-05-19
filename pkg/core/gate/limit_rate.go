package gate

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate/limitrate"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
)

func createLimitRate(cfg *config.Config) (*limitRate, error) {
	if cfg.Gate.LimitRate.Enable {
		n := &limitRate{
			b: limitrate.NewBucket(cfg.Gate.LimitRate.RatePerSecond),
		}
		n.b.Start()
		logger.Logger.Infof("[gate.limitRrate]:%s", n.b.ShowTokenArg())
		return n, nil
	}
	logger.Logger.Infow("[gate.limitRrate] disabled")
	return nil, nil
}

type limitRate struct {
	b *limitrate.Bucket
}

func (v *limitRate) Handle(c *gin.Context, apiType apiType) error {
	if !v.b.GetToken() {
		return errcode.ErrSystemBusy.NewError()
	}
	return nil
}
