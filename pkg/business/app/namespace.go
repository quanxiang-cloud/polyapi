package app

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/core/action"
	"github.com/quanxiang-cloud/polyapi/pkg/core/auth"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"

	"github.com/quanxiang-cloud/cabin/logger"
)

// consts define
const (
	RootOfAPPs  = "/system/app"
	RootOfFaas  = "/system/faas"
	ServiceForm = "form"
)

// NamespaceDefine define a namespace
type NamespaceDefine struct {
	Path   string
	Parent string
	Name   string
	Title  string
}

// enum define
var (
	PathEnum = enumset.New(nil) // API authorize method

	PathFaasGlobal       = PathEnum.MustReg("faas.global") // faas in global
	PathRoot             = PathEnum.MustReg("root")
	PathPoly             = PathEnum.MustReg("poly")
	PathRawRoot          = PathEnum.MustReg("raw.root")
	PathFaas             = PathEnum.MustReg("faas") // faas in app
	PathRawInner         = PathEnum.MustReg("raw.inner")
	PathRaw3Party        = PathEnum.MustReg("raw.3party")
	PathRaw3PartyDefault = PathEnum.MustReg("raw.3party.default")

	PathServiceForm = PathEnum.MustReg("inner.form")

	/*
		APP path tree:
		/system/app/<APP_ID>
			/poly
			/raw
				/faas
					/[func-name]
						- func_version
				/customer
					/default
				/inner
					/form
						/table-id
						 - opation_api
					/[flow and-so-on]...
	*/
	predefinedPath = map[enumset.Enum]*nsDefine{
		PathFaasGlobal: {
			NamespaceDefine: NamespaceDefine{
				Name:  "",
				Title: "函数服务",
				Path:  RootOfFaas,
			},
			paths: []enumset.Enum{},
		},
		PathRoot: {
			NamespaceDefine: NamespaceDefine{
				Name:  "",
				Title: "应用",
				Path:  RootOfAPPs,
			},
			paths: []enumset.Enum{},
		},
		PathPoly: {
			NamespaceDefine: NamespaceDefine{
				Name:  "poly",
				Title: "API编排",
			},
			paths: []enumset.Enum{},
		},
		PathRawRoot: {
			NamespaceDefine: NamespaceDefine{
				Name:  "raw",
				Title: "原生API",
			},
			paths: []enumset.Enum{},
		},
		PathFaas: {
			NamespaceDefine: NamespaceDefine{
				Name:  "faas",
				Title: "函数服务",
			},
			paths: []enumset.Enum{PathRawRoot},
		},
		PathRaw3Party: {
			NamespaceDefine: NamespaceDefine{
				Name:  "customer",
				Title: "代理第三方API",
			},
			paths: []enumset.Enum{PathRawRoot},
		},
		PathRaw3PartyDefault: {
			NamespaceDefine: NamespaceDefine{
				Name:  "default",
				Title: "默认分组",
			},
			paths: []enumset.Enum{PathRawRoot, PathRaw3Party},
		},
		PathRawInner: {
			NamespaceDefine: NamespaceDefine{
				Name:  "inner",
				Title: "平台API",
			},
			paths: []enumset.Enum{PathRawRoot},
		},
		PathServiceForm: {
			NamespaceDefine: NamespaceDefine{
				Name:  ServiceForm,
				Title: "表单模型API",
			},
			paths: []enumset.Enum{PathRawRoot, PathRawInner},
		},
	}
	/*
		/system/app/appX
		/system/app/appX/poly
		/system/app/appX/raw
		/system/app/appX/raw/faas
		/system/app/appX/raw/customer
		/system/app/appX/raw/customer/default
		/system/app/appX/raw/inner
		/system/app/appX/raw/inner/form
	*/
	createAppPathOrder = []enumset.Enum{
		PathRoot,
		PathPoly,
		PathRawRoot,
		PathFaas,
		PathRaw3Party,
		PathRaw3PartyDefault,
		PathRawInner,
		PathServiceForm,
	}
)

// RootPath get root path of app
func RootPath(appID string) string {
	return fmt.Sprintf("%s/%s", appRoot, appID)
}

// PathInfo get info of path type
func PathInfo(typ string) NamespaceDefine {
	if p, ok := predefinedPath[enumset.Enum(typ)]; ok {
		return p.NamespaceDefine
	}

	return NamespaceDefine{}
}

// GetCreateAppPaths get the namespace info to create
func GetCreateAppPaths(appID string) []NamespaceDefine {
	if appID != "" {
		ret := make([]NamespaceDefine, len(createAppPathOrder))
		for i, v := range createAppPathOrder {
			p := &ret[i]
			p.Path = Path(appID, v.String())
			p.Parent, p.Name = apipath.Split(p.Path)
			if f, ok := predefinedPath[v]; ok {
				p.Title = f.Title
			}
		}
		return ret
	}
	return nil
}

// Path get sub path of app
func Path(appID string, typ string) string {
	e := enumset.Enum(typ)
	switch e {
	case PathRoot:
		return RootPath(appID)
	case PathFaasGlobal: // /system/faas
		if p, ok := predefinedPath[e]; ok {
			return p.Path
		}
	default:
		if p, ok := predefinedPath[e]; ok {
			return fmt.Sprintf("%s/%s", RootPath(appID), p.Path)
		}
	}

	return ""
}

//------------------------------------------------------------------------------
var (
	appRoot          = predefinedPath[PathRoot].Path
	relPathRaw3Party = "" // raw/customer
	relPathPoly      = "" // poly
)

func init() {
	for _, v := range predefinedPath {
		if v.Path == "" {
			parent := v.mergeRelativePath()
			v.Path = path.Join(parent, v.Name)
		}
	}
	relPathRaw3Party = predefinedPath[PathRaw3Party].Path
	relPathPoly = predefinedPath[PathPoly].Path
}

type nsDefine struct {
	NamespaceDefine
	paths []enumset.Enum //path of names
}

func (n *nsDefine) mergeRelativePath() string {
	buf := bytes.NewBuffer(nil)
	for _, v := range n.paths {
		p, ok := predefinedPath[v]
		if !ok {
			panic(v)
		}
		buf.WriteString(p.Name)
		buf.WriteByte('/')
	}
	return buf.String()
}

// /system/app/:appID/:sub/...
var (
	exprApp    = fmt.Sprintf(`(?sm:^%s/(?P<appID>[\w-]+)(?:/(?P<sub>[-\w\./]*))?$)`, RootOfAPPs)
	regAPPPath = regexp.MustCompile(exprApp)
)

// SplitAsAppPath split a full namespace as app/sub path
func SplitAsAppPath(fullPath string) (string, string, error) {
	elems := regAPPPath.FindAllStringSubmatch(fullPath, 1)
	if len(elems) > 0 {
		appID, sub := elems[0][1], elems[0][2]
		return appID, sub, nil
	}
	return "", "", fmt.Errorf("invalid app path [%s]", fullPath)
}

func validAppSubPath(op adaptor.Operation, sub string) bool {
	switch op {
	case adaptor.OpAddPolyAPI:
		return strings.HasPrefix(sub, relPathPoly)
	case adaptor.OpAddRawAPI, adaptor.OpAddService:
		return strings.HasPrefix(sub, relPathRaw3Party)
	case adaptor.OpAddSub:
		return strings.HasPrefix(sub, relPathPoly) ||
			strings.HasPrefix(sub, relPathRaw3Party)
	}
	return true
}

// IsSystemOwner check if a owner is system
func isSystemOwner(owner string) bool {
	return owner == "" || owner == consts.SystemName
}

// ValidateNamespace verify namespace when some operation
func ValidateNamespace(owner string, op adaptor.Operation, fullPath string) error {
	if isSystemOwner(owner) {
		return nil
	}

	formated := apipath.Format(fullPath)
	_, sub, err := SplitAsAppPath(formated)
	if err != nil {
		return err
	}

	if !validAppSubPath(op, sub) {
		return fmt.Errorf("invalid app sub path [%s]", fullPath)
	}

	return nil
}

// ValidateServicePath verify the service path
func ValidateServicePath(owner string, op adaptor.Operation, servicePath string) error {
	if op != adaptor.OpCreate {
		return nil
	}

	ns, svsName := apipath.Split(servicePath)
	if err := ValidateNamespace(owner, adaptor.OpAddRawAPI, ns); err != nil {
		return err
	}

	if !isSystemOwner(owner) {
		_, nsName := apipath.Split(ns)
		if nsName != svsName {
			return fmt.Errorf("invalid app service path [%s]", servicePath)
		}
	}

	return nil
}

// ValidateAPIPath verify if refer api path from different APP
func ValidateAPIPath(owner, callerPath, referPath string) error {
	if isSystemOwner(owner) {
		return nil
	}
	apCaller, _, err := SplitAsAppPath(callerPath)
	if err != nil {
		return err
	}
	appRefer, _, err := SplitAsAppPath(referPath)
	if err != nil {
		return err
	}
	if apCaller != appRefer {
		return errcode.ErrUseDifferentApp.FmtError(referPath)
	}
	return nil
}

// ValidateAppAccessPermit verify if user has app acccess right
func ValidateAppAccessPermit(fullNamespace string, header http.Header, writterAPI bool) error {
	if fullNamespace == "" { //no namespace API
		return nil
	}
	if oper := adaptor.GetAppCenterServerOper(); oper != nil {
		userID := auth.GetAuthUserID(header)
		if isSystemOwner(userID) {
			return nil
		}
		app, _, _ := SplitAsAppPath(fullNamespace)
		if app == "" {
			if !writterAPI && strings.HasPrefix(fullNamespace, RootOfFaas) {
				return nil
			}
			logger.Logger.Warnf("[no-app-access-right] ns=%s app=%s", fullNamespace, app)
			//return nil
			return errcode.ErrNoNamespacePermit.NewError()
		}
		isSuper := auth.IsSuperManager(header)
		checkAdmin := writterAPI
		depID := auth.GetAuthDepartmentID(header)
		ok, err := oper.Check(context.Background(), userID, depID, app, isSuper, checkAdmin)
		if err != nil {
			logger.Logger.Warnf("[no-app-access-right] ns=%s app=%s user=%s super=%v checkAdmin=%v err=%v",
				fullNamespace, app, userID, isSuper, checkAdmin, err)
			//return nil
			return err
		}
		if !ok {
			logger.Logger.Warnf("[no-app-access-right] ns=%s app=%s user=%s super=%v checkAdmin=%v OK=%v",
				fullNamespace, app, userID, isSuper, checkAdmin, ok)
			//return nil
			return errcode.ErrNoNamespacePermit.NewError()
		}
	}
	return nil
}

// InitAppPath create namespaces within app
func InitAppPath(appID string, owner, ownerName string, header http.Header) error {
	if err := ValidateAppAccessPermit(Path(appID, PathRaw3Party.String()), header, true); err != nil {
		return err
	}

	var fails strings.Builder
	create := func(d NamespaceDefine) {
		if _, err := action.CreateNamespace(d.Parent, d.Name, d.Title, owner, ownerName, true); err != nil {
			fails.WriteString(fmt.Sprintf("%s\n", err.Error()))
		}
	}
	paths := GetCreateAppPaths(appID)
	for _, v := range paths {
		create(v)
	}
	if fails.Len() > 0 {
		return errcode.ErrInitAppPath.FmtError(fails.String())
	}

	return nil
}

// MakeRequestPath return path of request
func MakeRequestPath(namespacePath string) string {
	if len(namespacePath) > 0 && namespacePath[0] == '/' {
		namespacePath = namespacePath[1:]
	}
	return fmt.Sprintf(consts.APIRequestPath, namespacePath)
}
