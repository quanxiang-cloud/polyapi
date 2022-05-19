package models

import (
	"github.com/quanxiang-cloud/polyapi/pkg/lib/rediscache"
)

// CacheType exports
type CacheType = rediscache.CacheType

// RedisCache RedisCache
type RedisCache interface {
	PutCache(p interface{}, fullPath string, ty CacheType) error
	DeleteCache(fullPath string, ty CacheType, includeList bool) error
	DeleteCacheInBatch(fullPath []string, ty CacheType) error
	DeletePatternCache(pattern string, ty CacheType, includeList bool) error

	QueryRaw(fullPath string) (*RawAPICore, error)
	QueryRawDoc(fullPath string) (*RawAPIDoc, error)
	QueryRawList(fullPath string) (*RawAPIList, error)
	QueryPoly(fullPath string) (*PolyAPIScript, error)
	QueryPolyDoc(fullPath string) (*PolyAPIDoc, error)
	QueryPolyList(fullPath string) (*PolyAPIList, error)
	QueryNamespace(fullPath string) (*APINamespace, error)
	QueryNamespaceList(fullPath string) (*APINamespaceList, error)
	QueryService(fullPath string) (*APIService, error)
	QueryServiceList(fullPath string) (*APIServiceList, error)
	QuerySchema(fullPath string) (*APISchemaFull, error)
	QuerySchemaList(fullPath string) (*APISchemaList, error)

	IncAPIStat(apiPath string, raw bool) (int64, error)
	GetAPIStat(apiPath string, raw bool) (int64, error)
}
