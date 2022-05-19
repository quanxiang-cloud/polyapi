package adaptor

import (
	"context"
)

// InnerDelServiceByPrefixPathReq InnerDelServiceByPrefixPathReq
type InnerDelServiceByPrefixPathReq struct {
	NamespacePath string `json:"-"`
}

// InnerDelServiceByPrefixPathResp InnerDelServiceByPrefixPathResp
type InnerDelServiceByPrefixPathResp struct {
}

// InnerImportServiceReq InnerImportServiceReq
type InnerImportServiceReq struct {
	List []*APIService `json:"list"`
}

// InnerImportServiceResp InnerImportServiceResp
type InnerImportServiceResp struct {
}

// ListServiceByPrefixReq ListServiceByPrefixReq
type ListServiceByPrefixReq struct {
	NamespacePath string `json:"-"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
}

// ListServiceByPrefixResp ListServiceByPrefixResp
type ListServiceByPrefixResp struct {
	Total int           `json:"total"`
	Page  int           `json:"page"`
	List  []*APIService `json:"list"`
}

// ServicesResp ServicesResp
type ServicesResp struct {
	ID          string `json:"id"`
	Owner       string `json:"owner"`
	OwnerName   string `json:"owenrName"`
	FullPath    string `json:"fullPath"`
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	Active      uint   `json:"active"`
	Schema      string `json:"schema"`
	Host        string `json:"host"`
	AuthType    string `json:"authType"`
	AuthContent string `json:"authContent"`
	CreateAt    int64  `json:"createAt"`
	UpdateAt    int64  `json:"updateAt"`
}

// ListServiceReq ListServiceReq
type ListServiceReq struct {
	NamespacePath string `json:"-"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
}

// ListServiceResp ListServiceResp
type ListServiceResp struct {
	Total int             `json:"total"`
	Page  int             `json:"page"`
	List  []*ServicesResp `json:"list"`
}

// DeleteServiceReq DeleteServiceReq
type DeleteServiceReq struct {
	ServicePath string `json:"-"`
	Owner       string `json:"-"`
}

// DeleteServiceResp DelteServiceResp
type DeleteServiceResp struct {
	FullPath string `json:"fullPath"`
}

// InnerDeleteServiceReq InnerDeleteServiceReq
type InnerDeleteServiceReq struct {
	NamespacePath string   `json:"-"`
	Names         []string `json:"names"`
}

// InnerDeleteServiceResp InnerDeleteServiceResp
type InnerDeleteServiceResp struct {
}

// APIService Schema Objects
type APIService struct {
	ID        string
	Owner     string
	OwnerName string
	Namespace string
	Name      string
	Title     string
	Desc      string

	Access uint
	Active uint

	Schema    string
	Host      string
	AuthType  string
	Authorize string

	CreateAt int64
	UpdateAt int64
	DeleteAt *int64
}

// ServiceOper ServiceOper
type ServiceOper interface {
	Check(c context.Context, service string, owner string, op Operation) (*ServicesResp, error)
	List(c context.Context, req *ListServiceReq) (*ListServiceResp, error)
	Delete(c context.Context, req *DeleteServiceReq) (*DeleteServiceResp, error)
	Query(c context.Context, service string) (*ServicesResp, error)
	InnerDelete(c context.Context, req *InnerDeleteServiceReq) (*InnerDeleteServiceResp, error)
	ListByPrefixPath(c context.Context, req *ListServiceByPrefixReq) (*ListServiceByPrefixResp, error)
	InnerImport(c context.Context, req *InnerImportServiceReq) (*InnerImportServiceResp, error)
	InnerDelByPrefixPath(ctx context.Context, req *InnerDelServiceByPrefixPathReq) (*InnerDelServiceByPrefixPathResp, error)
}

// SetServiceOper set the instance of service oper
func SetServiceOper(f ServiceOper) ServiceOper {
	i := getInst()
	old := i.serviceOper
	i.serviceOper = f
	return old
}

// GetServiceOper get the instance of service oper
func GetServiceOper() ServiceOper {
	i := getInst()
	return i.serviceOper
}
