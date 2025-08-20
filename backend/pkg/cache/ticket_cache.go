package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/lipeichen/ticket-getter/internal/models"
)

const (
	// 票券類型快取鍵前綴
	ticketTypeCacheKeyPrefix = "ticket_type:"
	
	// 活動票券類型快取鍵前綴
	eventTicketTypesCacheKeyPrefix = "event_ticket_types:"
	
	// 票券快取鍵前綴
	ticketCacheKeyPrefix = "ticket:"
	
	// 票券類型快取過期時間 (15分鐘)
	ticketTypeCacheTTL = 15 * time.Minute
	
	// 活動票券類型快取過期時間 (10分鐘)
	eventTicketTypesCacheTTL = 10 * time.Minute
	
	// 票券快取過期時間 (30分鐘)
	ticketCacheTTL = 30 * time.Minute
)

// TicketCache 處理票券相關快取
type TicketCache struct {
	cache *RedisCache
}

// NewTicketCache 創建新的 TicketCache 實例
func NewTicketCache(redisCache *RedisCache) *TicketCache {
	return &TicketCache{
		cache: redisCache,
	}
}

// GetTicketType 從快取獲取票券類型
func (c *TicketCache) GetTicketType(ctx context.Context, ticketTypeID string) (*models.TicketType, error) {
	key := fmt.Sprintf("%s%s", ticketTypeCacheKeyPrefix, ticketTypeID)
	var ticketType models.TicketType
	
	err := c.cache.Get(ctx, key, &ticketType)
	if err != nil {
		return nil, err
	}
	
	return &ticketType, nil
}

// SetTicketType 將票券類型存入快取
func (c *TicketCache) SetTicketType(ctx context.Context, ticketType *models.TicketType) error {
	key := fmt.Sprintf("%s%s", ticketTypeCacheKeyPrefix, ticketType.ID)
	return c.cache.Set(ctx, key, ticketType, ticketTypeCacheTTL)
}

// DeleteTicketType 從快取刪除票券類型
func (c *TicketCache) DeleteTicketType(ctx context.Context, ticketTypeID string) error {
	key := fmt.Sprintf("%s%s", ticketTypeCacheKeyPrefix, ticketTypeID)
	return c.cache.Delete(ctx, key)
}

// GetEventTicketTypes 從快取獲取活動票券類型
func (c *TicketCache) GetEventTicketTypes(ctx context.Context, eventID string) ([]*models.TicketType, error) {
	key := fmt.Sprintf("%s%s", eventTicketTypesCacheKeyPrefix, eventID)
	var ticketTypes []*models.TicketType
	
	err := c.cache.Get(ctx, key, &ticketTypes)
	if err != nil {
		return nil, err
	}
	
	return ticketTypes, nil
}

// SetEventTicketTypes 將活動票券類型存入快取
func (c *TicketCache) SetEventTicketTypes(ctx context.Context, eventID string, ticketTypes []*models.TicketType) error {
	key := fmt.Sprintf("%s%s", eventTicketTypesCacheKeyPrefix, eventID)
	return c.cache.Set(ctx, key, ticketTypes, eventTicketTypesCacheTTL)
}

// DeleteEventTicketTypes 從快取刪除活動票券類型
func (c *TicketCache) DeleteEventTicketTypes(ctx context.Context, eventID string) error {
	key := fmt.Sprintf("%s%s", eventTicketTypesCacheKeyPrefix, eventID)
	return c.cache.Delete(ctx, key)
}

// GetTicket 從快取獲取票券
func (c *TicketCache) GetTicket(ctx context.Context, ticketCode string) (*models.Ticket, error) {
	key := fmt.Sprintf("%s%s", ticketCacheKeyPrefix, ticketCode)
	var ticket models.Ticket
	
	err := c.cache.Get(ctx, key, &ticket)
	if err != nil {
		return nil, err
	}
	
	return &ticket, nil
}

// SetTicket 將票券存入快取
func (c *TicketCache) SetTicket(ctx context.Context, ticket *models.Ticket) error {
	key := fmt.Sprintf("%s%s", ticketCacheKeyPrefix, ticket.TicketCode)
	return c.cache.Set(ctx, key, ticket, ticketCacheTTL)
}

// DeleteTicket 從快取刪除票券
func (c *TicketCache) DeleteTicket(ctx context.Context, ticketCode string) error {
	key := fmt.Sprintf("%s%s", ticketCacheKeyPrefix, ticketCode)
	return c.cache.Delete(ctx, key)
}
