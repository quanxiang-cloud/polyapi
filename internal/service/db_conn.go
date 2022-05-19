package service

import (
	"github.com/quanxiang-cloud/cabin/logger"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	myredis "github.com/quanxiang-cloud/polyapi/internal/models/redis"
	"github.com/quanxiang-cloud/polyapi/pkg/config"

	"gorm.io/gorm"
)

var (
	mysqlDBInst     *gorm.DB
	redisClientInst *myredis.Client
)

func createMysqlConn(conf *config.Config) (*gorm.DB, error) {
	if mysqlDBInst == nil {
		db, err := mysql2.New(conf.Mysql, logger.Logger)
		if err != nil {
			return nil, err
		}
		mysqlDBInst = db
	}
	return mysqlDBInst, nil
}

func createRedisConn(conf *config.Config) (*myredis.Client, error) {
	if redisClientInst == nil {
		c, err := redis2.NewClient(conf.Redis)
		if err != nil {
			return nil, err
		}
		redisCache := myredis.NewRedisClient(c)
		redisClientInst = redisCache
	}
	return redisClientInst, nil
}
