package store

import (
	"dqueue/setting"
	"github.com/go-redis/redis/v8"
)

const (
	JobInfoKeyPrefix = "dq:jobInfo:"
	QueueKey         = "dq:queueTask"
	ReadyKey         = "dq:readyTask"
)

var RedisCli *redis.Client

func RedisStoreSetup() {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:         setting.RedisSetting.Host,
		Password:     setting.RedisSetting.Password, // no password set
		DB:           setting.RedisSetting.DB,       // use default DB
		PoolSize:     setting.RedisSetting.PoolSize, // 连接池大小
		MinIdleConns: setting.RedisSetting.MinIdleConns,
		DialTimeout:  setting.RedisSetting.DialTimeout,
		ReadTimeout:  setting.RedisSetting.ReadTimeout,
		WriteTimeout: setting.RedisSetting.WriteTimeout,
	})
}
