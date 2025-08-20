package dto

import "time"

// 創建事件請求
type CreateEventRequest struct {
	Title       string    `json:"title" binding:"required" example:"2024 台北音樂節"`
	Description string    `json:"description" binding:"required" example:"年度最大音樂節，超過50組藝人演出"`
	Location    string    `json:"location" binding:"required" example:"台北市立體育場"`
	StartTime   time.Time `json:"start_time" binding:"required" example:"2024-08-15T18:00:00+08:00"`
	EndTime     time.Time `json:"end_time" binding:"required" example:"2024-08-15T22:00:00+08:00"`
}

// 更新事件請求
type UpdateEventRequest struct {
	Title       string    `json:"title" example:"2024 台北音樂節 (更新)"`
	Description string    `json:"description" example:"年度最大音樂節，超過50組藝人演出，包括國際巨星！"`
	Location    string    `json:"location" example:"台北市立體育場"`
	StartTime   time.Time `json:"start_time" example:"2024-08-15T18:00:00+08:00"`
	EndTime     time.Time `json:"end_time" example:"2024-08-15T22:00:00+08:00"`
}

// 事件列表查詢參數
type EventQueryParams struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}

// 事件搜索參數
type EventSearchParams struct {
	Query string `form:"q" binding:"required"`
	Page  int    `form:"page,default=1" binding:"min=1"`
	Limit int    `form:"limit,default=10" binding:"min=1,max=100"`
}
