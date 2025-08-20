package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lipeichen/ticket-getter/internal/services"
	"github.com/lipeichen/ticket-getter/pkg/utils"
)

// TicketController 處理票券相關 HTTP 請求
type TicketController struct {
	TicketService *services.TicketService
}

// NewTicketController 創建新的 TicketController 實例
func NewTicketController(ticketService *services.TicketService) *TicketController {
	return &TicketController{
		TicketService: ticketService,
	}
}

// CheckAvailability 檢查票券可用性
// @Summary 檢查票券可用性
// @Description 檢查指定票券類型是否可用
// @Tags 票券
// @Accept json
// @Produce json
// @Param ticket_type_id path string true "票券類型 ID"
// @Param quantity query int false "需要的數量，默認為 1"
// @Success 200 {object} map[string]interface{} "票券可用狀態"
// @Failure 400 {object} map[string]string "無效的請求"
// @Failure 404 {object} map[string]string "票券類型不存在"
// @Failure 409 {object} map[string]string "票券不可用"
// @Router /tickets/check-availability/{ticket_type_id} [get]
func (c *TicketController) CheckAvailability(ctx *gin.Context) {
	ticketTypeID := ctx.Param("ticket_type_id")
	
	// 獲取請求數量
	quantity := 1
	if quantityStr := ctx.Query("quantity"); quantityStr != "" {
		if _, err := fmt.Sscanf(quantityStr, "%d", &quantity); err != nil || quantity <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "無效的數量"})
			return
		}
	}
	
	// 檢查可用性
	available, err := c.TicketService.CheckAvailability(ticketTypeID, quantity)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"available": available,
		"quantity":  quantity,
	})
}

// CheckFingerprint 檢查指紋是否已購買特定票券
// @Summary 檢查指紋購買狀態
// @Description 檢查當前設備是否已購買特定票券
// @Tags 票券
// @Accept json
// @Produce json
// @Param ticket_type_id path string true "票券類型 ID"
// @Success 200 {object} map[string]interface{} "指紋檢查結果"
// @Failure 400 {object} map[string]string "無效的請求"
// @Router /tickets/check-fingerprint/{ticket_type_id} [get]
func (c *TicketController) CheckFingerprint(ctx *gin.Context) {
	ticketTypeID := ctx.Param("ticket_type_id")
	
	// 提取 TLS 指紋
	fingerprint := utils.GetClientIPFingerprint(ctx.Request)
	
	// 檢查是否已購買
	alreadyPurchased, err := c.TicketService.CheckFingerprint(ctx, fingerprint, ticketTypeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "檢查失敗"})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"already_purchased": alreadyPurchased,
	})
}

// ValidateTicket 驗證票券有效性
// @Summary 驗證票券
// @Description 驗證票券碼是否有效
// @Tags 票券
// @Accept json
// @Produce json
// @Param ticket_code path string true "票券碼"
// @Success 200 {object} map[string]interface{} "票券驗證結果"
// @Failure 400 {object} map[string]string "無效的請求"
// @Failure 404 {object} map[string]string "票券不存在"
// @Router /tickets/validate/{ticket_code} [get]
func (c *TicketController) ValidateTicket(ctx *gin.Context) {
	// TODO: 實現票券驗證邏輯
	ctx.JSON(http.StatusOK, gin.H{
		"message": "票券驗證功能尚未實現",
	})
}

// UseTicket 使用票券
// @Summary 使用票券
// @Description 標記票券為已使用狀態
// @Tags 票券
// @Accept json
// @Produce json
// @Param ticket_code path string true "票券碼"
// @Success 200 {object} map[string]interface{} "票券使用結果"
// @Failure 400 {object} map[string]string "無效的請求"
// @Failure 404 {object} map[string]string "票券不存在"
// @Failure 409 {object} map[string]string "票券已被使用"
// @Router /tickets/use/{ticket_code} [post]
func (c *TicketController) UseTicket(ctx *gin.Context) {
	// TODO: 實現使用票券邏輯
	ctx.JSON(http.StatusOK, gin.H{
		"message": "使用票券功能尚未實現",
	})
}
