package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lipeichen/ticket-getter/internal/controllers"
	"github.com/lipeichen/ticket-getter/internal/models"
	"github.com/lipeichen/ticket-getter/internal/services"
	"github.com/lipeichen/ticket-getter/pkg/cache"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 設置測試環境
func setupTestEnv(t *testing.T) (*gin.Engine, *gorm.DB) {
	// 使用 SQLite 內存數據庫進行測試
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("無法連接到測試數據庫: %v", err)
	}

	// 自動遷移表結構
	err = db.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.TicketType{},
		&models.Order{},
		&models.OrderItem{},
		&models.Ticket{},
	)
	if err != nil {
		t.Fatalf("自動遷移失敗: %v", err)
	}

	// 創建測試數據
	createTestData(t, db)

	// 創建 Gin 引擎
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// 創建一個模擬的 Redis 客戶端
	mockRedisClient := &MockRedisClient{}

	// 創建快取服務
	redisCache := cache.NewRedisCache(mockRedisClient)
	eventCache := cache.NewEventCache(redisCache)
	ticketCache := cache.NewTicketCache(redisCache)

	// 創建服務
	eventService := services.NewEventService(db, eventCache, ticketCache)

	// 創建控制器
	eventController := controllers.NewEventController(eventService)

	// 註冊路由
	apiV1 := router.Group("/api/v1")
	apiV1.GET("/events", eventController.GetEvents)
	apiV1.GET("/events/:id", eventController.GetEvent)
	apiV1.GET("/events/search", eventController.SearchEvents)
	apiV1.GET("/events/:id/ticket-types", eventController.GetEventTicketTypes)
	apiV1.GET("/events/featured", eventController.GetFeaturedEvents)

	return router, db
}

// 創建測試數據
func createTestData(t *testing.T, db *gorm.DB) {
	// 創建測試用戶
	user := models.User{
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		Name:         "測試用戶",
		Role:         "admin",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("創建用戶失敗: %v", err)
	}

	// 創建測試事件
	event := models.Event{
		Title:       "測試活動",
		Description: "這是一個測試活動",
		Location:    "測試地點",
		StartTime:   time.Now().Add(24 * time.Hour),
		EndTime:     time.Now().Add(48 * time.Hour),
		CreatedBy:   user.ID,
	}
	if err := db.Create(&event).Error; err != nil {
		t.Fatalf("創建事件失敗: %v", err)
	}

	// 創建測試票種
	ticketType := models.TicketType{
		EventID:          event.ID,
		Name:             "測試票種",
		Price:            100,
		TotalQuantity:    100,
		AvailableQuantity: 100,
		SaleStart:        time.Now().Add(-24 * time.Hour),
		SaleEnd:          time.Now().Add(48 * time.Hour),
	}
	if err := db.Create(&ticketType).Error; err != nil {
		t.Fatalf("創建票種失敗: %v", err)
	}
}

// 模擬 Redis 客戶端
type MockRedisClient struct{}

// 實現必要的接口方法
func (m *MockRedisClient) Get(ctx interface{}, key string) *redis.StringCmd {
	return redis.NewStringCmd(ctx)
}

func (m *MockRedisClient) Set(ctx interface{}, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redis.NewStatusCmd(ctx)
}

func (m *MockRedisClient) Del(ctx interface{}, keys ...string) *redis.IntCmd {
	return redis.NewIntCmd(ctx)
}

func (m *MockRedisClient) Exists(ctx interface{}, keys ...string) *redis.IntCmd {
	return redis.NewIntCmd(ctx)
}

func (m *MockRedisClient) HSet(ctx interface{}, key string, values ...interface{}) *redis.IntCmd {
	return redis.NewIntCmd(ctx)
}

func (m *MockRedisClient) HGet(ctx interface{}, key, field string) *redis.StringCmd {
	return redis.NewStringCmd(ctx)
}

func (m *MockRedisClient) HDel(ctx interface{}, key string, fields ...string) *redis.IntCmd {
	return redis.NewIntCmd(ctx)
}

func (m *MockRedisClient) Incr(ctx interface{}, key string) *redis.IntCmd {
	return redis.NewIntCmd(ctx)
}

func (m *MockRedisClient) IncrBy(ctx interface{}, key string, value int64) *redis.IntCmd {
	return redis.NewIntCmd(ctx)
}

func (m *MockRedisClient) Expire(ctx interface{}, key string, expiration time.Duration) *redis.BoolCmd {
	return redis.NewBoolCmd(ctx)
}

func (m *MockRedisClient) TTL(ctx interface{}, key string) *redis.DurationCmd {
	return redis.NewDurationCmd(ctx)
}

func (m *MockRedisClient) Scan(ctx interface{}, cursor uint64, match string, count int64) *redis.ScanCmd {
	return redis.NewScanCmd(ctx)
}

// 測試獲取事件列表 API
func TestGetEvents(t *testing.T) {
	router, _ := setupTestEnv(t)

	// 發送請求
	req := httptest.NewRequest("GET", "/api/v1/events", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 驗證響應
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 驗證響應數據
	events, ok := response["events"].([]interface{})
	assert.True(t, ok)
	assert.GreaterOrEqual(t, len(events), 1)
}

// 測試獲取單個事件 API
func TestGetEvent(t *testing.T) {
	router, db := setupTestEnv(t)

	// 獲取事件 ID
	var event models.Event
	err := db.First(&event).Error
	assert.NoError(t, err)

	// 發送請求
	req := httptest.NewRequest("GET", "/api/v1/events/"+event.ID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 驗證響應
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 驗證響應數據
	assert.Equal(t, event.Title, response["title"])
	assert.Equal(t, event.Description, response["description"])
}

// 測試搜索事件 API
func TestSearchEvents(t *testing.T) {
	router, _ := setupTestEnv(t)

	// 發送請求
	req := httptest.NewRequest("GET", "/api/v1/events/search?q=測試", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 驗證響應
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 驗證響應數據
	events, ok := response["events"].([]interface{})
	assert.True(t, ok)
	assert.GreaterOrEqual(t, len(events), 1)
}
