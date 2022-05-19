package rediscache

import (
	"fmt"
	"time"
)

var (
	keyPrefix      = "polyapi"
	defaultTimeout = time.Minute * 5
)

// CacheType represents the type of cache
type CacheType int

// InvalidCacheType exports
const InvalidCacheType CacheType = 0

// Name gets regist name of this type
func (t CacheType) Name() string {
	if n, ok := cacheName[t]; ok {
		return n
	}
	return ""
}

// Int get CacheType as int
func (t CacheType) Int() int {
	return int(t)
}

// GetCacheKey get key of the fullPath
func (t CacheType) GetCacheKey(fullPath string) string {
	name := t.Name()
	if name == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s:%s", keyPrefix, name, fullPath)
}

// MustRegCacheType regist cache type with name
func MustRegCacheType(id int, name string) CacheType {
	t, err := regCacheType(id, name)
	if err != nil {
		panic(err)
	}
	return t
}

func regCacheType(id int, name string) (CacheType, error) {
	t := CacheType(len(cacheName) + 1)
	if id > 0 {
		t = CacheType(id)
		if n := t.Name(); n != "" {
			return InvalidCacheType, fmt.Errorf("duplicate cache type %d of name %s and %s)", t, n, name)
		}
	}
	if name == "" {
		return InvalidCacheType, fmt.Errorf("missing name of cache type %d", t)
	}
	if ot, ok := cacheType[name]; ok {
		return InvalidCacheType, fmt.Errorf("duplicate cache name %s of type %d and %d)", name, ot, t)
	}
	cacheName[t] = name
	cacheType[name] = t
	return t, nil
}

var (
	cacheName = map[CacheType]string{}
	cacheType = map[string]CacheType{}
)
