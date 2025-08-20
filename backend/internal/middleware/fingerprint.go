package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// CheckTLSFingerprint 檢查 TLS 指紋，防止重複購買
func CheckTLSFingerprint(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 僅針對購票相關操作檢查指紋
		if !strings.Contains(c.Request.URL.Path, "/tickets/purchase") && 
		   !strings.Contains(c.Request.URL.Path, "/orders") {
			c.Next()
			return
		}
		
		// 獲取 TLS 指紋
		// 在實際應用中，這可能是從請求標頭或客戶端證書中提取的
		fingerprint := c.GetHeader("X-TLS-Fingerprint")
		if fingerprint == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 TLS 指紋識別"})
			c.Abort()
			return
		}
		
		// 檢查當前操作的資源（如票券 ID）
		resourceID := c.Param("id")
		if resourceID == "" {
			// 如果 URL 中沒有 ID，嘗試從請求體中獲取
			resourceID = c.PostForm("ticket_type_id")
		}
		
		if resourceID == "" {
			// 無法識別資源 ID，跳過檢查
			c.Next()
			return
		}
		
		// 建立 Redis key
		key := fmt.Sprintf("fingerprint:%s:%s", fingerprint, resourceID)
		
		ctx := context.Background()
		
		// 檢查是否已存在記錄
		exists, err := client.Exists(ctx, key).Result()
		if err != nil {
			fmt.Printf("Redis 錯誤: %v\n", err)
			c.Next() // 出錯時繼續處理
			return
		}
		
		// 如果指紋已存在，表示可能是重複購買
		if exists > 0 {
			// 獲取相關信息
			userInfo, _ := client.Get(ctx, key).Result()
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "您已經購買過此票券，請勿重複購買",
				"info":  userInfo,
			})
			c.Abort()
			return
		}
		
		// 設置一個過期時間較長的標記，以便後續檢查
		// 在實際應用中，可能需要根據業務需求調整過期時間
		userID, _ := c.Get("userID")
		userIDStr, _ := userID.(string)
		
		// 存儲信息
		info := fmt.Sprintf("購買者ID:%s,時間:%s", userIDStr, time.Now().Format(time.RFC3339))
		_, err = client.Set(ctx, key, info, 24*time.Hour).Result()
		if err != nil {
			fmt.Printf("Redis 錯誤: %v\n", err)
		}
		
		c.Next()
	}
}
