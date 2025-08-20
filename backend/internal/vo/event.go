package vo

import (
	"time"

	"github.com/google/uuid"
)

// EventResponse 事件回應
type EventResponse struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string    `json:"title" example:"2024 台北音樂節"`
	Description string    `json:"description" example:"年度最大音樂節，超過50組藝人演出"`
	Location    string    `json:"location" example:"台北市立體育場"`
	StartTime   time.Time `json:"start_time" example:"2024-08-15T18:00:00+08:00"`
	EndTime     time.Time `json:"end_time" example:"2024-08-15T22:00:00+08:00"`
	CreatedAt   time.Time `json:"created_at" example:"2024-06-01T10:30:00+08:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-06-01T10:30:00+08:00"`
}

// EventDetailResponse 事件詳情回應
type EventDetailResponse struct {
	ID          uuid.UUID            `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string               `json:"title" example:"2024 台北音樂節"`
	Description string               `json:"description" example:"年度最大音樂節，超過50組藝人演出"`
	Location    string               `json:"location" example:"台北市立體育場"`
	StartTime   time.Time            `json:"start_time" example:"2024-08-15T18:00:00+08:00"`
	EndTime     time.Time            `json:"end_time" example:"2024-08-15T22:00:00+08:00"`
	CreatedAt   time.Time            `json:"created_at" example:"2024-06-01T10:30:00+08:00"`
	UpdatedAt   time.Time            `json:"updated_at" example:"2024-06-01T10:30:00+08:00"`
	TicketTypes []TicketTypeResponse `json:"ticket_types,omitempty"`
}

// EventListResponse 事件列表回應
type EventListResponse struct {
	Events []EventResponse `json:"events"`
	Total  int64           `json:"total" example:"42"`
	Page   int             `json:"page" example:"1"`
	Limit  int             `json:"limit" example:"10"`
}
