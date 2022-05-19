package jsvm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	pkgtimestamp "github.com/quanxiang-cloud/polyapi/pkg/basic/timestamp"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/encoding"
)

const debugScript = false

// vmSetup setup some predefined object to VM
func (vm *JsVM) vmSetup() {
	// js defined apis:
	// pdFromJson(jsonStr) jsObj
	// pdToJson(jsObj) jsonStr
	// pdToJsonP(jsObj) jsonStr
	// pdMergeObjs(jsObj...) obj
	// pdFiltObject(input, config, filtFunc)
	// pdSelect(cond, yes, no) yes|no
	// pdToJsobj("encoding", strData) obj
	// sel(cond, yes, no) yes|no

	vm.Set(consts.PDHttpRequest, httpRequest)                  // pdHttpRequest(url, method, data, header) (string, error)
	vm.Set(consts.PDAddHTTPHeader, addHTTPHeader)              // pdAddHttpHeader(header, k, v)
	vm.Set(consts.PDUpdateReferPath, updateReferPath)          // pdUpdateReferPath(header, namespacePath)
	vm.Set(consts.PDFromXML, encoding.FromXML)                 // pdFromXml(xmlStr)(obj, error)
	vm.Set(consts.PDToXML, encoding.ToXML)                     // pdToXml(obj)(xmlStr, error)
	vm.Set(consts.PDFromYAML, encoding.FromYAML)               // pdFromYaml(yamlStr)(obj, error)
	vm.Set(consts.PDToYAML, encoding.ToYAML)                   // pdToYaml(obj)(yamlStr, error)
	vm.Set(consts.PDQueryUser, vm.QueryUser)                   // pdQueryUser(realDo)uid
	vm.Set(consts.PDCreateNS, vm.createNamespace)              // pdCreateNS(parent,name,title) string
	vm.Set(consts.PDAppendAuth, vm.appendAuth)                 // pdAppendAuth(keyID, authType, header, body) string
	vm.Set(consts.PDQueryGrantedAPIKey, vm.queryGrantedAPIKey) // pdQueryGrantedAPIKey(owner, service, authType, keyID) string
	vm.Set(consts.PDNewHTTPHeader, vm.NewHTTPHeader)           // pdNewHttpHeader() http.Header

	vm.Set(consts.PDtimestamp, timestamp) // timestamp() string , UTC timestamp of "YYYY-MM-DDThh:mm:ssZ"
	vm.Set(consts.PDformat, fmt.Sprintf)  // format(fmt, ...) string

	if debugScript {
		file := `./script/predef.js`
		v, err := vm.vmSetupScriptFile(file)
		if err != nil {
			fmt.Println("vmSetupScriptFile fail", file, v, err)
			panic(err)
		}
	} else {
		v, err := vm.RunString(predefJsScript)
		if err != nil {
			fmt.Println("vmSetup fail", v, err)
			panic(err)
		}
	}
}

// timestamp get current timestamp of UTC
func timestamp(format string) string {
	return pkgtimestamp.Timestamp(format)
}

//addHTTPHeader add a header parameter
func addHTTPHeader(h http.Header, k, v string) {
	h.Add(k, v)
}

var polyRequestRoot = fmt.Sprintf(consts.APIRequestPath, "")

// UpdateReferPath update a refer path
func updateReferPath(h http.Header, apiNamespacePath string) bool {
	if refer := h.Get(consts.HeaderRefer); refer != "" {
		if idx := strings.Index(refer, polyRequestRoot); idx > 0 {
			prefix := refer[:idx+len(polyRequestRoot)-1]
			n := prefix + apiNamespacePath
			h.Set(consts.HeaderRefer, n)
			return true
		}
	}
	return false
}
