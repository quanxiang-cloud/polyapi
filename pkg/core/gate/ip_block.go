package gate

import (
	"net/http"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
)

func createIPBlock(cfg *config.Config) (*ipBlock, error) {
	if cfg.Gate.IPBlock.Enable {
		n := &ipBlock{
			config: cfg,
			cfg:    &cfg.Gate.IPBlock,
		}
		n.tire.Init()
		n.tire.BatchInsert(n.cfg.White, white)
		n.tire.BatchInsert(n.cfg.Black, black)
		logger.Logger.Infof("[gate.ipBlock] config:%+v", cfg.Gate.IPBlock)
		return n, nil
	}
	logger.Logger.Infow("[gate.ipBlock] disabled")
	return nil, nil
}

type ipBlock struct {
	config *config.Config
	cfg    *config.GateIPBlock
	tire   TireTree
}

func (v *ipBlock) match(ip string) (byte, bool) {
	return v.tire.Match(ip)
}

func (v *ipBlock) Handle(c *gin.Context, apiType apiType) (err error) {
	ip := GetHTTPClientIPX(c.Request)
	if ip == "" {
		ip = "0.0.0.0"
	}
	defer func() {
		logger.Logger.Debugf("[gate.ipBlock] client ip '%s' ", ip)
		if err != nil {
			logger.Logger.Warnf("[gate.ipBlock] ip blocked '%s' ", ip)
		}
	}()

	if wb, ok := v.match(ip); ok {
		switch wb {
		case white:
			return nil
		case black:
			return errcode.ErrGateBlockedIP.NewError()
		}
	}
	if len(v.cfg.White) > 0 {
		return errcode.ErrGateBlockedIP.NewError()
	}
	return ValidateRequestIP(c.Request)
}

//------------------------------------------------------------------------------

// GetHTTPClientIPX get client IP with "x-real-ip" & "x-forwarded-for" header check
func GetHTTPClientIPX(r *http.Request) string {
	return GetHTTPClientIP(r, "x-real-ip", "x-forwarded-for")
}

// GetHTTPClientIP return the client IP of a http request, by the order
// of a given list of headers. If no ip is found in headers, then return request's
// RemoteAddr. This is useful when there are proxy servers between the client and the backend server.
// Example:
// GetHTTPClientIP(r, "x-real-ip", "x-forwarded-for"), will first check header  "x-real-ip"
// if it exists, then split it by  "," and return the first part. Otherwise, it will check
// the header  "x-forwarded-for" if it exists, then split it by  "," and return the first part.
// Otherwise it will return request's RemoteAddr.
//
func GetHTTPClientIP(r *http.Request, headers ...string) string {
	for _, header := range headers {
		ip := r.Header.Get(header)
		if ip != " " {
			return getFirstPart(ip, ",")
		}
	}
	return getFirstPart(r.RemoteAddr, ":")
}

func getFirstPart(s string, sep string) string {
	if index := strings.Index(s, sep); index >= 0 {
		return s[:index]
	}
	return s
}

// ValidateRequestIP verify the request IP
func ValidateRequestIP(r *http.Request) error {
	// TODO:
	return nil
}
