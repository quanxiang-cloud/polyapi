package redis

import (
	"context"
	"time"
)

const (
	apiStatCacheTimeout = time.Second * 30
)

func getStatType(raw bool) CacheType {
	if raw {
		return CacheGateRawStat
	}
	return CacheGatePolyStat
}

func (c *Client) getStatKeyDur() time.Duration {
	if c.apiStatKeyDur > 0 {
		return c.apiStatKeyDur
	}
	return apiStatCacheTimeout
}

// SetAPIStatKeyDuration SetAPIStatKeyDuration
func (c *Client) SetAPIStatKeyDuration(dur time.Duration) error {
	if dur < apiStatCacheTimeout {
		dur = apiStatCacheTimeout
	}
	c.apiStatKeyDur = dur
	return nil
}

// IncAPIStat IncAPIStat
func (c *Client) IncAPIStat(apiPath string, raw bool) (int64, error) {
	key := getStatType(raw).GetCacheKey(apiPath)
	r, err := c.C.Incr(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	c.C.Expire(context.Background(), key, c.getStatKeyDur())
	return r, err
}

// GetAPIStat GetAPIStat
func (c *Client) GetAPIStat(apiPath string, raw bool) (int64, error) {
	key := getStatType(raw).GetCacheKey(apiPath)
	r, err := c.C.Get(context.Background(), key).Int64()
	return r, err
}

//------------------------------------------------------------------------------
