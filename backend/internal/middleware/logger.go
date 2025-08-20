package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日誌中間件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 請求前
		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 繼續處理請求
		c.Next()

		// 請求後
		endTime := time.Now()
		latency := endTime.Sub(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// 日誌格式
		logFormat := fmt.Sprintf("[%s] | %3d | %13v | %15s | %-7s %s",
			endTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)

		// 根據狀態碼選擇日誌級別
		switch {
		case statusCode >= 500:
			fmt.Println(logFormat, "| 伺服器錯誤")
		case statusCode >= 400:
			fmt.Println(logFormat, "| 客戶端錯誤")
		default:
			fmt.Println(logFormat)
		}
	}
}
