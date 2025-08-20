package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter Redis 實現的頻率限制器
type RateLimiter struct {
	redisClient *redis.Client
	prefix      string
}

// NewRateLimiter 創建新的 RateLimiter 實例
func NewRateLimiter(redisClient *redis.Client, prefix string) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		prefix:      prefix,
	}
}

// Allow 檢查是否允許請求通過
// key: 唯一標識符（通常為 IP 或 用戶 ID）
// limit: 在時間窗口內允許的最大請求數
// window: 時間窗口（例如 1 分鐘）
func (l *RateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, int, time.Duration, error) {
	// 建立 Redis key
	redisKey := fmt.Sprintf("%s:%s", l.prefix, key)
	
	// 獲取目前計數
	count, err := l.redisClient.Get(ctx, redisKey).Int()
	
	// 如果 key 不存在，初始化為 1 並設置過期時間
	if err == redis.Nil {
		_, err = l.redisClient.Set(ctx, redisKey, 1, window).Result()
		if err != nil {
			return false, 0, 0, err
		}
		return true, limit - 1, window, nil
	} else if err != nil {
		return false, 0, 0, err
	}
	
	// 如果計數超過限制，拒絕請求
	if count >= limit {
		// 獲取剩餘過期時間
		ttl, err := l.redisClient.TTL(ctx, redisKey).Result()
		if err != nil {
			return false, 0, 0, err
		}
		return false, 0, ttl, nil
	}
	
	// 增加計數
	_, err = l.redisClient.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, 0, 0, err
	}
	
	// 獲取剩餘過期時間
	ttl, err := l.redisClient.TTL(ctx, redisKey).Result()
	if err != nil {
		return false, 0, 0, err
	}
	
	// 如果 TTL 為負數，重新設置過期時間
	if ttl < 0 {
		l.redisClient.Expire(ctx, redisKey, window)
		ttl = window
	}
	
	// 返回允許通過、剩餘可用數量和 TTL
	return true, limit - count - 1, ttl, nil
}

// Reset 重置特定鍵的限制
func (l *RateLimiter) Reset(ctx context.Context, key string) error {
	redisKey := fmt.Sprintf("%s:%s", l.prefix, key)
	return l.redisClient.Del(ctx, redisKey).Err()
}

// AllowFunc 特定功能的限制器
type AllowFunc func(ctx context.Context, key string) (bool, int, time.Duration, error)

// ForIPAndPath 基於 IP 和路徑的限制器
func (l *RateLimiter) ForIPAndPath(limit int, window time.Duration) AllowFunc {
	return func(ctx context.Context, key string) (bool, int, time.Duration, error) {
		return l.Allow(ctx, key, limit, window)
	}
}

// ForUser 基於用戶 ID 的限制器
func (l *RateLimiter) ForUser(limit int, window time.Duration) AllowFunc {
	return func(ctx context.Context, userID string) (bool, int, time.Duration, error) {
		key := fmt.Sprintf("user:%s", userID)
		return l.Allow(ctx, key, limit, window)
	}
}
