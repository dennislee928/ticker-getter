package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lipeichen/ticket-getter/config"
	"github.com/lipeichen/ticket-getter/internal/dto"
	"github.com/lipeichen/ticket-getter/internal/models"
	"github.com/lipeichen/ticket-getter/internal/vo"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 處理認證相關邏輯
type AuthService struct {
	DB     *gorm.DB
	Config *config.Config
}

// NewAuthService 創建新的 AuthService 實例
func NewAuthService(db *gorm.DB, config *config.Config) *AuthService {
	return &AuthService{
		DB:     db,
		Config: config,
	}
}

// Login 處理使用者登入
func (s *AuthService) Login(req dto.LoginRequest) (*vo.LoginResponse, error) {
	var user models.User
	
	// 查找使用者
	result := s.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("使用者不存在")
		}
		return nil, result.Error
	}
	
	// 驗證密碼
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("密碼不正確")
	}
	
	// 生成令牌
	token, refreshToken, err := s.generateTokens(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}
	
	// 建立回應
	response := &vo.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: vo.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Role:  user.Role,
		},
	}
	
	return response, nil
}

// Register 處理使用者註冊
func (s *AuthService) Register(req dto.RegisterRequest) (*vo.RegisterResponse, error) {
	// 檢查使用者是否已存在
	var existingUser models.User
	result := s.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("使用者已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	
	// 加密密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	// 創建新使用者
	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Phone:        req.Phone,
		Role:         "user", // 預設為普通用戶
	}
	
	// 儲存使用者
	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	
	// 生成令牌
	token, refreshToken, err := s.generateTokens(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}
	
	// 建立回應
	response := &vo.RegisterResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: vo.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Role:  user.Role,
		},
	}
	
	return response, nil
}

// GetUserByID 根據 ID 獲取使用者
func (s *AuthService) GetUserByID(id string) (*vo.UserResponse, error) {
	var user models.User
	
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("無效的使用者 ID")
	}
	
	result := s.DB.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("使用者不存在")
		}
		return nil, result.Error
	}
	
	response := &vo.UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Role:  user.Role,
	}
	
	return response, nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	// 解析刷新令牌
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// 檢查簽名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("無效的簽名方法: %v", token.Header["alg"])
		}
		return []byte(s.Config.JWTSecret), nil
	})
	
	if err != nil {
		return "", "", errors.New("無效的刷新令牌")
	}
	
	// 檢查令牌有效性
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 獲取使用者 ID 和角色
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", "", errors.New("無效的令牌聲明")
		}
		
		role, ok := claims["role"].(string)
		if !ok {
			return "", "", errors.New("無效的令牌聲明")
		}
		
		// 生成新的令牌對
		newToken, newRefreshToken, err := s.generateTokens(userID, role)
		if err != nil {
			return "", "", err
		}
		
		return newToken, newRefreshToken, nil
	}
	
	return "", "", errors.New("無效的刷新令牌")
}

// 生成 JWT 令牌
func (s *AuthService) generateTokens(userID string, role string) (string, string, error) {
	// 設置令牌過期時間
	expiresAt := time.Now().Add(time.Hour * time.Duration(s.Config.JWTExpiryHours))
	refreshExpiresAt := time.Now().Add(time.Hour * 24 * 7) // 一週
	
	// 創建令牌聲明
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     expiresAt.Unix(),
	}
	
	// 創建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.Config.JWTSecret))
	if err != nil {
		return "", "", err
	}
	
	// 創建刷新令牌聲明
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     refreshExpiresAt.Unix(),
	}
	
	// 創建刷新令牌
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.Config.JWTSecret))
	if err != nil {
		return "", "", err
	}
	
	return tokenString, refreshTokenString, nil
}
