package service

import (
	"context"

	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/internal/models/mysql"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/hash"

	id2 "github.com/quanxiang-cloud/cabin/id"
	"gorm.io/gorm"
)

// RawPoly is the raw_poly operator
type RawPoly interface {
	CreateInBatches(ctx context.Context, req *CreateRawPolyReq) (*CreateRawPolyResp, error)
	DeleteByRawAPI(ctx context.Context, req *DeleteByRawAPIReq) (*DeleteByRawAPIResp, error)
	DeleteByPolyAPI(ctx context.Context, req *DeleteByPolyAPIReq) (*DeleteByPolyAPIResp, error)
	QueryByRawAPI(ctx context.Context, req *QueryByRawAPIReq) (*QueryByRawAPIResp, error)
	QueryByPolyAPI(ctx context.Context, req *QueryByPolyAPIReq) (*QueryByPolyAPIResp, error)
	UpdateRawPoly(ctx context.Context, req *UpdateRawPolyReq) (*UpdateRawPolyResp, error)
}

// CreateRawPoly create raw_poly operator
func CreateRawPoly(conf *config.Config) (RawPoly, error) {
	db, err := createMysqlConn(conf)
	if err != nil {
		return nil, err
	}
	redisCache, err := createRedisConn(conf)
	if err != nil {
		return nil, err
	}

	rp := &rawPoly{
		conf:        conf,
		db:          db,
		redisAPI:    redisCache,
		rawPolyRepo: mysql.NewRawPolyRepo(),
	}

	adaptor.SetRawPolyOper(rp)

	return rp, nil
}

type rawPoly struct {
	conf        *config.Config
	db          *gorm.DB
	rawPolyRepo models.RawPolyRepo
	redisAPI    models.RedisCache
}

// RawPolyResp RawPolyResp
type RawPolyResp = adaptor.RawPolyResp

// CreateRawPolyReq CreateRawPolyReq
type CreateRawPolyReq = adaptor.CreateRawPolyReq

// CreateRawPolyResp CreateRawPolyResp
type CreateRawPolyResp = adaptor.CreateRawPolyResp

// CreateInBatches CreateInBatches
func (rp *rawPoly) CreateInBatches(ctx context.Context, req *CreateRawPolyReq) (*CreateRawPolyResp, error) {
	items := make([]*models.RawPoly, 0, len(req.RawAPIList))
	for _, raw := range req.RawAPIList {
		items = append(items, &models.RawPoly{
			ID:      hash.GenID("rp"),
			RawAPI:  raw,
			PolyAPI: req.PolyAPI,
		})
	}
	err := rp.rawPolyRepo.CreateInBatches(rp.db, items)
	return &CreateRawPolyResp{}, err
}

// DeleteByRawAPIReq DeleteByRawAPIReq
type DeleteByRawAPIReq = adaptor.DeleteByRawAPIReq

// DeleteByRawAPIResp DeleteByRawAPIResp
type DeleteByRawAPIResp = adaptor.DeleteByRawAPIResp

// DeleteByRawAPI DeleteByRawAPI
func (rp *rawPoly) DeleteByRawAPI(ctx context.Context, req *DeleteByRawAPIReq) (*DeleteByRawAPIResp, error) {
	err := rp.rawPolyRepo.DeleteByRawAPI(rp.db, req.RawAPI)
	return &DeleteByRawAPIResp{}, err
}

// DeleteByPolyAPIReq DeleteByPolyAPIReq
type DeleteByPolyAPIReq = adaptor.DeleteByPolyAPIReq

// DeleteByPolyAPIResp DeleteByPolyAPIResp
type DeleteByPolyAPIResp = adaptor.DeleteByPolyAPIResp

// DeleteByPolyAPI DeleteByPolyAPI
func (rp *rawPoly) DeleteByPolyAPI(ctx context.Context, req *DeleteByPolyAPIReq) (*DeleteByPolyAPIResp, error) {
	err := rp.rawPolyRepo.DeleteByPolyAPI(rp.db, req.PolyAPI)
	return &DeleteByPolyAPIResp{}, err
}

// DeleteByPolyAPIInBatchesReq DeleteByPolyAPIInBatchesReq
type DeleteByPolyAPIInBatchesReq = adaptor.DeleteByPolyAPIInBatchesReq

// DeleteByPolyAPIInBatchesResp DeleteByPolyAPIInBatchesResp
type DeleteByPolyAPIInBatchesResp = adaptor.DeleteByPolyAPIInBatchesResp

func (rp *rawPoly) DeleteByPolyAPIInBatches(ctx context.Context, req *DeleteByPolyAPIInBatchesReq) (*DeleteByPolyAPIInBatchesResp, error) {
	err := rp.rawPolyRepo.DeleteByPolyAPIInBatches(rp.db, req.PolyAPIPath)
	return &DeleteByPolyAPIInBatchesResp{}, err
}

// QueryByRawAPIReq QueryByRawAPIReq
type QueryByRawAPIReq = adaptor.QueryByRawAPIReq

// QueryByRawAPIResp QueryByRawAPIResp
type QueryByRawAPIResp = adaptor.QueryByRawAPIResp

// QueryByRawAPI QueryByRawAPI
func (rp *rawPoly) QueryByRawAPI(ctx context.Context, req *QueryByRawAPIReq) (*QueryByRawAPIResp, error) {
	data, err := rp.rawPolyRepo.QueryByRawAPI(rp.db, req.RawAPI)
	if err != nil {
		return nil, err
	}
	list := serializeRawPoly(data)
	return &QueryByRawAPIResp{
		Total: data.Total,
		List:  list,
	}, nil
}

// QueryByPolyAPIReq QueryByPolyAPIReq
type QueryByPolyAPIReq = adaptor.QueryByPolyAPIReq

// QueryByPolyAPIResp QueryByPolyAPIResp
type QueryByPolyAPIResp = adaptor.QueryByPolyAPIResp

// QueryByPolyAPI QueryByPolyAPI
func (rp *rawPoly) QueryByPolyAPI(ctx context.Context, req *QueryByPolyAPIReq) (*QueryByPolyAPIResp, error) {
	data, err := rp.rawPolyRepo.QueryByPolyAPI(rp.db, req.PolyAPI)
	if err != nil {
		return nil, err
	}
	list := serializeRawPoly(data)
	return &QueryByPolyAPIResp{
		Total: data.Total,
		List:  list,
	}, nil
}

func serializeRawPoly(data *models.RawPolyList) []*RawPolyResp {
	list := make([]*RawPolyResp, 0, len(data.List))
	for _, v := range data.List {
		list = append(list, &RawPolyResp{
			ID:      v.ID,
			RawAPI:  v.RawAPI,
			PolyAPI: v.PolyAPI,
		})
	}
	return list
}

// UpdateRawPolyReq UpdateRawPolyReq
type UpdateRawPolyReq = adaptor.UpdateRawPolyReq

// UpdateRawPolyResp UpdateRawPolyResp
type UpdateRawPolyResp = adaptor.UpdateRawPolyResp

// UpdateRawPoly UpdateRawPoly
func (rp *rawPoly) UpdateRawPoly(ctx context.Context, req *UpdateRawPolyReq) (*UpdateRawPolyResp, error) {
	// BUG: crash when list is empty
	if len(req.RawAPIList) == 0 {
		return &UpdateRawPolyResp{}, nil
	}
	items := make([]*models.RawPoly, 0, len(req.RawAPIList))
	for _, raw := range req.RawAPIList {
		items = append(items, &models.RawPoly{
			ID:      id2.WithPrefix(id2.ShortID(12), "rp_"),
			RawAPI:  raw,
			PolyAPI: req.PolyAPI,
		})
	}
	err := rp.rawPolyRepo.Update(rp.db, items)
	return &UpdateRawPolyResp{}, err
}

// ListRawPolyByPrefixPathReq ListRawPolyByPrefixPathReq
type ListRawPolyByPrefixPathReq = adaptor.ListRawPolyByPrefixPathReq

// ListRawPolyByPrefixPathResp ListRawPolyByPrefixPathResp
type ListRawPolyByPrefixPathResp = adaptor.ListRawPolyByPrefixPathResp

// ListByPrefixPath ListByPrefixPath
func (rp *rawPoly) ListByPrefixPath(ctx context.Context, req *ListRawPolyByPrefixPathReq) (*ListRawPolyByPrefixPathResp, error) {
	data, err := rp.rawPolyRepo.ListByPrefixPath(rp.db, req.Path)
	if err != nil {
		return nil, err
	}
	return &ListRawPolyByPrefixPathResp{
		List: data.List,
	}, nil
}

// InnerImportRawPolyReq InnerImportRawPolyReq
type InnerImportRawPolyReq = adaptor.InnerImportRawPolyReq

// InnerImportRawPolyResp InnerImportRawPolyResp
type InnerImportRawPolyResp = adaptor.InnerImportRawPolyResp

// InnerImport InnerImport
func (rp *rawPoly) InnerImport(ctx context.Context, req *InnerImportRawPolyReq) (*InnerImportRawPolyResp, error) {
	for _, v := range req.List {
		v.ID = hash.GenID("rp")
	}
	err := rp.rawPolyRepo.CreateInBatches(rp.db, req.List)
	return &InnerImportRawPolyResp{}, err
}

// InnerDelRawPolyByPrefixPathReq InnerDelRawPolyByPrefixPathReq
type InnerDelRawPolyByPrefixPathReq = adaptor.InnerDelRawPolyByPrefixPathReq

// InnerDelRawPolyByPrefixPathResp InnerDelRawPolyByPrefixPathResp
type InnerDelRawPolyByPrefixPathResp = adaptor.InnerDelRawPolyByPrefixPathResp

// InnerDelRawPolyByPrefixPath InnerDelRawPolyByPrefixPath
func (rp *rawPoly) InnerDelByPrefixPath(ctx context.Context, req *InnerDelRawPolyByPrefixPathReq) (*InnerDelRawPolyByPrefixPathResp, error) {
	err := rp.rawPolyRepo.DelByPrefixPath(rp.db, req.NamespacePath)
	return &InnerDelRawPolyByPrefixPathResp{}, err
}
