package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用戶模型
type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email         string         `gorm:"type:varchar(255);not null;unique"`
	PasswordHash  string         `gorm:"type:varchar(255);not null"`
	Name          string         `gorm:"type:varchar(100);not null"`
	Phone         string         `gorm:"type:varchar(20)"`
	Role          string         `gorm:"type:varchar(20);not null;default:'user'"` // user 或 admin
	TLSFingerprint string        `gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `gorm:"not null;default:now()"`
	UpdatedAt     time.Time      `gorm:"not null;default:now()"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Events        []Event        `gorm:"foreignKey:CreatedBy"`
	Orders        []Order        `gorm:"foreignKey:UserID"`
}

// BeforeCreate 在創建前生成 UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
