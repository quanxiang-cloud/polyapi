package models

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"

	"gorm.io/gorm"
)

// PolyAPIFull export
type PolyAPIFull = adaptor.PolyAPIFull

// PolyTableName is table name
const PolyTableName = "api_poly"

// PolyAPIArrange is view of poly arrange
type PolyAPIArrange struct {
	ID        string `gorm:"primarykey"` // uuid
	Owner     string // owner
	OwnerName string // owner name
	Namespace string
	Name      string // name
	Title     string
	Desc      string
	Active    uint
	Valid     uint
	Access    uint
	Method    string
	Arrange   string // arrange info
	CreateAt  int64  // create time
	UpdateAt  int64  // update time
	DeleteAt  *int64 // delete time
	BuildAt   int64  // build time
}

// PolyAPIList is list of poly api
type PolyAPIList struct {
	Total int64
	List  []*PolyAPIArrange
}

// PolyBuildResult is view of poly build result
type PolyBuildResult struct {
	ID        string `gorm:"primarykey"` // uuid
	Namespace string
	Name      string // name
	Script    string // script
	Doc       string // API doc
	BuildAt   int64  // build time
}

// PolyAPIScript is view of poly script
type PolyAPIScript struct {
	ID        string `gorm:"primarykey"` // uuid
	Namespace string
	Name      string // name
	Active    uint
	Valid     uint
	Method    string
	Script    string // script
	Owner     string
	BuildAt   int64 // build time
}

// PolyAPIDoc is view of poly doc
type PolyAPIDoc struct {
	ID        string `gorm:"primarykey"` // uuid
	Namespace string
	Name      string // name
	Title     string
	Desc      string
	Doc       *APIDoc // API doc
	BuildAt   int64   // build time
}

// TableName TableName
func (v *PolyAPIArrange) TableName() string {
	return PolyTableName
}

// TableName TableName
func (v *PolyAPIScript) TableName() string {
	return PolyTableName
}

// PolyAPIRepo PolyAPIRepo
type PolyAPIRepo interface {
	Create(db *gorm.DB, info *PolyAPIArrange) error
	Delete(db *gorm.DB, path string, name []string) error
	UpdateArrange(db *gorm.DB, info *PolyAPIArrange) error
	UpdateScript(db *gorm.DB, info *PolyBuildResult) error
	GetArrange(db *gorm.DB, path, name string) (*PolyAPIArrange, error)
	GetScript(db *gorm.DB, path, name string) (*PolyAPIScript, error)
	GetDoc(db *gorm.DB, path, name string) (*PolyAPIDoc, error)
	GetDocInBatches(db *gorm.DB, path [][2]string) ([]*PolyAPIDoc, error)
	List(db *gorm.DB, namespace string, active, page, pageSize int) (*PolyAPIList, error)
	UpdateActive(db *gorm.DB, namespace, name string, active uint) error
	Search(db *gorm.DB, namespace, name, title string, active, page, pageSize int, withSub bool) (*PolyAPIList, error)

	UpdateValid(db *gorm.DB, path [][2]string, valid uint) error
	UpdateValidByPrefixPath(db *gorm.DB, namespace string, valid uint) error
	CreateInBatches(db *gorm.DB, items []*PolyAPIFull) error
	ListByPrefixPath(db *gorm.DB, namespace string, active, page, pageSize int) ([]*PolyAPIFull, int64, error)
	DelByPrefixPath(db *gorm.DB, path string) error
}
