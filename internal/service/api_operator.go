package service

import (
	"context"
	"encoding/json"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/polyhost"
	"github.com/quanxiang-cloud/polyapi/pkg/business/rule"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/swagger"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
)

// APIOperator APIOperator
type APIOperator interface {
	QuerySwaggerInBatches(ctx context.Context, req *QueryAPISwaggerReq) (*QueryAPISwaggerResp, error)
}

// NewAPIOperator create apiOperator
func NewAPIOperator(cfg *config.Config) (APIOperator, error) {
	return &apiOperator{
		cfg: cfg,
	}, nil
}

type apiOperator struct {
	cfg *config.Config
}

// QueryAPISwaggerReq QueryAPISwaggerReq
type QueryAPISwaggerReq struct {
	APIPath []string `json:"apiPath"`
}

// QueryAPISwaggerResp QueryAPISwaggerResp
type QueryAPISwaggerResp struct {
	Swag *swagger.SwagDoc `json:"swag"`
}

// QuerySwaggerInBatches query api swagger
func (p *apiOperator) QuerySwaggerInBatches(ctx context.Context, req *QueryAPISwaggerReq) (*QueryAPISwaggerResp, error) {
	if err := rule.CheckArrayLength(len(req.APIPath)); err != nil {
		return nil, err
	}

	raw := req.APIPath[:0]
	poly := make([]string, 0, len(req.APIPath)/2)
	for _, v := range req.APIPath {
		t := apipath.GetAPIType(v)
		switch t {
		case "r":
			raw = append(raw, v)
		case "p":
			poly = append(poly, v)
		default:
		}
	}

	var rawSwag = &swagger.SwagDoc{
		Paths: make(map[string]swagger.SwagMethods),
	}
	if op := adaptor.GetRawAPIOper(); op != nil && len(raw) > 0 {
		resp, err := op.QuerySwagger(ctx, &adaptor.QueryRawSwaggerReq{
			APIPath: raw,
		})
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(resp.Swagger, rawSwag)
		if err != nil {
			return nil, err
		}
	}

	if op := adaptor.GetPolyOper(); op != nil && len(poly) > 0 {
		resp, err := op.QuerySwagger(ctx, &adaptor.QueryPolySwaggerReq{
			APIPath: poly,
		})
		if err != nil {
			return nil, err
		}
		var polySwag = &swagger.SwagDoc{}
		err = json.Unmarshal(resp.Swagger, polySwag)
		if err != nil {
			return nil, err
		}

		for k, v := range polySwag.Paths {
			rawSwag.Paths[k] = v
		}
	}

	rawSwag.BasePath = "/api/v1/polyapi/request"
	rawSwag.Schemes = []string{polyhost.GetSchema()}
	rawSwag.Host = polyhost.GetHost()

	return &QueryAPISwaggerResp{
		Swag: rawSwag,
	}, nil
}
