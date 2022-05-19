package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	//NOTE: init arrange.PolyDocGennerator
	_ "github.com/quanxiang-cloud/polyapi/polycore/pkg/core/polydoc"

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
	"github.com/quanxiang-cloud/polyapi/pkg/core/auth"
	"github.com/quanxiang-cloud/polyapi/pkg/core/docview"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"
	"github.com/quanxiang-cloud/polyapi/pkg/core/swagger"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/factory"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/hash"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/httputil"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/arrange"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/exprx"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/jsvm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"gorm.io/gorm"
)

// PolyAPI represents the service of poly api operator
type PolyAPI interface {
	Create(c context.Context, req *PolyCreateReq) (*PolyCreateResp, error)
	Delete(c context.Context, req *PolyDeleteReq) (*PolyDeleteResp, error)
	// InnerDelete(c context.Context, req *PolyDeleteReq) (*PolyDeleteResp, error)
	UpdateArrange(c context.Context, req *PolyUpdateArrangeReq) (*PolyUpdateArrangeResp, error)
	UpdateScript(c context.Context, req *PolyUpdateScriptReq) (*PolyUpdateScriptResp, error)
	GetArrange(c context.Context, req *PolyGetArrangeReq) (*PolyGetArrangeResp, error)
	GetScript(c context.Context, req *PolyGetScriptReq) (*PolyGetScriptResp, error)
	Build(c context.Context, req *PolyBuildReq) (*PolyBuildResp, error)
	ShowEnum(c context.Context, req *PolyEnumReq) (*PolyEnumResp, error)
	List(c context.Context, req *PolyListReq) (*PolyListResp, error)
	Active(c context.Context, req *PolyActiveReq) (*PolyActiveResp, error)
	Search(c context.Context, req *SearchPolyReq) (*SearchPolyResp, error)

	QuerySwagger(ctx context.Context, req *QueryPolySwaggerReq) (*QueryPolySwaggerResp, error)

	apiprovider.Provider
}

var instPoly PolyAPI

// CreatePoly create poly api operator
func CreatePoly(conf *config.Config) (PolyAPI, error) {
	if instPoly != nil {
		return instPoly, nil
	}

	db, err := createMysqlConn(conf)
	if err != nil {
		return nil, err
	}
	redisCache, err := createRedisConn(conf)
	if err != nil {
		return nil, err
	}

	logger.Logger.Infof("db connect ok: host=%s db=%s user=%s",
		conf.Mysql.Host, conf.Mysql.DB, conf.Mysql.User)

	p := &polyAPI{
		conf:     conf,
		db:       db,
		dbRepo:   mysql.NewPolyDbAPI(),
		redisAPI: redisCache,
	}

	if err := apiprovider.RegistAPIProvider(p); err != nil {
		return nil, err
	}
	adaptor.SetPolyOper(p)
	adaptor.SetEvalerOper(jsvm.NewEvalerCreator()) // set Evaler
	instPoly = p

	return p, nil
}

//------------------------------------------------------------------------------

type polyAPI struct {
	conf     *config.Config
	db       *gorm.DB
	dbRepo   models.PolyAPIRepo
	redisAPI models.RedisCache
}

// PolyCreateReq PolyCreateReq
type PolyCreateReq struct {
	Namespace string `json:"-" binding:"-"`
	Owner     string `json:"-" binding:"-"` // owner
	OwnerName string `json:"-" binding:"-"` // owner

	Name  string `json:"name" binding:"required,max=128"` // name
	Title string `json:"title" binding:"max=64"`
	Desc  string `json:"desc" binding:"max=4096"`
	//Access          []string `json:"access" binding:"max=64"`
	Method          string `json:"method" binding:"max=64"`
	TemplateAPIPath string `json:"templateAPIPath" binding:"max=512"` // template poly API path
}

// PolyCreateResp PolyCreateResp
type PolyCreateResp struct {
	ID      string `json:"id"`      // uuid
	APIPath string `json:"apiPath"` //
	//Access    []string  `json:"access"`
	Active    uint   `json:"active"`
	Method    string `json:"method"`
	Owner     string `json:"owner"`
	OwnerName string `json:"ownerName"`
	Arrange   string `json:"arrange"`  // arrange info
	CreateAt  int64  `json:"createAt"` // create time
	UpdateAt  int64  `json:"updateAt"` // update time
}

// Create create a new poly api
func (p *polyAPI) Create(c context.Context, req *PolyCreateReq) (*PolyCreateResp, error) {
	if err := rule.CheckCharSet(req.Title, req.Desc); err != nil {
		return nil, err
	}
	if err := rule.CheckDescLength(req.Desc); err != nil {
		return nil, err
	}

	if err := rule.ValidateName(req.Name, rule.MaxNameLength-2, false); err != nil {
		return nil, err
	}

	// check namespace
	if oper := adaptor.GetNamespaceOper(); oper != nil {
		if _, err := oper.Check(c, req.Namespace, req.Owner, OpAddPolyAPI); err != nil {
			return nil, err
		}
	}
	if _, err := p.check(c, nil, apipath.Join(req.Namespace, req.Name), req.Owner, OpCreate); err != nil {
		return nil, err
	}

	// access, err := permission.ParsePermits(req.Access)
	// if err != nil {
	// 	return nil, err
	// }

	// FIXME: check desc length
	d := &models.PolyAPIArrange{
		ID:        "?",
		Namespace: req.Namespace,
		Name:      apipath.GenerateAPIName(req.Name, p.APIType()),
		Title:     req.Title,
		Desc:      req.Desc,
		Owner:     req.Owner,
		OwnerName: req.OwnerName,
		//Access:    access.Uint(),
		Method: req.Method,
		Active: rule.ActiveDisable,
		Valid:  rule.Valid,
	}
	if t := req.TemplateAPIPath; t != "" {
		if err := app.ValidateAPIPath(req.Owner, req.Namespace, req.TemplateAPIPath); err != nil {
			return nil, err
		}
		ret, err := p.GetArrange(c, &PolyGetArrangeReq{APIPath: t})
		if err != nil {
			return nil, err
		}
		d.Arrange = ret.Arrange
	}
	info := arrange.Arrange{
		Info: &arrange.APIInfo{
			Namespace: d.Namespace,
			Name:      d.Name,
		},
	}
	t, err := arrange.InitArrange(d.Arrange, &info)
	if err != nil {
		return nil, err
	}
	d.Arrange = t

	for i := 0; i < hash.MaxHashConflict; i++ { // avoid hash conflict
		d.ID = hash.Default("poly", i, d.Namespace, d.Name)
		if err = p.dbRepo.Create(p.db, d); err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	p.redisAPI.DeleteCache(req.Namespace, myredis.CachePolyList, false)

	resp := &PolyCreateResp{
		ID:        d.ID,
		APIPath:   apipath.Join(d.Namespace, d.Name),
		Owner:     d.Owner,
		OwnerName: d.OwnerName,
		//Access:    permission.ToPermitList(d.Access),
		Method:   d.Method,
		Arrange:  d.Arrange,
		CreateAt: d.CreateAt,
		UpdateAt: d.UpdateAt,
	}
	return resp, nil
}

//------------------------------------------------------------------------------

// PolyDeleteReq PolyDeleteReq
type PolyDeleteReq = adaptor.PolyDeleteReq

// PolyDeleteResp PolyDeleteResp
type PolyDeleteResp = adaptor.PolyDeleteResp

// Delete delete a poly api
func (p *polyAPI) Delete(c context.Context, req *PolyDeleteReq) (*PolyDeleteResp, error) {
	for _, name := range req.Names {
		if _, err := p.check(c, nil, apipath.Join(req.NamespacePath, name), req.Owner, OpDelete); err != nil {
			return nil, err
		}
	}
	return p.InnerDelete(c, req)
}

func (p *polyAPI) InnerDelete(c context.Context, req *PolyDeleteReq) (*PolyDeleteResp, error) {
	err := p.dbRepo.Delete(p.db, req.NamespacePath, req.Names)
	if err != nil {
		return nil, err
	}

	// cache operation
	p.redisAPI.DeleteCache(req.NamespacePath, myredis.CachePolyList, false)
	for _, name := range req.Names {
		p.redisAPI.DeleteCache(apipath.Join(req.NamespacePath, name), myredis.CachePoly, false)
	}

	if op := adaptor.GetRawPolyOper(); op != nil {
		paths := make([]string, 0, len(req.Names))
		for _, v := range req.Names {
			paths = append(paths, apipath.Join(req.NamespacePath, v))
		}

		_, err := op.DeleteByPolyAPIInBatches(c, &adaptor.DeleteByPolyAPIInBatchesReq{
			PolyAPIPath: paths,
		})
		if err != nil {
			return nil, err
		}
	}

	resp := &PolyDeleteResp{}
	return resp, nil
}

//------------------------------------------------------------------------------

// PolyUpdateArrangeReq PolyUpdateArrangeReq
type PolyUpdateArrangeReq struct {
	Owner   string `json:"-"`
	APIPath string `json:"-" binding:"-"` //
	Title   string `json:"title" binding:"max=64"`
	Desc    string `json:"desc" binding:"max=4096"`
	Arrange string `json:"arrange" binding:"required"` // arrange info
}

// PolyUpdateArrangeResp PolyUpdateArrangeResp
type PolyUpdateArrangeResp struct {
	APIPath string `json:"apiPath"`
}

// UpdateArrange update a poly api arrange json
func (p *polyAPI) UpdateArrange(c context.Context, req *PolyUpdateArrangeReq) (*PolyUpdateArrangeResp, error) {
	if err := rule.CheckCharSet(req.Title, req.Desc); err != nil {
		return nil, err
	}
	if err := rule.CheckDescLength(req.Desc); err != nil {
		return nil, err
	}

	// check state
	if _, err := p.check(c, nil, req.APIPath, req.Owner, OpEdit); err != nil {
		return nil, err
	}

	ns, name := apipath.Split(req.APIPath)
	u := &models.PolyAPIArrange{
		Namespace: ns,
		Name:      name,
		Arrange:   req.Arrange,
		Title:     req.Title,
		Desc:      req.Desc,
	}
	if err := p.dbRepo.UpdateArrange(p.db, u); err != nil {
		return nil, err
	}

	resp := &PolyUpdateArrangeResp{
		APIPath: req.APIPath,
	}
	return resp, nil
}

//------------------------------------------------------------------------------

// PolyUpdateScriptReq PolyUpdateScriptReq
type PolyUpdateScriptReq struct {
	APIPath string `json:"-" binding:"-"`             //
	Script  string `json:"script" binding:"required"` // arrange info
	Doc     string `json:"doc"`                       // api doc
}

// PolyUpdateScriptResp PolyUpdateScriptResp
type PolyUpdateScriptResp struct {
	APIPath string `json:"apiPath"`
}

// UpdateScript save the build result
func (p *polyAPI) UpdateScript(c context.Context, req *PolyUpdateScriptReq) (*PolyUpdateScriptResp, error) {
	ns, name := apipath.Split(req.APIPath)
	u := &models.PolyBuildResult{
		Namespace: ns,
		Name:      name,
		Script:    req.Script,
		Doc:       req.Doc,
	}
	if err := p.dbRepo.UpdateScript(p.db, u); err != nil {
		return nil, err
	}

	p.redisAPI.DeleteCache(req.APIPath, myredis.CachePoly, false) // cache operation

	resp := &PolyUpdateScriptResp{
		APIPath: req.APIPath,
	}
	return resp, nil
}

//------------------------------------------------------------------------------

// PolyGetArrangeReq PolyGetArrangeReq
type PolyGetArrangeReq struct {
	APIPath string `json:"-" binding:"-"` //
}

// PolyGetArrangeResp PolyGetArrangeResp
type PolyGetArrangeResp struct {
	ID        string `json:"id"`        // uuid
	Namespace string `json:"namespace"` //
	Name      string `json:"name"`      //
	Title     string `json:"title"`     //
	Desc      string `json:"desc"`      //
	Owner     string `json:"owner"`
	OwnerName string `json:"ownerName"`
	//Access    []string       `json:"access"` //
	Active   uint   `json:"active"`
	Method   string `json:"method"`   //
	Arrange  string `json:"arrange"`  // arrange info
	CreateAt int64  `json:"createAt"` // create time
	UpdateAt int64  `json:"updateAt"` // update time
	BuildAt  int64  `json:"buildAt"`  // build time
}

func (p *polyAPI) query(apiPath string) (*models.PolyAPIArrange, error) {
	ns, name := apipath.Split(apiPath)
	item, err := p.dbRepo.GetArrange(p.db, ns, name)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// GetArrange return a poly api arrange json
func (p *polyAPI) GetArrange(c context.Context, req *PolyGetArrangeReq) (*PolyGetArrangeResp, error) {
	item, err := p.query(req.APIPath)
	if err != nil {
		return nil, err
	}

	resp := &PolyGetArrangeResp{
		ID:        item.ID,
		Name:      item.Name,
		Namespace: item.Namespace,
		Title:     item.Title,
		Desc:      item.Desc,
		//Access:    permission.ToPermitList(item.Access),
		Active:    item.Active,
		Method:    item.Method,
		Owner:     item.Owner,
		OwnerName: item.OwnerName,
		Arrange:   item.Arrange,
		CreateAt:  item.CreateAt,
		UpdateAt:  item.UpdateAt,
		BuildAt:   item.BuildAt,
	}
	return resp, nil
}

//------------------------------------------------------------------------------

// PolyGetScriptReq PolyGetScriptReq
type PolyGetScriptReq struct {
	APIPath string `json:"-" binding:"-"` //
}

// PolyGetScriptResp PolyGetScriptResp
type PolyGetScriptResp struct {
	ID      string `json:"id"`      // uuid
	APIPath string `json:"apiPath"` //
	Script  string `json:"-"`       // js script
	BuildAt int64  `json:"buildAt"` // build time
}

func (p *polyAPI) checkScript(c context.Context, item *models.PolyAPIScript, apiPath string, owner string, op Operation) (*models.PolyAPIScript, error) {
	if !op.RequireScriptReady() {
		return nil, errNotSupport
	}
	if item == nil {
		q, err := p.queryScript(apiPath)
		switch {
		case err != nil || q.Script == "":
			if err != nil {
				logger.Logger.Errorf("[polyapi.CheckScript(%s) error]: %s", apiPath, err.Error())
			}
			return nil, errcode.ErrPolyNotBuild.FmtError()
		default:
			item = q
		}
	}

	if err := rule.CheckValid(item.Valid, op, enums.ObjectPoly); err != nil {
		return nil, err
	}

	if err := rule.ValidateActive(item.Active, op, enums.ObjectPoly); err != nil {
		switch {
		case op == OpRequest && item.Owner == owner:
			//NOTE: allow owner request the disabled poly api for testing
			// do nothing
		default:
			return nil, err
		}
	}

	return item, nil
}

func (p *polyAPI) check(c context.Context, item *models.PolyAPIArrange, apiPath string, owner string, op Operation) (*models.PolyAPIArrange, error) {
	if op.RequireScriptReady() {
		return nil, errNotSupport
	}
	if item == nil {
		q, err := p.query(apiPath)
		switch {
		case op == OpCreate && err == nil:
			return nil, errcode.ErrCreateExistsPoly.NewError()
		case op == OpCreate && err != nil:
			return nil, nil
		case err != nil:
			logger.Logger.Debugf("polyapi.check error op=%v path=%v err=%s", op, apiPath, err.Error())
			msg := fmt.Sprintf("invalid poly api: %s", apiPath)
			return nil, error2.NewErrorWithString(error2.ErrParams, msg)
		default:
			item = q
		}
	}

	if err := rule.ValidateActive(item.Active, op, enums.ObjectPoly); err != nil {
		switch {
		case op == OpRequest && item.Owner == owner:
			//NOTE: allow owner request the disabled poly api for testing
			// do nothing
		default:
			return nil, err
		}
	}

	return item, nil
}

func (p *polyAPI) queryScript(apiPath string) (*models.PolyAPIScript, error) {
	script, err := p.redisAPI.QueryPoly(apiPath) // cache operation
	if err != nil {
		ns, name := apipath.Split(apiPath)
		script, err = p.dbRepo.GetScript(p.db, ns, name)
		if err != nil {
			return nil, err
		}
		p.redisAPI.PutCache(script, apiPath, myredis.CachePoly) // cache operation
	}
	return script, err
}

// GetScript get the build result
func (p *polyAPI) GetScript(c context.Context, req *PolyGetScriptReq) (*PolyGetScriptResp, error) {
	script, err := p.queryScript(req.APIPath)
	if err != nil {
		return nil, err
	}

	resp := &PolyGetScriptResp{
		ID:      script.ID,
		APIPath: req.APIPath,
		Script:  script.Script,
		BuildAt: script.BuildAt,
	}
	return resp, nil
}

func (p *polyAPI) queryDoc(apiPath string) (*models.PolyAPIDoc, error) {
	var doc *models.PolyAPIDoc
	var err error
	doc, err = p.redisAPI.QueryPolyDoc(apiPath) // cache operation
	if err != nil {
		ns, name := apipath.Split(apiPath)
		doc, err = p.dbRepo.GetDoc(p.db, ns, name)
		if err != nil {
			return nil, err
		}
		p.redisAPI.PutCache(doc, apiPath, myredis.CachePolyDoc) // cache operation
	}

	if err != nil || doc.ID == "" {
		return nil, errNotFound
	}
	return doc, err
}

//------------------------------------------------------------------------------

// PolyBuildReq PolyBuildReq
type PolyBuildReq struct {
	Owner   string `json:"-"`
	APIPath string `json:"-" binding:"-"`               //
	Arrange string `json:"arrange"  binding:"required"` // arrange info
}

// PolyBuildResp PolyBuildResp
type PolyBuildResp struct {
	APIPath string `json:"apiPath"` //
}

// Build build a poly api arrange json to JS code
func (p *polyAPI) Build(c context.Context, req *PolyBuildReq) (*PolyBuildResp, error) {
	arrangeTxt := req.Arrange

	// check state
	item, err := p.check(c, nil, req.APIPath, req.Owner, OpBuild)
	if err != nil {
		return nil, err
	}

	if req.Arrange != "" {
		// save arrange JSON
		u := &PolyUpdateArrangeReq{
			APIPath: req.APIPath,
			Arrange: req.Arrange,
		}
		_, err := p.UpdateArrange(c, u)
		if err != nil {
			return nil, err
		}
	} else {
		arrangeTxt = item.Arrange
	}

	info := arrange.APIInfo{
		Namespace: item.Namespace,
		Name:      item.Name,
		Title:     item.Title,
		Desc:      item.Desc,
		Method:    item.Method,
	}

	// build script
	script, doc, rawList, err := arrange.BuildJsScript(&info, arrangeTxt, req.Owner)
	if err != nil {
		return nil, err
	}

	// verify apis from the same app
	for _, v := range rawList {
		if err := app.ValidateAPIPath(req.Owner, req.APIPath, v); err != nil {
			return nil, err
		}
	}

	// BUG: it crash when rawList is empty
	if len(rawList) > 0 {
		if oper := adaptor.GetRawPolyOper(); oper != nil {
			req := &adaptor.UpdateRawPolyReq{
				PolyAPI:    req.APIPath,
				RawAPIList: rawList,
			}
			if _, err := oper.UpdateRawPoly(c, req); err != nil {
				return nil, err
			}
		}
	}

	// save script
	save := &PolyUpdateScriptReq{
		APIPath: req.APIPath,
		Script:  script,
		Doc:     doc,
	}
	if _, err := p.UpdateScript(c, save); err != nil {
		return nil, err
	}

	resp := &PolyBuildResp{
		APIPath: req.APIPath,
	}
	return resp, nil
}

//------------------------------------------------------------------------------

// PolyEnumReq PolyEnumReq
type PolyEnumReq struct {
	Type   string `json:"-" binding:"-"`
	Sample bool   `json:"sample"`
}

// PolyEnumElem PolyEnumElem
type PolyEnumElem struct {
	Name   string      `json:"name"`
	View   string      `json:"view"`
	Sample interface{} `json:"sample"`
}

// PolyEnumResp PolyEnumResp
type PolyEnumResp struct {
	EnumType string         `json:"enumType"`
	List     []PolyEnumElem `json:"list"`
}

// ShowEnum request a poly api by Name and Parent
func (p *polyAPI) ShowEnum(c context.Context, req *PolyEnumReq) (*PolyEnumResp, error) {
	out := &PolyEnumResp{
		EnumType: req.Type,
	}

	switch req.Type {
	case expr.EnumNode.String():
		showEnum(req, arrange.NodeTypeEnum, arrange.GetFactory(), out)
	case expr.EnumValue.String():
		showEnum(req, exprx.ExprTypeEnum.EnumSet, expr.GetFactory(), out)
	case expr.EnumOper.String():
		showEnum(req, exprx.OpEnum.EnumSet, expr.GetFactory(), out)
	case expr.EnumCond.String():
		showEnum(req, exprx.CondEnum.EnumSet, expr.GetFactory(), out)
	case expr.EnumCmp.String():
		showEnum(req, exprx.CmpEnum.EnumSet, expr.GetFactory(), out)
	case expr.EnumIn.String():
		showEnum(req, exprx.ParaTypeEnum.EnumSet, expr.GetFactory(), out)
	case expr.EnumAuth.String():
		showEnum(req, auth.AuthTypeEnum, auth.GetFactory(), out)
	default:
		msg := fmt.Sprintf("unsupported enum <%s>, valid: %v", req.Type, expr.EnumTypesEnum.GetAll())
		return nil, error2.NewErrorWithString(error2.ErrParams, msg)
	}
	return out, nil
}

func showEnum(req *PolyEnumReq, s *enumset.EnumSet, factory *factory.FlexObjFactory, out *PolyEnumResp) {
	for _, v := range s.GetAll() {
		var d interface{}
		if req.Sample {
			d, _ = factory.CreateSample(v)
		}

		view, _ := exprx.ConvertOp(v)
		e := PolyEnumElem{
			Name:   v,
			View:   view,
			Sample: d,
		}
		out.List = append(out.List, e)
	}
}

// PolyListReq polyListReq
type PolyListReq = adaptor.PolyListReq

// PolyListResp PolyListResp
type PolyListResp = adaptor.PolyListResp

// PolyListNode node
type PolyListNode = adaptor.PolyListNode

// List List
func (p *polyAPI) List(c context.Context, req *PolyListReq) (*PolyListResp, error) {
	cache := req.PageSize <= 0
	var data *models.PolyAPIList
	var err error
	if cache {
		data, err = p.redisAPI.QueryPolyList(req.NamespacePath)
	}
	if !cache || err != nil {
		data, err = p.dbRepo.List(p.db, req.NamespacePath, req.Active, req.Page, req.PageSize)
		if err != nil {
			return nil, err
		}
		if cache {
			p.redisAPI.PutCache(data, req.NamespacePath, myredis.CachePolyList)
		}
	}

	// pageSize <= 0
	if cache {
		p.activeFilter(data, req.Active)
	}
	list := serializePolyList(data.List)
	return &PolyListResp{
		Total: int(data.Total),
		Page:  req.Page,
		List:  list,
	}, nil
}

func (p *polyAPI) activeFilter(src *models.PolyAPIList, active int) {
	if active >= 0 {
		uactive := uint(active)
		list := make([]*models.PolyAPIArrange, 0, len(src.List))
		for _, v := range src.List {
			if v.Active == uactive {
				list = append(list, v)
			}
		}
		src.List = list
	}
}

// PolyActiveReq PolyActiveReq
type PolyActiveReq struct {
	APIPath string
	Active  uint `json:"active"`
}

// PolyActiveResp PolyActiveResp
type PolyActiveResp struct {
	FullPath string `json:"fullPath"`
	Active   uint   `json:"active"`
}

func (p *polyAPI) Active(c context.Context, req *PolyActiveReq) (*PolyActiveResp, error) {
	if req.Active == rule.ActiveEnable {
		if _, err := p.checkScript(c, nil, req.APIPath, "", OpPublish); err != nil {
			return nil, err
		}
	}

	namespace, name := apipath.Split(req.APIPath)
	err := p.dbRepo.UpdateActive(p.db, namespace, name, req.Active)

	p.redisAPI.DeleteCache(req.APIPath, myredis.CachePoly, true) // cache operation

	return &PolyActiveResp{
		FullPath: req.APIPath,
		Active:   req.Active,
	}, err
}

// SearchPolyReq SearchPolyReq
type SearchPolyReq struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Active    int    `json:"active"`
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
	WithSub   bool   `json:"withSub"`
}

// SearchPolyResp SearchPolyResp
type SearchPolyResp struct {
	Total int             `json:"total"`
	Page  int             `json:"page"`
	List  []*PolyListNode `json:"list"`
}

// Search Search
func (p *polyAPI) Search(c context.Context, req *SearchPolyReq) (*SearchPolyResp, error) {
	data, err := p.dbRepo.Search(p.db, req.Namespace, req.Name, req.Title, req.Active, req.Page, req.PageSize, req.WithSub)
	if err != nil {
		return nil, err
	}
	list := serializePolyList(data.List)
	return &SearchPolyResp{
		Total: int(data.Total),
		Page:  req.Page,
		List:  list,
	}, nil
}

func serializePolyList(data []*models.PolyAPIArrange) []*PolyListNode {
	list := make([]*PolyListNode, 0, len(data))
	for _, v := range data {
		list = append(list, serializePoly(v))
	}
	return list
}

func serializePoly(data *models.PolyAPIArrange) *PolyListNode {
	namespacePath := apipath.Join(data.Namespace, data.Name)
	return &PolyListNode{
		ID:         data.ID,
		Owner:      data.Owner,
		OwnerName:  data.OwnerName,
		FullPath:   namespacePath,
		Method:     data.Method,
		Name:       data.Name,
		Title:      data.Title,
		Desc:       data.Desc,
		Active:     data.Active,
		Valid:      data.Valid,
		CreateAt:   data.CreateAt,
		UpdateAt:   data.UpdateAt,
		AccessPath: app.MakeRequestPath(namespacePath),
	}
}

// PolyValidReq PolyValidReq
type PolyValidReq struct {
	APIPath string `json:"-"`
	Valid   uint   `json:"valid"`
}

// PolyValidResp PolyValidResp
type PolyValidResp struct {
	FullPath string `json:"fullPath"`
	Valid    uint   `json:"valid"`
}

// Valid Valid
func (p *polyAPI) Valid(ctx context.Context, req *PolyValidReq) (*PolyValidResp, error) {
	if err := checkRawValid(ctx, req.APIPath, req.Valid); err != nil {
		return nil, err
	}

	_, err := p.ValidInBatches(ctx, &PolyValidInBatchesReq{
		APIPath: []string{
			req.APIPath,
		},
		Valid: req.Valid,
	})
	return &PolyValidResp{
		FullPath: req.APIPath,
		Valid:    req.Valid,
	}, err
}

func checkRawValid(ctx context.Context, polyPath string, valid uint) error {
	if valid == rule.Valid {
		polyRawList, err := adaptor.GetRawPolyOper().QueryByPolyAPI(ctx, &adaptor.QueryByPolyAPIReq{
			PolyAPI: polyPath,
		})
		if err != nil {
			return err
		}

		rawPathList := make([]string, 0, len(polyRawList.List))
		for _, rawPath := range polyRawList.List {
			rawPathList = append(rawPathList, rawPath.RawAPI)
		}

		rawList, err := adaptor.GetRawAPIOper().QueryInBatches(ctx, &adaptor.QueryRawAPIInBatchesReq{
			APIPathList: rawPathList,
		})
		if err != nil {
			return err
		}

		for _, raw := range rawList.List {
			if raw.Valid != valid {
				return errcode.ErrAPIInvalid.NewError()
			}
		}
	}
	return nil
}

// PolyValidInBatchesReq PolyValidInBatchesReq
type PolyValidInBatchesReq = adaptor.PolyValidInBatchesReq

// PolyValidInBatchesResp PolyValidInBatchesResp
type PolyValidInBatchesResp = adaptor.PolyValidInBatchesResp

// ValidInBatches ValidInBatches
func (p *polyAPI) ValidInBatches(ctx context.Context, req *PolyValidInBatchesReq) (*PolyValidInBatchesResp, error) {
	path := make([][2]string, 0, len(req.APIPath))
	for _, v := range req.APIPath {
		namespace, name := apipath.Split(v)
		path = append(path, [2]string{namespace, name})
	}

	err := p.dbRepo.UpdateValid(p.db, path, req.Valid)
	if err != nil {
		return nil, err
	}

	for _, v := range req.APIPath {
		p.redisAPI.DeleteCache(v, myredis.CachePoly, true)
	}
	return &PolyValidInBatchesResp{}, nil
}

// PolyValidByPrefixPathReq PolyValidByPrefixPathReq
type PolyValidByPrefixPathReq = adaptor.PolyValidByPrefixPathReq

// PolyValidByPrefixPathResp PolyValidByPrefixPathResp
type PolyValidByPrefixPathResp = adaptor.PolyValidByPrefixPathResp

// ValidByPrefixPath ValidByPrefixPath
func (p *polyAPI) ValidByPrefixPath(ctx context.Context, req *PolyValidByPrefixPathReq) (*PolyValidByPrefixPathResp, error) {
	if err := p.dbRepo.UpdateValidByPrefixPath(p.db, req.NamespacePath, req.Valid); err != nil {
		return nil, err
	}
	err := p.redisAPI.DeletePatternCache(apipath.FormatPrefix(req.NamespacePath), myredis.CachePoly, true)
	return &PolyValidByPrefixPathResp{}, err
}

// ListPolyByPrefixPathReq ListPolyByPrefixPathReq
type ListPolyByPrefixPathReq = adaptor.ListPolyByPrefixPathReq

// ListPolyByPrefixPathResp ListPolyByPrefixPathResp
type ListPolyByPrefixPathResp = adaptor.ListPolyByPrefixPathResp

// ListByPrefixPath ListByPrefixPath
func (p *polyAPI) ListByPrefixPath(ctx context.Context, req *ListPolyByPrefixPathReq) (*ListPolyByPrefixPathResp, error) {
	data, total, err := p.dbRepo.ListByPrefixPath(p.db, req.NamespacePath, req.Active, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	return &ListPolyByPrefixPathResp{
		List:  data,
		Total: total,
		Page:  req.Page,
	}, nil
}

// InnerImportPolyReq InnerImportPolyReq
type InnerImportPolyReq = adaptor.InnerImportPolyReq

// InnerImportPolyResp InnerImportPolyResp
type InnerImportPolyResp = adaptor.InnerImportPolyResp

func (p *polyAPI) InnerImport(ctx context.Context, req *InnerImportPolyReq) (*InnerImportPolyResp, error) {
	for _, v := range req.List {
		v.ID = hash.GenID("poly")
	}
	err := p.dbRepo.CreateInBatches(p.db, req.List)
	return &InnerImportPolyResp{}, err
}

// InnerDelPolyByPrefixPathReq InnerDelPolyByPrefixPathReq
type InnerDelPolyByPrefixPathReq = adaptor.InnerDelPolyByPrefixPathReq

// InnerDelPolyByPrefixPathResp InnerDelPolyByPrefixPathResp
type InnerDelPolyByPrefixPathResp = adaptor.InnerDelPolyByPrefixPathResp

// InnerDelByPrefixPath InnerDelByPrefixPath
func (p *polyAPI) InnerDelByPrefixPath(ctx context.Context, req *InnerDelPolyByPrefixPathReq) (*InnerDelPolyByPrefixPathResp, error) {
	err := p.dbRepo.DelByPrefixPath(p.db, req.NamespacePath)
	if err != nil {
		return nil, err
	}
	p.redisAPI.DeletePatternCache(apipath.FormatPrefix(req.NamespacePath), myredis.CachePoly, true)
	return &InnerDelPolyByPrefixPathResp{}, nil
}

//------------------------------------------------------------------------------
// API provider

// Request request a poly api by http
func (p *polyAPI) Request(c context.Context, req *apiprovider.RequestReq) (*apiprovider.RequestResp, error) {
	if gate.APIStatIsBlocked(c, req.APIPath, false) { // shortly blocked api
		return nil, errcode.ErrGateBlockedAPI.NewError()
	}

	// check state
	api, err := p.checkScript(c, nil, req.APIPath, req.Owner, OpRequest)
	if err != nil {
		return nil, err
	}

	if !httputil.AllowMethod(api.Method, req.Method) {
		return &apiprovider.RequestResp{
			APIPath:    req.APIPath,
			Status:     http.StatusText(http.StatusNotFound),
			StatusCode: http.StatusNotFound,
		}, nil
	}

	resp := &apiprovider.RequestResp{
		APIPath: req.APIPath,
	}
	start := time.Now()
	if jsRet, err := jsvm.RunJsString(api.Script, req.Body, req.Header); err == nil {
		resp.Header = http.Header{} // TODO:
		resp.Header.Set(consts.HeaderContentType, consts.MIMEJSON)
		resp.Response = json.RawMessage(jsRet)
		resp.StatusCode = http.StatusOK
	} else {
		return nil, err
	}
	dur := time.Now().Sub(start)
	gate.APIStatAddTimeStat(c, req.APIPath, false, http.StatusOK, dur)

	return resp, nil
}

// QueryDoc return API doc by path & name
func (p *polyAPI) QueryDoc(c context.Context, req *apiprovider.QueryDocReq) (*apiprovider.QueryDocResp, error) {
	api, err := p.queryDoc(req.APIPath)
	if err != nil {
		return nil, err
	}

	doc, err := docview.GetAPIDocView(api.Doc, polyhost.GetSchemaHost(), req.DocType, req.TitleFirst)
	if err != nil {
		return nil, err
	}
	resp := &apiprovider.QueryDocResp{
		DocType: req.DocType,
		APIPath: req.APIPath,
		Doc:     doc,
		Title:   api.Title,
	}
	return resp, nil
}

// APIType is the type of api provider
func (p *polyAPI) APIType() string {
	return "p"
}

// QueryPolySwaggerReq QueryPolySwaggerReq
type QueryPolySwaggerReq = adaptor.QueryPolySwaggerReq

// QueryPolySwaggerResp QueryPolySwaggerResp
type QueryPolySwaggerResp = adaptor.QueryPolySwaggerResp

// QuerySwagger query poly swagger
func (p *polyAPI) QuerySwagger(ctx context.Context, req *QueryPolySwaggerReq) (*QueryPolySwaggerResp, error) {
	path := make([][2]string, 0, len(req.APIPath))
	for _, v := range req.APIPath {
		ns, name := apipath.Split(v)
		path = append(path, [2]string{ns, name})
	}
	polys, err := p.dbRepo.GetDocInBatches(p.db, path)
	if err != nil {
		return nil, err
	}

	var swags = &swagger.SwagDoc{}
	swags.Paths = make(map[string]swagger.SwagMethods)
	swags.Host = polyhost.GetHost()
	swags.Info.Title = "auto generated"
	swags.Schemes = []string{polyhost.GetSchema()}

	for _, v := range polys {
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

	return &QueryPolySwaggerResp{
		Swagger: b,
	}, nil
}
