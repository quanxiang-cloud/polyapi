package models

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"

	"gorm.io/gorm"
)

// APIService Schema Objects
type APIService = adaptor.APIService

// APIServiceList is list of namespace
type APIServiceList struct {
	List  []*APIService
	Total int64
}

// APIServiceRepo persistence layer interface
type APIServiceRepo interface {
	Create(db *gorm.DB, group *APIService) error
	CreateInBatches(db *gorm.DB, groups []*APIService) error
	Delete(db *gorm.DB, namespace, name string) error
	DeleteInBatch(db *gorm.DB, namespace string, names []string) error
	Update(db *gorm.DB, group *APIService) error
	UpdateProperty(db *gorm.DB, group *APIService) error
	UpdateActive(db *gorm.DB, group *APIService) error
	Query(db *gorm.DB, namespace, name string) (*APIService, error)
	List(db *gorm.DB, namespace string, page, pageSize int, withSub bool) (*APIServiceList, error)
	Count(db *gorm.DB, namespace string) (int, error)

	DelByPrefixPath(db *gorm.DB, path string) error
}
