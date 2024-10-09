package common

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type StateConfiguration struct {
	Database *DatabaseConfiguration
	Redis    *RedisConfiguration
	DataPath string
}

func NewStateConfiguration() *StateConfiguration {
	return &StateConfiguration{
		Database: &DatabaseConfiguration{},
		Redis:    &RedisConfiguration{},
		DataPath: ".data",
	}
}

type State struct {
	Database     *gorm.DB
	Redis        *redis.Client
	RedisContext *context.Context
	Storage      Storage
}

func NewState(config *StateConfiguration) (*State, error) {
	db, err := CreateDatabaseSession(config.Database)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	rdb, err := CreateRedisSession(ctx, config.Redis)
	if err != nil {
		return nil, err
	}

	storage := NewFileStorage(config.DataPath)
	err = storage.EnsureDefaultAvatar()
	if err != nil {
		return nil, err
	}

	return &State{
		Database:     db,
		Storage:      storage,
		Redis:        rdb,
		RedisContext: &ctx,
	}, nil
}
