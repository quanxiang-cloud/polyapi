package models

import (
	"gorm.io/gorm"
)

// APIPermitGroup Schema Objects
type APIPermitGroup struct {
	ID        string
	Owner     string
	OwnerName string

	Namespace string
	Name      string
	Title     string
	Access    uint
	Desc      string
	Active    uint

	CreateAt int64
	UpdateAt int64
	DeleteAt int64
}

// ListPermitGroupResp ListPermitGroupResp
type ListPermitGroupResp struct {
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	List  []*APIPermitGroup `json:"list"`
}

// APIPermitGroupRepo persistence layer interface
type APIPermitGroupRepo interface {
	Create(db *gorm.DB, group *APIPermitGroup) error
	Delete(db *gorm.DB, path, name string) error
	Update(db *gorm.DB, group *APIPermitGroup) error
	UpdateActive(db *gorm.DB, group *APIPermitGroup) error
	Query(db *gorm.DB, path, name string) (*APIPermitGroup, error)
	List(db *gorm.DB, namespace string, page, pageSize int) (*ListPermitGroupResp, error)
	Count(db *gorm.DB, path string) (int, error)
}

//------------------------------------------------------------------------------

// APIPermitElem Schema Objects
type APIPermitElem struct {
	ID        string
	Owner     string
	OwnerName string

	GroupPath string
	ElemType  string
	ElemID    string
	ElemPath  string
	Desc      string
	ElemPri   uint
	Content   string
	Active    uint

	CreateAt int64
	UpdateAt int64
	DeleteAt int64
}

// ListPermitElemResp ListPermitElemResp
type ListPermitElemResp struct {
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	List  []*APIPermitElem `json:"list"`
}

// APIPermitElemRepo persistence layer interface
type APIPermitElemRepo interface {
	Create(db *gorm.DB, group *APIPermitElem) error
	Delete(db *gorm.DB, id string) error
	Update(db *gorm.DB, group *APIPermitElem) error
	UpdateActive(db *gorm.DB, group *APIPermitElem) error
	Query(db *gorm.DB, id string) (*APIPermitElem, error)
	List(db *gorm.DB, groupPath string, page, pageSize int) (*ListPermitElemResp, error)
	Count(db *gorm.DB, groupID string) (int, error)
}

//------------------------------------------------------------------------------

// APIPermitGrant Schema Objects
type APIPermitGrant struct {
	ID        string
	Owner     string
	OwnerName string

	GroupPath string
	GrantType string
	GrantID   string
	GrantName string
	GrantPri  uint
	Active    uint
	Desc      string

	CreateAt int64
	UpdateAt int64
	DeleteAt int64
}

// ListPermitGrantResp ListPermitGrantResp
type ListPermitGrantResp struct {
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	List  []*APIPermitGrant `json:"list"`
}

// APIPermitGrantRepo persistence layer interface
type APIPermitGrantRepo interface {
	Create(db *gorm.DB, group *APIPermitGrant) error
	Delete(db *gorm.DB, id string) error
	Update(db *gorm.DB, group *APIPermitGrant) error
	UpdateActive(db *gorm.DB, group *APIPermitGrant) error
	Query(db *gorm.DB, id string) (*APIPermitGrant, error)
	QueryGrant(db *gorm.DB, groupPath, grantType, grantID string) (*APIPermitGrant, error)
	List(db *gorm.DB, groupPath string, page, pageSize int) (*ListPermitGrantResp, error)
	ListGroups(db *gorm.DB, grantType, grantID string) ([]string, error)
	Count(db *gorm.DB, groupID string) (int, error)
}
