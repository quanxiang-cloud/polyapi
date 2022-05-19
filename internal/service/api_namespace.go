package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/internal/models/mysql"
	myredis "github.com/quanxiang-cloud/polyapi/internal/models/redis"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/enums"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/business/app"
	"github.com/quanxiang-cloud/polyapi/pkg/business/rule"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/daprevent"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/hash"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"gorm.io/gorm"
)

// NamespaceAPI represents the namespace api operator
type NamespaceAPI interface {
	Create(c context.Context, req *CreateNsReq) (*CreateNsResp, error)
	Delete(c context.Context, req *DeleteNsReq) (*DeleteNsResp, error)
	InnerDelete(c context.Context, req *DeleteNsReq) (*DeleteNsResp, error)
	Update(c context.Context, req *UpdateNsReq) (*UpdateNsResp, error)
	Active(c context.Context, req *ActiveNsReq) (*ActiveNsResp, error)
	List(c context.Context, req *ListNsReq) (*ListNsResp, error)
	Query(c context.Context, namespace string) (*NamespaceResp, error)
	Check(c context.Context, namespace string, owner string, op Operation) (*NamespaceResp, error)
	APPPath(c context.Context, req *AppPathReq) (*AppPathResp, error)
	InitAPPPath(c context.Context, req *InitAppPathReq) (*InitAppPathResp, error)
	Search(c context.Context, req *SearchNsReq) (*SearchNsResp, error)
	Tree(c context.Context, req *NsTreeReq) (*NsTreeResp, error)
}

// CreateNamespaceOper create a namespace API operater
func CreateNamespaceOper(conf *config.Config) (NamespaceAPI, error) {
	db, err := createMysqlConn(conf)
	if err != nil {
		return nil, err
	}
	redisCache, err := createRedisConn(conf)
	if err != nil {
		return nil, err
	}

	p := &namespace{
		conf:     conf,
		db:       db,
		nsDB:     mysql.NewAPINamespaceRepo(),
		redisAPI: redisCache,
	}

	adaptor.SetNamespaceOper(p)
	return p, nil
}

type namespace struct {
	conf     *config.Config
	db       *gorm.DB
	redisAPI models.RedisCache
	nsDB     models.APINamespaceRepo
}

// CreateNsReq export
type CreateNsReq = adaptor.CreateNsReq

// CreateNsResp export
type CreateNsResp = adaptor.NamespaceResp

// NamespaceResp export
type NamespaceResp = adaptor.NamespaceResp

// Create Create
func (ns *namespace) Create(c context.Context, req *CreateNsReq) (*CreateNsResp, error) {
	if err := rule.CheckCharSet(req.Title, req.Desc); err != nil {
		return nil, err
	}
	if err := rule.CheckDescLength(req.Desc); err != nil {
		return nil, err
	}

	if err := rule.ValidateName(req.Name, rule.MaxNameLength, false); err != nil {
		return nil, err
	}

	ownerCheck := req.Owner
	if req.IgnoreAccessCheck { // use "system" to check access right
		ownerCheck = consts.SystemName
	}

	if _, err := ns.check(c, nil, req.Namespace, ownerCheck, OpAddSub); err != nil {
		return nil, err
	}

	if _, err := ns.check(c, nil, apipath.Join(req.Namespace, req.Name), ownerCheck, OpCreate); err != nil {
		return nil, err
	}

	apiNs := &models.APINamespace{
		ID:        hash.GenID("ns"),
		Owner:     req.Owner,
		OwnerName: req.OwnerName,
		Parent:    req.Namespace,
		Namespace: req.Name,
		Desc:      req.Desc,
		Access:    0,
		SubCount:  0,
		Active:    rule.ActiveEnable,
		Valid:     rule.Valid,
		Title:     req.Title,
	}

	err := ns.nsDB.Create(ns.db, apiNs)
	if err != nil {
		return nil, err
	}

	fullPath := apipath.Join(req.Namespace, req.Name)
	ns.redisAPI.DeleteCache(fullPath, myredis.CacheNS, true)
	ns.redisAPI.DeleteCache(req.Namespace, myredis.CacheNS, true)

	return &CreateNsResp{
		ID:        apiNs.ID,
		Owner:     apiNs.Owner,
		OwnerName: apiNs.OwnerName,
		Parent:    apiNs.Parent,
		Name:      apiNs.Namespace,
		SubCount:  apiNs.SubCount,
		Title:     apiNs.Title,
		Desc:      apiNs.Desc,
		Active:    apiNs.Active,
		CreateAt:  apiNs.CreateAt,
		UpdateAt:  apiNs.UpdateAt,
	}, nil
}

// DeleteNsReq DeleteNsReq
type DeleteNsReq = adaptor.DeleteNsReq

// DeleteNsResp DeleteNsResp
type DeleteNsResp = adaptor.DeleteNsResp

// Delete delete
func (ns *namespace) Delete(c context.Context, req *DeleteNsReq) (*DeleteNsResp, error) {
	if p, err := ns.check(c, nil, req.NamespacePath, req.Owner, OpDelete); err == nil {
		req.Pointer = p
	} else {
		return nil, err
	}

	return ns.InnerDelete(c, req)
}

func (ns *namespace) InnerDelete(c context.Context, req *DeleteNsReq) (*DeleteNsResp, error) {
	needCheck := false
	if needCheck = req.Pointer == nil; needCheck {
		if p, err := ns.query(c, req.NamespacePath, OpDelete); err == nil {
			req.Pointer = p
		} else {
			return nil, err
		}
	}
	if req.Pointer.SubCount > 0 {
		return nil, errcode.ErrNSWithSubs.NewError()
	}

	if req.ForceDelAPI && req.Owner == consts.SystemName {
		if oper := adaptor.GetRawAPIOper(); oper != nil {
			r := &adaptor.InnerDelRawByPrefixPathReq{
				NamespacePath: req.NamespacePath,
			}
			if _, err := oper.InnerDelByPrefixPath(c, r); err != nil {
				return nil, err
			}
		}
	}

	if needCheck {
		if _, err := ns.check(c, req.Pointer, req.NamespacePath, req.Owner, OpDelete); err != nil {
			return nil, err
		}
	}

	path, name := apipath.Split(req.NamespacePath)
	err := ns.nsDB.Delete(ns.db, path, name)
	if err != nil {
		return nil, err
	}

	ns.redisAPI.DeleteCache(req.NamespacePath, myredis.CacheNS, true)
	parent, _ := apipath.Split(req.NamespacePath)
	ns.redisAPI.DeleteCache(parent, myredis.CacheNS, true)

	return &DeleteNsResp{
		FullPath: req.NamespacePath,
	}, nil
}

// UpdateNsReq UpdateNsReq
type UpdateNsReq = adaptor.UpdateNsReq

// UpdateNsResp UpdateNsResp
type UpdateNsResp = adaptor.UpdateNsResp

// Update Update
func (ns *namespace) Update(c context.Context, req *UpdateNsReq) (*UpdateNsResp, error) {
	if err := rule.CheckCharSet(req.Title, req.Desc); err != nil {
		return nil, err
	}

	path, name := apipath.Split(req.NamespacePath)
	apiNs := &models.APINamespace{
		Parent:    path,
		Namespace: name,
		Title:     req.Title,
		Desc:      req.Desc,
	}
	if err := ns.nsDB.Update(ns.db, apiNs); err != nil {
		return nil, err
	}

	ns.redisAPI.DeleteCache(req.NamespacePath, myredis.CacheNS, true)

	return &UpdateNsResp{
		FullPath: req.NamespacePath,
	}, nil
}

// ActiveNsReq ActiveNsReq
type ActiveNsReq struct {
	NamespacePath string `json:"-"`
	Active        uint   `json:"active"`
}

// ActiveNsResp ActiveNsResp
type ActiveNsResp struct {
	FullPath string `json:"fullPath"`
	Active   uint   `json:"active"`
}

// Active Active
func (ns *namespace) Active(c context.Context, req *ActiveNsReq) (*ActiveNsResp, error) {
	path, name := apipath.Split(req.NamespacePath)

	apiNs := &models.APINamespace{
		Parent:    path,
		Namespace: name,
		Active:    req.Active,
	}

	if err := ns.nsDB.UpdateActive(ns.db, apiNs); err != nil {
		return nil, err
	}
	ns.redisAPI.DeleteCache(req.NamespacePath, myredis.CacheNS, true)

	return &ActiveNsResp{
		FullPath: req.NamespacePath,
		Active:   req.Active,
	}, nil
}

// ListNsReq ListNsReq
type ListNsReq = adaptor.ListNsReq

// ListNsResp ListNsResp
type ListNsResp = adaptor.ListNsResp

// List List
func (ns *namespace) List(c context.Context, req *ListNsReq) (*ListNsResp, error) {
	if _, err := ns.check(c, nil, req.NamespacePath, consts.SystemName, OpQuery); err != nil {
		return nil, err
	}

	cache := req.PageSize <= 0
	var data *models.APINamespaceList
	var err error
	if cache {
		data, err = ns.redisAPI.QueryNamespaceList(req.NamespacePath)
	}
	if !cache || err != nil {
		data, err = ns.nsDB.List(ns.db, req.NamespacePath, req.Active, req.Page, req.PageSize)
		if err != nil {
			return nil, err
		}
		if cache {
			ns.redisAPI.PutCache(data, req.NamespacePath, myredis.CacheNSList)
		}
	}
	// pageSize <= 0
	if cache {
		ns.activeFilter(data, req.Active)
	}
	listItems := make([]*NamespaceResp, len(data.List))
	for index, item := range data.List {
		listItems[index] = copyNamespaceResp(item)
	}

	return &ListNsResp{
		Total: int(data.Total),
		Page:  req.Page,
		List:  listItems,
	}, err
}

func (ns *namespace) activeFilter(src *models.APINamespaceList, active int) {
	if active >= 0 {
		uactive := uint(active)
		list := make([]*models.APINamespace, 0, len(src.List))
		for _, v := range src.List {
			if v.Active == uactive {
				list = append(list, v)
			}
		}
		src.List = list
	}
}

// SearchNsReq SearchNsReq
type SearchNsReq struct {
	NamespacePath string `json:"-"`
	Namespace     string `json:"namespace"`
	Title         string `json:"title"`
	Active        int    `json:"active"`
	WithSub       bool   `json:"withSub"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
}

// SearchNsResp SearchNsResp
type SearchNsResp struct {
	Total int              `json:"total"`
	Page  int              `json:"page"`
	List  []*NamespaceResp `json:"list"`
}

func (ns *namespace) search(req *SearchNsReq) (*models.APINamespaceList, error) {
	data, err := ns.nsDB.Search(ns.db, req.NamespacePath, req.Namespace, req.Title, req.Active, req.Page, req.PageSize, req.WithSub)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Search Search
func (ns *namespace) Search(c context.Context, req *SearchNsReq) (*SearchNsResp, error) {
	data, err := ns.search(req)
	if err != nil {
		return nil, err
	}

	list := make([]*NamespaceResp, 0, len(data.List))
	for _, v := range data.List {
		list = append(list, copyNamespaceResp(v))
	}
	return &SearchNsResp{
		Total: int(data.Total),
		Page:  req.Page,
		List:  list,
	}, nil
}

func (ns *namespace) makeTree(list *models.APINamespaceList) (*NsTreeResp, error) {
	return newTreeMaker(list).makeTree()
}

func newTreeMaker(list *models.APINamespaceList) *treeMaker {
	ls := list.List
	for _, v := range ls {
		v.FullPath = apipath.Join(v.Parent, v.Namespace)
	}
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].FullPath < ls[j].FullPath
	})

	resp := &NsTreeResp{}
	if len(ls) > 0 {
		resp.Root.FullPath, _ = apipath.Split(ls[0].FullPath)
		resp.Root.Parent, resp.Root.Name = apipath.Split(resp.Root.FullPath)
	}

	mt := &treeMaker{
		list:    ls,
		parents: []*NsTreeNode{&resp.Root},
		lastAdd: nil,
		resp:    resp,
	}
	return mt
}

type treeMaker struct {
	list    []*models.APINamespace
	parents []*NsTreeNode
	lastAdd *NsTreeNode
	resp    *NsTreeResp
}

func (m *treeMaker) makeTree() (*NsTreeResp, error) {
	for len(m.list) > 0 {
		child := m.list[0]
	innerLoop:
		for len(m.parents) > 0 {
			parent := m.parents[len(m.parents)-1]
			switch {
			case child.Parent == parent.FullPath: // push child
				m.lastAdd = &NsTreeNode{
					NamespaceResp: *copyNamespaceResp(child),
				}
				parent.Children = append(parent.Children, m.lastAdd)
				m.list = m.list[1:]

				break innerLoop //NOTE: break loop
			case len(child.Parent) > len(parent.FullPath) &&
				strings.HasPrefix(child.Parent, parent.FullPath): //push last add
				if m.lastAdd != nil && m.lastAdd.FullPath != parent.FullPath {
					m.parents = append(m.parents, m.lastAdd)
					if len(m.parents) > 100 {
						return nil, fmt.Errorf("add path(%s) to tree fail", child.FullPath)
					}
				} else {
					return nil, fmt.Errorf("path '%s' is not child of '%s'",
						child.FullPath, parent.FullPath)
				}

			default: // pop parent
				if parent == &m.resp.Root {
					return nil, fmt.Errorf("path '%s' out of root '%s'",
						child.FullPath, parent.FullPath)
				}
				m.parents = m.parents[:len(m.parents)-1]
			}
		}
	}
	return m.resp, nil
}

// NsTreeReq NsTreeReq
type NsTreeReq struct {
	RootPath string `json:"-"`
	Active   int    `json:"active"`
}

// NsTreeNode NsTreeNode
type NsTreeNode struct {
	NamespaceResp
	Children []*NsTreeNode `json:"children"`
}

// NsTreeResp NsTreeResp
type NsTreeResp struct {
	Root NsTreeNode `json:"root"`
}

// Search Search
func (ns *namespace) Tree(c context.Context, req *NsTreeReq) (*NsTreeResp, error) {
	searchReq := &SearchNsReq{
		NamespacePath: req.RootPath,
		Namespace:     "",
		Title:         "",
		Active:        req.Active,
		WithSub:       true,
		Page:          -1,
		PageSize:      -1,
	}
	list, err := ns.search(searchReq)
	if err != nil {
		return nil, err
	}
	return ns.makeTree(list)
}

func (ns *namespace) check(c context.Context, item *models.APINamespace, namespace string, owner string, op Operation) (*models.APINamespace, error) {
	if item == nil {
		q, err := ns.query(c, namespace, op)
		switch {
		case op == OpCreate && err == nil:
			return nil, errcode.ErrCreateExistsNS.FmtError(namespace)
		case op == OpCreate && err != nil:
			return nil, nil
		case err != nil:
			logger.Logger.Debugf("namespace.check error op=%v path=%v err=%s", op, namespace, err.Error())
			msg := fmt.Sprintf("invalid namespace: %s", namespace)
			return nil, error2.NewErrorWithString(error2.ErrParams, msg)
		default:
			item = q
		}
	}

	if op == OpDelete {
		if err := ns.checkEmptyNamespace(c, namespace, item, owner); err != nil {
			return nil, err
		}
	}

	if err := app.ValidateNamespace(owner, op, namespace); err != nil {
		return nil, err
	}

	if err := rule.ValidateActive(item.Active, op, enums.ObjectNamespace); err != nil {
		return nil, err
	}

	if err := ns.validForAPI(item, owner, op); err != nil {
		return nil, err
	}
	return item, nil
}

func (ns *namespace) checkEmptyNamespace(c context.Context, fullPath string, item *models.APINamespace, owner string) error {
	if item.SubCount > 0 {
		return errcode.ErrNSWithSubs.NewError()
	}

	if op := adaptor.GetRawAPIOper(); op != nil {
		req := &adaptor.RawListReq{
			NamespacePath: fullPath,
			Page:          1,
			PageSize:      1,
			Active:        -1,
		}
		if l, err := op.List(c, req); err == nil {
			if l.Total > 0 {
				return errcode.ErrNSWithRaw.NewError()
			}
		} else {
			return err
		}
	}

	if op := adaptor.GetPolyOper(); op != nil {
		req := &adaptor.PolyListReq{
			NamespacePath: fullPath,
			Page:          1,
			PageSize:      1,
			Active:        -1,
		}
		if l, err := op.List(c, req); err == nil {
			if l.Total > 0 {
				return errcode.ErrNSWithPoly.NewError()
			}
		} else {
			return err
		}
	}

	if svsOper := adaptor.GetServiceOper(); svsOper != nil {
		req := &adaptor.ListServiceReq{
			NamespacePath: fullPath,
			Page:          1,
			PageSize:      1,
		}
		if l, err := svsOper.List(c, req); err == nil {
			switch {
			case l.Total > 1:
				return errcode.ErrNSWithService.NewError()
			case l.Total == 1 && len(l.List) > 0: // check default service and customer keys
				svs := l.List[0]
				if svs.FullPath != apipath.Join(fullPath, item.Namespace) {
					// the only services without default name with path
					return errcode.ErrNSWithService.NewError()
				}

				// try delete the default service
				req := &adaptor.DeleteServiceReq{
					ServicePath: svs.FullPath,
					Owner:       owner,
				}
				if _, err := svsOper.Delete(c, req); err != nil {
					return err
				}
			}
		} else {
			return err
		}
	}

	return nil
}

// Check verify if namespace is valid for create API
func (ns *namespace) Check(c context.Context, namespace string, owner string, op Operation) (*NamespaceResp, error) {
	item, err := ns.check(c, nil, namespace, owner, op)
	return copyNamespaceResp(item), err
}

func (ns *namespace) validForAPI(item *models.APINamespace, owner string, op Operation) error {
	// TODO: remove
	return nil

	return nil
}

func copyNamespaceResp(item *models.APINamespace) *NamespaceResp {
	if item == nil {
		return nil
	}
	resp := &NamespaceResp{
		ID:        item.ID,
		Owner:     item.Owner,
		OwnerName: item.OwnerName,
		Parent:    item.Parent,
		Name:      item.Namespace,
		SubCount:  item.SubCount,
		Title:     item.Title,
		Desc:      item.Desc,
		Active:    item.Active,
		CreateAt:  item.CreateAt,
		UpdateAt:  item.UpdateAt,
		FullPath:  item.FullPath,
	}
	return resp
}

func (ns *namespace) query(c context.Context, namespace string, op Operation) (*models.APINamespace, error) {
	var item *models.APINamespace
	var err error
	noCache := rule.IgnoreCache(op)
	if !noCache {
		item, err = ns.redisAPI.QueryNamespace(namespace)
	}
	if noCache || err != nil {
		path, name := apipath.Split(namespace)
		item, err = ns.nsDB.Query(ns.db, path, name)
		if !noCache && err == nil {
			ns.redisAPI.PutCache(item, namespace, myredis.CacheNS)
		}
	}

	return item, err
}

// Query query namespace item from DB
func (ns *namespace) Query(c context.Context, namespace string) (*NamespaceResp, error) {
	item, err := ns.query(c, namespace, OpQuery)
	return copyNamespaceResp(item), err
}

// AppPathReq AppPathReq
type AppPathReq struct {
	APPID    string `json:"appID" binding:"required,max=64"`
	PathType string `json:"pathType" binding:"required,max=64"`
}

// AppPathResp AppPathResp
type AppPathResp struct {
	APPID    string `json:"appID"`
	PathType string `json:"pathType"`
	APPPath  string `json:"appPath"`
}

// APPPath query namespace of app path
func (ns *namespace) APPPath(c context.Context, req *AppPathReq) (*AppPathResp, error) {
	if !app.PathEnum.Verify(req.PathType) {
		return nil, errcode.ErrInvalidAppPathType.FmtError(req.PathType, app.PathEnum.GetAll())
	}

	resp := &AppPathResp{
		APPID:    req.APPID,
		PathType: req.PathType,
		APPPath:  app.Path(req.APPID, req.PathType),
	}
	return resp, nil
}

// InitAppPathReq InitAppPathReq
type InitAppPathReq = daprevent.EventInitAppPath

// InitAppPathResp InitAppPathResp
type InitAppPathResp struct {
	APPID string `json:"appID"`
}

// InitAPPPath init app path
func (ns *namespace) InitAPPPath(c context.Context, req *InitAppPathReq) (*InitAppPathResp, error) {
	d := &req.Data
	if err := app.InitAppPath(d.APPID, d.Owner, d.OwnerName, d.Header); err != nil {
		return nil, err
	}
	return &InitAppPathResp{APPID: d.APPID}, nil
}

// UpdateNsValidReq UpdateNsValidReq
type UpdateNsValidReq = adaptor.UpdateNsValidReq

// UpdateNsValidResp UpdateNsValidResp
type UpdateNsValidResp = adaptor.UpdateNsValidResp

// ValidWithSub ValidWithSub
func (ns *namespace) ValidWithSub(c context.Context, req *UpdateNsValidReq) (*UpdateNsValidResp, error) {
	err := ns.nsDB.UpdateValidWithSub(ns.db, req.NamespacePath, req.Valid)
	if err != nil {
		return nil, err
	}
	err = ns.redisAPI.DeletePatternCache(apipath.FormatPrefix(req.NamespacePath), myredis.CacheNS, true)
	return &UpdateNsValidResp{}, err
}

// ListNsByPrefixPathReq ListNsByPrefixPathReq
type ListNsByPrefixPathReq = adaptor.ListNsByPrefixPathReq

// ListNsByPrefixPathResp ListNsByPrefixPathResp
type ListNsByPrefixPathResp = adaptor.ListNsByPrefixPathResp

// ListByPrefixPath ListByPrefixPath
func (ns *namespace) ListByPrefixPath(c context.Context, req *ListNsByPrefixPathReq) (*ListNsByPrefixPathResp, error) {
	data, err := ns.nsDB.Search(ns.db, req.NamespacePath, "", "", req.Active, req.Page, req.PageSize, true)
	if err != nil {
		return nil, err
	}

	return &ListNsByPrefixPathResp{
		List:  data.List,
		Total: int(data.Total),
		Page:  req.Page,
	}, nil
}

// InnerImportNsReq InnerImportNsReq
type InnerImportNsReq = adaptor.InnerImportNsReq

// InnerImportNsResp InnerImportNsResp
type InnerImportNsResp = adaptor.InnerImportNsResp

// InnerImport InnerImport
func (ns *namespace) InnerImport(c context.Context, req *InnerImportNsReq) (*InnerImportNsResp, error) {
	for _, v := range req.List {
		v.ID = hash.GenID("ns")
	}
	err := ns.nsDB.CreateInBatches(ns.db, req.List)
	return &InnerImportNsResp{}, err
}

// InnerDelNsByPrefixPathReq InnerDelNsByPrefixPathReq
type InnerDelNsByPrefixPathReq = adaptor.InnerDelNsByPrefixPathReq

// InnerDelNsByPrefixPathResp InnerDelNsByPrefixPathResp
type InnerDelNsByPrefixPathResp = adaptor.InnerDelNsByPrefixPathResp

// InnerDelByPrefixPath InnerDelByPrefixPath
func (ns *namespace) InnerDelByPrefixPath(c context.Context, req *InnerDelNsByPrefixPathReq) (*InnerDelNsByPrefixPathResp, error) {
	err := ns.nsDB.DelByPrefixPath(ns.db, req.NamespacePath)
	if err != nil {
		return nil, err
	}
	ns.redisAPI.DeletePatternCache(apipath.FormatPrefix(req.NamespacePath), myredis.CacheNS, true)
	return &InnerDelNsByPrefixPathResp{}, nil
}
