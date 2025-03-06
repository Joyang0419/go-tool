package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-tool/pkg/cache"
)

func TestManager_AddCache_And_Caches(t *testing.T) {
	// 初始化 manager
	manager := cache.NewManager()

	// 創建第一個緩存實例
	dataLoader1 := &testDataLoader{}
	expiration1 := 1 * time.Second
	config1 := cache.BigCacheConfig(expiration1)
	cache1, err := cache.NewBigCache[string, string](config1, "cache1", dataLoader1)
	if err != nil {
		t.Fatalf("failed to create first BigCache: %v", err)
	}

	// 創建第二個緩存實例
	dataLoader2 := &testDataLoader{}
	expiration2 := 2 * time.Second
	config2 := cache.BigCacheConfig(expiration2)
	cache2, err := cache.NewBigCache[string, string](config2, "cache2", dataLoader2)
	if err != nil {
		t.Fatalf("failed to create second BigCache: %v", err)
	}

	// 測試 AddCache
	manager.AddCache(cache1)
	manager.AddCache(cache2)

	// 測試 Caches
	caches := manager.Caches()
	t.Log(caches)
	assert.Equal(t, 2, len(caches))
	assert.Contains(t, caches, "cache1")
	assert.Contains(t, caches, "cache2")
}

func TestManager_ResetCache(t *testing.T) {
	// 初始化 manager
	manager := cache.NewManager()

	// 創建緩存實例
	dataLoader := &testDataLoader{}
	expiration := 1 * time.Second
	config := cache.BigCacheConfig(expiration)
	c, err := cache.NewBigCache[string, string](config, "testCache", dataLoader)
	if err != nil {
		t.Fatalf("failed to create BigCache: %v", err)
	}

	// 添加緩存到管理器
	manager.AddCache(c)

	// 設置一個值
	key := "test_key"
	value := "test_value"
	err = c.Set(context.Background(), key, value)
	if err != nil {
		t.Fatalf("failed to set value: %v", err)
	}

	// 重置緩存
	err = manager.ResetCache(context.Background(), "testCache")
	if err != nil {
		t.Fatalf("failed to reset cache: %v", err)
	}

	// 確認值已被重置
	val, err := c.Get(context.Background(), key)
	assert.Error(t, err, "應該獲取不到值")
	assert.Nil(t, val, "重置後應該為 nil")
}

func TestManager_ResetAllCaches(t *testing.T) {
	// 初始化 manager
	manager := cache.NewManager()

	// 創建兩個緩存實例
	dataLoader1 := &testDataLoader{}
	expiration1 := 1 * time.Second
	config1 := cache.BigCacheConfig(expiration1)
	cache1, _ := cache.NewBigCache[string, string](config1, "cache1", dataLoader1)

	dataLoader2 := &testDataLoader{}
	expiration2 := 2 * time.Second
	config2 := cache.BigCacheConfig(expiration2)
	cache2, _ := cache.NewBigCache[string, string](config2, "cache2", dataLoader2)

	// 添加緩存到管理器
	manager.AddCache(cache1)
	manager.AddCache(cache2)

	// 設置值
	key1 := "key1"
	value1 := "value1"
	key2 := "key2"
	value2 := "value2"

	_ = cache1.Set(context.Background(), key1, value1)
	_ = cache2.Set(context.Background(), key2, value2)

	// 重置所有緩存
	err := manager.ResetAllCaches(context.Background())
	assert.NoError(t, err)

	// 確認所有值都被重置
	val1, err1 := cache1.Get(context.Background(), key1)
	val2, err2 := cache2.Get(context.Background(), key2)

	assert.Error(t, err1, "cache1 應該獲取不到值")
	assert.Error(t, err2, "cache2 應該獲取不到值")
	assert.Nil(t, val1, "cache1 重置後應該沒有值")
	assert.Nil(t, val2, "cache2 重置後應該沒有值")
}
