package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/polyapi/internal/service"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"
)

// NewAPIOperator new
func NewAPIOperator(cfg *config.Config) (*APIOperator, error) {
	op, err := service.NewAPIOperator(cfg)
	if err != nil {
		return nil, err
	}
	return &APIOperator{
		op:   op,
		gate: gate.NewGate(cfg),
	}, nil
}

// APIOperator APIOperator
type APIOperator struct {
	op   service.APIOperator
	gate *gate.Gate
}

// QuerySwagger QuerySwagger
func (p *APIOperator) QuerySwagger(c *gin.Context) {
	if err := p.gate.Filt(c, APIRead); err != nil {
		return
	}

	req := &service.QueryAPISwaggerReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(p.op.QuerySwaggerInBatches(c, req)).Context(c)
}
