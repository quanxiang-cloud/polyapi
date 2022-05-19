package rediscache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"

	redis "github.com/go-redis/redis/v8"
)

// RedisCache RedisCache
type RedisCache interface {
	PutCache(p interface{}, fullName string, t CacheType) error
	DeleteCache(fullName string, t CacheType, includeList bool) error
	DeleteCacheInBatch(fullName []string, t CacheType) error
	DeletePatternCache(pattern string, t CacheType, includeList bool) error
	QueryCache(fullPath string, t CacheType) ([]byte, error)
	GetDocKey(t CacheType, fullPath string) string
	GetDocKeys(t CacheType, fullPaths []string) []string
	GetListKey(t CacheType, fullPath string) string
}

// NewClient create an Client
func NewClient(client *redis.ClusterClient) *Client {
	return &Client{
		C:          client,
		KeyTimeout: defaultTimeout,
	}
}

// Client is redis client
type Client struct {
	C                *redis.ClusterClient
	KeyTimeout       time.Duration
	GetDocCacheType  func(CacheType) CacheType
	GetListCacheType func(CacheType) CacheType
}

// GetDocKey get key of doc
func (c *Client) GetDocKey(t CacheType, fullPath string) string {
	if fn := c.GetDocCacheType; fn != nil {
		return fn(t).GetCacheKey(fullPath)
	}
	return ""
}

// GetDocKeys get keys of doc
func (c *Client) GetDocKeys(t CacheType, fullPaths []string) []string {
	if fn := c.GetDocCacheType; fn != nil {
		tDoc := fn(t)
		ret := make([]string, 0, len(fullPaths))
		for _, v := range fullPaths {
			if k := tDoc.GetCacheKey(v); k != "" {
				ret = append(ret, k)
			}
		}
		return ret
	}
	return nil
}

// GetListKey get key of list
func (c *Client) GetListKey(t CacheType, fullPath string) string {
	if fn := c.GetListCacheType; fn != nil {
		ns, _ := apipath.Split(fullPath)
		return fn(t).GetCacheKey(ns)
	}
	return ""
}

// PutCache store cache by path
func (c *Client) PutCache(p interface{}, fullPath string, t CacheType) error {
	key := t.GetCacheKey(fullPath)
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	err = c.C.Set(context.Background(), key, b, c.KeyTimeout).Err()
	return err
}

// DeleteCache delete cache by path
func (c *Client) DeleteCache(fullPath string, t CacheType, includeList bool) error {
	key := t.GetCacheKey(fullPath)
	err := c.C.Del(context.Background(), key).Err()

	if key := c.GetDocKey(t, fullPath); key != "" {
		err2 := c.C.Del(context.Background(), key).Err()
		if err2 == nil {
			err = err2
		}
	}

	if includeList {
		if key := c.GetListKey(t, fullPath); key != "" {
			err2 := c.C.Del(context.Background(), key).Err()
			if err2 == nil {
				err = err2
			}
		}
	}

	return err
}

// DeleteCacheInBatch delete cache in batch
func (c *Client) DeleteCacheInBatch(fullPaths []string, t CacheType) error {
	keys := make([]string, 0, len(fullPaths))
	for _, v := range fullPaths {
		keys = append(keys, t.GetCacheKey(v))
	}
	err := c.C.Del(context.Background(), keys...).Err()

	if keys := c.GetDocKeys(t, fullPaths); keys != nil {
		err2 := c.C.Del(context.Background(), keys...).Err()
		if err2 != nil {
			err = err2
		}
	}

	return err
}

// DeletePatternCache delete cache by pattern eg. /system/app/:appid/*
func (c *Client) DeletePatternCache(pattern string, t CacheType, includeList bool) error {
	patternKey := t.GetCacheKey(pattern)
	keys, err := c.C.Keys(context.Background(), patternKey).Result()
	if err != nil {
		return err
	}
	c.C.Del(context.Background(), keys...)

	if patternKey := c.GetDocKey(t, pattern); patternKey != "" {
		keys, err := c.C.Keys(context.Background(), patternKey).Result()
		if err != nil {
			return err
		}
		c.C.Del(context.Background(), keys...)
	}

	if includeList {
		if patternKey := c.GetListKey(t, pattern); patternKey != "" {
			keys, err := c.C.Keys(context.Background(), patternKey).Result()
			if err != nil {
				return err
			}
			c.C.Del(context.Background(), keys...)
		}
	}
	return err
}

// QueryCache query cache bytes
func (c *Client) QueryCache(fullPath string, t CacheType) ([]byte, error) {
	key := t.GetCacheKey(fullPath)
	content, err := c.C.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	c.C.Expire(context.Background(), key, c.KeyTimeout)
	return []byte(content), err
}
