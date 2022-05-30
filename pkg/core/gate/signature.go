package gate

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	apisecret "github.com/quanxiang-cloud/polyapi/pkg/dependence/remote"

	"github.com/gin-gonic/gin"
)

func createAPISignature(conf *config.Config) (*apiSignature, error) {
	authClient, err := apisecret.NewAuthorizeClient(&conf.Authorize)
	if err != nil {
		return nil, err
	}

	obj := &apiSignature{
		authClient: authClient,

		enableSignature: true, //NOTE: always enable signature check
		openRegister:    true, //NOTE: always enable register
	}

	return obj, nil
}

// APISignature is the object for judge signature
type apiSignature struct {
	authClient      *apisecret.APIAuthVerifier
	enableSignature bool
	openRegister    bool
}

func (s *apiSignature) setFromInnerFlag(c *gin.Context, inner bool) (err error) {
	return SetFromInnerFlag(c, inner)
}

func (s *apiSignature) Handle(c *gin.Context, apiType apiType) error {
	return s.verifyAPISignature(c, apiType)
}

func (s *apiSignature) verifyAPISignature(c *gin.Context, apiType apiType) (err error) {
	if s.enableSignature && apiType.IsWriter() && !s.openRegister { // not open for outter
		return errcode.ErrAPINotOpen.NewError()
	}
	return nil
}
