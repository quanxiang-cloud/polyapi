package restful

import (
	"errors"

	"github.com/quanxiang-cloud/polyapi/internal/service"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

var errInternal = errors.New("internal error")

// PolyAPI is the route for poly api
type PolyAPI struct {
	svs  service.PolyAPI
	gate *gate.Gate
}

// NewPolyAPI create a poly api provider
func NewPolyAPI(conf *config.Config) (*PolyAPI, error) {
	poly, err := service.CreatePoly(conf)
	if err != nil {
		return nil, err
	}

	serv := &PolyAPI{
		svs:  poly,
		gate: gate.NewGate(conf),
	}

	return serv, nil
}

//------------------------------------------------------------------------------

// Create create a new poly api
func (s *PolyAPI) Create(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.PolyCreateReq{}
	req.Owner, req.OwnerName = s.gate.GetAuthOwnerPair(c)
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.Namespace, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}

	resp.Format(s.svs.Create(c, req)).Context(c)
}

//------------------------------------------------------------------------------

// Delete delete a poly api
func (s *PolyAPI) Delete(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.PolyDeleteReq{}
	req.Owner = s.gate.GetAuthOwner(c)
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}

	resp.Format(s.svs.Delete(c, req)).Context(c)
}

//------------------------------------------------------------------------------

// UpdateArrange update a poly api arrange json
func (s *PolyAPI) UpdateArrange(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.PolyUpdateArrangeReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.APIPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	req.Owner = s.gate.GetAuthOwner(c)

	resp.Format(s.svs.UpdateArrange(c, req)).Context(c)
}

//------------------------------------------------------------------------------

// GetArrange return a poly api arrange json
func (s *PolyAPI) GetArrange(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.PolyGetArrangeReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.APIPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}

	resp.Format(s.svs.GetArrange(c, req)).Context(c)
}

//------------------------------------------------------------------------------

// Build build a poly api arrange json to JS code
func (s *PolyAPI) Build(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.PolyBuildReq{}
	req.Owner = s.gate.GetAuthOwner(c)
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.APIPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}

	// don't hide build error
	resp.Format(s.svs.Build(c, req)).Context(c)
}

//------------------------------------------------------------------------------

// List list poly api in namespace
func (s *PolyAPI) List(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.PolyListReq{
		Active: -1,
	}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.svs.List(c, req)).Context(c)
}

// Active reverse active state of a poly api
func (s *PolyAPI) Active(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.PolyActiveReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.APIPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.svs.Active(c, req)).Context(c)
}

// Search search poly api by title or path, allow to find poly api at subdirectory
func (s *PolyAPI) Search(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.SearchPolyReq{
		Active: -1,
	}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.Namespace, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.svs.Search(c, req)).Context(c)
}

//------------------------------------------------------------------------------

// ShowEnum show enums of poly api
// @Summary list enums
// @Description list enums
// @Produce json
// @Param body body service.PolyEnumReq true "body parameters"
// @Success 200 {object} resp.R{data=service.PolyEnumResp} {} "return list of specify enum"
// @Router /api/v1/polyapi/poly/enums [post]
func (s *PolyAPI) ShowEnum(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.PolyEnumReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargType, &req.Type, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}

	resp.Format(s.svs.ShowEnum(c, req)).Context(c)
}
