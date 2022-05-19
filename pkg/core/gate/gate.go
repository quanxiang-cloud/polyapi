package gate

import (
	"net/http"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/auth"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

var gateInst *Gate

// NewGate create Gate object
func NewGate(cfg *config.Config) *Gate {
	if gateInst == nil {
		n := &Gate{}
		// if err := n.add(createIPBlock(cfg)); err != nil {
		// 	panic(err)
		// }
		// if err := n.add(createLimitRate(cfg)); err != nil {
		// 	panic(err)
		// }
		if err := n.add(createAPISignature(cfg)); err != nil {
			panic(err)
		}
		if err := n.add(createAppPermit(cfg)); err != nil { //NOTE: must after createAPISignature
			panic(err)
		}
		gateInst = n
	}

	return gateInst
}

type handler interface {
	Handle(c *gin.Context, apiType apiType) error
}

// Gate object
type Gate struct {
	h []handler
}

// Filt do gate file action
func (v *Gate) Filt(c *gin.Context, apiType apiType) (err error) {
	defer func() {
		if err != nil {
			//c.AbortWithError(http.StatusUnauthorized, err)
			resp.Format(nil, err).Context(c, http.StatusForbidden)
		}
	}()

	SetFromInnerFlag(c, apiType.IsInner())
	if !apiType.IsInner() {
		if len(v.h) > 0 {
			for _, v := range v.h {
				if err := v.Handle(c, apiType); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// SetFromInnerFlag set HeaderRequestFromInner header
func (v *Gate) SetFromInnerFlag(c *gin.Context, inner bool) error {
	return SetFromInnerFlag(c, inner)
}

// GetAuthUser get the authed user id
func (v *Gate) GetAuthUser(c *gin.Context) string {
	user := auth.GetAuthUserID(c.Request.Header)
	return user
}

// GetAuthOwner get the authed owner id
func (v *Gate) GetAuthOwner(c *gin.Context) string {
	owner := auth.GetAuthOwner(c.Request.Header)
	return owner
}

// GetAuthOwnerPair get the owner id & name
func (v *Gate) GetAuthOwnerPair(c *gin.Context) (string, string) {
	owner := auth.GetAuthOwner(c.Request.Header)
	ownerName := auth.GetAuthOwnerName(c.Request.Header)
	return owner, ownerName
}

// GetNamespacePath get path arg :namespacePath
func (v *Gate) GetNamespacePath(c *gin.Context) string {
	return GetNamespacePath(c)
}

func (v *Gate) add(h handler, err error) error {
	if err != nil {
		return err
	}
	if h != nil {
		v.h = append(v.h, h)
	}
	return nil
}

//------------------------------------------------------------------------------

// SetFromInnerFlag set HeaderRequestFromInner header
func SetFromInnerFlag(c *gin.Context, inner bool) error {
	if inner {
		c.Request.Header.Set(consts.HeaderRequestFromInner, "true")
	} else {
		delete(c.Request.Header, consts.HeaderRequestFromInner)
	}
	return nil
}

// GetNamespacePath get path arg :namespacePath
func GetNamespacePath(c *gin.Context) string {
	return c.Param(consts.PathArgNamespacePath)
}

//------------------------------------------------------------------------------

// APIType is the kind of API
type apiType uint8

// APIType enum
const (
	APIInner apiType = iota + 1
	APIRead
	APIWrite
)

// IsInner check if it is an inner API
func (t apiType) IsInner() bool {
	return t == APIInner
}

// IsReader check if it is a reader API
func (t apiType) IsReader() bool {
	return t == APIRead
}

// IsWriter check if it is a writer API
func (t apiType) IsWriter() bool {
	return t == APIWrite
}
