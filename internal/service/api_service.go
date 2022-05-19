package service

import (
	"context"
	"fmt"

	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/internal/models/mysql"
	myredis "github.com/quanxiang-cloud/polyapi/internal/models/redis"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/enums"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/business/app"
	"github.com/quanxiang-cloud/polyapi/pkg/business/rule"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/hash"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"gorm.io/gorm"
)

// APIService represents the api service operator
type APIService interface {
	Delete(c context.Context, req *DeleteServiceReq) (*DeleteServiceResp, error)
	Update(c context.Context, req *UpdateServiceReq) (*UpdateServiceResp, error)
	UpdateProperty(c context.Context, req *UpdatePropertyReq) (*UpdatePropertyResp, error)
	Active(c context.Context, req *ActiveServiceReq) (*ActiveServiceResp, error)
	List(c context.Context, req *ListServiceReq) (*ListServiceResp, error)
	Query(c context.Context, service string) (*ServicesResp, error)
	Create(c context.Context, req *CreateServiceReq) (*CreateServiceResp, error)
	Check(c context.Context, service string, owner string, op Operation) (*ServicesResp, error)
}

// CreateServiceOper create a namespace API operater
func CreateServiceOper(conf *config.Config) (APIService, error) {
	db, err := createMysqlConn(conf)
	if err != nil {
		return nil, err
	}
	redisCache, err := createRedisConn(conf)
	if err != nil {
		return nil, err
	}

	p := &serviceAPI{
		conf:       conf,
		db:         db,
		servicesDB: mysql.NewAPIServiceRepo(),
		redisAPI:   redisCache,
	}

	adaptor.SetServiceOper(p)
	return p, nil
}

type serviceAPI struct {
	conf       *config.Config
	db         *gorm.DB
	redisAPI   models.RedisCache
	servicesDB models.APIServiceRepo
}

// CreateServiceReq CreateServiceReq
type CreateServiceReq struct {
	NamespacePath string `json:"-"`
	Owner         string `json:"-"`
	OwnerName     string `json:"-"`

	Name      string `json:"name"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Schema    string `json:"schema"`
	Host      string `json:"host"`
	AuthType  string `json:"authType"`
	Authorize string `json:"authorize"`
}

// CreateServiceResp CreateServiceResp
type CreateServiceResp struct {
	ID        string `json:"id"`
	Owner     string `json:"owner"`
	OwnerName string `json:"ownerName"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Active    uint   `json:"active"`
	Schema    string `json:"schema"`
	Host      string `json:"host"`
	AuthType  string `json:"authType"`
	Authorize string `json:"authorize"`
	CreateAt  int64  `json:"createAt"`
	UpdateAt  int64  `json:"updateAt"`
}

// Create Create
func (s *serviceAPI) Create(c context.Context, req *CreateServiceReq) (*CreateServiceResp, error) {
	if err := rule.ValidateName(req.Name, rule.MaxNameLength, false); err != nil {
		return nil, err
	}
	if err := rule.ValidateHost(req.Host); err != nil {
		return nil, err
	}

	var servicePath = apipath.Join(req.NamespacePath, req.Name)
	if op := adaptor.GetKMSOper(); op != nil {
		req := &adaptor.CheckAuthReq{
			AuthType:    req.AuthType,
			AuthContent: req.Authorize,
			ServicePath: servicePath,
		}
		if _, err := op.CheckAuth(c, req); err != nil {
			return nil, errcode.ErrServiceAuth.FmtError(err.Error())
		}
	}

	// check namespace
	if oper := adaptor.GetNamespaceOper(); oper != nil {
		if _, err := oper.Check(c, req.NamespacePath, req.Owner, OpAddService); err != nil {
			return nil, err
		}
	}
	if _, err := s.check(c, nil, servicePath, req.Owner, OpCreate); err != nil {
		return nil, err
	}

	apiSc := &models.APIService{
		ID:        hash.GenID("svs"),
		Owner:     req.Owner,
		OwnerName: req.OwnerName,
		Namespace: req.NamespacePath,
		Name:      req.Name,
		Title:     req.Title,
		Desc:      req.Desc,
		Schema:    req.Schema,
		Host:      req.Host,
		AuthType:  req.AuthType,
		Authorize: req.Authorize,
		Access:    0,
		Active:    rule.ActiveEnable,
	}

	if err := s.servicesDB.Create(s.db, apiSc); err != nil {
		return nil, err
	}

	fullPath := apipath.Join(req.NamespacePath, req.Name)
	s.redisAPI.DeleteCache(fullPath, myredis.CacheService, true)

	return &CreateServiceResp{
		ID:        apiSc.ID,
		Owner:     apiSc.Owner,
		OwnerName: apiSc.OwnerName,
		Namespace: apiSc.Namespace,
		Name:      apiSc.Name,
		Title:     apiSc.Title,
		Desc:      apiSc.Desc,
		Active:    apiSc.Active,
		Schema:    apiSc.Schema,
		Host:      apiSc.Host,
		AuthType:  apiSc.AuthType,
		Authorize: apiSc.Authorize,
		CreateAt:  apiSc.CreateAt,
		UpdateAt:  apiSc.UpdateAt,
	}, nil
}

// ServicesResp exports
type ServicesResp = adaptor.ServicesResp

func (s *serviceAPI) check(c context.Context, item *models.APIService, servicePath string, owner string, op Operation) (*models.APIService, error) {
	if item == nil {
		q, err := s.query(c, servicePath, op)
		switch {
		case op == OpCreate && err == nil:
			return nil, errcode.ErrCreateExistsService.NewError()
		case op == OpCreate && err != nil:
			return nil, nil
		case err != nil:
			logger.Logger.Debugf("service.check error op=%v path=%v err=%s", op, servicePath, err.Error())
			msg := fmt.Sprintf("invalid service: %s", servicePath)
			return nil, error2.NewErrorWithString(error2.ErrParams, msg)
		default:
			item = q
		}
	}

	if op == OpDelete {
		if err := s.checkEmptyService(c, servicePath, item); err != nil {
			return nil, err
		}
	}

	if err := app.ValidateServicePath(owner, op, servicePath); err != nil {
		return nil, err
	}

	if err := rule.ValidateActive(item.Active, op, enums.ObjectService); err != nil {
		return nil, err
	}

	if err := s.validFoAPI(item, owner, op); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *serviceAPI) checkEmptyService(c context.Context, servicePath string, item *models.APIService) error {
	// if contains raw api
	if op := adaptor.GetRawAPIOper(); op != nil {
		req := &adaptor.ListInServiceReq{
			ServicePath: servicePath,
			Page:        1,
			PageSize:    1,
		}
		if l, err := op.ListInService(c, req); err == nil {
			if l.Total > 0 {
				return errcode.ErrServiceWithRaw.NewError()
			}
		} else {
			return err
		}
	} else {
		return errcode.ErrServiceWithRaw.NewError()
	}

	// if contains customer api keys
	if op := adaptor.GetKMSOper(); op != nil {
		req := &adaptor.ListKMSCustomerKeyReq{
			Service:  servicePath,
			Page:     1,
			PageSize: 1,
			Owner:    "", //NOTE: keep empty
		}
		l, err := op.ListCustomerAPIKey(c, req)
		if err != nil {
			return err
		}
		if l.Total > 0 {
			return errcode.ErrServiceWithKeys.NewError()
		}
	} else {
		return errcode.ErrServiceWithKeys.NewError()
	}

	return nil
}

// Check verify if service is valid for create API
func (s *serviceAPI) Check(c context.Context, servicePath string, owner string, op Operation) (*ServicesResp, error) {
	item, err := s.check(c, nil, servicePath, owner, op)
	return s.copyResp(item), err
}

func (s *serviceAPI) validFoAPI(item *models.APIService, owner string, op Operation) error {
	// TODO: remove
	return nil

	return nil
}

func (s *serviceAPI) copyResp(item *models.APIService) *ServicesResp {
	if item == nil {
		return nil
	}
	resp := &ServicesResp{
		FullPath:    apipath.Join(item.Namespace, item.Name),
		Schema:      item.Schema,
		Host:        item.Host,
		AuthType:    item.AuthType,
		AuthContent: item.Authorize,
		Owner:       item.Owner,
		OwnerName:   item.OwnerName,
		ID:          item.ID,
		Active:      item.Active,
		CreateAt:    item.CreateAt,
		UpdateAt:    item.UpdateAt,
		Title:       item.Title,
		Desc:        item.Desc,
	}
	return resp
}

func (s *serviceAPI) query(c context.Context, servicePath string, op Operation) (*models.APIService, error) {
	var item *models.APIService
	var err error
	noCache := rule.IgnoreCache(op)
	if !noCache {
		item, err = s.redisAPI.QueryService(servicePath)
	}
	if noCache || err != nil {
		namespace, name := apipath.Split(servicePath)
		item, err = s.servicesDB.Query(s.db, namespace, name)
		if err == nil {
			s.redisAPI.PutCache(item, servicePath, myredis.CacheService)
		}
	}

	return item, err
}

// QueryService query service item from DB
func (s *serviceAPI) Query(c context.Context, servicePath string) (*ServicesResp, error) {
	f, err := s.query(c, servicePath, OpQuery)
	return s.copyResp(f), err
}

// DeleteServiceReq DeleteServiceReq
type DeleteServiceReq = adaptor.DeleteServiceReq

// DeleteServiceResp DelteServiceResp
type DeleteServiceResp = adaptor.DeleteServiceResp

// Delete delete
func (s *serviceAPI) Delete(c context.Context, req *DeleteServiceReq) (*DeleteServiceResp, error) {
	if _, err := s.check(c, nil, req.ServicePath, req.Owner, OpDelete); err != nil {
		return nil, err
	}

	path, name := apipath.Split(req.ServicePath)
	if err := s.servicesDB.Delete(s.db, path, name); err != nil {
		return nil, err
	}

	s.redisAPI.DeleteCache(req.ServicePath, myredis.CacheService, true)

	return &DeleteServiceResp{
		FullPath: req.ServicePath,
	}, nil
}

// InnerDeleteServiceReq InnerDeleteServiceReq
type InnerDeleteServiceReq = adaptor.InnerDeleteServiceReq

// InnerDeleteServiceResp InnerDeleteServiceResp
type InnerDeleteServiceResp = adaptor.InnerDeleteServiceResp

func (s *serviceAPI) InnerDelete(c context.Context, req *InnerDeleteServiceReq) (*InnerDeleteServiceResp, error) {
	err := s.servicesDB.DeleteInBatch(s.db, req.NamespacePath, req.Names)
	if err != nil {
		return nil, err
	}
	for _, name := range req.Names {
		servicePath := apipath.Join(req.NamespacePath, name)
		s.redisAPI.DeleteCache(servicePath, myredis.CacheService, true)
	}
	return &InnerDeleteServiceResp{}, nil
}

// UpdateServiceReq UpdateServiceReq
type UpdateServiceReq struct {
	ServicePath string `json:"-"`
	Title       string `json:"title"`
	Desc        string `json:"desc"`
}

// UpdateServiceResp UpdateServiceResp
type UpdateServiceResp struct {
	FullPath string `json:"fullPath"`
}

// Update update title and desc
func (s *serviceAPI) Update(c context.Context, req *UpdateServiceReq) (*UpdateServiceResp, error) {
	ns, name := apipath.Split(req.ServicePath)
	apiSc := &models.APIService{
		Namespace: ns,
		Name:      name,
		Title:     req.Title,
		Desc:      req.Desc,
	}
	if err := s.servicesDB.Update(s.db, apiSc); err != nil {
		return nil, err
	}
	s.redisAPI.DeleteCache(req.ServicePath, myredis.CacheService, true)

	return &UpdateServiceResp{
		FullPath: req.ServicePath,
	}, nil
}

// UpdatePropertyReq UpdatePropertyReq
type UpdatePropertyReq struct {
	ServicePath string `json:"-"`
	Host        string `json:"host"`
	Schema      string `json:"schema"`
	AuthType    string `json:"authType"`
	Authorize   string `json:"authorize"`
}

// UpdatePropertyResp UpdatePropertyResp
type UpdatePropertyResp struct {
}

func (s *serviceAPI) UpdateProperty(c context.Context, req *UpdatePropertyReq) (*UpdatePropertyResp, error) {
	if err := rule.ValidateHost(req.Host); err != nil {
		return nil, err
	}
	if op := adaptor.GetKMSOper(); op != nil {
		req := &adaptor.CheckAuthReq{
			AuthType:    req.AuthType,
			AuthContent: req.Authorize,
			ServicePath: req.ServicePath,
		}
		if _, err := op.CheckAuth(c, req); err != nil {
			return nil, errcode.ErrServiceAuth.FmtError(err.Error())
		}
	}

	ns, name := apipath.Split(req.ServicePath)
	apiSc := &models.APIService{
		Namespace: ns,
		Name:      name,
		Host:      req.Host,
		Schema:    req.Schema,
		AuthType:  req.AuthType,
		Authorize: req.Authorize,
	}
	if err := s.servicesDB.UpdateProperty(s.db, apiSc); err != nil {
		return nil, err
	}
	s.redisAPI.DeleteCache(req.ServicePath, myredis.CacheService, true)

	if op := adaptor.GetRawAPIOper(); op != nil {
		_, err := op.InnerUpdateRawInBatch(c, &adaptor.InnerUpdateRawInBatchReq{
			Namespace: ns,
			Service:   req.ServicePath,
			Host:      apiSc.Host,
			Schema:    apiSc.Schema,
			AuthType:  apiSc.AuthType,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := updateKMSKey(c, apiSc.Host, req.ServicePath, apiSc.AuthType, apiSc.Authorize); err != nil {
		return nil, err
	}
	return &UpdatePropertyResp{}, nil
}

func updateKMSKey(c context.Context, host, service, authType, authorize string) error {
	if op := adaptor.GetKMSOper(); op != nil {
		req := &adaptor.UpdateCustomerKeyInBatchReq{
			Host:        host,
			Service:     service,
			AuthType:    authType,
			AuthContent: authorize,
		}
		_, err := op.UpdateCustomerKeyInBatch(c, req)
		if err != nil {
			return err
		}
		return nil
	}
	return errcode.ErrServiceWithKeys.NewError()
}

// ActiveServiceReq ActiveServiceReq
type ActiveServiceReq struct {
	ServicePath string `json:"-"`
	Active      uint   `json:"active"`
}

// ActiveServiceResp ActiveServiceResp
type ActiveServiceResp struct {
	FullPath string `json:"fullPath"`
	Active   uint   `json:"active"`
}

func (s *serviceAPI) Active(c context.Context, req *ActiveServiceReq) (*ActiveServiceResp, error) {
	ns, name := apipath.Split(req.ServicePath)
	apiSc := &models.APIService{
		Namespace: ns,
		Name:      name,
		Active:    req.Active,
	}
	if err := s.servicesDB.UpdateActive(s.db, apiSc); err != nil {
		return nil, err
	}
	s.redisAPI.DeleteCache(req.ServicePath, myredis.CacheService, true)

	return &ActiveServiceResp{
		FullPath: req.ServicePath,
		Active:   req.Active,
	}, nil
}

// ListServiceReq ListServiceReq
type ListServiceReq = adaptor.ListServiceReq

// ListServiceResp ListServiceResp
type ListServiceResp = adaptor.ListServiceResp

// List List
func (s *serviceAPI) List(c context.Context, req *ListServiceReq) (*ListServiceResp, error) {
	list, err := s.redisAPI.QueryServiceList(req.NamespacePath)
	if err != nil {
		list, err = s.servicesDB.List(s.db, req.NamespacePath, req.Page, req.PageSize, false)
		if err != nil {
			return nil, err
		}
		s.redisAPI.PutCache(list, req.NamespacePath, myredis.CacheServiceList)
	}

	listItems := make([]*ServicesResp, len(list.List))
	for index, item := range list.List {
		listItems[index] = s.copyResp(item)
	}

	return &ListServiceResp{
		Total: int(list.Total),
		Page:  req.Page,
		List:  listItems,
	}, err
}

// ListServiceByPrefixReq ListServiceByPrefixReq
type ListServiceByPrefixReq = adaptor.ListServiceByPrefixReq

// ListServiceByPrefixResp ListServiceByPrefixResp
type ListServiceByPrefixResp = adaptor.ListServiceByPrefixResp

func (s *serviceAPI) ListByPrefixPath(c context.Context, req *ListServiceByPrefixReq) (*ListServiceByPrefixResp, error) {
	data, err := s.servicesDB.List(s.db, req.NamespacePath, req.Page, req.PageSize, true)
	if err != nil {
		return nil, err
	}
	return &ListServiceByPrefixResp{
		List:  data.List,
		Total: int(data.Total),
		Page:  req.Page,
	}, nil
}

// InnerImportServiceReq InnerImportServiceReq
type InnerImportServiceReq = adaptor.InnerImportServiceReq

// InnerImportServiceResp InnerImportServiceResp
type InnerImportServiceResp = adaptor.InnerImportServiceResp

// InnerImport InnerImport
func (s *serviceAPI) InnerImport(c context.Context, req *InnerImportServiceReq) (*InnerImportServiceResp, error) {
	for _, v := range req.List {
		v.ID = hash.GenID("svs")
	}
	err := s.servicesDB.CreateInBatches(s.db, req.List)
	return &InnerImportServiceResp{}, err
}

// InnerDelServiceByPrefixPathReq InnerDelServiceByPrefixPathReq
type InnerDelServiceByPrefixPathReq = adaptor.InnerDelServiceByPrefixPathReq

// InnerDelServiceByPrefixPathResp InnerDelServiceByPrefixPathResp
type InnerDelServiceByPrefixPathResp = adaptor.InnerDelServiceByPrefixPathResp

// InnerDelServiceByPrefixPath InnerDelServiceByPrefixPath
func (s *serviceAPI) InnerDelByPrefixPath(c context.Context, req *InnerDelServiceByPrefixPathReq) (*InnerDelServiceByPrefixPathResp, error) {
	err := s.servicesDB.DelByPrefixPath(s.db, req.NamespacePath)
	if err != nil {
		return nil, err
	}
	s.redisAPI.DeletePatternCache(apipath.FormatPrefix(req.NamespacePath), myredis.CacheService, true)
	return &InnerDelServiceByPrefixPathResp{}, nil
}
