package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lipeichen/ticket-getter/internal/dto"
	"github.com/lipeichen/ticket-getter/internal/services"
)

// EventController 處理事件相關 HTTP 請求
type EventController struct {
	EventService *services.EventService
}

// NewEventController 創建新的 EventController 實例
func NewEventController(eventService *services.EventService) *EventController {
	return &EventController{
		EventService: eventService,
	}
}

// GetEvents 獲取事件列表
// @Summary 獲取事件列表
// @Description 獲取所有事件的分頁列表
// @Tags 事件
// @Accept json
// @Produce json
// @Param page query int false "頁碼，默認為 1"
// @Param limit query int false "每頁數量，默認為 10"
// @Success 200 {object} vo.EventListResponse "事件列表"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Router /events [get]
func (c *EventController) GetEvents(ctx *gin.Context) {
	var params dto.EventQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	events, total, err := c.EventService.GetEvents(params.Page, params.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "獲取事件列表失敗"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"events": events,
		"total":  total,
		"page":   params.Page,
		"limit":  params.Limit,
	})
}

// GetEvent 獲取單個事件詳情
// @Summary 獲取事件詳情
// @Description 根據ID獲取事件詳情
// @Tags 事件
// @Accept json
// @Produce json
// @Param id path string true "事件 ID"
// @Success 200 {object} vo.EventDetailResponse "事件詳情"
// @Failure 400 {object} map[string]string "無效的 ID"
// @Failure 404 {object} map[string]string "事件不存在"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Router /events/{id} [get]
func (c *EventController) GetEvent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無效的事件 ID"})
		return
	}

	event, err := c.EventService.GetEventByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "事件不存在"})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

// SearchEvents 搜索事件
// @Summary 搜索事件
// @Description 根據關鍵字搜索事件
// @Tags 事件
// @Accept json
// @Produce json
// @Param q query string true "搜索關鍵字"
// @Param page query int false "頁碼，默認為 1"
// @Param limit query int false "每頁數量，默認為 10"
// @Success 200 {object} vo.EventListResponse "搜索結果"
// @Failure 400 {object} map[string]string "無效的請求"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Router /events/search [get]
func (c *EventController) SearchEvents(ctx *gin.Context) {
	var params dto.EventSearchParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	events, total, err := c.EventService.SearchEvents(params.Query, params.Page, params.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "搜索事件失敗"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"events": events,
		"total":  total,
		"page":   params.Page,
		"limit":  params.Limit,
	})
}

// GetEventTicketTypes 獲取事件的票種列表
// @Summary 獲取事件票種
// @Description 獲取指定事件的所有票種
// @Tags 事件
// @Accept json
// @Produce json
// @Param id path string true "事件 ID"
// @Success 200 {object} vo.TicketTypeListResponse "票種列表"
// @Failure 400 {object} map[string]string "無效的 ID"
// @Failure 404 {object} map[string]string "事件不存在"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Router /events/{id}/ticket-types [get]
func (c *EventController) GetEventTicketTypes(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無效的事件 ID"})
		return
	}

	ticketTypes, err := c.EventService.GetEventTicketTypes(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "獲取票種失敗"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"ticket_types": ticketTypes,
	})
}

// GetFeaturedEvents 獲取精選事件
// @Summary 獲取精選事件
// @Description 獲取精選或熱門事件列表
// @Tags 事件
// @Accept json
// @Produce json
// @Param limit query int false "數量限制，默認為 6"
// @Success 200 {array} vo.EventResponse "精選事件列表"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Router /events/featured [get]
func (c *EventController) GetFeaturedEvents(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "6")
	limit := 6
	ctx.ScanJSON(&limit, limitStr)

	events, err := c.EventService.GetFeaturedEvents(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "獲取精選事件失敗"})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
