package adaptor

import (
	"context"
	"time"
)

// APIStatOper is the interface for api stat
type APIStatOper interface {
	IncTimeStat(c context.Context, apiPath string, raw bool, dur time.Duration) error
	//	GetStat(c context.Context, apiPath string, raw bool) (int64, error)
	IsBlocked(c context.Context, apiPath string, raw bool) bool
}

// SetAPIStatOper set the instance of api stat oper
func SetAPIStatOper(f APIStatOper) APIStatOper {
	i := getInst()
	old := i.apiStatOper
	i.apiStatOper = f
	return old
}

// GetAPIStatOper get the instance of api stat oper
func GetAPIStatOper() APIStatOper {
	i := getInst()
	if oper := i.apiStatOper; oper != nil {
		return i.apiStatOper
	}
	panic("got nil")
}
