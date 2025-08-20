package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/lipeichen/ticket-getter/internal/models"
)

const (
	// 活動快取鍵前綴
	eventCacheKeyPrefix = "event:"
	
	// 活動列表快取鍵
	eventListCacheKey = "event:list"
	
	// 活動快取過期時間 (10分鐘)
	eventCacheTTL = 10 * time.Minute
	
	// 活動列表快取過期時間 (5分鐘)
	eventListCacheTTL = 5 * time.Minute
)

// EventCache 處理活動相關快取
type EventCache struct {
	cache *RedisCache
}

// NewEventCache 創建新的 EventCache 實例
func NewEventCache(redisCache *RedisCache) *EventCache {
	return &EventCache{
		cache: redisCache,
	}
}

// GetEvent 從快取獲取活動
func (c *EventCache) GetEvent(ctx context.Context, eventID string) (*models.Event, error) {
	key := fmt.Sprintf("%s%s", eventCacheKeyPrefix, eventID)
	var event models.Event
	
	err := c.cache.Get(ctx, key, &event)
	if err != nil {
		return nil, err
	}
	
	return &event, nil
}

// SetEvent 將活動存入快取
func (c *EventCache) SetEvent(ctx context.Context, event *models.Event) error {
	key := fmt.Sprintf("%s%s", eventCacheKeyPrefix, event.ID)
	return c.cache.Set(ctx, key, event, eventCacheTTL)
}

// DeleteEvent 從快取刪除活動
func (c *EventCache) DeleteEvent(ctx context.Context, eventID string) error {
	key := fmt.Sprintf("%s%s", eventCacheKeyPrefix, eventID)
	return c.cache.Delete(ctx, key)
}

// GetEventList 從快取獲取活動列表
func (c *EventCache) GetEventList(ctx context.Context) ([]*models.Event, error) {
	var events []*models.Event
	
	err := c.cache.Get(ctx, eventListCacheKey, &events)
	if err != nil {
		return nil, err
	}
	
	return events, nil
}

// SetEventList 將活動列表存入快取
func (c *EventCache) SetEventList(ctx context.Context, events []*models.Event) error {
	return c.cache.Set(ctx, eventListCacheKey, events, eventListCacheTTL)
}

// DeleteEventList 從快取刪除活動列表
func (c *EventCache) DeleteEventList(ctx context.Context) error {
	return c.cache.Delete(ctx, eventListCacheKey)
}
