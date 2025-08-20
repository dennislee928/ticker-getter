package dto

// 登入請求 DTO
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// 註冊請求 DTO
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2" example:"張三"`
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Phone    string `json:"phone" binding:"required" example:"0912345678"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}

// 重設密碼請求 DTO
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// 重設密碼確認 DTO
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required" example:"1234567890abcdef"`
	Password string `json:"password" binding:"required,min=8" example:"newPassword123"`
}
