package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lipeichen/ticket-getter/config"
	"github.com/lipeichen/ticket-getter/internal/controllers"
	"github.com/lipeichen/ticket-getter/internal/middleware"
	"github.com/lipeichen/ticket-getter/internal/services"
	"github.com/lipeichen/ticket-getter/pkg/cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// RegisterRoutes 注冊所有 API 路由
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client) {
	cfg := config.LoadConfig()

	// 初始化快取
	redisCache := cache.NewRedisCache(redisClient)
	eventCache := cache.NewEventCache(redisCache)
	ticketCache := cache.NewTicketCache(redisCache)

	// 初始化服務
	authService := services.NewAuthService(db, cfg)
	ticketService := services.NewTicketService(db, redisClient)

	// 初始化控制器
	authController := controllers.NewAuthController(authService)
	ticketController := controllers.NewTicketController(ticketService)

	// 公開路由
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/refresh", authController.RefreshToken)
	}

	// 票券可用性檢查（公開路由）
	ticketRoutes := router.Group("/tickets")
	{
		ticketRoutes.GET("/check-availability/:ticket_type_id", ticketController.CheckAvailability)
		ticketRoutes.GET("/check-fingerprint/:ticket_type_id", ticketController.CheckFingerprint)
	}

	// 需要認證的路由
	authenticatedRoutes := router.Group("")
	authenticatedRoutes.Use(middleware.AuthRequired())
	{
		// 認證相關路由
		authRoutes := authenticatedRoutes.Group("/auth")
		{
			authRoutes.GET("/me", authController.Me)
			authRoutes.POST("/logout", authController.Logout)
		}

		// 用戶相關路由 (後續添加)
		userRoutes := authenticatedRoutes.Group("/users")
		{
			// TODO: 添加用戶相關路由
		}

		// 活動相關路由 (後續添加)
		eventRoutes := authenticatedRoutes.Group("/events")
		{
			// TODO: 添加活動相關路由
		}

		// 票種相關路由 (後續添加)
		ticketTypeRoutes := authenticatedRoutes.Group("/ticket-types")
		{
			// TODO: 添加票種相關路由
		}

		// 訂單相關路由 (後續添加)
		orderRoutes := authenticatedRoutes.Group("/orders")
		{
			// 購買票券 (添加使用者限流中間件)
			orderRoutes.Use(middleware.UserRateLimiter(redisClient))
			
			// TODO: 添加訂單相關路由
		}

		// 票券相關路由
		ticketAuthRoutes := authenticatedRoutes.Group("/tickets")
		{
			ticketAuthRoutes.GET("/validate/:ticket_code", ticketController.ValidateTicket)
			ticketAuthRoutes.POST("/use/:ticket_code", ticketController.UseTicket)
		}
	}

	// 需要管理員權限的路由
	adminRoutes := router.Group("")
	adminRoutes.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		// 管理員相關路由 (後續添加)
		adminEventRoutes := adminRoutes.Group("/admin/events")
		{
			// TODO: 添加管理員活動相關路由
		}

		adminUserRoutes := adminRoutes.Group("/admin/users")
		{
			// TODO: 添加管理員用戶相關路由
		}

		adminOrderRoutes := adminRoutes.Group("/admin/orders")
		{
			// TODO: 添加管理員訂單相關路由
		}
	}
}