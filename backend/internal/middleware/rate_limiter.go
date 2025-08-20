package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lipeichen/ticket-getter/pkg/limiter"
	"github.com/redis/go-redis/v9"
)

// RateLimiter 基於 Redis 的請求頻率限制中間件
func RateLimiter(client *redis.Client) gin.HandlerFunc {
	rateLimiter := limiter.NewRateLimiter(client, "rate_limit")
	
	return func(c *gin.Context) {
		// 獲取客戶端 IP 作為識別
		clientIP := c.ClientIP()
		
		// 根據路徑分類限制
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		
		// 鍵包含 IP 和路徑
		key := fmt.Sprintf("%s:%s", clientIP, path)
		
		// 設定限制: 60秒內最多100個請求
		limit := 100
		window := 60 * time.Second
		
		// 對於敏感操作（例如購買票券）設置更嚴格的限制
		if c.FullPath() == "/api/v1/tickets/purchase" || 
		   c.FullPath() == "/api/v1/orders" {
			limit = 10
			window = 60 * time.Second
		}
		
		allowed, remaining, retryAfter, err := rateLimiter.Allow(context.Background(), key, limit, window)
		if err != nil {
			// Redis 錯誤，繼續處理請求
			fmt.Printf("Redis 錯誤: %v\n", err)
			c.Next()
			return
		}
		
		// 設定標頭
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		
		// 檢查是否超過限制
		if !allowed {
			c.Header("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "請求頻率過高，請稍後再試",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// UserRateLimiter 基於用戶 ID 的請求頻率限制中間件
func UserRateLimiter(client *redis.Client) gin.HandlerFunc {
	rateLimiter := limiter.NewRateLimiter(client, "user_rate_limit")
	
	return func(c *gin.Context) {
		// 從上下文中獲取用戶 ID（由 Auth 中間件設置）
		userID, exists := c.Get("userID")
		if !exists {
			// 如果沒有用戶 ID，跳過此中間件
			c.Next()
			return
		}
		
		userIDStr, ok := userID.(string)
		if !ok {
			c.Next()
			return
		}
		
		// 根據路徑分類限制
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		
		// 鍵包含用戶 ID 和路徑
		key := fmt.Sprintf("%s:%s", userIDStr, path)
		
		// 設定限制: 購買操作每分鐘最多 5 次
		limit := 5
		window := 60 * time.Second
		
		allowed, remaining, retryAfter, err := rateLimiter.Allow(context.Background(), key, limit, window)
		if err != nil {
			// Redis 錯誤，繼續處理請求
			fmt.Printf("Redis 錯誤: %v\n", err)
			c.Next()
			return
		}
		
		// 設定標頭
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		
		// 檢查是否超過限制
		if !allowed {
			c.Header("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "購買操作頻率過高，請稍後再試",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}