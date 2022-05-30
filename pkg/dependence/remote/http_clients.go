package apisecret

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	ginLogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/polysign"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/hash"
)

func cloneProfile(dst *http.Header, src http.Header) {
	dst.Set(consts.HeaderUserID, deepCopy(src.Values(consts.HeaderUserID)))
	dst.Set(consts.HeaderUserName, deepCopy(src.Values(consts.HeaderUserName)))
	dst.Set(consts.HeaderDepartmentID, deepCopy(src.Values(consts.HeaderDepartmentID)))
	dst.Set(consts.HeaderTenantID, deepCopy(src.Values(consts.HeaderTenantID)))
}

// DeepCopy copy the first none-empty header value
func deepCopy(src []string) string {
	for _, elem := range src {
		if elem != "" {
			return elem
		}
	}
	return ""
}

// NewAuthorizeClient creates authorize client
func NewAuthorizeClient(cfg *config.AuthorizeConfig) (*APIAuthVerifier, error) {
	r := &APIAuthVerifier{
		// oauthToken:      (*oauthTokenResolver)(newOathClient(&cfg.OauthToken)),
		// oauthKey:        (*oauthKeyResolver)(newOathClient(&cfg.OauthKey)),
		// goalie:          (*goalieResolver)(newOathClient(&cfg.Goalie)),
		//fileServer:      (*fileUploadResolver)(newOathClient(&cfg.FileServer)),
		appCenterServer: newAppClient(&cfg.AppAccess, &cfg.AppAdmin),
	}
	//adaptor.SetFileServerOper(r.fileServer) // set adaptor
	adaptor.SetAppCenterServerOper(r.appCenterServer)
	return r, nil
}

// APIAuthVerifier presents API auth verifier
type APIAuthVerifier struct {
	// oauthToken      resolver
	// oauthKey        resolver
	// goalie          resolver
	//fileServer      *fileUploadResolver // used for adaptor only
	appCenterServer *appCenterResolver
}

// Authorize verify authorize info for apis
func (v *APIAuthVerifier) Authorize(c *gin.Context, verifySignature bool) error {
	//
	// 	if token := c.GetHeader(polysign.XHeaderPolyAccessToken); token != "" {
	// 		if err := v.authByToken(c, token); err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		if err := v.authByKey(c, verifySignature); err != nil {
	// 			return err
	// 		}
	// 	}

	// 	if err := v.goalie.request(c, nil); err != nil {
	// 		return err
	// 	}

	// 	if c.GetHeader(consts.HeaderRequestID) == "" {
	// 		c.Request.Header.Add(consts.HeaderRequestID, hash.GenID("req"))
	// 	}
	//
	return nil
}

// func (v *APIAuthVerifier) authByToken(c *gin.Context, token string) error {
// 	if err := v.oauthToken.request(c, nil); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (v *APIAuthVerifier) authByKey(c *gin.Context, verifySignature bool) error {
// 	if err := v.verifyAPISignature(c, verifySignature); err != nil {
// 		return err
// 	}
// 	return nil
// }

//------------------------------------------------------------------------------
type requestArg struct {
	signature       string
	accessKeyID     string
	body            json.RawMessage
	verifySignature bool
}
type resolver interface {
	request(c *gin.Context, arg *requestArg) error
}

type oauthKeyResolver oauthClient

type oauthKeyReq struct {
	Key string `json:"key"`
}

type oauthKeyResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Secret   string `json:"secret"`
		UserInfo struct {
			UserID       string `json:"userID"`
			UserName     string `json:"userName"`
			DepartmentID string `json:"departmentID"`
		} `json:"userInfo"`
	} `json:"data"`
}

func (r *oauthKeyResolver) request(c *gin.Context, arg *requestArg) error {
	req := &http.Request{
		Method: http.MethodPost,
		Header: http.Header{},
		URL:    r.url,
		Body:   io.NopCloser(bytes.NewReader(arg.body)),
	}
	req.Header.Set(consts.HeaderContentType, consts.MIMEJSON)
	req.Header.Set(polysign.XHeaderPolySignKeyID, arg.accessKeyID)

	resp, err := r.client.Do(req)
	if err != nil {
		logger.Logger.Errorw("[oauthKey]", ginLogger.GetRequestID(c))
		return err
	}

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(resp.StatusCode)
		return errcode.ErrInternal.FmtError(resp.Status)
	}

	if arg.verifySignature {
		expect := resp.Header.Get(polysign.XInternalHeaderPolySignSignature)
		got := arg.signature
		if expect == "" || got != expect {
			//fmt.Printf("****expect signature: %s  body=%s\n", expect, string(arg.body))
			return errcode.ErrInputArgValidateMismatch.FmtError("body", polysign.XBodyPolySignSignature)
		}
	}

	cloneProfile(&c.Request.Header, resp.Header)
	c.Request.Header.Set(consts.HeaderAccessKeyID, consts.HeaderAccessKeyIDPrefix+arg.accessKeyID) // auth by access key
	delete(c.Request.Header, polysign.XHeaderPolySignKeyID)
	delete(c.Request.Header, polysign.XHeaderPolySignMethod)
	delete(c.Request.Header, polysign.XHeaderPolySignVersion)
	delete(c.Request.Header, polysign.XHeaderPolySignTimestamp)

	return nil
}

//------------------------------------------------------------------------------

type oauthTokenResolver oauthClient

func (r *oauthTokenResolver) request(c *gin.Context, arg *requestArg) error {
	// uri := ctx.Request.URL.String()
	// if o.inWhiteList(uri) {
	// 	return false, nil
	// }

	req := &http.Request{
		Header: c.Request.Header.Clone(),
		URL:    r.url,
	}

	req.Header.Set(consts.HeaderContentType, consts.MIMEJSON)

	resp, err := r.client.Do(req)
	if err != nil {
		logger.Logger.Errorw("[auth-key]", ginLogger.GetRequestID(c))
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(resp.StatusCode)
		return errcode.ErrInternal.FmtError(resp.Status)
	}

	delete(c.Request.Header, consts.HeaderAccessToken)
	cloneProfile(&c.Request.Header, resp.Header)

	return nil
}

//------------------------------------------------------------------------------

type goalieResolver oauthClient

func (r *goalieResolver) request(c *gin.Context, arg *requestArg) error {
	req := &http.Request{
		Header: c.Request.Header.Clone(),
		URL:    r.url,
	}

	req.Header.Set(consts.HeaderContentType, consts.MIMEJSON)

	resp, err := r.client.Do(req)
	if err != nil {
		logger.Logger.Errorw("[goalie]", ginLogger.GetRequestID(c))
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(resp.StatusCode)
		return errcode.ErrInternal.FmtError(resp.Status)
	}

	return nil
}

//------------------------------------------------------------------------------

type appCenterResolver struct {
	appAccess *oauthClient
	appAdmin  *oauthClient
}

func newAppClient(appAccessCfg, appAdminCfg *config.OauthConfig) *appCenterResolver {
	return &appCenterResolver{
		appAccess: newOathClient(appAccessCfg),
		appAdmin:  newOathClient(appAdminCfg),
	}
}

// CheckReq CheckReq
type CheckReq struct {
	UserID  string `json:"userID"`
	DepID   string `json:"depID"`
	AppID   string `json:"appID"`
	IsSuper bool   `json:"is_super"`
}

func (r appCenterResolver) Check(c context.Context, userID, depID, appID string, isSuper, admin bool) (bool, error) {
	if !admin {
		if ok, err := r.CheckAccess(c, userID, depID, appID); err == nil && ok {
			return ok, err
		}
	}
	return r.CheckAdmin(c, userID, depID, appID, isSuper)
}

func (r appCenterResolver) CheckAdmin(c context.Context, userID, depID, appID string, isSuper bool) (bool, error) {
	body, err := r.buildBody(userID, depID, appID, isSuper)
	if err != nil {
		return false, err
	}

	req := &http.Request{
		URL:    r.appAdmin.url,
		Header: http.Header{},
		Body:   body,
		Method: http.MethodPost,
	}

	req.Header.Set(consts.HeaderContentType, consts.MIMEJSON)
	req.Header.Set(consts.HeaderRequestID, hash.GenID("req"))

	resp, err := r.appAdmin.client.Do(req)
	if err != nil {
		return false, err
	}

	appAdmin := resp.Header.Get("X-App-Admin")
	return appAdmin == "true", nil
}

func (r appCenterResolver) CheckAccess(c context.Context, userID, depID, appID string) (bool, error) {
	body, err := r.buildBody(userID, depID, appID, false)
	if err != nil {
		return false, err
	}

	req := &http.Request{
		URL:    r.appAccess.url,
		Header: http.Header{},
		Body:   body,
		Method: http.MethodPost,
	}

	req.Header.Set(consts.HeaderContentType, consts.MIMEJSON)
	req.Header.Set(consts.HeaderRequestID, hash.GenID("req"))

	resp, err := r.appAccess.client.Do(req)
	if err != nil {
		return false, err
	}

	appAccess := resp.Header.Get("X-App-Access")
	return appAccess == "true", nil
}

func (r appCenterResolver) buildBody(userID, depID, appID string, isSuper bool) (io.ReadCloser, error) {
	data := &CheckReq{
		UserID:  userID,
		DepID:   depID,
		AppID:   appID,
		IsSuper: isSuper,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	byteBuf := bytes.NewBuffer(b)
	return io.NopCloser(byteBuf), nil
}

type oauthClient struct {
	url    *url.URL
	client http.Client
}

func newOathClient(cfg *config.OauthConfig) *oauthClient {
	uri, err := url.ParseRequestURI(cfg.Addr)
	if err != nil {
		panic(err)
	}

	c := &oauthClient{
		url: uri,
		client: http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					deadline := time.Now().Add(cfg.Timeout * time.Second)
					c, err := net.DialTimeout(netw, addr, time.Second*cfg.Timeout)
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
				MaxIdleConns:      cfg.MaxIdleConns,
				DisableKeepAlives: false,
			},
		},
	}
	initFileServer()
	return c
}
