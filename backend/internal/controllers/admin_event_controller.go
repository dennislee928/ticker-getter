package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lipeichen/ticket-getter/internal/dto"
	"github.com/lipeichen/ticket-getter/internal/models"
	"github.com/lipeichen/ticket-getter/internal/services"
)

// AdminEventController 處理管理員事件相關 HTTP 請求
type AdminEventController struct {
	EventService *services.EventService
}

// NewAdminEventController 創建新的 AdminEventController 實例
func NewAdminEventController(eventService *services.EventService) *AdminEventController {
	return &AdminEventController{
		EventService: eventService,
	}
}

// GetEvents 獲取所有事件
// @Summary 管理員獲取事件列表
// @Description 管理員獲取所有事件的分頁列表
// @Tags 管理員-事件
// @Accept json
// @Produce json
// @Param page query int false "頁碼，默認為 1"
// @Param limit query int false "每頁數量，默認為 10"
// @Success 200 {object} vo.EventListResponse "事件列表"
// @Failure 401 {object} map[string]string "未授權"
// @Failure 403 {object} map[string]string "禁止訪問"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Security BearerAuth
// @Router /admin/events [get]
func (c *AdminEventController) GetEvents(ctx *gin.Context) {
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

// CreateEvent 創建新事件
// @Summary 創建事件
// @Description 管理員創建新事件
// @Tags 管理員-事件
// @Accept json
// @Produce json
// @Param event body dto.CreateEventRequest true "事件信息"
// @Success 201 {object} vo.EventResponse "創建成功"
// @Failure 400 {object} map[string]string "無效的輸入"
// @Failure 401 {object} map[string]string "未授權"
// @Failure 403 {object} map[string]string "禁止訪問"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Security BearerAuth
// @Router /admin/events [post]
func (c *AdminEventController) CreateEvent(ctx *gin.Context) {
	var req dto.CreateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 從上下文中獲取用戶 ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未認證"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "無效的用戶 ID 類型"})
		return
	}

	createdByID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "無效的用戶 ID"})
		return
	}

	// 創建事件模型
	event := models.Event{
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		CreatedBy:   createdByID,
	}

	// 創建事件
	createdEvent, err := c.EventService.CreateEvent(&event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "創建事件失敗"})
		return
	}

	ctx.JSON(http.StatusCreated, createdEvent)
}

// UpdateEvent 更新事件
// @Summary 更新事件
// @Description 管理員更新現有事件
// @Tags 管理員-事件
// @Accept json
// @Produce json
// @Param id path string true "事件 ID"
// @Param event body dto.UpdateEventRequest true "更新的事件信息"
// @Success 200 {object} vo.EventResponse "更新成功"
// @Failure 400 {object} map[string]string "無效的輸入"
// @Failure 401 {object} map[string]string "未授權"
// @Failure 403 {object} map[string]string "禁止訪問"
// @Failure 404 {object} map[string]string "事件不存在"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Security BearerAuth
// @Router /admin/events/{id} [put]
func (c *AdminEventController) UpdateEvent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無效的事件 ID"})
		return
	}

	var req dto.UpdateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 檢查事件是否存在
	_, err = c.EventService.GetEventByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "事件不存在"})
		return
	}

	// 更新事件
	event := models.Event{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	}

	updatedEvent, err := c.EventService.UpdateEvent(&event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "更新事件失敗"})
		return
	}

	ctx.JSON(http.StatusOK, updatedEvent)
}

// DeleteEvent 刪除事件
// @Summary 刪除事件
// @Description 管理員刪除事件
// @Tags 管理員-事件
// @Accept json
// @Produce json
// @Param id path string true "事件 ID"
// @Success 200 {object} map[string]string "刪除成功"
// @Failure 400 {object} map[string]string "無效的輸入"
// @Failure 401 {object} map[string]string "未授權"
// @Failure 403 {object} map[string]string "禁止訪問"
// @Failure 404 {object} map[string]string "事件不存在"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Security BearerAuth
// @Router /admin/events/{id} [delete]
func (c *AdminEventController) DeleteEvent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無效的事件 ID"})
		return
	}

	// 檢查事件是否存在
	_, err = c.EventService.GetEventByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "事件不存在"})
		return
	}

	// 刪除事件
	if err := c.EventService.DeleteEvent(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "刪除事件失敗"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "事件已成功刪除"})
}

// CreateTicketType 為事件創建票種
// @Summary 創建票種
// @Description 為指定事件創建新票種
// @Tags 管理員-票種
// @Accept json
// @Produce json
// @Param id path string true "事件 ID"
// @Param ticketType body dto.CreateTicketTypeRequest true "票種信息"
// @Success 201 {object} vo.TicketTypeResponse "創建成功"
// @Failure 400 {object} map[string]string "無效的輸入"
// @Failure 401 {object} map[string]string "未授權"
// @Failure 403 {object} map[string]string "禁止訪問"
// @Failure 404 {object} map[string]string "事件不存在"
// @Failure 500 {object} map[string]string "內部服務器錯誤"
// @Security BearerAuth
// @Router /admin/events/{id}/ticket-types [post]
func (c *AdminEventController) CreateTicketType(ctx *gin.Context) {
	idStr := ctx.Param("id")
	eventID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無效的事件 ID"})
		return
	}

	// 檢查事件是否存在
	_, err = c.EventService.GetEventByID(eventID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "事件不存在"})
		return
	}

	var req dto.CreateTicketTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果沒有提供可用數量，則設為總數量
	availableQuantity := req.AvailableQuantity
	if availableQuantity == 0 {
		availableQuantity = req.TotalQuantity
	}

	// 創建票種模型
	ticketType := models.TicketType{
		EventID:          eventID,
		Name:             req.Name,
		Price:            req.Price,
		TotalQuantity:    req.TotalQuantity,
		AvailableQuantity: availableQuantity,
		SaleStart:        req.SaleStart,
		SaleEnd:          req.SaleEnd,
	}

	// 創建票種
	createdTicketType, err := c.EventService.CreateTicketType(&ticketType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "創建票種失敗"})
		return
	}

	ctx.JSON(http.StatusCreated, createdTicketType)
}
