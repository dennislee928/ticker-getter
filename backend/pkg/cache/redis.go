package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache 提供 Redis 快取功能
type RedisCache struct {
	Client *redis.Client
}

// NewRedisCache 創建新的 RedisCache 實例
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		Client: client,
	}
}

// Get 從快取獲取數據
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("key not found")
		}
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Set 設置快取數據
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.Client.Set(ctx, key, json, expiration).Err()
}

// Delete 刪除快取數據
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

// Exists 檢查key是否存在
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	val, err := c.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// HashSet 設置哈希表字段值
func (c *RedisCache) HashSet(ctx context.Context, key, field string, value interface{}) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Client.HSet(ctx, key, field, json).Err()
}

// HashGet 獲取哈希表字段值
func (c *RedisCache) HashGet(ctx context.Context, key, field string, dest interface{}) error {
	val, err := c.Client.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errors.New("field not found")
		}
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// HashDelete 刪除哈希表字段
func (c *RedisCache) HashDelete(ctx context.Context, key string, fields ...string) error {
	return c.Client.HDel(ctx, key, fields...).Err()
}

// Increment 遞增值
func (c *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return c.Client.Incr(ctx, key).Result()
}

// IncrementBy 按指定值遞增
func (c *RedisCache) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.Client.IncrBy(ctx, key, value).Result()
}

// SetExpire 設置過期時間
func (c *RedisCache) SetExpire(ctx context.Context, key string, expiration time.Duration) error {
	return c.Client.Expire(ctx, key, expiration).Err()
}

// ClearByPattern 清除符合模式的所有鍵
func (c *RedisCache) ClearByPattern(ctx context.Context, pattern string) error {
	iter := c.Client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.Client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
