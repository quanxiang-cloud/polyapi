package models

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"

	"gorm.io/gorm"
)

// RawPoly the relationship between raw api and poly api
type RawPoly = adaptor.RawPoly

// RawPolyList is list of RawPoly
type RawPolyList struct {
	List  []*RawPoly
	Total int64
}

// RawPolyRepo is the interface for the relationship between raw and poly api db operator
type RawPolyRepo interface {
	Create(db *gorm.DB, item *RawPoly) error
	CreateInBatches(db *gorm.DB, items []*RawPoly) error
	DeleteByRawAPI(db *gorm.DB, rawPath string) error
	DeleteByPolyAPI(db *gorm.DB, polyPath string) error
	DeleteByPolyAPIInBatches(db *gorm.DB, polyPath []string) error
	QueryByRawAPI(db *gorm.DB, rawPath []string) (*RawPolyList, error)
	QueryByPolyAPI(db *gorm.DB, polyPath string) (*RawPolyList, error)
	Update(db *gorm.DB, item []*RawPoly) error

	ListByPrefixPath(db *gorm.DB, path string) (*RawPolyList, error)
	DelByPrefixPath(db *gorm.DB, path string) error
}
