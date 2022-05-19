package jsvm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"unsafe"

	"github.com/quanxiang-cloud/polyapi/pkg/core/action"
	"github.com/quanxiang-cloud/polyapi/pkg/core/auth"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/httputil"

	"github.com/quanxiang-cloud/cabin/logger"
)

// unsafeByteString convert []byte to string without copy
// the origin []byte **MUST NOT** accessed after that
func unsafeByteString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// unsafeStringBytes return GoString's buffer slice
// ** NEVER modify returned []byte **
func unsafeStringBytes(s string) []byte {
	var bh reflect.SliceHeader
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// httpRequest was the JS api pdHttpRequest implication
func httpRequest(reqURL, method, data string, header http.Header, owner string) (string, error) {
	body, _, err := httputil.HTTPRequest(reqURL, method, data, header, owner)
	if err != nil {
		logger.Logger.PutError(err, "httpRequest")
	}

	// if body is not json, format it
	var d json.RawMessage
	if err := json.Unmarshal(unsafeStringBytes(body), &d); err != nil {
		logger.Logger.PutError(err, "httpRequest", "body", body)
		body = fmt.Sprintf(`{"$":%q}`, body)
	}
	return body, err
}

func (vm *JsVM) vmSetupScriptFile(path string) (sysVMValue, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return emptySysVMValue, err
	}
	script := unsafeByteString(b)
	return vm.RunString(script)
}

func (vm *JsVM) createNamespace(parent, name, title string) string {
	fullPath, err := action.CreateNamespace(parent, name, title,
		vm.QueryUser(true), vm.QueryUserName(), false)
	if err != nil {
		return fmt.Sprintf("error: %s %s", err.Error(), fullPath)
	}
	return fullPath
}

func (vm *JsVM) appendAuth(keyID, authType string, header http.Header, body string) string {
	return auth.AppendAuth(keyID, authType, vm.header, header, body)
}

func (vm *JsVM) queryGrantedAPIKey(owner, service, keyID, authType string) string {
	return auth.QueryGrantedAPIKey(owner, service, keyID, authType)
}
