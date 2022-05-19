package models

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"

	"gorm.io/gorm"
)

// APINamespace Schema Objects
type APINamespace = adaptor.APINamespace

// APINamespaceList is list of namespace
type APINamespaceList struct {
	List  []*APINamespace
	Total int64
}

// APINamespaceRepo persistence layer interface
type APINamespaceRepo interface {
	Create(db *gorm.DB, item *APINamespace) error
	Delete(db *gorm.DB, path, name string) error
	Update(db *gorm.DB, item *APINamespace) error
	UpdateActive(db *gorm.DB, item *APINamespace) error
	Query(db *gorm.DB, path, name string) (*APINamespace, error)
	List(db *gorm.DB, path string, active, page, pageSize int) (*APINamespaceList, error)
	Count(db *gorm.DB, path string) (int, error)
	Search(db *gorm.DB, parent, namespace, title string, active, page, pageSize int, withSub bool) (*APINamespaceList, error)

	CreateInBatches(db *gorm.DB, items []*APINamespace) error
	UpdateValidWithSub(db *gorm.DB, path string, valid uint) error
	DelByPrefixPath(db *gorm.DB, path string) error
}
