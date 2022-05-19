package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"
	"unsafe"

	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/internal/models/mysql"
	myredis "github.com/quanxiang-cloud/polyapi/internal/models/redis"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/apiprovider"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/enums"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/polyhost"
	"github.com/quanxiang-cloud/polyapi/pkg/business/app"
	"github.com/quanxiang-cloud/polyapi/pkg/business/rule"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/action"
	"github.com/quanxiang-cloud/polyapi/pkg/core/auth"
	"github.com/quanxiang-cloud/polyapi/pkg/core/docview"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"
	"github.com/quanxiang-cloud/polyapi/pkg/core/swagger"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/encoding"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/hash"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/httputil"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/xsvc"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"gorm.io/gorm"
)

// RawAPI represents the raw api operator
type RawAPI interface {
	RegSwagger(c context.Context, req *RegReq) (*RegResp, error)
	Del(c context.Context, req *DelReq) (*DelResp, error)
	// InnerDel(c context.Context, req *DelReq) (*DelResp, error)
	List(c context.Context, req *RawListReq) (*RawListResp, error)
	ListInService(c context.Context, req *ListInServiceReq) (*ListInServiceResp, error)
	Active(c context.Context, req *ActiveReq) (*ActiveResp, error)
	Query(c context.Context, req *QueryReq) (*QueryResp, error)
	Search(c context.Context, req *SearchRawReq) (*SearchRawResp, error)

	QuerySwagger(ctx context.Context, req *QueryRawSwaggerReq) (*QueryRawSwaggerResp, error)

	apiprovider.Provider
}

// RegReq RegReq
type RegReq struct {
	Service   string `json:"-" binding:"-"`
	Owner     string `json:"-" binding:"-"`
	OwnerName string `json:"-" binding:"-"`

	Namespace  string `json:"namespace" binding:"max=384"`
	Swagger    string `json:"swagger" binding:"max=4096000"`
	Version    string `json:"version" binding:"max=32"`
	SingleName string `json:"singleName" binding:"max=64"` // api name when has only 1 api in swagger

	Schema   string `json:"schema" binding:"max=64"`
	Host     string `json:"host" binding:"max=64"`
	AuthType string `json:"authType" binding:"max=64"` //system, none, signature

	// auto create leaf namespace
	AutoCreateNamespaceTitle string `json:"autoCreateNamespaceTitle" binding:"max=64"`
}

// RegID RegID
type RegID struct {
	Path string `json:"path"`
}

// RegResp RegResp
type RegResp struct {
	List []RegID `json:"list"`
}

// DelReq DelReq
type DelReq = adaptor.DelReq

// DelResp DelResp
type DelResp = adaptor.DelResp

// QueryReq QueryReq
type QueryReq = adaptor.QueryRawAPIReq

// QueryResp QueryResp
type QueryResp = adaptor.QueryRawAPIResp

// QueryInBatchesReq QueryInBatchesReq
type QueryInBatchesReq = adaptor.QueryRawAPIInBatchesReq

// QueryInBatchesResp QueryInBatchesResp
type QueryInBatchesResp = adaptor.QueryRawAPIInBatchesResp

var instRaw RawAPI

// CreateRaw create a raw API operater
func CreateRaw(conf *config.Config) (RawAPI, error) {
	if instRaw != nil {
		return instRaw, nil
	}

	if _, err := CreateAPIStatOper(conf); err != nil { // init adaptor
		return nil, err
	}

	db, err := createMysqlConn(conf)
	if err != nil {
		return nil, err
	}
	redisCache, err := createRedisConn(conf)
	if err != nil {
		return nil, err
	}

	p := &raw{
		conf:     conf,
		db:       db,
		rawRepo:  mysql.NewRawAPIRepo(),
		redisAPI: redisCache,
	}

	if err := apiprovider.RegistAPIProvider(p); err != nil {
		return nil, err
	}
	adaptor.SetRawAPIOper(p)
	instRaw = p

	return p, nil

}

type raw struct {
	conf     *config.Config
	db       *gorm.DB
	rawRepo  models.RawAPIRepo
	redisAPI models.RedisCache
}

func uploadFileName(owner string) string {
	now := time.Now()
	return fmt.Sprintf("polyapi/%s_%s_swagger.json", owner, now.Format("20060102_150405"))
}

// RegSwagger RegSwagger
func (r *raw) RegSwagger(c context.Context, req *RegReq) (*RegResp, error) {
	if err := rule.CheckCharSet(req.Swagger); err != nil {
		return nil, err
	}

	uploadedURL := ""
	var err error
	if fs := adaptor.GetFileServerOper(); fs != nil {
		fn := uploadFileName(req.OwnerName)
		uploadedURL, err = fs.UploadFile(c, fn, unsafeStringBytes(req.Swagger))
		if err != nil {
			logger.Logger.Error("UploadSwagger fail", fn, req.Namespace, err.Error(), req.Swagger)
		} else {
			logger.Logger.Infof("UploadSwagger ok url=%s %s,%s,%s", uploadedURL, req.Owner, req.OwnerName, req.Namespace)
		}
	} else {
		logger.Logger.Info("UploadSwagger unsupport", req.Owner, req.Namespace, req.Swagger)
	}

	ns, _ := apipath.Split(req.Service)
	if req.Namespace != "" { // use namespace in request firstly, else use namespace of service
		ns = req.Namespace
	}

	// check namespace
	if oper := adaptor.GetNamespaceOper(); oper != nil {
		if nsObj, err := oper.Check(c, ns, req.Owner, OpAddRawAPI); err != nil {
			// try auto create leaf namespace in inner API
			if req.AutoCreateNamespaceTitle != "" && req.Owner == consts.SystemName {
				p, n := apipath.Split(ns)
				nsReq := &adaptor.CreateNsReq{
					Owner:             req.Owner,
					OwnerName:         req.OwnerName,
					IgnoreAccessCheck: true,
					Namespace:         p,
					Name:              n,
					Title:             req.AutoCreateNamespaceTitle,
					Desc:              req.AutoCreateNamespaceTitle,
				}
				if _, e := oper.Create(c, nsReq); e == nil {
					err = nil
				}
			}
			if err != nil {
				return nil, err
			}
		} else {
			// update Title if changed
			if req.AutoCreateNamespaceTitle != "" && req.Owner == consts.SystemName &&
				nsObj.Title != req.AutoCreateNamespaceTitle {
				nsReq := &UpdateNsReq{
					NamespacePath: ns,
					Title:         req.AutoCreateNamespaceTitle,
				}
				if _, err := oper.Update(c, nsReq); err != nil {
					logger.Logger.PutError(err, "ns.Update", "ns", ns)
				}
			}
		}
	}

	cfg := &swagger.APIServiceConfig{
		Host:       req.Host,
		Schema:     req.Schema,
		AuthType:   req.AuthType,
		Service:    req.Service,
		Namespace:  ns,
		StoredURL:  uploadedURL,
		SingleName: req.SingleName,
		Version:    req.Version,
	}

	if req.Service != "" {
		if oper := adaptor.GetServiceOper(); oper != nil {
			svs, err := oper.Check(c, req.Service, req.Owner, OpQuery)
			if err != nil {
				return nil, err
			}

			// use service config firstly
			cfg.Host = svs.Host
			cfg.Schema = svs.Schema
			cfg.AuthType = svs.AuthType
		}
	} else {
		if auth.RequireAPIKey(cfg.AuthType) {
			return nil, errcode.ErrIsolateAuthType.NewError()
		}
	}

	// check auth type
	if err := auth.ValidateAuthType(cfg.AuthType); err != nil {
		return nil, err
	}

	list, err := swagger.ParseSwagger(unsafeStringBytes(req.Swagger), cfg)
	if err != nil {
		logger.Logger.Error("Swagger parse fail", uploadedURL, err.Error())
		return nil, error2.NewErrorWithString(error2.ErrParams, err.Error())
	}

	// TODO: host health check?

	resp := &RegResp{List: make([]RegID, 0, len(list))}
	for _, api := range list {
		api.Owner = req.Owner
		api.OwnerName = req.OwnerName
		api.Namespace = ns
		api.Service = req.Service
		api.Access = 0
		api.Active = 1
		api.Valid = rule.Valid

		err = r.createRawWithHashConflict(api)
		if err != nil {
			return nil, errcode.ErrCreateExistsRaw.NewError()
		}

		resp.List = append(resp.List, RegID{Path: apipath.Join(api.Namespace, api.Name)})
	}
	if len(resp.List) > 0 {
		r.redisAPI.DeleteCache(ns, myredis.CacheRawList, false) // cache operation
	}
	return resp, nil
}

func (r *raw) createRawWithHashConflict(api *adaptor.RawAPIFull) error {
	var err error
	originName := api.Name
	for i := 0; i < hash.MaxHashConflict; i++ { // avoid hash conflict
		if originName == "" { // set name as explicit hash-random-name if not set outside
			api.Name = hash.HShortID(0, i, api.Namespace, api.Path, api.Action, api.Method, api.Version)
		}
		id := hash.Default("raw", i, api.Namespace, api.Name)

		api.Name = apipath.GenerateAPIName(api.Name, r.APIType())
		api.ID = id
		api.Content.ID = id
		apiPath := apipath.Join(api.Namespace, api.Name)
		api.Doc.FmtInOut.SetAccessURL(apiPath)

		if err = r.rawRepo.Create(r.db, api); err == nil {
			r.redisAPI.DeleteCache(apiPath, myredis.CacheRaw, false) // cache operation
			break
		}
	}
	return err
}

// Del delete the raw API
func (r *raw) Del(c context.Context, req *DelReq) (*DelResp, error) {
	paths := make([]string, 0, len(req.Names))
	for _, name := range req.Names {
		path := apipath.Join(req.NamespacePath, name)
		if _, err := r.check(c, nil, path, req.Owner, OpDelete); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	if err := checkPoly(c, paths); err != nil {
		return nil, err
	}

	return r.InnerDel(c, req)
}

func (r *raw) InnerDel(c context.Context, req *DelReq) (*DelResp, error) {
	err := r.rawRepo.Del(r.db, req.NamespacePath, req.Names)
	if err != nil {
		return nil, err
	}

	for _, name := range req.Names {
		r.redisAPI.DeleteCache(apipath.Join(req.NamespacePath, name), myredis.CacheRaw, false)
	}
	r.redisAPI.DeleteCache(req.NamespacePath, myredis.CacheRawList, false)

	return &DelResp{}, nil
}

func checkPoly(ctx context.Context, paths []string) error {
	if op := adaptor.GetRawPolyOper(); op != nil {
		resp, err := op.QueryByRawAPI(ctx, &adaptor.QueryByRawAPIReq{
			RawAPI: paths,
		})
		if err != nil {
			return err
		}

		if resp.Total != 0 {
			var set = make(map[string]bool)
			var ret = make([]string, 0, resp.Total)
			for _, v := range resp.List {
				if _, ok := set[v.RawAPI]; !ok {
					ret = append(ret, v.RawAPI)
				}
			}
			return errcode.ErrExistPoly.FmtError(ret)
		}
	}
	return nil
}

func (r *raw) check(c context.Context, item *models.RawAPICore, apiPath string, owner string, op Operation) (*models.RawAPICore, error) {
	if item == nil {
		q, err := r.query(apiPath)
		switch {
		case op == OpCreate && err == nil:
			return nil, errcode.ErrCreateExistsRaw.NewError()
		case op == OpCreate && err != nil:
			return nil, nil
		case err != nil:
			logger.Logger.Debugf("rawapi.check error op=%v path=%v err=%s", op, apiPath, err.Error())
			msg := fmt.Sprintf("invalid raw api: %s", apiPath)
			return nil, error2.NewErrorWithString(error2.ErrParams, msg)
		default:
			item = q
		}
	}

	if err := rule.CheckValid(item.Valid, op, enums.ObjectRaw); err != nil {
		return nil, err
	}

	if err := rule.ValidateActive(item.Active, op, enums.ObjectRaw); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *raw) query(apiPath string) (*models.RawAPICore, error) {
	var rawAPI *models.RawAPICore
	var err error
	rawAPI, err = r.redisAPI.QueryRaw(apiPath) // cache operation
	if err != nil {
		ns, name := apipath.Split(apiPath)
		rawAPI, err = r.rawRepo.Get(r.db, ns, name)
		if err != nil {
			return nil, err
		}
		r.redisAPI.PutCache(rawAPI, apiPath, myredis.CacheRaw) // cache operation
	}

	if err != nil || rawAPI.ID == "" {
		return nil, errNotFound
	}
	return rawAPI, err
}

func (r *raw) queryDoc(req *QueryReq) (*models.RawAPIDoc, error) {
	var doc *models.RawAPIDoc
	var err error
	doc, err = r.redisAPI.QueryRawDoc(req.APIPath) // cache operation
	if err != nil {
		ns, name := apipath.Split(req.APIPath)
		doc, err = r.rawRepo.GetDoc(r.db, ns, name)
		if err != nil {
			return nil, err
		}
		r.redisAPI.PutCache(doc, req.APIPath, myredis.CacheRawDoc) // cache operation
	}

	if err != nil || doc.ID == "" {
		return nil, errNotFound
	}
	return doc, err
}

// Query query the raw API info
func (r *raw) Query(c context.Context, req *QueryReq) (*QueryResp, error) {
	rawAPI, err := r.query(req.APIPath)
	if err != nil {
		return nil, err
	}

	resp := &QueryResp{
		ID:        rawAPI.ID,
		Content:   rawAPI.Content,
		URL:       rawAPI.URL,
		Method:    rawAPI.Method,
		Namespace: rawAPI.Namespace,
		Owner:     rawAPI.Owner,
		OwnerName: rawAPI.OwnerName,
		Service:   rawAPI.Service,
		Schema:    rawAPI.Schema,
		Host:      rawAPI.Host,
		Active:    rawAPI.Active,
		Valid:     rawAPI.Valid,
		Name:      rawAPI.Name,
		Title:     rawAPI.Title,
		UpdateAt:  rawAPI.UpdateAt,
	}
	return resp, nil
}

func (r *raw) QueryInBatches(c context.Context, req *QueryInBatchesReq) (*QueryInBatchesResp, error) {
	path := make([][2]string, 0, len(req.APIPathList))
	for _, v := range req.APIPathList {
		namespace, name := apipath.Split(v)
		path = append(path, [2]string{
			namespace,
			name,
		})
	}
	data, err := r.rawRepo.GetInBatches(r.db, path)
	if err != nil {
		return nil, err
	}
	list := serializeList(data.List)
	return &QueryInBatchesResp{
		List: list,
	}, nil
}

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

// RawListNode RawListNode
type RawListNode = adaptor.RawListNode

// RawListReq RawListReq
type RawListReq = adaptor.RawListReq

// RawListResp RawListResp
type RawListResp = adaptor.RawListResp

// List list
func (r *raw) List(c context.Context, req *RawListReq) (*RawListResp, error) {
	list, total, err := r.list(req.NamespacePath, "", req.Active, req.Page, req.PageSize)
	return &RawListResp{
		Total: total,
		Page:  req.Page,
		List:  list,
	}, err
}

// ListInServiceReq ListInServiceReq
type ListInServiceReq = adaptor.ListInServiceReq

// ListInServiceResp ListInServiceResp
type ListInServiceResp = adaptor.ListInServiceResp

// ListInService ListInservice
func (r *raw) ListInService(c context.Context, req *ListInServiceReq) (*ListInServiceResp, error) {
	list, total, err := r.list("", req.ServicePath, req.Active, req.Page, req.PageSize)
	return &ListInServiceResp{
		Total: total,
		Page:  req.Page,
		List:  list,
	}, err
}

func (r *raw) list(namespace, service string, active, page, pageSize int) ([]*RawListNode, int, error) {
	cache := namespace != "" && pageSize <= 0
	var data *models.RawAPIList
	var err error
	if cache {
		data, err = r.redisAPI.QueryRawList(namespace)
	}
	if !cache || err != nil {
		data, err = r.rawRepo.List(r.db, namespace, service, active, page, pageSize)
		if err != nil {
			return nil, 0, err
		}
		if cache {
			r.redisAPI.PutCache(data, namespace, myredis.CacheRawList)
		}
	}
	// active filter when pageSize <= 0
	if pageSize <= 0 {
		r.activeFilter(data, active)
	}

	list := serializeList(data.List)
	return list, int(data.Total), nil
}

func (r *raw) activeFilter(src *models.RawAPIList, active int) {
	if active >= 0 {
		uactive := uint(active)
		list := make([]*models.RawAPICore, 0, len(src.List))
		for _, v := range src.List {
			if v.Active == uactive {
				list = append(list, v)
			}
		}
		src.List = list
	}
}

// ListRawByPrefixPathReq ListRawByPrefixPathReq
type ListRawByPrefixPathReq = adaptor.ListRawByPrefixPathReq

// ListRawByPrefixPathResp ListRawByPrefixPathResp
type ListRawByPrefixPathResp = adaptor.ListRawByPrefixPathResp

func (r *raw) ListByPrefixPath(c context.Context, req *ListRawByPrefixPathReq) (*ListRawByPrefixPathResp, error) {
	data, total, err := r.rawRepo.ListByPrefixPath(r.db, req.NamespacePath, req.Active, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &ListRawByPrefixPathResp{
		List:  data,
		Total: total,
		Page:  req.Page,
	}, nil
}

// ActiveReq ActiveReq
type ActiveReq struct {
	APIPath string
	Active  uint `json:"active"`
}

// ActiveResp ActiveResp
type ActiveResp struct {
	FullPath string `json:"fullPath"`
	Active   uint   `json:"active"`
}

// Active Active
func (r *raw) Active(c context.Context, req *ActiveReq) (*ActiveResp, error) {
	namespace, name := apipath.Split(req.APIPath)
	err := r.rawRepo.UpdateActive(r.db, namespace, name, req.Active)

	r.redisAPI.DeleteCache(req.APIPath, myredis.CacheRaw, true) // cache operation

	return &ActiveResp{
		FullPath: req.APIPath,
		Active:   req.Active,
	}, err
}

// SearchRawReq SearchRawReq
type SearchRawReq struct {
	NamespacePath string `json:"-"`
	Name          string `json:"name"`
	Title         string `json:"title"`
	Active        int    `json:"active"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
	WithSub       bool   `json:"withSub"`
}

// SearchRawResp SearchRawResp
type SearchRawResp struct {
	Total int            `json:"total"`
	Page  int            `json:"page"`
	List  []*RawListNode `json:"list"`
}

// Search Search
func (r *raw) Search(c context.Context, req *SearchRawReq) (*SearchRawResp, error) {
	data, err := r.rawRepo.Search(r.db, req.NamespacePath, req.Name,
		req.Title, req.Active, req.Page, req.PageSize, req.WithSub)
	if err != nil {
		return nil, err
	}
	list := serializeList(data.List)
	return &SearchRawResp{
		Total: int(data.Total),
		Page:  req.Page,
		List:  list,
	}, nil
}

func serializeList(data []*models.RawAPICore) []*RawListNode {
	list := make([]*RawListNode, 0, len(data))
	for _, v := range data {
		list = append(list, serializeRaw(v))
	}
	return list
}

func serializeRaw(data *models.RawAPICore) *RawListNode {
	namespacePath := apipath.Join(data.Namespace, data.Name)
	return &RawListNode{
		ID:         data.ID,
		Owner:      data.Owner,
		OwnerName:  data.OwnerName,
		Name:       data.Name,
		Title:      data.Title,
		Desc:       data.Desc,
		FullPath:   namespacePath,
		URL:        data.URL,
		Version:    data.Version,
		Method:     data.Method,
		Action:     data.Action,
		Active:     data.Active,
		Valid:      data.Valid,
		CreateAt:   data.CreateAt,
		UpdateAt:   data.UpdateAt,
		URI:        data.Path,
		AccessPath: app.MakeRequestPath(namespacePath),
	}
}

// InnerUpdateRawInBatchReq InnerUpdateRawInBatchReq
type InnerUpdateRawInBatchReq = adaptor.InnerUpdateRawInBatchReq

// InnerUpdateRawInBatchResp InnerUpdateRawInBatchResp
type InnerUpdateRawInBatchResp = adaptor.InnerUpdateRawInBatchResp

func (r *raw) InnerUpdateRawInBatch(ctx context.Context, req *InnerUpdateRawInBatchReq) (*InnerUpdateRawInBatchResp, error) {
	list, err := r.ListInService(ctx, &adaptor.ListInServiceReq{
		ServicePath: req.Service,
		Active:      -1,
		Page:        1,
		PageSize:    -1,
	})
	if err != nil {
		return nil, err
	}

	err = r.rawRepo.UpdateInBatch(r.db, req.Namespace, req.Service, req.Host, req.Schema, req.AuthType)
	if err != nil {
		return nil, err
	}

	fullPath := make([]string, 0, len(list.List))
	for _, raw := range list.List {
		fullPath = append(fullPath, raw.FullPath)
	}
	r.redisAPI.DeleteCacheInBatch(fullPath, myredis.CacheRaw)
	r.redisAPI.DeleteCache(req.Namespace, myredis.CacheRawList, false)

	return &adaptor.InnerUpdateRawInBatchResp{}, nil
}

// RawValidReq RawValidReq
type RawValidReq = struct {
	APIPath string
	Valid   uint `json:"valid"`
}

// RawValidResp RawValidResp
type RawValidResp struct {
	FullPath string `json:"fullPath"`
	Valid    uint   `json:"valid"`
}

// Valid Valid
func (r *raw) Valid(ctx context.Context, req *RawValidReq) (*RawValidResp, error) {
	_, err := r.ValidInBatches(ctx, &RawValidInBatchesReq{
		APIPath: []string{req.APIPath},
		Valid:   req.Valid,
	})
	if err != nil {
		return nil, err
	}

	// update poly if valid = 0
	if req.Valid == rule.Invalid {
		queryPolyResp, err := adaptor.GetRawPolyOper().QueryByRawAPI(ctx, &adaptor.QueryByRawAPIReq{
			RawAPI: []string{req.APIPath},
		})
		if err != nil {
			return nil, err
		}

		polyList := make([]string, 0, queryPolyResp.Total)
		for _, v := range queryPolyResp.List {
			polyList = append(polyList, v.PolyAPI)
		}

		_, err = adaptor.GetPolyOper().ValidInBatches(ctx, &adaptor.PolyValidInBatchesReq{
			APIPath: polyList,
			Valid:   req.Valid,
		})
		if err != nil {
			return nil, err
		}
	}

	return &RawValidResp{
		FullPath: req.APIPath,
		Valid:    req.Valid,
	}, nil
}

// RawValidInBatchesReq RawValidInBatchesReq
type RawValidInBatchesReq = adaptor.RawValidInBatchesReq

// RawValidInBatchesResp RawValidInBatchesResp
type RawValidInBatchesResp = adaptor.RawValidInBatchesResp

func (r *raw) ValidInBatches(ctx context.Context, req *RawValidInBatchesReq) (*RawValidInBatchesResp, error) {
	path := make([][2]string, 0, len(req.APIPath))
	for _, v := range req.APIPath {
		namespace, name := apipath.Split(v)
		path = append(path, [2]string{namespace, name})
	}

	err := r.rawRepo.UpdateValid(r.db, path, req.Valid)
	if err != nil {
		return nil, err
	}

	for _, v := range req.APIPath {
		r.redisAPI.DeleteCache(v, myredis.CacheRaw, true)
	}
	return &RawValidInBatchesResp{}, nil
}

// RawValidByPrefixPathReq RawValidByPrefixPathReq
type RawValidByPrefixPathReq = adaptor.RawValidByPrefixPathReq

// RawValidByPrefixPathResp RawValidByPrefixPathResp
type RawValidByPrefixPathResp = adaptor.RawValidByPrefixPathResp

// ValidByPrefixPath ValidByPrefixPath
func (r *raw) ValidByPrefixPath(ctx context.Context, req *RawValidByPrefixPathReq) (*RawValidByPrefixPathResp, error) {
	if err := r.rawRepo.UpdateValidByPrefixPath(r.db, req.NamespacePath, req.Valid); err != nil {
		return nil, err
	}
	err := r.redisAPI.DeletePatternCache(apipath.FormatPrefix(req.NamespacePath), myredis.CacheRaw, true)
	return &RawValidByPrefixPathResp{}, err
}

// InnerImportRawReq InnerImportRawReq
type InnerImportRawReq = adaptor.InnerImportRawReq

// InnerImportRawResp InnerImportRawResp
type InnerImportRawResp = adaptor.InnerImportRawResp

// InnerImport InnerImport
func (r *raw) InnerImport(ctx context.Context, req *InnerImportRawReq) (*InnerImportRawResp, error) {
	for _, v := range req.List {
		v.ID = hash.GenID("raw")
	}
	r.rawRepo.CreateInBatches(r.db, req.List)
	return &InnerImportRawResp{}, nil
}

// InnerDelRawByPrefixPathReq InnerDelRawByPrefixPathReq
type InnerDelRawByPrefixPathReq = adaptor.InnerDelRawByPrefixPathReq

// InnerDelRawByPrefixPathResp InnerDelRawByPrefixPathResp
type InnerDelRawByPrefixPathResp = adaptor.InnerDelRawByPrefixPathResp

// InnerDelRawByPrefixPath InnerDelRawByPrefixPath
func (r *raw) InnerDelByPrefixPath(ctx context.Context, req *InnerDelRawByPrefixPathReq) (*InnerDelRawByPrefixPathResp, error) {
	err := r.rawRepo.DelByPrefixPath(r.db, req.NamespacePath)
	if err != nil {
		return nil, err
	}
	r.redisAPI.DeletePatternCache(apipath.FormatPrefix(req.NamespacePath), myredis.CacheRaw, true)
	return &InnerDelRawByPrefixPathResp{}, nil
}

//------------------------------------------------------------------------------
// API provider

// Request call the API directly
func (r *raw) Request(c context.Context, req *apiprovider.RequestReq) (*apiprovider.RequestResp, error) {
	if gate.APIStatIsBlocked(c, req.APIPath, true) { // shortly blocked api
		return nil, errcode.ErrGateBlockedAPI.NewError()
	}

	api, err := r.query(req.APIPath)
	if err != nil {
		return nil, err
	}
	if api.ID == "" {
		return nil, errNotFound
	}

	// check method
	if !httputil.AllowMethod(api.Method, req.Method) {
		return &apiprovider.RequestResp{
			APIPath:    req.APIPath,
			Status:     http.StatusText(http.StatusNotFound),
			StatusCode: http.StatusNotFound,
		}, nil
	}
	// check state
	if _, err := r.check(c, api, req.APIPath, req.Owner, OpRequest); err != nil {
		return nil, err
	}

	args := &expr.RequestArgs{
		Name:         api.Name,
		URL:          api.URL,
		Body:         req.Body,
		Header:       req.Header,
		EncodingIn:   api.Content.EncodingIn,
		EncodingPoly: expr.PolyEncoding,
		Action:       api.Content.Action,
	}

	authType := api.AuthType
	services := api.Service
	if req.APIService != "" {
		svs, err := action.QueryService(req.APIService)
		if err != nil {
			return nil, err
		}
		args.URL, authType = action.GetServiceInfo(svs, api.Path)
	}
	if req.APIServiceArgs != "" {
		xArgs, err := xsvc.Unmarshal(req.APIServiceArgs)
		if err != nil {
			return nil, err
		}
		args.URL, _ = action.GetServiceInfo(&adaptor.ServicesResp{
			Schema: xArgs.Schema,
			Host:   xArgs.Host,
			// Note: ignore authType in XPolyServiceArgs
			// AuthType: xArgs.AuthType,
			// AuthContent: xArgs.AuthContent,
		}, api.Path)
	}

	if err := api.Content.Consts.PrepareRequest(args); err != nil {
		return nil, err
	}

	keyID, err := auth.QueryGrantedAPIKeyWithError(req.Owner, services, api.AuthType, "")
	if err != nil {
		return nil, err
	}
	body, err := auth.AppendAuthWithError(keyID, authType, args.Header, args.Header, unsafeByteString(args.Body))
	if err != nil {
		return nil, err
	}
	args.Body = unsafeStringBytes(body)

	switch api.Schema {
	case consts.SchemaHTTP, consts.SchemaHTTPS:
		t := api.Content
		start := time.Now()
		ret, resp, err := httputil.HTTPRequest(args.URL, t.Method, unsafeByteString(args.Body), args.Header, req.Owner)
		if err != nil {
			return nil, err
		}
		dur := time.Now().Sub(start)
		gate.APIStatAddTimeStat(c, req.APIPath, true, resp.StatusCode, dur)
		body, err := encoding.ConvertEncoding(api.Content.EncodingOut, ret, expr.PolyEncoding, true)
		if err != nil {
			logger.Logger.PutError(err, "raw.request")
			return nil, err
		}
		r := &apiprovider.RequestResp{
			APIPath:    req.APIPath,
			Response:   json.RawMessage(body),
			Header:     resp.Header,
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
		}
		return r, nil
	default:
		return nil, fmt.Errorf("unsupported schema %s, api path: %s", api.Schema, req.APIPath)
	}
}

// QueryAPIDoc return API doc by path & name
func (r *raw) QueryDoc(c context.Context, req *apiprovider.QueryDocReq) (*apiprovider.QueryDocResp, error) {
	q := &QueryReq{
		APIPath: req.APIPath,
	}

	rawAPI, err := r.queryDoc(q)
	if err != nil {
		return nil, err
	}

	doc, err := docview.GetAPIDocView(rawAPI.Doc, polyhost.GetSchemaHost(), req.DocType, req.TitleFirst)
	if err != nil {
		return nil, err
	}
	resp := &apiprovider.QueryDocResp{
		DocType: req.DocType,
		APIPath: req.APIPath,
		Doc:     doc,
		Title:   rawAPI.Title,
	}
	return resp, nil
}

// APIType is the type of api provider
func (r *raw) APIType() string {
	return "r"
}

// QueryRawSwaggerReq QueryRawSwaggerReq
type QueryRawSwaggerReq = adaptor.QueryRawSwaggerReq

// QueryRawSwaggerResp QueryRawSwaggerResp
type QueryRawSwaggerResp = adaptor.QueryRawSwaggerResp

// QuerySwagger query raw swagger
func (r *raw) QuerySwagger(ctx context.Context, req *QueryRawSwaggerReq) (*QueryRawSwaggerResp, error) {
	path := make([][2]string, 0, len(req.APIPath))
	for _, v := range req.APIPath {
		ns, name := apipath.Split(v)
		path = append(path, [2]string{ns, name})
	}

	raws, err := r.rawRepo.GetDocInBatches(r.db, path)
	if err != nil {
		return nil, err
	}

	var swags = &swagger.SwagDoc{}
	swags.Paths = make(map[string]swagger.SwagMethods)
	swags.Host = polyhost.GetHost()
	swags.Info.Title = "auto generated"
	swags.Schemes = []string{polyhost.GetSchema()}

	for _, v := range raws {
		swag := &swagger.SwagDoc{}
		err := json.Unmarshal(v.Doc.Swagger, swag)
		if err != nil {
			return nil, err
		}

		path := apipath.Join(v.Namespace, v.Name)
		for _, v := range swag.Paths {
			swags.Paths[path] = v
		}
	}

	b, err := json.Marshal(swags)
	if err != nil {
		return nil, err
	}

	return &QueryRawSwaggerResp{
		Swagger: b,
	}, nil
}
