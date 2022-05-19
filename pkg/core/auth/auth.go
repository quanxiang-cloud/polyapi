package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/business/rule"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/hash"

	"github.com/quanxiang-cloud/cabin/logger"
)

// ValidateAPIKeyPermission check if owenr has permission of keyUUID
func ValidateAPIKeyPermission(keyUUID string, service string, owner string) error {
	if keyUUID == "" {
		return nil
	}
	if op := adaptor.GetKMSOper(); op != nil {
		key, err := op.QueryCustomerAPIKey(context.Background(), keyUUID)
		if err != nil {
			return err
		}
		if key.Service != service {
			return errcode.ErrAPIKeyServiceMismatch.NewError()
		}
	}

	// TODO: check api key permission...

	return nil
}

// AppendAuth append auth info for an http request
func AppendAuth(keyUUID, authType string, orginHeader, reqHeader http.Header, body string) string {
	bd, err := AppendAuthWithError(keyUUID, authType, orginHeader, reqHeader, body)
	if err != nil {
		logger.Logger.Errorf("AppendAuth error: keyUUID=%s, authType=%s err=%s",
			keyUUID, authType, err.Error())
	}
	return bd
}

// AppendAuthWithError append auth info for an http request with error handling
func AppendAuthWithError(keyUUID, authType string, orginHeader, reqHeader http.Header, body string) (string, error) {
	if requestID := getOriginRequestID(orginHeader); requestID != "" {
		reqHeader.Set(consts.HeaderRequestID, requestID)
	}
	switch a := Enum(authType); a {
	case AuthNone:
		//do nothing
	case AuthSystem:
		copySystemHeader(orginHeader, &reqHeader)
	case AuthSignature, AuthCookie:
		if op := adaptor.GetKMSOper(); op != nil {
			resp, err := op.Authorize(context.Background(), keyUUID, json.RawMessage(body), orginHeader) // generate token from kms
			if err != nil {
				return body, err
			}
			toBody := []*adaptor.KMSAuthorizeRespItem{}
			for _, v := range resp.Token {
				switch v.In {
				case consts.ParaInBody, consts.ParaInQuery:
					toBody = append(toBody, v)
				case consts.ParaInHeader:
					reqHeader.Set(v.Name, v.Value)
				default:
					return body, fmt.Errorf("unsupport parameter in [%s] from KMS.Authorize", v.In)
				}
			}
			if len(toBody) > 0 {
				newBody, err := appendToBody(body, toBody)
				if err != nil {
					logger.Logger.Warnf("appendToBody('%s', %+v) fail: %s", body, toBody, err.Error())
					return body, err
				}
				body = newBody
			}
		}
	}

	if authType != AuthSystem.String() {
		removeSystemHeader(&reqHeader)
	}

	return body, nil
}

func appendToBody(body string, vals []*adaptor.KMSAuthorizeRespItem) (string, error) {
	var tmp map[string]interface{}
	if err := json.Unmarshal([]byte(body), &tmp); err != nil {
		return body, err
	}
	for _, v := range vals {
		if err := appendToMap(tmp, v.Name, v.Value, 0); err != nil {
			return body, err
		}
	}
	b, err := json.MarshalIndent(tmp, "", "  ")
	if err != nil {
		return body, err
	}
	return string(b), nil
}

// a.b.c will be treat as recursion
func appendToMap(m map[string]interface{}, name string, val string, depth int) error {
	if depth > 10 {
		return fmt.Errorf("recursion too deep(>10)")
	}
	if name == "" {
		return fmt.Errorf("missing name for value: %s depth:%d", val, depth)
	}

	index := strings.Index(name, ".")
	if index < 0 { // single filed
		if _, ok := m[name]; ok {
			logger.Logger.Warnf("append existing field %s when add child", name)
		}
		m[name] = val
	} else { // sub field
		childName := name[:index]
		f := m[childName]
		if f == nil { // add a sub {}
			f = map[string]interface{}{}
			m[childName] = f
		}
		if sub, ok := f.(map[string]interface{}); ok {
			subName := name[index+1:]
			if err := appendToMap(sub, subName, val, depth+1); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("refer non object field %s when add child", childName)
		}
	}

	return nil
}

// QueryGrantedAPIKey query granted api key for owner
func QueryGrantedAPIKey(owner, service, authType, keyUUID string) string {
	key, err := QueryGrantedAPIKeyWithError(owner, service, keyUUID, authType)
	if err != nil {
		logger.Logger.Errorf("QueryGrantedAPIKey error: owner=%s service=%s, authType=%s err=%s",
			owner, service, authType, authType, err.Error())
	}
	return key
}

// QueryGrantedAPIKeyWithError query granted api key for owner with error handling
func QueryGrantedAPIKeyWithError(owner, service, authType, keyUUID string) (string, error) {
	switch {
	case RequireAPIKey(authType):
		if keyUUID == "" {
			if service == "" {
				return keyUUID, fmt.Errorf(" authType '%s' require apikey but missing service", authType)
			}

			// Query key from service.
			// TODO: query from key grant
			if op := adaptor.GetKMSOper(); op != nil {
				req := &adaptor.ListKMSCustomerKeyReq{
					Service:  service,
					Page:     1,
					PageSize: 1,
					Owner:    "",                //NOTE: keep empty
					Active:   rule.ActiveEnable, // NOTE: list active key only
				}
				resp, err := op.ListCustomerAPIKey(context.Background(), req)
				if err != nil {
					return "", err
				}
				if len(resp.Keys) == 0 {
					return "", errcode.ErrServiceWithoutKeys.NewError()
				}
				keyUUID = resp.Keys[0].ID
			}
		}
	}
	return keyUUID, nil
}

// GetAsAuthType convert auth type if valid
func GetAsAuthType(authType string) Enum {
	if AuthTypeEnum.Verify(authType) {
		return Enum(authType)
	}
	return AuthNone
}

// GetAuthDepartmentID get department id from header
func GetAuthDepartmentID(header http.Header) string {
	return header.Get(consts.HeaderDepartmentID)
}

// GetAuthUserID get userID from http header
func GetAuthUserID(header http.Header) string {
	if s := header.Get(consts.HeaderRequestFromInner); s != "" { // from inner => system
		return consts.SystemName
	}
	if s := header.Get(consts.HeaderUserID); s != "" {
		return s
	}
	return ""
}

// GetAuthOwner get ownerID from http header
func GetAuthOwner(header http.Header) string {
	if s := header.Get(consts.HeaderRequestFromInner); s != "" { // from inner => system
		return consts.SystemName
	}
	if s := header.Get(consts.HeaderAccessKeyID); s != "" {
		return consts.HeaderAccessKeyIDPrefix + s
	}
	return GetAuthUserID(header)
}

// GetAuthOwnerName get ownerName from http header
func GetAuthOwnerName(header http.Header) string {
	if s := header.Get(consts.HeaderUserName); s != "" {
		return s
	}
	return ""
}

// IsSuperManager verify if a http requester is super manager
func IsSuperManager(header http.Header) bool {
	roles := header.Get(consts.HeaderRole)
	return strings.Contains(roles, "super")
}

//------------------------------------------------------------------------------

func getOriginRequestID(header http.Header) string {
	// from web, gen suffix
	if id := header.Get(consts.HeaderXRequestID); id != "" {
		return id + hash.GenID("-req")
	}
	// from inner, use it directly
	if id := header.Get(consts.HeaderRequestID); id != "" {
		return id
	}
	// gen new one
	return hash.GenID("req")
}

func copySystemHeader(orginHeader http.Header, reqHeader *http.Header) {
	copyHeader(orginHeader, reqHeader, consts.HeaderUserID)
	copyHeader(orginHeader, reqHeader, consts.HeaderUserName)
	copyHeader(orginHeader, reqHeader, consts.HeaderDepartmentID)
	copyHeader(orginHeader, reqHeader, consts.HeaderRole)
	copyHeader(orginHeader, reqHeader, consts.HeaderTenantID)

}

func removeSystemHeader(reqHeader *http.Header) {
	delete(*reqHeader, consts.HeaderUserID)
	delete(*reqHeader, consts.HeaderUserName)
	delete(*reqHeader, consts.HeaderDepartmentID)
	delete(*reqHeader, consts.HeaderRole)
	delete(*reqHeader, consts.HeaderTenantID)
	delete(*reqHeader, consts.HeaderAccessKeyID)
	delete(*reqHeader, consts.HeaderRequestFromInner)
}

func copyHeader(orginHeader http.Header, reqHeader *http.Header, name string) {
	if h := orginHeader.Get(name); h != "" {
		reqHeader.Set(name, h)
	}
}
