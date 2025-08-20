package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lipeichen/ticket-getter/config"
)

// JWTClaims 定義 JWT 聲明結構
type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthRequired 驗證 JWT token
func AuthRequired() gin.HandlerFunc {
	cfg := config.LoadConfig()
	
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供認證令牌"})
			c.Abort()
			return
		}
		
		// 從 Bearer Token 中提取令牌
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		
		// 解析令牌
		token, err := jwt.ParseWithClaims(
			tokenString,
			&JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				// 檢查簽名方法
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("無效的簽名方法: %v", token.Header["alg"])
				}
				return []byte(cfg.JWTSecret), nil
			},
		)
		
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的認證令牌"})
			c.Abort()
			return
		}
		
		// 驗證令牌並提取聲明
		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			// 將聲明信息設置在上下文中
			c.Set("userID", claims.UserID)
			c.Set("role", claims.Role)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認證令牌已過期或無效"})
			c.Abort()
			return
		}
	}
}

// AdminRequired 檢查用戶是否為管理員角色
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 檢查是否已經通過了 AuthRequired 中間件的認證
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未經認證"})
			c.Abort()
			return
		}
		
		// 檢查用戶角色是否為管理員
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理員權限"})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
