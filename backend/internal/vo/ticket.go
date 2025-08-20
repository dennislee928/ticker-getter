package vo

import (
	"time"

	"github.com/google/uuid"
)

// TicketTypeResponse 票種回應
type TicketTypeResponse struct {
	ID               uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	EventID          uuid.UUID `json:"event_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name             string    `json:"name" example:"VIP票"`
	Price            float64   `json:"price" example:"2000"`
	TotalQuantity    int       `json:"total_quantity" example:"100"`
	AvailableQuantity int      `json:"available_quantity" example:"75"`
	SaleStart        time.Time `json:"sale_start" example:"2024-07-01T10:00:00+08:00"`
	SaleEnd          time.Time `json:"sale_end" example:"2024-08-14T23:59:59+08:00"`
	CreatedAt        time.Time `json:"created_at" example:"2024-06-01T10:30:00+08:00"`
	UpdatedAt        time.Time `json:"updated_at" example:"2024-06-01T10:30:00+08:00"`
}

// TicketResponse 票券回應
type TicketResponse struct {
	ID           uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OrderItemID  uuid.UUID  `json:"order_item_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	TicketCode   string     `json:"ticket_code" example:"TICKET123456"`
	IsUsed       bool       `json:"is_used" example:"false"`
	UsedAt       *time.Time `json:"used_at,omitempty" example:"2024-08-15T19:30:00+08:00"`
	CreatedAt    time.Time  `json:"created_at" example:"2024-06-01T10:30:00+08:00"`
	UpdatedAt    time.Time  `json:"updated_at" example:"2024-06-01T10:30:00+08:00"`
	EventTitle   string     `json:"event_title" example:"2024 台北音樂節"`
	EventTime    time.Time  `json:"event_time" example:"2024-08-15T18:00:00+08:00"`
	EventLocation string    `json:"event_location" example:"台北市立體育場"`
	TicketTypeName string   `json:"ticket_type_name" example:"VIP票"`
}

// TicketTypeListResponse 票種列表回應
type TicketTypeListResponse struct {
	TicketTypes []TicketTypeResponse `json:"ticket_types"`
}

// TicketAvailabilityResponse 票券可用性回應
type TicketAvailabilityResponse struct {
	Available bool `json:"available" example:"true"`
	Quantity  int  `json:"quantity" example:"2"`
}

// TicketFingerprintResponse 票券指紋檢查回應
type TicketFingerprintResponse struct {
	AlreadyPurchased bool `json:"already_purchased" example:"false"`
}
