package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements the Cache interface using a Redis backend.
// It provides methods to get and set key-value pairs in Redis with optional expiration.
type RedisCache struct {
	redis *redis.Client   // Redis client instance
	ctx   context.Context // Context for Redis operations
}

// NewRedisCache creates a new RedisCache with the given Redis client.
func NewRedisCache(redis *redis.Client) *RedisCache {
	return &RedisCache{
		redis: redis,
		ctx:   context.Background(),
	}
}

// Get retrieves the value for a given key from Redis.
// Returns the value as a string or an error if the key does not exist or on failure.
func (r *RedisCache) Get(key string) (string, error) {
	return r.redis.Get(r.ctx, key).Result()
}

// Set stores a key-value pair in Redis with the specified expiration duration.
// If expire is 0, the key does not expire.
func (r *RedisCache) Set(key string, value string, expire time.Duration) {
	r.redis.Set(r.ctx, key, value, expire)
}

// Incr atomically increments the integer value of a key by one in Redis.
// Returns the new value as int64 or an error if the operation fails.
func (r *RedisCache) Incr(key string) (int64, error) {
	return r.redis.Incr(r.ctx, key).Result()
}

// Expire sets a timeout on a key in Redis. After the timeout, the key will be automatically deleted.
// If expire is 0, the key will not expire.
func (r *RedisCache) Expire(key string, expire time.Duration) {
	r.redis.Expire(r.ctx, key, expire).Result()
}
