package redis

import (
	"time"

	"go.uber.org/zap"

	"github.com/gomodule/redigo/redis"
	"github.com/wt5858/go-ddd-api/infrastructure/conf"
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/log"
	"go.uber.org/fx"
)

var (
	redisPool *redis.Pool
)

var Module = fx.Provide(func(cfg *conf.Config, logger *log.Logger) *redis.Pool {
	redisPool = &redis.Pool{
		MaxIdle:     cfg.RedisConf.MaxIdle,                   // 最初链接数量
		MaxActive:   cfg.RedisConf.MaxActive,                 // 最大链接数量
		IdleTimeout: cfg.RedisConf.IdleTimeout * time.Second, // 超时时间
		Dial: func() (redis.Conn, error) { // 链接Redis
			setPassword := redis.DialPassword(cfg.RedisConf.Auth)
			dial, err := redis.Dial(cfg.RedisConf.Protocol, cfg.RedisConf.Host+":"+cfg.RedisConf.Port, setPassword)
			if err != nil {
				logger.ZapLogger.Error(
					"[redis-conn-error]",
					zap.Any("module", "redis"),
					zap.Any("type", "cache"),
					zap.Any("info", cfg.RedisConf.Protocol+cfg.RedisConf.Host+":"+cfg.RedisConf.Port),
				)
			}
			dial.Do("SELECT", cfg.RedisConf.Db) // 选择数据库
			return dial, err
		},
	}
	return redisPool
})
