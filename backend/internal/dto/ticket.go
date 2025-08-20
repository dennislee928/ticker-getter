package dto

import "time"

// 創建票種請求
type CreateTicketTypeRequest struct {
	Name             string    `json:"name" binding:"required" example:"VIP票"`
	Price            float64   `json:"price" binding:"required,min=0" example:"2000"`
	TotalQuantity    int       `json:"total_quantity" binding:"required,min=1" example:"100"`
	AvailableQuantity int      `json:"available_quantity" binding:"omitempty,min=0" example:"100"`
	SaleStart        time.Time `json:"sale_start" binding:"required" example:"2024-07-01T10:00:00+08:00"`
	SaleEnd          time.Time `json:"sale_end" binding:"required" example:"2024-08-14T23:59:59+08:00"`
}

// 更新票種請求
type UpdateTicketTypeRequest struct {
	Name             string    `json:"name" example:"VIP票 (更新)"`
	Price            float64   `json:"price" binding:"omitempty,min=0" example:"2200"`
	TotalQuantity    int       `json:"total_quantity" binding:"omitempty,min=1" example:"120"`
	AvailableQuantity int      `json:"available_quantity" binding:"omitempty,min=0" example:"120"`
	SaleStart        time.Time `json:"sale_start" example:"2024-07-01T10:00:00+08:00"`
	SaleEnd          time.Time `json:"sale_end" example:"2024-08-14T23:59:59+08:00"`
}

// 票券查詢參數
type TicketQueryParams struct {
	Code string `form:"code" binding:"required" example:"TICKET123456"`
}

// 購買票券請求
type PurchaseTicketRequest struct {
	TicketTypeID string `json:"ticket_type_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Quantity     int    `json:"quantity" binding:"required,min=1,max=10" example:"2"`
}

// 訂單創建請求
type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,dive,required"`
}

// 訂單項目請求
type OrderItemRequest struct {
	TicketTypeID string `json:"ticket_type_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Quantity     int    `json:"quantity" binding:"required,min=1,max=10" example:"2"`
}
