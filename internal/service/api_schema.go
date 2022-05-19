package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/internal/models/mysql"
	myredis "github.com/quanxiang-cloud/polyapi/internal/models/redis"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/apiprovider"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/pkg/core/schema"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/hash"
	"gorm.io/gorm"
)

// APISchemaOper represents the api schema operator
type APISchemaOper interface {
	GenSchema(ctx context.Context, req *GenSchemaReq) (*GenSchemaResp, error)
	QuerySchema(ctx context.Context, req *QuerySchemaReq) (*QuerySchemaResp, error)
	ListSchema(ctx context.Context, req *ListSchemaReq) (*ListSchemaResp, error)
	DeleteSchema(ctx context.Context, req *DeleteSchemaReq) (*DeleteSchemaResp, error)
	Request(ctx context.Context, req *SchemaRequestReq) (*SchemaRequestResp, error)
}

type apiSchema struct {
	conf     *config.Config
	db       *gorm.DB
	dbRepo   models.APISchemaRepo
	redisAPI *myredis.Client
}

// CreateSchemaOper new
func CreateSchemaOper(cfg *config.Config) (APISchemaOper, error) {
	db, err := createMysqlConn(cfg)
	if err != nil {
		return nil, err
	}
	redisCache, err := createRedisConn(cfg)
	if err != nil {
		return nil, err
	}

	logger.Logger.Infof("db connect ok: host=%s db=%s user=%s",
		cfg.Mysql.Host, cfg.Mysql.DB, cfg.Mysql.User)

	s := &apiSchema{
		conf:     cfg,
		db:       db,
		dbRepo:   mysql.NewAPISchemaRepo(),
		redisAPI: redisCache,
	}
	return s, nil
}

// GenSchemaReq GenSchemaReq
type GenSchemaReq struct {
	Owner     string `json:"-"`
	OwnerName string `json:"-"`
	APIPath   string `json:"namespace"`
	Title     string `json:"title"`
}

// GenSchemaResp GenSchemaResp
type GenSchemaResp struct {
}

func (s *apiSchema) GenSchema(ctx context.Context, req *GenSchemaReq) (*GenSchemaResp, error) {
	docRaw, err := apiprovider.APIQueryDoc(ctx, &apiprovider.QueryDocReq{
		APIPath:    req.APIPath,
		DocType:    apiprovider.DocTypeRaw.String(),
		TitleFirst: false,
	})
	if err != nil {
		return nil, err
	}

	var input expr.InputNodeDetail
	var output expr.OutputNodeDetail
	apiDoc := &apiprovider.APIDoc{
		Input:  &input,
		Output: &output,
	}
	if err := json.Unmarshal(docRaw.Doc, apiDoc); err != nil {
		return nil, err
	}

	inout := &expr.FmtAPIInOut{
		Input:  input,
		Output: output,
	}

	ns, name := apipath.Split(docRaw.APIPath)
	schema, err := schema.Gen(name, inout)
	if err != nil {
		return nil, err
	}

	entity := &models.APISchemaFull{
		ID:        "",
		Namespace: ns,
		Name:      name,
		Title:     req.Title,
		Schema:    schema,
	}

	for i := 0; i < hash.MaxHashConflict; i++ {
		entity.ID = hash.Default("schema", i, ns, name)
		if err = s.dbRepo.Create(s.db, entity); err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	s.redisAPI.DeleteCache(req.APIPath, myredis.CacheSchema, true)

	return &GenSchemaResp{}, nil
}

// QuerySchemaReq QuerySchemaReq
type QuerySchemaReq struct {
	APIPath string `json:"-"`
}

// QuerySchemaResp QuerySchemaResp
type QuerySchemaResp struct {
	ID           string          `json:"id"`
	Namespace    string          `json:"namespace"`
	Service      string          `json:"serivce"`
	Name         string          `json:"name"`
	Title        string          `json:"title"`
	Desc         string          `json:"desc"`
	InputSchema  json.RawMessage `json:"inputSchema"`
	OutputSchema json.RawMessage `json:"outputSchema"`
	CreateAt     int64           `json:"createAt"`
	UpdateAt     int64           `json:"updateAt"`
}

func (s *apiSchema) QuerySchema(ctx context.Context, req *QuerySchemaReq) (*QuerySchemaResp, error) {
	schema, err := s.redisAPI.QuerySchema(req.APIPath)
	if err != nil {
		ns, name := apipath.Split(req.APIPath)
		schema, err = s.dbRepo.Query(s.db, ns, name)
		if err != nil {
			return nil, err
		}
		s.redisAPI.PutCache(schema, req.APIPath, myredis.CacheSchema)
	}

	return &QuerySchemaResp{
		ID:           schema.ID,
		Namespace:    schema.Namespace,
		Name:         schema.Name,
		Title:        schema.Title,
		Desc:         schema.Desc,
		InputSchema:  schema.Schema.Input,
		OutputSchema: schema.Schema.Output,
		CreateAt:     schema.CreateAt,
		UpdateAt:     schema.UpdateAt,
	}, nil
}

// ListSchemaReq ListSchemaReq
type ListSchemaReq struct {
	NamespacePath string `json:"-"`
	WithSub       bool   `json:"withSub"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
}

// ListSchemaResp ListSchemaResp
type ListSchemaResp struct {
	Total int64             `json:"total"`
	List  []*ListSchemaNode `json:"list"`
}

// ListSchemaNode ListSchemaNode
type ListSchemaNode = adaptor.APISchema

func (s *apiSchema) ListSchema(ctx context.Context, req *ListSchemaReq) (*ListSchemaResp, error) {
	cache := req.PageSize <= 0 && !req.WithSub
	var ret *models.APISchemaList
	var err error
	if cache {
		ret, err = s.redisAPI.QuerySchemaList(req.NamespacePath)
	}
	if !cache || err != nil {
		ret, err = s.dbRepo.List(s.db, req.NamespacePath, req.WithSub, req.Page, req.PageSize)
		if err != nil {
			return nil, err
		}
		if cache {
			s.redisAPI.PutCache(req, req.NamespacePath, myredis.CacheSchemaList)
		}
	}
	return &ListSchemaResp{
		Total: ret.Total,
		List:  ret.List,
	}, nil
}

// DeleteSchemaReq DeleteSchemaReq
type DeleteSchemaReq struct {
	APIPath string `json:"-"`
}

// DeleteSchemaResp DeleteSchemaResp
type DeleteSchemaResp struct {
}

func (s *apiSchema) DeleteSchema(ctx context.Context, req *DeleteSchemaReq) (*DeleteSchemaResp, error) {
	ns, name := apipath.Split(req.APIPath)
	err := s.dbRepo.Delete(s.db, ns, name)
	if err != nil {
		return nil, err
	}

	s.redisAPI.DeleteCache(req.APIPath, myredis.CacheSchema, true)
	return &DeleteSchemaResp{}, nil
}

// SchemaRequestReq SchemaRequestReq
type SchemaRequestReq struct {
	Owner   string          `json:"-"`
	APIPath string          `json:"-"`
	APIType string          `json:"-"`
	Method  string          `json:"-"`
	Entity  json.RawMessage `json:"entity"`
	Header  http.Header     `json:"-"`
}

// SchemaRequestResp SchemaRequestResp
type SchemaRequestResp = apiprovider.RequestResp

func (s *apiSchema) Request(ctx context.Context, req *SchemaRequestReq) (*SchemaRequestResp, error) {
	body, err := schema.ParseRequest(req.Entity, req.Header)
	if err != nil {
		return nil, err
	}

	resp, err := apiprovider.APIRequest(ctx, &apiprovider.RequestReq{
		Owner:   req.Owner,
		APIPath: req.APIPath,
		Method:  req.Method,
		Body:    body,
		Header:  req.Header,
	})
	return resp, err
}
