package cache

import (
	"github.com/christo-andrew/haven/pkg/config"
	"github.com/redis/go-redis/v9"
)

func RedisClient(config config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
}
