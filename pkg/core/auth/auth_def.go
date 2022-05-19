package auth

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
)

// Enum export
type Enum = enumset.Enum

// AuthType enum define
var (
	AuthTypeEnum = enumset.New(nil) // API authorize method

	AuthNone      = AuthTypeEnum.MustReg("none")      //none authorize
	AuthSystem    = AuthTypeEnum.MustReg("system")    //system authorize
	AuthSignature = AuthTypeEnum.MustReg("signature") //signature authorize
	AuthCookie    = AuthTypeEnum.MustReg("cookie")    //represents cookie authorize
	// AuthOauth2 = AuthTypeEnum.MustReg("oauth2")
)

// ValidateAuthType verify the parameter of authType
func ValidateAuthType(at string) error {
	if !AuthTypeEnum.Verify(at) {
		return errcode.ErrInvalidAuthType.FmtError(at, AuthTypeEnum.GetAll())
	}
	return nil
}

// RequireAPIKey check if an authType require API key
func RequireAPIKey(authType string) bool {
	switch a := Enum(authType); a {
	case AuthNone, AuthSystem, "":
		return false
	case AuthSignature:
		return true
	}
	return true
}

//------------------------------------------------------------------------------

type authNone struct{}

func (t authNone) TypeName() string { return AuthNone.String() }

//------------------------------------------------------------------------------

type authSystem struct{}

func (t authSystem) TypeName() string { return AuthSystem.String() }

//------------------------------------------------------------------------------

type authSignature struct {
	Cmds string `json:"cmds"`
}

func (t authSignature) TypeName() string { return AuthSignature.String() }
