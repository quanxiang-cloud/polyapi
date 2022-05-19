package restful

import (
	"github.com/quanxiang-cloud/polyapi/internal/service"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

// APIPermit is the route for api permit
type APIPermit struct {
	s    service.APIPermit
	gate *gate.Gate
}

// NewAPIPermit create a api permit provider
func NewAPIPermit(conf *config.Config) (*APIPermit, error) {
	svs, err := service.CreateAPIPermit(conf)
	if err != nil {
		return nil, err
	}
	r := &APIPermit{
		gate: gate.NewGate(conf),
		s:    svs,
	}
	return r, nil
}

// CreateGroup create a api permit group
func (s *APIPermit) CreateGroup(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.AddPermitGroupReq{}
	req.Owner, req.OwnerName = s.gate.GetAuthOwnerPair(c)
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.Namespace, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.AddPermitGroup(c, req)).Context(c)
}

// DeleteGroup delete a api permit group
func (s *APIPermit) DeleteGroup(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.DelPermitGroupReq{}
	var err error
	if !getPathArg(c, PathargAPIPath, &req.GroupPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.DelPermitGroup(c, req)).Context(c)
}

// UpdateGroup update a api permit group
func (s *APIPermit) UpdateGroup(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.UpdatePermitGroupReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.GroupPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.UpdatePermitGroup(c, req)).Context(c)
}

// ActiveGroup reverse active state of a api permit group
func (s *APIPermit) ActiveGroup(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.ActivePermitGroupReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.GroupPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.ActivePermitGroup(c, req)).Context(c)
}

// ListGroup list children api permit group within a namespace
func (s *APIPermit) ListGroup(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.ListPermitGroupReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.Namespace, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.ListPermitGroup(c, req)).Context(c)
}

// QueryGroup query a api permit group info
func (s *APIPermit) QueryGroup(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.QueryPermitGroupReq{}
	var err error
	if !getPathArg(c, PathargAPIPath, &req.GroupPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.QueryPermitGroup(c, req)).Context(c)
}

// AddElem add a permit element to a permit group
func (s *APIPermit) AddElem(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.AddPermitElemReq{}
	req.Owner, req.OwnerName = s.gate.GetAuthOwnerPair(c)
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.GroupPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.AddPermitElem(c, req)).Context(c)
}

// DeleteElem delete a permit element to a permit group
func (s *APIPermit) DeleteElem(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.DelPermitElemReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.DelPermitElem(c, req)).Context(c)
}

// UpdateElem update a permit element to a permit group
func (s *APIPermit) UpdateElem(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.UpdatePermitElemReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.UpdatePermitElem(c, req)).Context(c)
}

// ActiveElem reverse active state of a api permit element
func (s *APIPermit) ActiveElem(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.ActivePermitElemReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.ActivePermitElem(c, req)).Context(c)
}

// ListElem list permit elements within a permit group
func (s *APIPermit) ListElem(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.ListPermitElemReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.GroupPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.ListPermitElem(c, req)).Context(c)
}

// QueryElem query a permit element within a permit group
func (s *APIPermit) QueryElem(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.QueryPermitElemReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.QueryPermitElem(c, req)).Context(c)
}

// AddGrant add a permit grant to a permit group
func (s *APIPermit) AddGrant(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.AddPermitGrantReq{}
	req.Owner, req.OwnerName = s.gate.GetAuthOwnerPair(c)
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.GroupPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.AddPermitGrant(c, req)).Context(c)
}

// DeleteGrant delete a permit grant from a permit group
func (s *APIPermit) DeleteGrant(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.DelPermitGrantReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.DelPermitGrant(c, req)).Context(c)
}

// UpdateGrant update a permit grant from a permit group
func (s *APIPermit) UpdateGrant(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.UpdatePermitGrantReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.UpdatePermitGrant(c, req)).Context(c)
}

// ActiveGrant reverse active state of a api permit grant
func (s *APIPermit) ActiveGrant(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.ActivePermitGrantReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.ActivePermitGrant(c, req)).Context(c)
}

// ListGrant list permit grants within a permit group
func (s *APIPermit) ListGrant(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.ListPermitGrantReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.GroupPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.ListPermitGrant(c, req)).Context(c)
}

// QueryGrant query a permit grant within a permit group
func (s *APIPermit) QueryGrant(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.QueryPermitGrantReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.QueryPermitGrant(c, req)).Context(c)
}
