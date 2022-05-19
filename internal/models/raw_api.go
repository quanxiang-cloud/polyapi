package models

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"

	"gorm.io/gorm"
)

// RawAPIContent export
type RawAPIContent = adaptor.RawAPIContent

// APIDoc export
type APIDoc = adaptor.APIDoc

// RawAPIFull export
type RawAPIFull = adaptor.RawAPIFull

// RawAPICore is the raw api core scheme
type RawAPICore struct {
	ID        string
	Owner     string
	OwnerName string
	Namespace string
	Service   string
	Name      string
	Title     string
	Desc      string
	Path      string
	URL       string
	Action    string // action
	Method    string // method GET|POST|...
	Version   string
	Access    uint
	Active    uint
	Valid     uint

	Schema   string
	Host     string
	AuthType string

	CreateAt int64  // create time
	UpdateAt int64  // update time
	DeleteAt *int64 // delete time

	Content *RawAPIContent
}

// RawAPIList is list of raw api
type RawAPIList struct {
	Total int64
	List  []*RawAPICore
}

// RawAPIDoc is the raw api doc scheme
type RawAPIDoc struct {
	ID        string
	Owner     string
	OwnerName string
	Namespace string
	Service   string
	Name      string
	Title     string
	Desc      string
	CreateAt  int64  // create time
	UpdateAt  int64  // update time
	DeleteAt  *int64 // delete time
	Doc       *APIDoc
}

// RawAPIRepo is the interface for raw api db operator
type RawAPIRepo interface {
	Create(db *gorm.DB, raw *RawAPIFull) error
	CreateInBatches(db *gorm.DB, items []*RawAPIFull) error
	Del(db *gorm.DB, namespace string, names []string) error
	Get(db *gorm.DB, path, name string) (*RawAPICore, error)
	GetDoc(db *gorm.DB, path, name string) (*RawAPIDoc, error)
	GetDocInBatches(db *gorm.DB, path [][2]string) ([]*RawAPIDoc, error)
	GetByID(db *gorm.DB, id string) (*RawAPICore, error)
	GetInBatches(db *gorm.DB, path [][2]string) (*RawAPIList, error)
	List(db *gorm.DB, namespace, service string, active, page, pageSize int) (*RawAPIList, error)
	ListByPrefixPath(db *gorm.DB, path string, active, page, pageSize int) ([]*RawAPIFull, int64, error)
	UpdateActive(db *gorm.DB, namespace, name string, active uint) error
	UpdateValid(db *gorm.DB, path [][2]string, valid uint) error
	UpdateValidByPrefixPath(db *gorm.DB, path string, valid uint) error
	UpdateInBatch(db *gorm.DB, namespace, service, host, schema, authType string) error
	Search(db *gorm.DB, namespace, name, title string, active, page, pageSize int, withSub bool) (*RawAPIList, error)

	DelByPrefixPath(db *gorm.DB, path string) error
	// GetByCondition(db *gorm.DB, raw *RawAPI) (int, error)
	// UpdateParentByID(db *gorm.DB, id, parent string) error
	// SelectList(db *gorm.DB, page, limit int, raw *RawAPI) ([]*RawAPI, error)
}
