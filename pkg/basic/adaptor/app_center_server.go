package adaptor

import "context"

// AppCenterServerOper is the interface for app-center server proxy
type AppCenterServerOper interface {
	Check(c context.Context, userID, depID, appID string, isSuper, admin bool) (bool, error)
}

// SetAppCenterServerOper set the instance of app-center server oper
func SetAppCenterServerOper(c AppCenterServerOper) AppCenterServerOper {
	i := getInst()
	old := i.appCenterServerOper
	i.appCenterServerOper = c
	return old
}

// GetAppCenterServerOper get the instance of app-center server oper
func GetAppCenterServerOper() AppCenterServerOper {
	i := getInst()
	return i.appCenterServerOper
}
