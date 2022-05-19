package restful

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/quanxiang-cloud/polyapi/internal/service"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

// APIKeyHolding is the route for customer api key
type APIKeyHolding struct {
	s    service.KMSAgent
	gate *gate.Gate
}

// NewAPIKeyHolding create a customer api key provider
func NewAPIKeyHolding(conf *config.Config) (*APIKeyHolding, error) {
	svs, err := service.CreateKMSAgent(conf)
	if err != nil {
		return nil, err
	}
	r := &APIKeyHolding{
		gate: gate.NewGate(conf),
		s:    svs,
	}
	return r, nil
}

// UploadHoldKeyReq is the upload key request schema
type UploadHoldKeyReq struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	KeyID       string `json:"keyID"`
	KeySecret   string `json:"keySecret"`
	KeyContent  string `json:"keyContent"`
	Service     string `json:"service"`

	Host        string `json:"host" binding:"-"`
	AuthType    string `json:"authType" binding:"-"`
	AuthContent string `json:"authContent" binding:"-"`
}

// Upload upload a customer api key
func (s *APIKeyHolding) Upload(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}

	serviceOper := adaptor.GetServiceOper()
	if serviceOper == nil {
		resp.Format(nil, errInternal)
		return
	}

	var req UploadHoldKeyReq
	if err := bindBody(c, &req); err != nil {
		resp.Format(nil, paraErr(err.Error())).Context(c)
		return
	}

	var svs *service.ServicesResp
	var err error
	if req.Service != "" {
		svs, err = serviceOper.Check(c, req.Service, s.gate.GetAuthOwner(c), service.OpQuery)
	}
	if req.Service == "" || err != nil {
		resp.Format(nil, paraErr("invalid service")).Context(c)
		return
	}

	req.Host = svs.Host
	req.AuthType = svs.AuthType
	req.AuthContent = svs.AuthContent

	proxy := func(c *gin.Context) (io.ReadCloser, error) {
		b, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		r := bytes.NewReader(b)
		return io.NopCloser(r), nil
	}
	if err := s.s.Proxy(c, proxy); err != nil {
		resp.Format(nil, errInternal)
		return
	}
}

// Update create a customer api key info
func (s *APIKeyHolding) Update(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}
	s.s.Proxy(c, nil)
}

// Active update a customer api key active status
func (s *APIKeyHolding) Active(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}
	s.s.Proxy(c, nil)
}

// UpdateByService update holding key info when service update(service level only)
// func (s *APIKeyHolding) UpdateByService(c *gin.Context) {
// 	if err := s.sign.verifyAPISignature(c, readerAPI); err != nil { // verify signature
//
// 		return
// 	}
// }

// Delete delete a customer api key
func (s *APIKeyHolding) Delete(c *gin.Context) {
	if err := s.gate.Filt(c, APIWrite); err != nil { //gate filter
		return
	}
	s.s.Proxy(c, nil)
}

// List list customer api key in service
func (s *APIKeyHolding) List(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}
	s.s.Proxy(c, nil)
}

// Query query customer api key
func (s *APIKeyHolding) Query(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}
	s.s.Proxy(c, nil)
}

// AuthTypes query all support auth types
func (s *APIKeyHolding) AuthTypes(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil { //gate filter
		return
	}
	s.s.Proxy(c, nil)
}

// Sample get auth sample
func (s *APIKeyHolding) Sample(c *gin.Context) {
	if err := s.gate.Filt(c, APIRead); err != nil {
		return
	}
	s.s.Proxy(c, nil)
}
