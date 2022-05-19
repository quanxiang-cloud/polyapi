package adaptor

import (
	"context"

	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/protocol"
)

// QueryPolySwaggerReq QueryPolySwaggerReq
type QueryPolySwaggerReq struct {
	APIPath []string `json:"-"`
}

// QueryPolySwaggerResp QueryPolySwaggerResp
type QueryPolySwaggerResp struct {
	Swagger []byte `json:"swagger"`
}

// InnerDelPolyByPrefixPathReq InnerDelPolyByPrefixPathReq
type InnerDelPolyByPrefixPathReq struct {
	NamespacePath string `json:"-"`
}

// InnerDelPolyByPrefixPathResp InnerDelPolyByPrefixPathResp
type InnerDelPolyByPrefixPathResp struct {
}

// InnerImportPolyReq InnerImportPolyReq
type InnerImportPolyReq struct {
	List []*PolyAPIFull `json:"list"`
}

// InnerImportPolyResp InnerImportPolyResp
type InnerImportPolyResp struct {
}

// ListPolyByPrefixPathReq ListPolyByPrefixPathReq
type ListPolyByPrefixPathReq struct {
	NamespacePath string `json:"-"`
	Active        int    `json:"active"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
}

// ListPolyByPrefixPathResp ListPolyByPrefixPathResp
type ListPolyByPrefixPathResp struct {
	List  []*PolyAPIFull `json:"list"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
}

// PolyValidByPrefixPathReq PolyValidByPrefixPathReq
type PolyValidByPrefixPathReq struct {
	NamespacePath string `json:"-"`
	Valid         uint   `json:"valid"`
}

// PolyValidByPrefixPathResp PolyValidByPrefixPathResp
type PolyValidByPrefixPathResp struct {
}

// PolyValidInBatchesReq PolyValidInBatchesReq
type PolyValidInBatchesReq struct {
	APIPath []string `json:"-"`
	Valid   uint     `json:"valid"`
}

// PolyValidInBatchesResp PolyValidInBatchesResp
type PolyValidInBatchesResp struct{}

// PolyListReq polyListReq
type PolyListReq struct {
	NamespacePath string `json:"-"`
	Active        int    `json:"active"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
}

// PolyListResp PolyListResp
type PolyListResp struct {
	Total int             `json:"total"`
	Page  int             `json:"page"`
	List  []*PolyListNode `json:"list"`
}

// PolyDeleteReq PolyDeleteReq
type PolyDeleteReq struct {
	NamespacePath string   `json:"-" binding:"-"` //
	Names         []string `json:"names"`
	Owner         string   `json:"-"`
}

// PolyDeleteResp PolyDeleteResp
type PolyDeleteResp struct{}

// PolyListNode node
type PolyListNode struct {
	ID         string `json:"id"`
	Owner      string `json:"owner"`
	OwnerName  string `json:"ownerName"`
	FullPath   string `json:"fullPath"`
	Method     string `json:"method"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Active     uint   `json:"active"`
	Valid      uint   `json:"valid"`
	CreateAt   int64  `json:"createAt"`
	UpdateAt   int64  `json:"updateAt"`
	AccessPath string `json:"accessPath"`
}

// PolyAPIFull is the poly api db scheme
type PolyAPIFull struct {
	ID        string `gorm:"primarykey" json:"id"` // uuid
	Owner     string `json:"owner"`                // owner
	OwnerName string `json:"ownerName"`            // owner name
	Namespace string `json:"namespace"`
	Name      string `json:"name"` // name
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	Active    uint   `json:"active"`
	Method    string `json:"method"`
	Arrange   string `json:"arrange"` // arrange info
	Doc       string `json:"doc"`
	Script    string `json:"script"`
	Valid     uint   `json:"valid"`
	CreateAt  int64  `json:"createAt"` // create time
	UpdateAt  int64  `json:"updateAt"` // update time
	DeleteAt  *int64 `json:"deleteAt"` // delete time
	BuildAt   int64  `json:"buildAt"`  // build time
}

// PolyAPIOper is the interface for poly api operator
type PolyAPIOper interface {
	List(c context.Context, req *PolyListReq) (*PolyListResp, error)
	InnerDelete(c context.Context, req *PolyDeleteReq) (*PolyDeleteResp, error)
	ValidInBatches(ctx context.Context, req *PolyValidInBatchesReq) (*PolyValidInBatchesResp, error)
	ValidByPrefixPath(ctx context.Context, req *PolyValidByPrefixPathReq) (*PolyValidByPrefixPathResp, error)
	ListByPrefixPath(ctx context.Context, req *ListPolyByPrefixPathReq) (*ListPolyByPrefixPathResp, error)
	InnerImport(ctx context.Context, req *InnerImportPolyReq) (*InnerImportPolyResp, error)
	InnerDelByPrefixPath(ctx context.Context, req *InnerDelPolyByPrefixPathReq) (*InnerDelPolyByPrefixPathResp, error)
	QuerySwagger(ctx context.Context, req *QueryPolySwaggerReq) (*QueryPolySwaggerResp, error)
}

// SetPolyOper set the instance of poly oper
func SetPolyOper(f PolyAPIOper) PolyAPIOper {
	i := getInst()
	old := i.polyOper
	i.polyOper = f
	return old
}

// GetPolyOper get the instance of poly oper
func GetPolyOper() PolyAPIOper {
	i := getInst()
	return i.polyOper
}

// EvalerOper if Evaler creater
type EvalerOper interface {
	CreateEvaler() protocol.Evaler
}

// SetEvalerOper set the instance of Evaler oper
func SetEvalerOper(f EvalerOper) EvalerOper {
	i := getInst()
	old := i.evalerOper
	i.evalerOper = f
	return old
}

// GetEvalerOper get the instance of Evaler oper
func GetEvalerOper() EvalerOper {
	i := getInst()
	return i.evalerOper
}
