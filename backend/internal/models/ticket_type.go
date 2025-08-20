package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TicketType 票券類型模型
type TicketType struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EventID          uuid.UUID      `gorm:"type:uuid;not null"`
	Name             string         `gorm:"type:varchar(100);not null"`
	Price            float64        `gorm:"type:decimal(10,2);not null"`
	TotalQuantity    int            `gorm:"not null"`
	AvailableQuantity int           `gorm:"not null"`
	SaleStart        time.Time      `gorm:"not null"`
	SaleEnd          time.Time      `gorm:"not null"`
	CreatedAt        time.Time      `gorm:"not null;default:now()"`
	UpdatedAt        time.Time      `gorm:"not null;default:now()"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	OrderItems       []OrderItem    `gorm:"foreignKey:TicketTypeID"`
}

// BeforeCreate 在創建前生成 UUID
func (t *TicketType) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
