package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ticket 票券模型
type Ticket struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderItemID  uuid.UUID      `gorm:"type:uuid;not null"`
	TicketCode   string         `gorm:"type:varchar(100);not null;unique"`
	IsUsed       bool           `gorm:"not null;default:false"`
	UsedAt       *time.Time     `gorm:""`
	CreatedAt    time.Time      `gorm:"not null;default:now()"`
	UpdatedAt    time.Time      `gorm:"not null;default:now()"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate 在創建前生成 UUID 和票券碼
func (t *Ticket) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	
	// 如果沒有票券碼，自動生成一個
	if t.TicketCode == "" {
		t.TicketCode = uuid.New().String()
	}
	
	return nil
}
