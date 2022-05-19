package apiprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	ginLogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

var m = &apiProviderManager{
	apiTypeEnum: enumset.New(nil),
	providers:   make(map[string]Provider),
}

func init() {
	enumset.FinishReg()
}

// Provider is the interface that provide api service
type Provider interface {
	APIType() string
	QueryDoc(c context.Context, req *QueryDocReq) (*QueryDocResp, error)
	Request(c context.Context, req *RequestReq) (*RequestResp, error)
}

type apiProviderManager struct {
	apiTypeEnum *enumset.EnumSet
	providers   map[string]Provider
}

func (p *apiProviderManager) queryProvider(apiPath string) (Provider, error) {
	apiType := apipath.GetAPIType(apiPath)
	provider, ok := p.providers[apiType]
	if !ok {
		return nil, errcode.ErrUnrecognizedAPIType.FmtError(apiType)
	}
	return provider, nil
}

// RegistAPIProvider regist api provider
func RegistAPIProvider(provider Provider) error {
	if provider == nil {
		return fmt.Errorf("error:RegistAPIProvider missing  provider(%p)", provider)
	}
	apiType := provider.APIType()
	if _, err := m.apiTypeEnum.Reg(apiType); err != nil {
		return err
	}
	m.providers[apiType] = provider
	return nil
}

// APIRequest request an api
func APIRequest(c context.Context, req *RequestReq) (*RequestResp, error) {
	p, err := m.queryProvider(req.APIPath)
	if err != nil {
		return nil, err
	}
	req.APIType = p.APIType()
	return p.Request(c, req)
}

// APIQueryDoc query an api doc
func APIQueryDoc(c context.Context, req *QueryDocReq) (*QueryDocResp, error) {
	p, err := m.queryProvider(req.APIPath)
	if err != nil {
		return nil, err
	}
	req.APIType = p.APIType()
	return p.QueryDoc(c, req)
}

func copyHeader(dst, src http.Header) {
	for k := range src {
		dst.Add(k, src.Get(k))
	}
}

// DoHTTPRequestResp response a universal api request
func DoHTTPRequestResp(c *gin.Context, req *RequestReq,
	f func(c context.Context, req *RequestReq) (*RequestResp, error)) {
	rep, err := f(c, req)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}

	reps := rep.Response
	if len(reps) == 0 {
		reps = json.RawMessage("{}")
	}

	if rep.StatusCode != http.StatusOK {
		c.AbortWithStatusJSON(rep.StatusCode, reps)
		logger.Logger.Warnf("[api-request-fail] type=%s apipath=%s, status=%s, reqID=%v resp=%s", req.APIType, req.APIPath, rep.Status, ginLogger.GetRequestID(c).String, string(rep.Response))
		return
	}

	copyHeader(c.Writer.Header(), rep.Header)

	// BUG: copyHeader makes "Content-Length" mismatch
	//resp.Format(json.RawMessage(reps), err).Context(c)

	logger.Logger.Debugf("HTTPRequest %s resp: header=%v body=%v", req.APIPath, c.Writer.Header(), string(reps))

	c.Writer.Write(reps)
}
