package restful

import (
	"encoding/base64"
	"io"

	"github.com/quanxiang-cloud/polyapi/internal/service"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/apiprovider"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/business/app"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"
	"github.com/quanxiang-cloud/polyapi/pkg/core/swagger"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/httputil"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

const (
	// PathargType is type in path arg
	PathargType = "type"

	// PathargAPIPath is apiPath in path arg
	PathargAPIPath = consts.PathArgNamespacePath
)

// Raw is raw api router
type Raw struct {
	svs  service.RawAPI
	gate *gate.Gate
}

// NewRawAPI NewRawAPI
func NewRawAPI(conf *config.Config) (*Raw, error) {
	raw, err := service.CreateRaw(conf)
	if err != nil {
		return nil, err
	}
	return &Raw{
		svs:  raw,
		gate: gate.NewGate(conf),
	}, nil
}

// RegFile RegFile
func (s *Raw) RegFile(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	if file == nil {
		resp.Format(nil, paraErr("missing file")).Context(c)
		return
	}
	req := &service.RegReq{}
	open, _ := file.Open()
	defer open.Close()
	all, _ := io.ReadAll(open)
	if !getPathArg(c, PathargAPIPath, &req.Service, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}

	req.Namespace = c.PostForm("namespace")
	req.Version = c.PostForm("version")

	req.Swagger = unsafeByteString(all)
	req.Owner, req.OwnerName = s.gate.GetAuthOwnerPair(c)
	resp.Format(s.svs.RegSwagger(c, req)).Context(c)
}

func decodeSwagger(s *string, err *error) bool {
	b, e := base64.StdEncoding.DecodeString(*s)
	if e != nil {
		return true // TODO: remove
		// *err = e
		// return false
	}
	*s = string(b)
	return true
}

func (s *Raw) regSwagger(c *gin.Context) {
	req := &service.RegReq{}
	req.Owner, req.OwnerName = s.gate.GetAuthOwnerPair(c)
	if err := bindBody(c, req); err != nil ||
		!getPathArg(c, PathargAPIPath, &req.Service, &err) ||
		!decodeSwagger(&req.Swagger, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	if contentType := c.GetHeader("Content-Type"); contentType != "" {
		enc, err := consts.FromMIME(contentType)
		if err != nil {
			resp.Format(nil, paraErr(err.Error())).Context(c)
			return
		}
		jsonSwagger, err := swagger.ConvertToJSON(enc, req.Swagger, false)
		if err != nil {
			resp.Format(nil, paraErr(err.Error())).Context(c)
			return
		}
		req.Swagger = jsonSwagger
	}

	resp.Format(s.svs.RegSwagger(c, req)).Context(c)
}

// InnerRegSwagger InnerRegSwagger
func (s *Raw) InnerRegSwagger(c *gin.Context) {
	s.gate.SetFromInnerFlag(c, true)
	s.regSwagger(c)
}

// InnerRegSwaggerAlone reg swaager string that is without service
func (s *Raw) InnerRegSwaggerAlone(c *gin.Context) {
	s.gate.SetFromInnerFlag(c, true)
	req := &service.RegReq{}
	req.Owner, req.OwnerName = s.gate.GetAuthOwnerPair(c)
	if err := bindBody(c, req); err != nil ||
		!getPathArg(c, PathargAPIPath, &req.Namespace, &err) ||
		!decodeSwagger(&req.Swagger, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	if contentType := c.GetHeader("Content-Type"); contentType != "" {
		enc, err := consts.FromMIME(contentType)
		if err != nil {
			resp.Format(nil, paraErr(err.Error())).Context(c)
			return
		}
		jsonSwagger, err := swagger.ConvertToJSON(enc, req.Swagger, false)
		if err != nil {
			resp.Format(nil, paraErr(err.Error())).Context(c)
			return
		}
		req.Swagger = jsonSwagger
	}

	req.Service = "" // alone API, without service
	resp.Format(s.svs.RegSwagger(c, req)).Context(c)
}

// RegSwagger RegSwagger
func (s *Raw) RegSwagger(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}
	s.regSwagger(c)
}

// Del Del
func (s *Raw) Del(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.DelReq{}
	req.Owner = s.gate.GetAuthOwner(c)
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.svs.Del(c, req)).Context(c)
}

// Query Query
func (s *Raw) Query(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.QueryReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.APIPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.svs.Query(c, req)).Context(c)
}

// List list raw api in namespace
func (s *Raw) List(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.RawListReq{
		Active: -1,
	}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.svs.List(c, req)).Context(c)
}

// ListInService list raw api in service
func (s *Raw) ListInService(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.ListInServiceReq{
		Active: -1,
	}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.ServicePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.svs.ListInService(c, req)).Context(c)
}

// Active reverse active state of a raw api
func (s *Raw) Active(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	req := &service.ActiveReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.APIPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.svs.Active(c, req)).Context(c)
}

// Search seach raw by title or path, allow to find raw at subdirectory
func (s *Raw) Search(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &service.SearchRawReq{
		Active: -1,
	}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(s.svs.Search(c, req)).Context(c)
}

//------------------------------------------------------------------------------

func (s *Raw) apiProviderRequest(c *gin.Context) {
	req := &apiprovider.RequestReq{}
	if err := getRequestArgs(c, &req.Body); err != nil || !getPathArg(c, PathargAPIPath, &req.APIPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	req.Header = c.Request.Header.Clone()
	req.Header.Set(consts.HeaderRefer, httputil.MakeRefer(c.Request))

	req.Method = c.Request.Method
	req.Owner = s.gate.GetAuthOwner(c)
	req.APIType = ""
	req.APIService = c.GetHeader(consts.XHeaderServicePath)
	req.APIServiceArgs = c.GetHeader(consts.XHeaderServiceArgs)
	if req.APIService != "" {
		if err := app.ValidateAPIPath(req.Owner, req.APIPath, req.APIService); err != nil {
			resp.Format(nil, paraErr(err.Error())).Context(c)
			return
		}
	}

	apiprovider.DoHTTPRequestResp(c, req, apiprovider.APIRequest)
}

// APIProviderRequest request an api
func (s *Raw) APIProviderRequest(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	s.apiProviderRequest(c)
}

// InnerAPIProviderRequest request an api by inner
func (s *Raw) InnerAPIProviderRequest(c *gin.Context) {
	s.gate.SetFromInnerFlag(c, true)
	s.apiProviderRequest(c)
}

// APIProviderQueryDoc query an api doc
func (s *Raw) APIProviderQueryDoc(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}

	req := &apiprovider.QueryDocReq{}
	if err := bindBody(c, &req); err != nil || !getPathArg(c, PathargAPIPath, &req.APIPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	resp.Format(apiprovider.APIQueryDoc(c, req)).Context(c)
}
