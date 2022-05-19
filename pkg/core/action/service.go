package action

import (
	"context"
	"errors"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
)

var errNotFound = errors.New("not found")

// ServicesResp exports
type ServicesResp = adaptor.ServicesResp

// QueryService query service from DB
func QueryService(service string) (*adaptor.ServicesResp, error) {
	if op := adaptor.GetServiceOper(); op != nil {
		return op.Query(context.Background(), service)
	}
	return nil, errNotFound
}

// GetServiceInfo make URL and authType from service
func GetServiceInfo(svs *adaptor.ServicesResp, apiPath string) (url, authType string) {
	return apipath.MakeRequestURL(svs.Schema, svs.Host, apiPath), svs.AuthType
}
