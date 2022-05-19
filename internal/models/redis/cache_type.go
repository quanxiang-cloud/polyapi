package redis

import (
	"github.com/quanxiang-cloud/polyapi/pkg/lib/rediscache"
)

// CacheType exports
type CacheType = rediscache.CacheType

// cache types, NEVER CHANGE THE ORDER
const (
	CacheRaw CacheType = iota + 1
	CacheRawDoc
	CacheRawList

	CachePoly
	CachePolyDoc
	CachePolyList

	CacheNS     // namespace
	CacheNSList // namespace list

	CacheService
	CacheServiceList

	CacheGateRawStat
	CacheGatePolyStat

	CacheGateIPBlack
	CacheGateIPWhite

	CacheSchema
	CacheSchemaList
)

func init() {
	rediscache.MustRegCacheType(CacheRaw.Int(), "r")
	rediscache.MustRegCacheType(CacheRawDoc.Int(), "rd")
	rediscache.MustRegCacheType(CacheRawList.Int(), "rl")
	rediscache.MustRegCacheType(CachePoly.Int(), "p")
	rediscache.MustRegCacheType(CachePolyDoc.Int(), "pd")
	rediscache.MustRegCacheType(CachePolyList.Int(), "pl")
	rediscache.MustRegCacheType(CacheNS.Int(), "n")
	rediscache.MustRegCacheType(CacheNSList.Int(), "nl")
	rediscache.MustRegCacheType(CacheService.Int(), "s")
	rediscache.MustRegCacheType(CacheServiceList.Int(), "sl")
	rediscache.MustRegCacheType(CacheGateRawStat.Int(), "rst")
	rediscache.MustRegCacheType(CacheGatePolyStat.Int(), "pst")
	rediscache.MustRegCacheType(CacheSchema.Int(), "sc")
	rediscache.MustRegCacheType(CacheSchemaList.Int(), "scl")
}

func getDocKey(t CacheType) CacheType {
	var d CacheType = rediscache.InvalidCacheType
	switch t {
	case CacheRaw:
		d = CacheRawDoc
	case CachePoly:
		d = CachePolyDoc
	}
	return d
}

func getStatKey(t CacheType) CacheType {
	var d CacheType = rediscache.InvalidCacheType
	switch t {
	case CacheRaw:
		d = CacheGateRawStat
	case CachePoly:
		d = CacheGatePolyStat
	}

	return d
}

func getListKey(t CacheType) CacheType {
	var d CacheType = rediscache.InvalidCacheType
	switch t {
	case CacheRaw:
		d = CacheRawList
	case CachePoly:
		d = CachePolyList
	case CacheNS:
		d = CacheNSList
	case CacheService:
		d = CacheServiceList
	case CacheSchema:
		d = CacheSchemaList
	}
	return d
}
