package cache

import (
	"fmt"
	"strings"
)

// CacheType тип кеша
type CacheType string

const (
	RedisCacheType CacheType = "redis"
)

// Factory создает экземпляр кеша указанного типа
type Factory struct{}

// NewFactory создает новую фабрику кеша
func NewFactory() *Factory {
	return &Factory{}
}

// Create создает кеш указанного типа
func (f *Factory) Create(cacheType CacheType, config interface{}) (Cache, error) {
	switch strings.ToLower(string(cacheType)) {
	case string(RedisCacheType):
		return f.createRedisCache(config)
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cacheType)
	}
}

// createRedisCache создает Redis кеш
func (f *Factory) createRedisCache(config interface{}) (Cache, error) {
	redisConfig, ok := config.(RedisConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type for Redis cache")
	}

	return NewRedisCache(redisConfig)
}

// CreateFromString создает кеш из строки подключения
func (f *Factory) CreateFromString(cacheType CacheType, connectionString string, password string, db int) (Cache, error) {
	switch strings.ToLower(string(cacheType)) {
	case string(RedisCacheType):
		return NewRedisCacheFromString(connectionString, password, db)
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cacheType)
	}
}
