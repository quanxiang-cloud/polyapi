package gate

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
)

var errAPIStatNotSetup = errors.New("feature api stat is not setup")

// APIStatAddTimeStat stat time of api and add too slow api to block
func APIStatAddTimeStat(c context.Context, apiPath string, raw bool, statusCode int, dur time.Duration) error {
	if statusCode >= http.StatusInternalServerError { // server error
		dur = -1
	}
	if oper := adaptor.GetAPIStatOper(); oper != nil {
		return oper.IncTimeStat(c, apiPath, raw, dur)
	}
	return errAPIStatNotSetup
}

// APIStatIsBlocked check if an api is shortly blocked
func APIStatIsBlocked(c context.Context, apiPath string, raw bool) bool {
	if oper := adaptor.GetAPIStatOper(); oper != nil {
		return oper.IsBlocked(c, apiPath, raw)
	}
	return false
}
