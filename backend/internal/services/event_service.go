package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lipeichen/ticket-getter/internal/models"
	"github.com/lipeichen/ticket-getter/internal/vo"
	"github.com/lipeichen/ticket-getter/pkg/cache"
	"gorm.io/gorm"
)

// EventService 處理事件相關邏輯
type EventService struct {
	DB          *gorm.DB
	EventCache  *cache.EventCache
	TicketCache *cache.TicketCache
}

// NewEventService 創建新的 EventService 實例
func NewEventService(db *gorm.DB, eventCache *cache.EventCache, ticketCache *cache.TicketCache) *EventService {
	return &EventService{
		DB:          db,
		EventCache:  eventCache,
		TicketCache: ticketCache,
	}
}

// GetEvents 獲取事件列表
func (s *EventService) GetEvents(page, limit int) ([]vo.EventResponse, int64, error) {
	// 檢查是否有快取
	ctx := context.Background()
	cachedEvents, err := s.EventCache.GetEventList(ctx)
	if err == nil && len(cachedEvents) > 0 {
		// 計算分頁
		var total int64 = int64(len(cachedEvents))
		start := (page - 1) * limit
		end := start + limit
		if start >= len(cachedEvents) {
			return []vo.EventResponse{}, total, nil
		}
		if end > len(cachedEvents) {
			end = len(cachedEvents)
		}

		// 將模型轉換為 VO
		events := make([]vo.EventResponse, end-start)
		for i, event := range cachedEvents[start:end] {
			events[i] = vo.EventResponse{
				ID:          event.ID,
				Title:       event.Title,
				Description: event.Description,
				Location:    event.Location,
				StartTime:   event.StartTime,
				EndTime:     event.EndTime,
				CreatedAt:   event.CreatedAt,
				UpdatedAt:   event.UpdatedAt,
			}
		}

		return events, total, nil
	}

	// 沒有快取或快取錯誤，從數據庫獲取
	var events []models.Event
	var total int64

	// 獲取總數
	if err := s.DB.Model(&models.Event{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 獲取分頁數據
	if err := s.DB.
		Order("start_time ASC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&events).Error; err != nil {
		return nil, 0, err
	}

	// 將模型轉換為 VO
	eventResponses := make([]vo.EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = vo.EventResponse{
			ID:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			Location:    event.Location,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}
	}

	// 將事件列表存入快取
	if len(events) > 0 {
		go s.EventCache.SetEventList(ctx, events)
	}

	return eventResponses, total, nil
}

// GetEventByID 根據 ID 獲取事件
func (s *EventService) GetEventByID(id uuid.UUID) (*vo.EventDetailResponse, error) {
	// 檢查是否有快取
	ctx := context.Background()
	cachedEvent, err := s.EventCache.GetEvent(ctx, id.String())
	if err == nil {
		// 獲取票種信息
		ticketTypes, _ := s.GetEventTicketTypes(id)

		// 將模型轉換為 VO
		eventResponse := &vo.EventDetailResponse{
			ID:          cachedEvent.ID,
			Title:       cachedEvent.Title,
			Description: cachedEvent.Description,
			Location:    cachedEvent.Location,
			StartTime:   cachedEvent.StartTime,
			EndTime:     cachedEvent.EndTime,
			CreatedAt:   cachedEvent.CreatedAt,
			UpdatedAt:   cachedEvent.UpdatedAt,
			TicketTypes: ticketTypes,
		}

		return eventResponse, nil
	}

	// 沒有快取或快取錯誤，從數據庫獲取
	var event models.Event
	if err := s.DB.First(&event, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("事件不存在")
		}
		return nil, err
	}

	// 獲取票種信息
	ticketTypes, _ := s.GetEventTicketTypes(id)

	// 將模型轉換為 VO
	eventResponse := &vo.EventDetailResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
		TicketTypes: ticketTypes,
	}

	// 將事件存入快取
	go s.EventCache.SetEvent(ctx, &event)

	return eventResponse, nil
}

// CreateEvent 創建新事件
func (s *EventService) CreateEvent(event *models.Event) (*vo.EventResponse, error) {
	// 在數據庫中創建事件
	if err := s.DB.Create(event).Error; err != nil {
		return nil, err
	}

	// 清除事件列表快取
	ctx := context.Background()
	go s.EventCache.DeleteEventList(ctx)

	// 將模型轉換為 VO
	eventResponse := &vo.EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}

	return eventResponse, nil
}

// UpdateEvent 更新事件
func (s *EventService) UpdateEvent(event *models.Event) (*vo.EventResponse, error) {
	// 查找要更新的事件
	var existingEvent models.Event
	if err := s.DB.First(&existingEvent, event.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("事件不存在")
		}
		return nil, err
	}

	// 更新事件
	updateData := map[string]interface{}{
		"title":       event.Title,
		"description": event.Description,
		"location":    event.Location,
		"start_time":  event.StartTime,
		"end_time":    event.EndTime,
		"updated_at":  time.Now(),
	}

	if err := s.DB.Model(&existingEvent).Updates(updateData).Error; err != nil {
		return nil, err
	}

	// 刷新事件以獲取更新後的數據
	if err := s.DB.First(&existingEvent, event.ID).Error; err != nil {
		return nil, err
	}

	// 清除快取
	ctx := context.Background()
	go s.EventCache.DeleteEvent(ctx, event.ID.String())
	go s.EventCache.DeleteEventList(ctx)

	// 將模型轉換為 VO
	eventResponse := &vo.EventResponse{
		ID:          existingEvent.ID,
		Title:       existingEvent.Title,
		Description: existingEvent.Description,
		Location:    existingEvent.Location,
		StartTime:   existingEvent.StartTime,
		EndTime:     existingEvent.EndTime,
		CreatedAt:   existingEvent.CreatedAt,
		UpdatedAt:   existingEvent.UpdatedAt,
	}

	return eventResponse, nil
}

// DeleteEvent 刪除事件
func (s *EventService) DeleteEvent(id uuid.UUID) error {
	// 在事務中刪除事件及相關數據
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 查找要刪除的事件
		var event models.Event
		if err := tx.First(&event, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("事件不存在")
			}
			return err
		}

		// 獲取事件的票種
		var ticketTypes []models.TicketType
		if err := tx.Where("event_id = ?", id).Find(&ticketTypes).Error; err != nil {
			return err
		}

		// 逐一處理每個票種
		for _, tt := range ticketTypes {
			// 查找與該票種相關的訂單項
			var orderItems []models.OrderItem
			if err := tx.Where("ticket_type_id = ?", tt.ID).Find(&orderItems).Error; err != nil {
				return err
			}

			// 處理每個訂單項
			for _, item := range orderItems {
				// 刪除票券
				if err := tx.Where("order_item_id = ?", item.ID).Delete(&models.Ticket{}).Error; err != nil {
					return err
				}
			}

			// 刪除訂單項
			if err := tx.Where("ticket_type_id = ?", tt.ID).Delete(&models.OrderItem{}).Error; err != nil {
				return err
			}

			// 清除票種快取
			ctx := context.Background()
			go s.TicketCache.DeleteTicketType(ctx, tt.ID.String())
		}

		// 刪除票種
		if err := tx.Where("event_id = ?", id).Delete(&models.TicketType{}).Error; err != nil {
			return err
		}

		// 刪除事件
		if err := tx.Delete(&event).Error; err != nil {
			return err
		}

		// 清除快取
		ctx := context.Background()
		go s.EventCache.DeleteEvent(ctx, id.String())
		go s.EventCache.DeleteEventList(ctx)
		go s.TicketCache.DeleteEventTicketTypes(ctx, id.String())

		return nil
	})
}

// GetEventTicketTypes 獲取事件的票種
func (s *EventService) GetEventTicketTypes(eventID uuid.UUID) ([]vo.TicketTypeResponse, error) {
	// 檢查是否有快取
	ctx := context.Background()
	cachedTicketTypes, err := s.TicketCache.GetEventTicketTypes(ctx, eventID.String())
	if err == nil && len(cachedTicketTypes) > 0 {
		// 將模型轉換為 VO
		ticketTypeResponses := make([]vo.TicketTypeResponse, len(cachedTicketTypes))
		for i, tt := range cachedTicketTypes {
			ticketTypeResponses[i] = vo.TicketTypeResponse{
				ID:               tt.ID,
				EventID:          tt.EventID,
				Name:             tt.Name,
				Price:            tt.Price,
				TotalQuantity:    tt.TotalQuantity,
				AvailableQuantity: tt.AvailableQuantity,
				SaleStart:        tt.SaleStart,
				SaleEnd:          tt.SaleEnd,
				CreatedAt:        tt.CreatedAt,
				UpdatedAt:        tt.UpdatedAt,
			}
		}
		return ticketTypeResponses, nil
	}

	// 沒有快取或快取錯誤，從數據庫獲取
	var ticketTypes []models.TicketType
	if err := s.DB.Where("event_id = ?", eventID).Order("price ASC").Find(&ticketTypes).Error; err != nil {
		return nil, err
	}

	// 將模型轉換為 VO
	ticketTypeResponses := make([]vo.TicketTypeResponse, len(ticketTypes))
	for i, tt := range ticketTypes {
		ticketTypeResponses[i] = vo.TicketTypeResponse{
			ID:               tt.ID,
			EventID:          tt.EventID,
			Name:             tt.Name,
			Price:            tt.Price,
			TotalQuantity:    tt.TotalQuantity,
			AvailableQuantity: tt.AvailableQuantity,
			SaleStart:        tt.SaleStart,
			SaleEnd:          tt.SaleEnd,
			CreatedAt:        tt.CreatedAt,
			UpdatedAt:        tt.UpdatedAt,
		}
	}

	// 將票種列表存入快取
	if len(ticketTypes) > 0 {
		go s.TicketCache.SetEventTicketTypes(ctx, eventID.String(), ticketTypes)
	}

	return ticketTypeResponses, nil
}

// CreateTicketType 創建票種
func (s *EventService) CreateTicketType(ticketType *models.TicketType) (*vo.TicketTypeResponse, error) {
	// 在數據庫中創建票種
	if err := s.DB.Create(ticketType).Error; err != nil {
		return nil, err
	}

	// 清除相關快取
	ctx := context.Background()
	go s.TicketCache.DeleteEventTicketTypes(ctx, ticketType.EventID.String())

	// 將模型轉換為 VO
	ticketTypeResponse := &vo.TicketTypeResponse{
		ID:               ticketType.ID,
		EventID:          ticketType.EventID,
		Name:             ticketType.Name,
		Price:            ticketType.Price,
		TotalQuantity:    ticketType.TotalQuantity,
		AvailableQuantity: ticketType.AvailableQuantity,
		SaleStart:        ticketType.SaleStart,
		SaleEnd:          ticketType.SaleEnd,
		CreatedAt:        ticketType.CreatedAt,
		UpdatedAt:        ticketType.UpdatedAt,
	}

	return ticketTypeResponse, nil
}

// SearchEvents 搜索事件
func (s *EventService) SearchEvents(query string, page, limit int) ([]vo.EventResponse, int64, error) {
	var events []models.Event
	var total int64

	// 構建搜索查詢
	searchQuery := "%" + query + "%"
	
	// 獲取總數
	if err := s.DB.Model(&models.Event{}).
		Where("title LIKE ? OR description LIKE ? OR location LIKE ?", 
			searchQuery, searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 獲取分頁數據
	if err := s.DB.
		Where("title LIKE ? OR description LIKE ? OR location LIKE ?", 
			searchQuery, searchQuery, searchQuery).
		Order("start_time ASC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&events).Error; err != nil {
		return nil, 0, err
	}

	// 將模型轉換為 VO
	eventResponses := make([]vo.EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = vo.EventResponse{
			ID:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			Location:    event.Location,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}
	}

	return eventResponses, total, nil
}

// GetFeaturedEvents 獲取精選事件
func (s *EventService) GetFeaturedEvents(limit int) ([]vo.EventResponse, error) {
	var events []models.Event

	// 獲取即將開始的活動
	now := time.Now()
	if err := s.DB.
		Where("start_time > ?", now).
		Order("start_time ASC").
		Limit(limit).
		Find(&events).Error; err != nil {
		return nil, err
	}

	// 如果沒有足夠的即將開始的活動，獲取最近的活動
	if len(events) < limit {
		var recentEvents []models.Event
		if err := s.DB.
			Where("start_time <= ?", now).
			Order("start_time DESC").
			Limit(limit - len(events)).
			Find(&recentEvents).Error; err != nil {
			return nil, err
		}
		events = append(events, recentEvents...)
	}

	// 將模型轉換為 VO
	eventResponses := make([]vo.EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = vo.EventResponse{
			ID:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			Location:    event.Location,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}
	}

	return eventResponses, nil
}
