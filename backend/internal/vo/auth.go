package vo

// 使用者資訊 VO
type UserResponse struct {
	ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name  string `json:"name" example:"張三"`
	Email string `json:"email" example:"user@example.com"`
	Phone string `json:"phone" example:"0912345678"`
	Role  string `json:"role" example:"user"`
}

// 登入響應 VO
type LoginResponse struct {
	Token        string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User         UserResponse `json:"user"`
}

// 註冊響應 VO
type RegisterResponse struct {
	Token        string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User         UserResponse `json:"user"`
}

// 基本響應 VO
type BaseResponse struct {
	Message string `json:"message" example:"操作成功"`
}
