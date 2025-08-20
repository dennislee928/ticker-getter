package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderItem 訂單項目模型
type OrderItem struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID       uuid.UUID      `gorm:"type:uuid;not null"`
	TicketTypeID  uuid.UUID      `gorm:"type:uuid;not null"`
	Quantity      int            `gorm:"not null"`
	PricePerUnit  float64        `gorm:"type:decimal(10,2);not null"`
	CreatedAt     time.Time      `gorm:"not null;default:now()"`
	UpdatedAt     time.Time      `gorm:"not null;default:now()"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Tickets       []Ticket       `gorm:"foreignKey:OrderItemID"`
}

// BeforeCreate 在創建前生成 UUID
func (i *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}
