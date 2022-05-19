package restful

import (
	"github.com/quanxiang-cloud/polyapi/internal/service"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

// PathOperator PathOperator
type PathOperator struct {
	pathOpt service.PathOperator
}

// NewOperator NewOperator
func NewOperator(config *config.Config) (*PathOperator, error) {
	recycler, err := service.CreateOperator(config)
	return &PathOperator{
		pathOpt: recycler,
	}, err
}

// InnerDelApp InnerDelApp
func (opt *PathOperator) InnerDelApp(c *gin.Context) {
	gate.SetFromInnerFlag(c, true)
	req := &service.DelAppReq{}
	var err error
	if !getPathArg(c, PathargAPIPath, &req.AppID, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}

	resp.Format(opt.pathOpt.DelApp(c, req)).Context(c)
}

// InnerUpdateAppValid update namespace raw and poly valid by app root path
func (opt *PathOperator) InnerUpdateAppValid(c *gin.Context) {
	gate.SetFromInnerFlag(c, true)
	req := &service.UpdateAppValidReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.AppID, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(opt.pathOpt.UpdateAppValid(c, req)).Context(c)
}

// InnerExportApp export app
func (opt *PathOperator) InnerExportApp(c *gin.Context) {
	gate.SetFromInnerFlag(c, true)
	req := &service.ExportAppReq{}
	var err error
	if !getPathArg(c, PathargAPIPath, &req.AppID, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(opt.pathOpt.ExportApp(c, req)).Context(c)
}

// InnerImport import data
func (opt *PathOperator) InnerImport(c *gin.Context) {
	gate.SetFromInnerFlag(c, true)
	req := &service.ImportReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(opt.pathOpt.Import(c, req)).Context(c)
}
