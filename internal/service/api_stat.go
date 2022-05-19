package service

import (
	"context"
	"errors"
	"time"

	"github.com/quanxiang-cloud/polyapi/internal/models"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/config"

	"github.com/quanxiang-cloud/cabin/logger"
)

var errFeatureDisabled = errors.New("feature disabled")

// APIStat represents the api stat operator
type APIStat interface {
	IncTimeStat(c context.Context, apiPath string, raw bool, dur time.Duration) error
	//GetStat(c context.Context, apiPath string, raw bool) (int64, error)
	IsBlocked(c context.Context, apiPath string, raw bool) bool
}

// CreateAPIStatOper create a API stat operater
func CreateAPIStatOper(conf *config.Config) (APIStat, error) {
	redisCache, err := createRedisConn(conf)
	if err != nil {
		return nil, err
	}

	dur := time.Second * time.Duration(conf.Gate.APIBlock.BlockSeconds)
	if err := redisCache.SetAPIStatKeyDuration(dur); err != nil {
		return nil, err
	}

	p := &apiStat{
		conf:             conf,
		redisAPI:         redisCache,
		cfg:              &conf.Gate.APIBlock,
		timeoutThreshold: time.Duration(conf.Gate.APIBlock.APITimeoutMS) * time.Millisecond,
	}
	logger.Logger.Infof("[gate.apistat] config is %+v", p.cfg)

	adaptor.SetAPIStatOper(p)
	return p, nil
}

type apiStat struct {
	conf             *config.Config
	redisAPI         models.RedisCache
	cfg              *config.GateAPIBlock
	timeoutThreshold time.Duration
}

func (s *apiStat) isEnabled() bool {
	return s.cfg.Enable && s.cfg.MaxAllowError > 0
}

func (s *apiStat) IncTimeStat(c context.Context, apiPath string, raw bool, dur time.Duration) error {
	if !s.isEnabled() {
		return errFeatureDisabled
	}
	if dur < 0 || dur > s.timeoutThreshold {
		counter, err := s.redisAPI.IncAPIStat(apiPath, raw)
		if err == nil && counter >= s.cfg.MaxAllowError {
			logger.Logger.Infof("[gate.apistat] api=%t(raw),%s will blocked during next %ds", raw, apiPath, s.cfg.BlockSeconds)
		}
		return err
	}
	return nil
}

// func (s *apiStat) GetStat(c context.Context, apiPath string, raw bool) (int64, error) {
// 	if !s.isEnabled() {
// 		return 0, errFeatureDisabled
// 	}
// 	r, err := s.redisAPI.GetAPIStat(apiPath, raw)
// 	return r, err
// }

func (s *apiStat) IsBlocked(c context.Context, apiPath string, raw bool) bool {
	if !s.isEnabled() {
		return false
	}

	r, err := s.redisAPI.GetAPIStat(apiPath, raw)
	if err != nil {
		return false
	}

	return r >= s.cfg.MaxAllowError
}
