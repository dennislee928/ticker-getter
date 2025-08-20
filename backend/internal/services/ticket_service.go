package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lipeichen/ticket-getter/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// TicketService 處理票券相關業務邏輯
type TicketService struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

// NewTicketService 創建新的 TicketService 實例
func NewTicketService(db *gorm.DB, redisClient *redis.Client) *TicketService {
	return &TicketService{
		DB:          db,
		RedisClient: redisClient,
	}
}

// CheckFingerprint 檢查指紋是否已經購買過特定票券
func (s *TicketService) CheckFingerprint(ctx context.Context, fingerprint string, ticketTypeID string) (bool, error) {
	// 建立 Redis key
	key := fmt.Sprintf("fingerprint:%s:%s", fingerprint, ticketTypeID)
	
	// 檢查是否存在
	exists, err := s.RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	
	return exists > 0, nil
}

// RecordFingerprint 記錄指紋購買了特定票券
func (s *TicketService) RecordFingerprint(ctx context.Context, fingerprint string, ticketTypeID string, userID string) error {
	// 建立 Redis key
	key := fmt.Sprintf("fingerprint:%s:%s", fingerprint, ticketTypeID)
	
	// 存儲信息
	info := fmt.Sprintf("user_id=%s,timestamp=%s", userID, time.Now().Format(time.RFC3339))
	
	// 設置過期時間為 24 小時
	_, err := s.RedisClient.Set(ctx, key, info, 24*time.Hour).Result()
	return err
}

// CheckAvailability 檢查票券是否可用
func (s *TicketService) CheckAvailability(ticketTypeID string, quantity int) (bool, error) {
	var ticketType models.TicketType
	
	// 解析票券類型 ID
	id, err := uuid.Parse(ticketTypeID)
	if err != nil {
		return false, errors.New("無效的票券類型 ID")
	}
	
	// 查詢票券類型
	result := s.DB.First(&ticketType, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, errors.New("票券類型不存在")
		}
		return false, result.Error
	}
	
	// 檢查銷售時間
	now := time.Now()
	if now.Before(ticketType.SaleStart) {
		return false, errors.New("票券銷售尚未開始")
	}
	
	if now.After(ticketType.SaleEnd) {
		return false, errors.New("票券銷售已結束")
	}
	
	// 檢查剩餘數量
	if ticketType.AvailableQuantity < quantity {
		return false, errors.New("票券數量不足")
	}
	
	return true, nil
}

// UpdateAvailability 更新票券可用數量
func (s *TicketService) UpdateAvailability(ticketTypeID string, quantity int) error {
	// 解析票券類型 ID
	id, err := uuid.Parse(ticketTypeID)
	if err != nil {
		return errors.New("無效的票券類型 ID")
	}
	
	// 在事務中更新票券數量
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var ticketType models.TicketType
		
		// 查詢票券類型並鎖定行
		result := tx.Set("gorm:query_option", "FOR UPDATE").First(&ticketType, id)
		if result.Error != nil {
			return result.Error
		}
		
		// 檢查剩餘數量
		if ticketType.AvailableQuantity < quantity {
			return errors.New("票券數量不足")
		}
		
		// 更新剩餘數量
		ticketType.AvailableQuantity -= quantity
		
		// 保存更改
		return tx.Save(&ticketType).Error
	})
}

// GenerateTickets 生成票券
func (s *TicketService) GenerateTickets(orderItemID string, quantity int) error {
	// 解析訂單項目 ID
	id, err := uuid.Parse(orderItemID)
	if err != nil {
		return errors.New("無效的訂單項目 ID")
	}
	
	// 在事務中創建票券
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var orderItem models.OrderItem
		
		// 查詢訂單項目
		result := tx.First(&orderItem, id)
		if result.Error != nil {
			return result.Error
		}
		
		// 創建票券
		tickets := make([]models.Ticket, quantity)
		for i := 0; i < quantity; i++ {
			tickets[i] = models.Ticket{
				OrderItemID: orderItem.ID,
				IsUsed:      false,
			}
		}
		
		// 批量插入票券
		return tx.Create(&tickets).Error
	})
}
