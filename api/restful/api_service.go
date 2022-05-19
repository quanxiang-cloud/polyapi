package restful

import (
	"github.com/quanxiang-cloud/polyapi/internal/service"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

// APIService is the route for api namespace
type APIService struct {
	s    service.APIService
	gate *gate.Gate
}

// NewAPIService create a api namespace provider
func NewAPIService(conf *config.Config) (*APIService, error) {
	svs, err := service.CreateServiceOper(conf)
	if err != nil {
		return nil, err
	}
	r := &APIService{
		gate: gate.NewGate(conf),
		s:    svs,
	}
	return r, nil
}

func (s *APIService) create(c *gin.Context) {
	req := &service.CreateServiceReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	req.Owner, req.OwnerName = s.gate.GetAuthOwnerPair(c)

	resp.Format(s.s.Create(c, req)).Context(c)
}

// Create create a namespace
func (s *APIService) Create(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}
	s.create(c)
}

// InnerCreate create a namespace without authorize
func (s *APIService) InnerCreate(c *gin.Context) {
	s.gate.SetFromInnerFlag(c, true)
	s.create(c)
}

// Delete delete a namespace
func (s *APIService) Delete(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.DeleteServiceReq{}
	req.Owner = s.gate.GetAuthOwner(c)
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.ServicePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.Delete(c, req)).Context(c)
}

// Update updates a namespace info
func (s *APIService) Update(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.UpdateServiceReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.ServicePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.Update(c, req)).Context(c)
}

// UpdateProperty updateProperty
func (s *APIService) UpdateProperty(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.UpdatePropertyReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.ServicePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.UpdateProperty(c, req)).Context(c)
}

// Active reverse active state of a namespace
func (s *APIService) Active(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.ActiveServiceReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.ServicePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.Active(c, req)).Context(c)
}

// List list child namespaces
func (s *APIService) List(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.ListServiceReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.List(c, req)).Context(c)
}

// Query query namespace info
func (s *APIService) Query(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	var service string
	var err error
	if !getPathArg(c, PathargAPIPath, &service, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.Query(c, service)).Context(c)
}
