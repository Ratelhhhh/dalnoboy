package cache

import (
	"context"
	"time"
)

// Cache интерфейс для работы с кешем
type Cache interface {
	// Set устанавливает значение в кеш с TTL
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	
	// Get получает значение из кеша
	Get(ctx context.Context, key string) (string, error)
	
	// Delete удаляет ключ из кеша
	Delete(ctx context.Context, key string) error
	
	// Exists проверяет существование ключа
	Exists(ctx context.Context, key string) (bool, error)
	
	// SetNX устанавливает значение только если ключ не существует
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	
	// Incr увеличивает значение на 1
	Incr(ctx context.Context, key string) (int64, error)
	
	// IncrBy увеличивает значение на указанное число
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	
	// Expire устанавливает TTL для ключа
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	
	// TTL получает оставшееся время жизни ключа
	TTL(ctx context.Context, key string) (time.Duration, error)
	
	// Close закрывает соединение с кешем
	Close() error
} 