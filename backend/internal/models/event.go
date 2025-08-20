package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Event 活動模型
type Event struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Location    string         `gorm:"type:varchar(255);not null"`
	StartTime   time.Time      `gorm:"not null"`
	EndTime     time.Time      `gorm:"not null"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null"`
	CreatedAt   time.Time      `gorm:"not null;default:now()"`
	UpdatedAt   time.Time      `gorm:"not null;default:now()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	TicketTypes []TicketType   `gorm:"foreignKey:EventID"`
}

// BeforeCreate 在創建前生成 UUID
func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
