package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/httputil"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	ginLogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
)

// KMSAgent represents the kms agent
type KMSAgent interface {
	Proxy(c *gin.Context, proxy ProxyFunc) error
	Authorize(c context.Context, keyUUID string, body json.RawMessage, header http.Header) (*KMSAuthorizeResp, error)
	ListCustomerAPIKey(c context.Context, req *ListKMSCustomerKeyReq) (*ListKMSCustomerKeyResp, error)
	QueryCustomerAPIKey(c context.Context, keyUUID string) (*QueryKMSCustomerKeyResp, error)
	UpdateCustomerKeyInBatch(c context.Context, req *UpdateCustomerKeyInBatchReq) (*UpdateCustomerKeyInBatchResp, error)
	DeleteCustomerKeyInBatch(c context.Context, req *DeleteCustomerKeyInBatchReq) (*DeleteCustomerKeyInBatchResp, error)
	DeleteCustomerKeyByPrefix(c context.Context, req *DeleteCustomerKeyByPrefixReq) (*DeleteCustomerKeyByPrefixResp, error)
}

var kmsAgentInst *kmsAgent

// CreateKMSAgent CreateKMSAgent
func CreateKMSAgent(conf *config.Config) (KMSAgent, error) {
	if kmsAgentInst != nil {
		return kmsAgentInst, nil
	}

	kms := &conf.Authorize.OauthKey
	kmsHostBase, err := getKMSHostBase(kms.Addr)
	if err != nil {
		return nil, err
	}

	auth3URI, err := url.ParseRequestURI(kmsHostBase + "/api/v1/kms/ext/authorize")
	if err != nil {
		return nil, err
	}

	listCustomerKeyURI, err := url.ParseRequestURI(kmsHostBase + "/api/v1/kms/ext/list")
	if err != nil {
		return nil, err
	}

	queryCustomerKeyURI, err := url.ParseRequestURI(kmsHostBase + "/api/v1/kms/ext/query")
	if err != nil {
		return nil, err
	}

	updateCustomerKeyInBatchURI, err := url.ParseRequestURI(kmsHostBase + "/api/v1/kms/ext/updateInBatch")
	if err != nil {
		return nil, err
	}

	deleteCustomerKeyInBatchURI, err := url.ParseRequestURI(kmsHostBase + "/api/v1/kms/ext/deleteInBatch")
	if err != nil {
		return nil, err
	}

	deleteCustomerKeyByPrefixURI, err := url.ParseRequestURI(kmsHostBase + "/api/v1/kms/ext/deleteByPrefix")
	if err != nil {
		return nil, err
	}

	checkAuthURI, err := url.ParseRequestURI(kmsHostBase + "/api/v1/kms/ext/checkAuth")
	if err != nil {
		return nil, err
	}

	p := &kmsAgent{
		conf:        conf,
		kmsHostBase: kmsHostBase,
		kmsAPIs: map[string]string{
			"/api/v1/polyapi/apikey/create":        "/api/v1/kms/key/create",
			"/api/v1/polyapi/apikey/delete":        "/api/v1/kms/key/delete",
			"/api/v1/polyapi/apikey/update":        "/api/v1/kms/key/update",
			"/api/v1/polyapi/apikey/active":        "/api/v1/kms/key/active",
			"/api/v1/polyapi/apikey/list":          "/api/v1/kms/key/list",
			"/api/v1/polyapi/apikey/query":         "/api/v1/kms/key/query",
			"/api/v1/polyapi/holdingkey/upload":    "/api/v1/kms/ext/upload",
			"/api/v1/polyapi/holdingkey/delete":    "/api/v1/kms/ext/delete",
			"/api/v1/polyapi/holdingkey/update":    "/api/v1/kms/ext/update",
			"/api/v1/polyapi/holdingkey/active":    "/api/v1/kms/ext/active",
			"/api/v1/polyapi/holdingkey/list":      "/api/v1/kms/ext/list",
			"/api/v1/polyapi/holdingkey/query":     "/api/v1/kms/ext/query",
			"/api/v1/polyapi/holdingkey/authTypes": "/api/v1/kms/ext/authTypes",
			"/api/v1/polyapi/holdingkey/sample":    "/api/v1/kms/ext/sample",
		},
		client: http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					deadline := time.Now().Add(kms.Timeout * time.Second)
					c, err := net.DialTimeout(netw, addr, time.Second*kms.Timeout)
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
				MaxIdleConns:      kms.MaxIdleConns,
				DisableKeepAlives: false,
			},
		},
		auth3URL:                     auth3URI,
		listCustomerKeyURI:           listCustomerKeyURI,
		queryCustomerKeyURI:          queryCustomerKeyURI,
		updateCustomerKeyInBatchURI:  updateCustomerKeyInBatchURI,
		deleteCustomerKeyInBatchURI:  deleteCustomerKeyInBatchURI,
		deleteCustomerKeyByPrefixURI: deleteCustomerKeyByPrefixURI,
		checkAuthURI:                 checkAuthURI,
	}

	kmsAgentInst = p
	adaptor.SetKMSOper(p)
	return p, nil
}

type kmsAgent struct {
	conf                         *config.Config
	kmsHostBase                  string
	auth3URL                     *url.URL // "/api/v1/kms/ext/authorize"
	listCustomerKeyURI           *url.URL // "/api/v1/kms/ext/list"
	queryCustomerKeyURI          *url.URL // "/api/v1/kms/ext/query"
	updateCustomerKeyInBatchURI  *url.URL // "/api/v1/kms/ext/updateInBatch"
	deleteCustomerKeyInBatchURI  *url.URL // "/api/v1/kms/ext/deleteInBatch"
	deleteCustomerKeyByPrefixURI *url.URL // "/api/v1/kms/ext/deleteByPrefix"
	checkAuthURI                 *url.URL // "/api/v1/kms/ext/checkAuth"
	kmsAPIs                      map[string]string
	client                       http.Client
}

// http://kms/api/xxx => http://kms
func getKMSHostBase(url string) (string, error) {
	delimiterGot := false // found "//"
	var last byte
	for i := 0; i < len(url); i++ {
		v := url[i]
		if v == '/' {
			if !delimiterGot {
				if last == '/' {
					delimiterGot = true
				}
			} else {
				return url[:i], nil
			}
			last = v
		}
	}

	return "", fmt.Errorf("invalid kms base url: %s", url)
}

func (a *kmsAgent) redirectURL(c *gin.Context) (*url.URL, error) {
	// c.Request.URL.Path:"/api/v1/xxx"
	apiPath := c.Request.URL.Path
	kmsPath, ok := a.kmsAPIs[apiPath]
	if !ok {
		return nil, fmt.Errorf("kms api mapping not found: %s", apiPath)
	}

	uri := a.kmsHostBase + kmsPath

	return url.ParseRequestURI(uri)
}

// Proxy request kms api by proxy
func (a *kmsAgent) Proxy(c *gin.Context, proxy ProxyFunc) (e error) {
	// defer func() {
	// 	if e != nil {
	// 		fmt.Println("err:", e)
	// 	}
	// }()
	kmsURI, err := a.redirectURL(c)
	if err != nil {
		return err
	}

	if proxy == nil {
		proxy = DefultProxyFunc
	}
	reqBody, err := proxy(c)
	if err != nil {
		return err
	}

	req := &http.Request{
		Header: c.Request.Header.Clone(),
		URL:    kmsURI,
		Body:   reqBody,
		Method: c.Request.Method,
	}

	body, resp, err := a.doRequest(req)
	if err != nil {
		logger.Logger.Errorw("[kms]", ginLogger.GetRequestID(c), err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(resp.StatusCode)
		return nil
	}

	c.Writer.Write(body)
	return nil
}

type authorize3Req struct {
	ID   string      `json:"ID"`
	Body interface{} `json:"body"`
}

// KMSAuthorizeResp exports
type KMSAuthorizeResp = adaptor.KMSAuthorizeResp

// Authorize proxy the third party authorize
func (a *kmsAgent) Authorize(c context.Context, keyUUID string, body json.RawMessage, header http.Header) (*KMSAuthorizeResp, error) {
	req := authorize3Req{
		ID:   keyUUID,
		Body: body,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request := &http.Request{
		Header: header,
		URL:    a.auth3URL,
		Body:   io.NopCloser(bytes.NewReader(b)),
		Method: http.MethodPost,
	}

	respBody, resp, err := a.doRequest(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errcode.ErrInternal.FmtError(resp.Status)
	}

	var r struct {
		Code int                       `json:"code"`
		Msg  string                    `json:"msg"`
		Data *adaptor.KMSAuthorizeResp `json:"data"`
	}
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("internal error: %d,%s", r.Code, r.Msg)
	}
	return r.Data, nil
}

// ListKMSCustomerKeyReq req
type ListKMSCustomerKeyReq = adaptor.ListKMSCustomerKeyReq

// KMSCustomerKey key info except secret
type KMSCustomerKey = adaptor.KMSCustomerKey

// ListKMSCustomerKeyResp resp
type ListKMSCustomerKeyResp = adaptor.ListKMSCustomerKeyResp

// ListCustomerAPIKey list customer secret key by service
func (a *kmsAgent) ListCustomerAPIKey(c context.Context, req *ListKMSCustomerKeyReq) (*ListKMSCustomerKeyResp, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request := &http.Request{
		Header: http.Header{},
		URL:    a.listCustomerKeyURI,
		Body:   io.NopCloser(bytes.NewReader(b)),
		Method: http.MethodPost,
	}
	if req.Owner != "" {
		request.Header.Add(consts.HeaderUserID, req.Owner)
	}

	respBody, resp, err := a.doRequest(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errcode.ErrInternal.FmtError(resp.Status)
	}

	var r struct {
		Code int                             `json:"code"`
		Msg  string                          `json:"msg"`
		Data *adaptor.ListKMSCustomerKeyResp `json:"data"`
	}
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("internal error: %d,%s", r.Code, r.Msg)
	}
	return r.Data, nil
}

type queryKMSCustomerKeyReq struct {
	ID string `json:"id"`
}

// QueryKMSCustomerKeyResp resp
type QueryKMSCustomerKeyResp = adaptor.QueryKMSCustomerKeyResp

// QueryCustomerAPIKey query customer api key
func (a *kmsAgent) QueryCustomerAPIKey(c context.Context, keyUUID string) (*QueryKMSCustomerKeyResp, error) {
	req := queryKMSCustomerKeyReq{
		ID: keyUUID,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request := &http.Request{
		Header: http.Header{},
		URL:    a.queryCustomerKeyURI,
		Body:   io.NopCloser(bytes.NewReader(b)),
		Method: http.MethodPost,
	}

	respBody, resp, err := a.doRequest(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errcode.ErrInternal.FmtError(resp.Status)
	}

	var r struct {
		Code int                              `json:"code"`
		Msg  string                           `json:"msg"`
		Data *adaptor.QueryKMSCustomerKeyResp `json:"data"`
	}
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("internal error: %d,%s", r.Code, r.Msg)
	}
	return r.Data, nil
}

// UpdateCustomerKeyInBatchReq UpdateCustomerAPIKeyInBatchReq
type UpdateCustomerKeyInBatchReq = adaptor.UpdateCustomerKeyInBatchReq

// UpdateCustomerKeyInBatchResp UpdateCustomerAPIKeyInBatchResp
type UpdateCustomerKeyInBatchResp = adaptor.UpdateCustomerKeyInBatchResp

// UpdateCustomerKeyInBatch UpdateCustomerAPIKeyInBatch
func (a *kmsAgent) UpdateCustomerKeyInBatch(c context.Context, req *UpdateCustomerKeyInBatchReq) (*UpdateCustomerKeyInBatchResp, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request := &http.Request{
		Header: http.Header{},
		URL:    a.updateCustomerKeyInBatchURI,
		Body:   io.NopCloser(bytes.NewReader(b)),
		Method: http.MethodPost,
	}
	respBody, resp, err := a.doRequest(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errcode.ErrInternal.FmtError(resp.Status)
	}

	var r struct {
		Code int                                   `json:"code"`
		Msg  string                                `json:"msg"`
		Data *adaptor.UpdateCustomerKeyInBatchResp `json:"data"`
	}
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("internal error: %d,%s", r.Code, r.Msg)
	}
	return r.Data, nil
}

// DeleteCustomerKeyInBatchReq DeleteCustomerKeyInBatchReq
type DeleteCustomerKeyInBatchReq = adaptor.DeleteCustomerKeyInBatchReq

// DeleteCustomerKeyInBatchResp DeleteCustomerKeyInBatchResp
type DeleteCustomerKeyInBatchResp = adaptor.DeleteCustomerKeyInBatchResp

func (a *kmsAgent) DeleteCustomerKeyInBatch(c context.Context, req *DeleteCustomerKeyInBatchReq) (*DeleteCustomerKeyInBatchResp, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request := &http.Request{
		Header: http.Header{},
		URL:    a.deleteCustomerKeyInBatchURI,
		Body:   io.NopCloser(bytes.NewReader(b)),
		Method: http.MethodPost,
	}
	respBody, resp, err := a.doRequest(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errcode.ErrInternal.FmtError(resp.Status)
	}

	var r struct {
		Code int                                   `json:"code"`
		Msg  string                                `json:"msg"`
		Data *adaptor.DeleteCustomerKeyInBatchResp `json:"data"`
	}
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("internal error: %d,%s", r.Code, r.Msg)
	}
	return r.Data, nil
}

// DeleteCustomerKeyByPrefixReq DeleteCustomerKeyByPrefixReq
type DeleteCustomerKeyByPrefixReq = adaptor.DeleteCustomerKeyByPrefixReq

// DeleteCustomerKeyByPrefixResp DeleteCustomerKeyByPrefixResp
type DeleteCustomerKeyByPrefixResp = adaptor.DeleteCustomerKeyByPrefixResp

func (a *kmsAgent) DeleteCustomerKeyByPrefix(c context.Context, req *DeleteCustomerKeyByPrefixReq) (*DeleteCustomerKeyByPrefixResp, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request := &http.Request{
		Header: http.Header{},
		URL:    a.deleteCustomerKeyByPrefixURI,
		Body:   io.NopCloser(bytes.NewReader(b)),
		Method: http.MethodPost,
	}
	respBody, resp, err := a.doRequest(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errcode.ErrInternal.FmtError(resp.Status)
	}

	var r struct {
		Code int                                    `json:"code"`
		Msg  string                                 `json:"msg"`
		Data *adaptor.DeleteCustomerKeyByPrefixResp `json:"data"`
	}
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("internal error: %d,%s", r.Code, r.Msg)
	}
	return r.Data, nil
}

// CheckAuthReq CheckAuthReq
type CheckAuthReq = adaptor.CheckAuthReq

// CheckAuthResp CheckAuthResp
type CheckAuthResp = adaptor.CheckAuthResp

func (a *kmsAgent) CheckAuth(c context.Context, req *CheckAuthReq) (*CheckAuthResp, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request := &http.Request{
		Header: http.Header{},
		URL:    a.checkAuthURI,
		Body:   io.NopCloser(bytes.NewReader(b)),
		Method: http.MethodPost,
	}
	respBody, resp, err := a.doRequest(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errcode.ErrInternal.FmtError(resp.Status)
	}

	var r struct {
		Code int                    `json:"code"`
		Msg  string                 `json:"msg"`
		Data *adaptor.CheckAuthResp `json:"data"`
	}
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("%s", r.Msg)
	}
	return r.Data, nil
}

func (a *kmsAgent) doRequest(req *http.Request) ([]byte, *http.Response, error) {
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, nil, err
	}
	return body, resp, nil
}

// ProxyFunc is the function to make proxy body
type ProxyFunc func(c *gin.Context) (io.ReadCloser, error)

// DefultProxyFunc is default ProxyFunc
func DefultProxyFunc(c *gin.Context) (io.ReadCloser, error) {
	var body json.RawMessage
	if err := httputil.BindBody(c, &body); err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(body)), nil
}
