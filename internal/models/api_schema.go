package models

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"gorm.io/gorm"
)

// APISchemaFull describes api json schema
type APISchemaFull struct {
	ID        string
	Namespace string
	Name      string
	Title     string
	Desc      string
	Schema    *Schema
	CreateAt  int64
	UpdateAt  int64
}

// APISchema APISchema
type APISchema = adaptor.APISchema

// APISchemaList is the list of schema
type APISchemaList struct {
	Total int64
	List  []*APISchema
}

// Schema json schema
type Schema = adaptor.Schema

// APISchemaRepo is the interface for api schema db operator
type APISchemaRepo interface {
	Create(db *gorm.DB, schema *APISchemaFull) error
	Delete(db *gorm.DB, namespace, name string) error
	Query(db *gorm.DB, namespace, name string) (*APISchemaFull, error)
	List(db *gorm.DB, namespace string, withSub bool, page, pageSize int) (*APISchemaList, error)
}
