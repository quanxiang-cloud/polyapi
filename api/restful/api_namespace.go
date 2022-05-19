package restful

import (
	"github.com/quanxiang-cloud/polyapi/internal/service"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

// APINamespace is the route for api namespace
type APINamespace struct {
	s    service.NamespaceAPI
	gate *gate.Gate
}

// NewAPINamespace create a api namespace provider
func NewAPINamespace(conf *config.Config) (*APINamespace, error) {
	svs, err := service.CreateNamespaceOper(conf)
	if err != nil {
		return nil, err
	}
	r := &APINamespace{
		gate: gate.NewGate(conf),
		s:    svs,
	}
	return r, nil
}

func (s *APINamespace) create(c *gin.Context) {
	req := &service.CreateNsReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.Namespace, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	req.Owner, req.OwnerName = s.gate.GetAuthOwnerPair(c)
	resp.Format(s.s.Create(c, req)).Context(c)
}

// Create create a namespace
func (s *APINamespace) Create(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}
	s.create(c)
}

// InnerCreate create a namespace without authorize
func (s *APINamespace) InnerCreate(c *gin.Context) {
	s.gate.SetFromInnerFlag(c, true)
	s.create(c)
}

// Delete delete a namespace
func (s *APINamespace) Delete(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.DeleteNsReq{}
	req.Owner = s.gate.GetAuthOwner(c)
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.Delete(c, req)).Context(c)
}

// Update updates a namespace info
func (s *APINamespace) Update(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.UpdateNsReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.Update(c, req)).Context(c)
}

// Active reverse active state of a namespace
func (s *APINamespace) Active(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.ActiveNsReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.Active(c, req)).Context(c)
}

// List list child namespaces
func (s *APINamespace) List(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.ListNsReq{
		Active: -1,
	}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.List(c, req)).Context(c)
}

// Tree list namespace tree
func (s *APINamespace) Tree(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.NsTreeReq{
		Active: -1,
	}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.RootPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.Tree(c, req)).Context(c)
}

// Search search namespace by namespace title and active, allow to find sub namespace
func (s *APINamespace) Search(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.SearchNsReq{
		Active: -1,
	}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.Search(c, req)).Context(c)
}

// Query query namespace info
func (s *APINamespace) Query(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	var namespace string
	var err error
	if !getPathArg(c, PathargAPIPath, &namespace, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.Query(c, namespace)).Context(c)
}

// APPPath query namespace of app path
func (s *APINamespace) APPPath(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.AppPathReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.APPPath(c, req)).Context(c)
}

// InitAPPPath init app path
func (s *APINamespace) InitAPPPath(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.InitAppPathReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	d := &req.Data
	d.Owner, d.OwnerName = s.gate.GetAuthOwnerPair(c)
	d.Header = c.Request.Header
	resp.Format(s.s.InitAPPPath(c, req)).Context(c)
}

// InnerInitAPPPath init app path from inner
func (s *APINamespace) InnerInitAPPPath(c *gin.Context) {
	s.gate.SetFromInnerFlag(c, true)

	req := &service.InitAppPathReq{}
	if err := bindBody(c, req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	d := &req.Data
	d.Owner, d.OwnerName = s.gate.GetAuthOwnerPair(c)
	d.Header = c.Request.Header
	resp.Format(s.s.InitAPPPath(c, req)).Context(c)
}

// InnerDelete delete a namespace from inner
func (s *APINamespace) InnerDelete(c *gin.Context) {
	s.gate.SetFromInnerFlag(c, true)

	req := &service.DeleteNsReq{}
	req.Owner = s.gate.GetAuthOwner(c)
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.s.InnerDelete(c, req)).Context(c)
}
