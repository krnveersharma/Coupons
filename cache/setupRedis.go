package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	mu  sync.Mutex
	rdb *redis.Client
	ttl time.Duration
)

// SetupRedis initializes the global Redis client
func SetupRedis(redisAddress, redisPassword string, defaultTTL time.Duration) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       0,
	})
	ttl = defaultTTL

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	if err := rdb.ConfigSet(ctx, "maxmemory", "100mb").Err(); err != nil {
		fmt.Println("Warning: failed to set maxmemory:", err)
	}
	if err := rdb.ConfigSet(ctx, "maxmemory-policy", "allkeys-lru").Err(); err != nil {
		fmt.Println("Warning: failed to set maxmemory-policy:", err)
	}

	fmt.Println("Redis server configured for LRU")
	return nil
}

func Get(key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func Set(key string, value string) error {
	mu.Lock()
	defer mu.Unlock()
	return rdb.Set(ctx, key, value, ttl).Err()
}

func Delete(key string) error {
	return rdb.Del(ctx, key).Err()
}
