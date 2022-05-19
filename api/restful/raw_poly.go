package restful

import (
	"github.com/quanxiang-cloud/polyapi/internal/service"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
)

// RawPoly RawPoly
type RawPoly struct {
	rp service.RawPoly
}

// CreateRawPoly CreateRawPoly
func CreateRawPoly(config *config.Config) (*RawPoly, error) {
	rp, err := service.CreateRawPoly(config)
	if err != nil {
		return nil, err
	}
	return &RawPoly{
		rp: rp,
	}, nil
}
