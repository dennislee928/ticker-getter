package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter 基於 Redis 的請求頻率限制中間件
func RateLimiter(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 獲取客戶端 IP 作為識別
		clientIP := c.ClientIP()
		
		// 建立 Redis key
		key := fmt.Sprintf("rate_limit:%s", clientIP)
		
		// 設定限制: 60秒內最多100個請求
		limit := 100
		window := 60 * time.Second
		
		ctx := context.Background()
		
		// 獲取目前計數
		count, err := client.Get(ctx, key).Int()
		if err == redis.Nil {
			// key 不存在，設置新的計數器
			_, err = client.Set(ctx, key, 1, window).Result()
			if err != nil {
				fmt.Printf("Redis 錯誤: %v\n", err)
				c.Next()
				return
			}
			c.Header("X-RateLimit-Remaining", strconv.Itoa(limit-1))
			c.Next()
			return
		} else if err != nil {
			// Redis 錯誤，繼續處理請求
			fmt.Printf("Redis 錯誤: %v\n", err)
			c.Next()
			return
		}
		
		// 增加計數
		count++
		
		// 檢查是否超過限制
		if count > limit {
			c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
			c.Header("X-RateLimit-Remaining", "0")
			
			// 計算重試時間
			ttl, _ := client.TTL(ctx, key).Result()
			c.Header("Retry-After", strconv.Itoa(int(ttl.Seconds())))
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "請求頻率過高，請稍後再試",
			})
			c.Abort()
			return
		}
		
		// 更新 Redis 計數
		_, err = client.Set(ctx, key, count, window).Result()
		if err != nil {
			fmt.Printf("Redis 錯誤: %v\n", err)
		}
		
		// 設定剩餘請求數量的標頭
		remaining := limit - count
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		
		c.Next()
	}
}
