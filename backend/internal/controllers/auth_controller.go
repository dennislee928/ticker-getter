package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lipeichen/ticket-getter/internal/dto"
	"github.com/lipeichen/ticket-getter/internal/services"
)

// AuthController 處理認證相關 HTTP 請求
type AuthController struct {
	AuthService *services.AuthService
}

// NewAuthController 創建新的 AuthController 實例
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

// Register 處理用戶註冊
// @Summary 用戶註冊
// @Description 創建新用戶帳戶
// @Tags 認證
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "用戶註冊信息"
// @Success 201 {object} vo.RegisterResponse "註冊成功"
// @Failure 400 {object} map[string]string "無效的輸入"
// @Failure 409 {object} map[string]string "用戶已存在"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無效的輸入: " + err.Error()})
		return
	}
	
	response, err := c.AuthService.Register(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusCreated, response)
}

// Login 處理用戶登入
// @Summary 用戶登入
// @Description 用戶登入並獲取認證令牌
// @Tags 認證
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "用戶登入憑證"
// @Success 200 {object} vo.LoginResponse "登入成功"
// @Failure 400 {object} map[string]string "無效的輸入"
// @Failure 401 {object} map[string]string "登入失敗"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無效的輸入: " + err.Error()})
		return
	}
	
	response, err := c.AuthService.Login(req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, response)
}

// Me 獲取當前登入用戶信息
// @Summary 獲取當前用戶
// @Description 獲取當前登入用戶的信息
// @Tags 認證
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} vo.UserResponse "用戶信息"
// @Failure 401 {object} map[string]string "未認證"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Router /auth/me [get]
func (c *AuthController) Me(ctx *gin.Context) {
	// 從上下文中獲取用戶 ID（由 Auth 中間件設置）
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未經認證"})
		return
	}
	
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "無效的用戶 ID 類型"})
		return
	}
	
	user, err := c.AuthService.GetUserByID(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, user)
}

// RefreshToken 刷新認證令牌
// @Summary 刷新令牌
// @Description 使用刷新令牌獲取新的認證令牌
// @Tags 認證
// @Accept json
// @Produce json
// @Param refresh_token body map[string]string true "刷新令牌"
// @Success 200 {object} map[string]string "新的令牌對"
// @Failure 400 {object} map[string]string "無效的輸入"
// @Failure 401 {object} map[string]string "無效的刷新令牌"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req map[string]string
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無效的輸入"})
		return
	}
	
	refreshToken, exists := req["refresh_token"]
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "必須提供刷新令牌"})
		return
	}
	
	newToken, newRefreshToken, err := c.AuthService.RefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"token":         newToken,
		"refresh_token": newRefreshToken,
	})
}

// Logout 處理用戶登出
// @Summary 用戶登出
// @Description 使當前用戶登出
// @Tags 認證
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "登出成功"
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	// 在前端處理令牌清理，這裡只返回成功響應
	ctx.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}
