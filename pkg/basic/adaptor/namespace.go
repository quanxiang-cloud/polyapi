package adaptor

import (
	"context"
)

// InnerDelNsByPrefixPathReq InnerDelNsByPrefixPathReq
type InnerDelNsByPrefixPathReq struct {
	NamespacePath string `json:"-"`
}

// InnerDelNsByPrefixPathResp InnerDelNsByPrefixPathResp
type InnerDelNsByPrefixPathResp struct {
}

// InnerImportNsReq InnerImportNsReq
type InnerImportNsReq struct {
	List []*APINamespace `json:"list"`
}

// InnerImportNsResp InnerImportNsResp
type InnerImportNsResp struct {
}

// ListNsByPrefixPathReq ListNsByPrefixPathReq
type ListNsByPrefixPathReq struct {
	NamespacePath string `json:"-"`
	Active        int    `json:"active"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
}

// ListNsByPrefixPathResp ListNsByPrefixPathResp
type ListNsByPrefixPathResp struct {
	Total int             `json:"total"`
	Page  int             `json:"page"`
	List  []*APINamespace `json:"list"`
}

// CreateNsReq CreateNsReq
type CreateNsReq struct {
	Owner             string `json:"-"`
	OwnerName         string `json:"-"`
	IgnoreAccessCheck bool   `json:"-"` // NOTE: launch by inner
	Namespace         string `uri:"namespace"`
	Name              string `json:"name" binding:"max=64"`
	Title             string `json:"title"`
	Desc              string `json:"desc"`
}

// NamespaceResp NamespaceResp
type NamespaceResp struct {
	ID        string `json:"id"`
	Owner     string `json:"owner"`
	OwnerName string `json:"ownerName"`
	Parent    string `json:"parent"`
	Name      string `json:"name"`
	SubCount  uint   `json:"subCount"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	// Access   uint
	Active   uint   `json:"active"`
	CreateAt int64  `json:"createAt"`
	UpdateAt int64  `json:"updateAt"`
	FullPath string `json:"-"`
}

// ListNsReq ListNsReq
type ListNsReq struct {
	NamespacePath string `json:"-"`
	Active        int    `json:"active"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
}

// ListNsResp ListNsResp
type ListNsResp struct {
	Total int              `json:"total"`
	Page  int              `json:"page"`
	List  []*NamespaceResp `json:"list"`
}

// DeleteNsReq DeleteNsReq
type DeleteNsReq struct {
	NamespacePath string        `json:"-"`
	Owner         string        `json:"-"`
	Pointer       *APINamespace `json:"-"`
	ForceDelAPI   bool          `json:"forceDelAPI"`
}

// DeleteNsResp DeleteNsResp
type DeleteNsResp struct {
	FullPath string `json:"fullPath"`
}

// UpdateNsValidReq UpdateNsValidReq
type UpdateNsValidReq struct {
	NamespacePath string `json:"-"`
	Valid         uint   `json:"valid"`
}

// UpdateNsValidResp UpdateNsValidResp
type UpdateNsValidResp struct {
}

// APINamespace Schema Objects
type APINamespace struct {
	ID        string
	Owner     string
	OwnerName string
	Parent    string
	Namespace string
	SubCount  uint
	Title     string
	Desc      string
	Access    uint
	Active    uint
	Valid     uint
	CreateAt  int64
	UpdateAt  int64
	DeleteAt  *int64
	FullPath  string `json:"-" gorm:"-"`
}

// UpdateNsReq UpdateNsReq
type UpdateNsReq struct {
	NamespacePath string `json:"-"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
}

// UpdateNsResp UpdateNsResp
type UpdateNsResp struct {
	FullPath string `json:"fullPath"`
}

// NamespaceOper NamespaceOper
type NamespaceOper interface {
	Create(c context.Context, req *CreateNsReq) (*NamespaceResp, error)
	Check(c context.Context, namespace string, owner string, op Operation) (*NamespaceResp, error)
	List(c context.Context, req *ListNsReq) (*ListNsResp, error)
	InnerDelete(c context.Context, req *DeleteNsReq) (*DeleteNsResp, error)
	ValidWithSub(c context.Context, req *UpdateNsValidReq) (*UpdateNsValidResp, error)
	ListByPrefixPath(c context.Context, req *ListNsByPrefixPathReq) (*ListNsByPrefixPathResp, error)
	InnerImport(c context.Context, req *InnerImportNsReq) (*InnerImportNsResp, error)
	InnerDelByPrefixPath(ctx context.Context, req *InnerDelNsByPrefixPathReq) (*InnerDelNsByPrefixPathResp, error)
	Update(c context.Context, req *UpdateNsReq) (*UpdateNsResp, error)
}

// SetNamespaceOper set the instance of namespace oper
func SetNamespaceOper(f NamespaceOper) NamespaceOper {
	i := getInst()
	old := i.namespaceOper
	i.namespaceOper = f
	return old
}

// GetNamespaceOper get the instance of namespace oper
func GetNamespaceOper() NamespaceOper {
	i := getInst()
	return i.namespaceOper
}
