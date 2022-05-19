package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/rediscache"
)

// NewRedisClient create a redis client
func NewRedisClient(client *redis.ClusterClient) *Client {
	c := rediscache.NewClient(client)
	c.GetDocCacheType = getDocKey
	c.GetListCacheType = getListKey
	return &Client{
		Client: c,
	}
}

// Client represend the redis client
type Client struct {
	*rediscache.Client
	apiStatKeyDur time.Duration
}

// QueryRaw query cache of raw api by fullPath
func (c *Client) QueryRaw(fullPath string) (*models.RawAPICore, error) {
	content, err := c.QueryCache(fullPath, CacheRaw)
	if err != nil {
		return nil, err
	}
	ret := &models.RawAPICore{}
	if err := json.Unmarshal(content, ret); err != nil {
		return nil, err
	}
	if err := ret.Content.Consts.DelayedJSONDecode(); err != nil {
		return nil, err
	}
	return ret, nil
}

// QueryRawDoc query cache of raw api doc by fullPath
func (c *Client) QueryRawDoc(fullPath string) (*models.RawAPIDoc, error) {
	content, err := c.QueryCache(fullPath, CacheRawDoc)
	if err != nil {
		return nil, err
	}
	ret := &models.RawAPIDoc{}
	if err := json.Unmarshal(content, ret); err != nil {
		return nil, err
	}
	if err := ret.Doc.FmtInOut.DelayedJSONDecode(); err != nil {
		return nil, err
	}
	return ret, nil
}

// QueryRawList query cache of raw api list by fullPath
func (c *Client) QueryRawList(fullPath string) (*models.RawAPIList, error) {
	content, err := c.QueryCache(fullPath, CacheRawList)
	if err != nil {
		return nil, err
	}
	ret := &models.RawAPIList{}
	if err := json.Unmarshal(content, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// QueryPoly query poly api cache by fullPath
func (c *Client) QueryPoly(fullPath string) (*models.PolyAPIScript, error) {
	content, err := c.QueryCache(fullPath, CachePoly)
	if err != nil {
		return nil, err
	}
	ret := &models.PolyAPIScript{}
	err = json.Unmarshal(content, ret)
	return ret, err
}

// QueryPolyDoc query poly api doc cache by fullPath
func (c *Client) QueryPolyDoc(fullPath string) (*models.PolyAPIDoc, error) {
	content, err := c.QueryCache(fullPath, CachePolyDoc)
	if err != nil {
		return nil, err
	}
	ret := &models.PolyAPIDoc{}
	if err := json.Unmarshal(content, ret); err != nil {
		return nil, err
	}
	if err := ret.Doc.FmtInOut.DelayedJSONDecode(); err != nil {
		return nil, err
	}
	return ret, nil
}

// QueryPolyList query poly api list cache by fullPath
func (c *Client) QueryPolyList(fullPath string) (*models.PolyAPIList, error) {
	content, err := c.QueryCache(fullPath, CachePolyList)
	if err != nil {
		return nil, err
	}
	ret := &models.PolyAPIList{}
	if err := json.Unmarshal(content, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// QueryNamespace query namespace cache by fullPath
func (c *Client) QueryNamespace(fullPath string) (*models.APINamespace, error) {
	content, err := c.QueryCache(fullPath, CacheNS)
	if err != nil {
		return nil, err
	}
	ret := &models.APINamespace{}
	if err := json.Unmarshal(content, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// QueryNamespaceList query namespace list cache by fullPath
func (c *Client) QueryNamespaceList(fullPath string) (*models.APINamespaceList, error) {
	content, err := c.QueryCache(fullPath, CacheNSList)
	if err != nil {
		return nil, err
	}
	ret := &models.APINamespaceList{}
	if err := json.Unmarshal(content, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// QueryService query service cache by fullPath
func (c *Client) QueryService(fullPath string) (*models.APIService, error) {
	content, err := c.QueryCache(fullPath, CacheService)
	if err != nil {
		return nil, err
	}
	ret := &models.APIService{}
	if err := json.Unmarshal(content, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// QueryServiceList query service list cache by fullPath
func (c *Client) QueryServiceList(fullPath string) (*models.APIServiceList, error) {
	content, err := c.QueryCache(fullPath, CacheServiceList)
	if err != nil {
		return nil, err
	}
	ret := &models.APIServiceList{}
	if err := json.Unmarshal(content, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// QuerySchema query schema cache by fullPath
func (c *Client) QuerySchema(fullPath string) (*models.APISchemaFull, error) {
	content, err := c.QueryCache(fullPath, CacheSchema)
	if err != nil {
		return nil, err
	}
	ret := &models.APISchemaFull{}
	if err := json.Unmarshal(content, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// QuerySchemaList query schema list cache by fullPath
func (c *Client) QuerySchemaList(fullPath string) (*models.APISchemaList, error) {
	content, err := c.QueryCache(fullPath, CacheSchemaList)
	if err != nil {
		return nil, err
	}
	ret := &models.APISchemaList{}
	if err := json.Unmarshal(content, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
