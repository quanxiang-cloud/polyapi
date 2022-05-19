package adaptor

import (
	"context"
)

// RawPolyOper is the interface of rawPoly adaptor
type RawPolyOper interface {
	CreateInBatches(ctx context.Context, req *CreateRawPolyReq) (*CreateRawPolyResp, error)
	DeleteByRawAPI(ctx context.Context, req *DeleteByRawAPIReq) (*DeleteByRawAPIResp, error)
	DeleteByPolyAPI(ctx context.Context, req *DeleteByPolyAPIReq) (*DeleteByPolyAPIResp, error)
	DeleteByPolyAPIInBatches(ctx context.Context, req *DeleteByPolyAPIInBatchesReq) (*DeleteByPolyAPIInBatchesResp, error)
	QueryByRawAPI(ctx context.Context, req *QueryByRawAPIReq) (*QueryByRawAPIResp, error)
	QueryByPolyAPI(ctx context.Context, req *QueryByPolyAPIReq) (*QueryByPolyAPIResp, error)
	UpdateRawPoly(ctx context.Context, req *UpdateRawPolyReq) (*UpdateRawPolyResp, error)
	ListByPrefixPath(ctx context.Context, req *ListRawPolyByPrefixPathReq) (*ListRawPolyByPrefixPathResp, error)
	InnerImport(ctx context.Context, req *InnerImportRawPolyReq) (*InnerImportRawPolyResp, error)
	InnerDelByPrefixPath(ctx context.Context, req *InnerDelRawPolyByPrefixPathReq) (*InnerDelRawPolyByPrefixPathResp, error)
}

// DeleteByPolyAPIInBatchesReq DeleteByPolyAPIInBatchesReq
type DeleteByPolyAPIInBatchesReq struct {
	PolyAPIPath []string `json:"polyAPIPath"`
}

// DeleteByPolyAPIInBatchesResp DeleteByPolyAPIInBatchesResp
type DeleteByPolyAPIInBatchesResp struct {
}

// InnerDelRawPolyByPrefixPathReq InnerDelRawPolyByPrefixPathReq
type InnerDelRawPolyByPrefixPathReq struct {
	NamespacePath string `json:"-"`
}

// InnerDelRawPolyByPrefixPathResp InnerDelRawPolyByPrefixPathResp
type InnerDelRawPolyByPrefixPathResp struct {
}

// InnerImportRawPolyReq InnerImportRawPolyReq
type InnerImportRawPolyReq struct {
	List []*RawPoly `json:"list"`
}

// InnerImportRawPolyResp InnerImportRawPolyResp
type InnerImportRawPolyResp struct {
}

// ListRawPolyByPrefixPathReq ListRawPolyByPrefixPathReq
type ListRawPolyByPrefixPathReq struct {
	Path string `json:"-"`
}

// ListRawPolyByPrefixPathResp ListRawPolyByPrefixPathResp
type ListRawPolyByPrefixPathResp struct {
	List []*RawPoly `json:"list"`
}

// RawPolyResp RawPolyResp
type RawPolyResp struct {
	ID      string `json:"id"`
	RawAPI  string `json:"rawAPI"`
	PolyAPI string `json:"polyAPI"`
}

// CreateRawPolyReq CreateRawPolyReq
type CreateRawPolyReq struct {
	PolyAPI    string   `json:"polyAPI"`
	RawAPIList []string `json:"rawAPIList"`
}

// CreateRawPolyResp CreateRawPolyResp
type CreateRawPolyResp struct{}

// DeleteByRawAPIReq DeleteByRawAPIReq
type DeleteByRawAPIReq struct {
	RawAPI string `json:"rawAPI"`
}

// DeleteByRawAPIResp DeleteByRawAPIResp
type DeleteByRawAPIResp struct{}

// DeleteByPolyAPIReq DeleteByPolyAPIReq
type DeleteByPolyAPIReq struct {
	PolyAPI string `json:"polyAPI"`
}

// DeleteByPolyAPIResp DeleteByPolyAPIResp
type DeleteByPolyAPIResp struct{}

// QueryRawPolyResp QueryRawPolyResp
type QueryRawPolyResp struct {
	Total int64          `json:"total"`
	List  []*RawPolyResp `json:"list"`
}

// QueryByRawAPIReq QueryByRawAPIReq
type QueryByRawAPIReq struct {
	RawAPI []string `json:"rawAPI"`
}

// QueryByRawAPIResp QueryByRawAPIResp
type QueryByRawAPIResp = QueryRawPolyResp

// QueryByPolyAPIReq QueryByPolyAPIReq
type QueryByPolyAPIReq struct {
	PolyAPI string `json:"polyAPI"`
}

// QueryByPolyAPIResp QueryByPolyAPIResp
type QueryByPolyAPIResp = QueryRawPolyResp

// UpdateRawPolyReq UpdateRawPolyReq
type UpdateRawPolyReq struct {
	PolyAPI    string   `json:"polyAPI"`
	RawAPIList []string `json:"rawAPIList"`
}

// UpdateRawPolyResp UpdateRawPolyResp
type UpdateRawPolyResp struct{}

// RawPoly the relationship between raw api and poly api
type RawPoly struct {
	ID      string `json:"id"`
	RawAPI  string `gorm:"column:raw_api" json:"rawAPI"`
	PolyAPI string `gorm:"column:poly_api" json:"polyAPI"`
}

// SetRawPolyOper set the instance of RawPolyOper
func SetRawPolyOper(f RawPolyOper) RawPolyOper {
	i := getInst()
	old := i.rawPolyOper
	i.rawPolyOper = f
	return old
}

// GetRawPolyOper get the instance of RawPolyOper
func GetRawPolyOper() RawPolyOper {
	i := getInst()
	return i.rawPolyOper
}
