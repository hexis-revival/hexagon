package common

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type RedisConfiguration struct {
	Host     string
	Port     int
	Password string
	Database int
}

func CreateRedisSession(ctx context.Context, config *RedisConfiguration) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Database,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to redis")
	}

	return rdb, nil
}
