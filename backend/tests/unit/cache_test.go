package unit

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/lipeichen/ticket-getter/pkg/cache"
	"github.com/redis/go-redis/v9"
)

func TestRedisCache(t *testing.T) {
	// 使用 miniredis 模擬 Redis
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("無法啟動 miniredis: %v", err)
	}
	defer mr.Close()

	// 創建 Redis 客戶端
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	// 創建 RedisCache 實例
	redisCache := cache.NewRedisCache(client)

	// 測試數據
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	testData := TestData{
		Name:  "test",
		Value: 42,
	}
	key := "test_key"
	ctx := context.Background()

	// 測試 Set 和 Get
	t.Run("Set and Get", func(t *testing.T) {
		// 設置快取
		err := redisCache.Set(ctx, key, testData, 1*time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		// 獲取快取
		var result TestData
		err = redisCache.Get(ctx, key, &result)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		// 驗證結果
		if result.Name != testData.Name || result.Value != testData.Value {
			t.Errorf("Expected %v, got %v", testData, result)
		}
	})

	// 測試 Delete
	t.Run("Delete", func(t *testing.T) {
		// 設置快取
		err := redisCache.Set(ctx, key, testData, 1*time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		// 刪除快取
		err = redisCache.Delete(ctx, key)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		// 嘗試獲取已刪除的快取
		var result TestData
		err = redisCache.Get(ctx, key, &result)
		if err == nil {
			t.Error("Expected error when getting deleted key, got nil")
		}
	})

	// 測試 Exists
	t.Run("Exists", func(t *testing.T) {
		// 設置快取
		err := redisCache.Set(ctx, key, testData, 1*time.Minute)
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}

		// 檢查是否存在
		exists, err := redisCache.Exists(ctx, key)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if !exists {
			t.Error("Expected key to exist, got false")
		}

		// 刪除快取
		err = redisCache.Delete(ctx, key)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		// 檢查是否不存在
		exists, err = redisCache.Exists(ctx, key)
		if err != nil {
			t.Fatalf("Exists failed: %v", err)
		}
		if exists {
			t.Error("Expected key to not exist, got true")
		}
	})
}
