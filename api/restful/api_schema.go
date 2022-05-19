package restful

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	ginLogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/polyapi/internal/service"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"
)

// APISchema APISchema
type APISchema struct {
	s    service.APISchemaOper
	gate *gate.Gate
}

// NewAPISchema NewAPISchema
func NewAPISchema(cfg *config.Config) (*APISchema, error) {
	s, err := service.CreateSchemaOper(cfg)
	if err != nil {
		return nil, err
	}

	apiSchema := &APISchema{
		s:    s,
		gate: gate.NewGate(cfg),
	}
	return apiSchema, nil
}

// GenSchema gen api json schema
func (s *APISchema) GenSchema(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil {
		return
	}

	req := &service.GenSchemaReq{}
	req.Owner, req.OwnerName = s.gate.GetAuthOwnerPair(c)
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.APIPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}

	resp.Format(s.s.GenSchema(c, req)).Context(c)
}

// QuerySchema QuerySchema
func (s *APISchema) QuerySchema(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil {
		return
	}

	req := &service.QuerySchemaReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.APIPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}

	resp.Format(s.s.QuerySchema(c, req)).Context(c)
}

// ListSchema ListSchema
func (s *APISchema) ListSchema(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil {
		return
	}

	req := &service.ListSchemaReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.NamespacePath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}

	resp.Format(s.s.ListSchema(c, req)).Context(c)
}

// DeleteSchema DeleteSchema
func (s *APISchema) DeleteSchema(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil {
		return
	}

	req := &service.DeleteSchemaReq{}
	if err := bindBody(c, req); err != nil || !getPathArg(c, PathargAPIPath, &req.APIPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}

	resp.Format(s.s.DeleteSchema(c, req)).Context(c)
}

// APISchemaRequest APISchemaRequest
// TODO: copy from raw_api.go apiprovider request, maybe unite requestion func
func (s *APISchema) APISchemaRequest(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil {
		return
	}

	req := &service.SchemaRequestReq{}
	if err := getRequestArgs(c, &req); err != nil || !getPathArg(c, PathargAPIPath, &req.APIPath, &err) {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}
	req.Header = c.Request.Header.Clone()
	req.Method = c.Request.Method
	req.Owner = s.gate.GetAuthOwner(c)
	req.APIType = ""

	rep, err := s.s.Request(c, req)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}

	if rep.StatusCode != http.StatusOK {
		c.AbortWithStatus(rep.StatusCode)
		logger.Logger.Warnf("[api-request-fail] type=%s apipath=%s, status=%s, reqID=%v resp=%s", req.APIType, req.APIPath, rep.Status, ginLogger.GetRequestID(c), rep.Response)
		return
	}

	copyHeader(c.Writer.Header(), rep.Header)
	reps := rep.Response
	if len(reps) == 0 {
		reps = json.RawMessage("{}")
	}

	// BUG: copyHeader makes "Content-Length" mismatch
	//resp.Format(json.RawMessage(reps), err).Context(c)

	c.Writer.Write(reps)
}

func copyHeader(dst, src http.Header) {
	for k := range src {
		dst.Add(k, src.Get(k))
	}
}
