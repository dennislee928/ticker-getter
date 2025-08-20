package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lipeichen/ticket-getter/config"
	"github.com/lipeichen/ticket-getter/internal/middleware"
	"github.com/lipeichen/ticket-getter/internal/models"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Run 啟動 API 服務器
func Run() {
	// 載入配置
	cfg := config.LoadConfig()

	// 設置生產模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 連接數據庫
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("無法連接到數據庫: %v", err)
	}
	
	// 自動遷移（僅開發環境）
	if cfg.Environment == "development" {
		autoMigrateDB(db)
	}

	// 初始化 Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})
	
	// 檢查 Redis 連接
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("無法連接到 Redis: %v，繼續啟動服務但部分功能可能受限", err)
	}

	// 創建 Gin 引擎
	router := gin.Default()

	// CORS 設定
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 設置中間件
	router.Use(middleware.Logger())
	router.Use(middleware.RateLimiter(redisClient))

	// API 版本前綴
	apiV1 := router.Group("/api/v1")
	
	// 健康檢查端點
	apiV1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
	
	// 註冊路由
	RegisterRoutes(apiV1, db, redisClient)

	// Swagger 文檔
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 啟動服務器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// 在 goroutine 中啟動服務器
	go func() {
		log.Printf("服務器啟動於 :%s 端口\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服務器啟動失敗: %v\n", err)
		}
	}()

	// 優雅關閉設置
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("正在關閉服務器...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服務器關閉時發生錯誤: %v\n", err)
	}
	
	log.Println("服務器已優雅關閉")
}

// connectDB 連接到 PostgreSQL 數據庫
func connectDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	
	// 獲取底層 SQL DB 並設置連接池參數
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	
	// 設置連接池參數
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	return db, nil
}

// autoMigrateDB 自動遷移數據庫結構
func autoMigrateDB(db *gorm.DB) {
	log.Println("正在執行數據庫自動遷移...")
	
	err := db.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.TicketType{},
		&models.Order{},
		&models.OrderItem{},
		&models.Ticket{},
	)
	
	if err != nil {
		log.Fatalf("數據庫自動遷移失敗: %v", err)
	}
	
	log.Println("數據庫自動遷移完成")
}