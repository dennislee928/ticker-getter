package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Order 訂單模型
type Order struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null"`
	TotalAmount   float64        `gorm:"type:decimal(10,2);not null"`
	Status        string         `gorm:"type:varchar(20);not null;default:'pending'"` // pending, paid, cancelled
	PaymentMethod string         `gorm:"type:varchar(50)"`
	PaymentStatus string         `gorm:"type:varchar(20);default:'unpaid'"` // unpaid, paid, refunded
	CreatedAt     time.Time      `gorm:"not null;default:now()"`
	UpdatedAt     time.Time      `gorm:"not null;default:now()"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	OrderItems    []OrderItem    `gorm:"foreignKey:OrderID"`
}

// BeforeCreate 在創建前生成 UUID
func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
