package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"go-tool/pkg/cache"
)

var TestData = map[string]string{
	"preload_key": "value1",
}

// testDataLoader 是一個簡單的模擬數據加載器，用於測試
type testDataLoader struct{}

func (m *testDataLoader) Load(_ context.Context, key string) (*string, error) {
	if val, ok := TestData[key]; ok {
		return &val, nil
	}
	return nil, errors.New("data not found") // 模擬未找到時返回空字符串
}

func (m *testDataLoader) PreLoad(_ context.Context) (map[string]string, error) {
	return TestData, nil
}

func TestBigCache_Init(t *testing.T) {
	dataLoader := &testDataLoader{}
	expiration := 1 * time.Second
	config := cache.BigCacheConfig(expiration)
	name := "testCache"

	c, err := cache.NewBigCache(config, name, dataLoader)
	if err != nil {
		t.Fatalf("failed to create BigCache: %v", err)
	}

	if c.Name() != name {
		t.Errorf("expected name test, got %s", c.Name())
	}
}

func TestBigCache_Get(t *testing.T) {
	dataLoader := &testDataLoader{}
	expiration := 1 * time.Second
	config := cache.BigCacheConfig(expiration)
	name := "testCache"

	c, _ := cache.NewBigCache(config, name, dataLoader)

	key := "preload_key"
	val, err := c.Get(context.Background(), key)
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	if *val != TestData[key] {
		t.Errorf("expected value %s, got %v", TestData[key], val)
	}
}

func TestBigCache_Set(t *testing.T) {
	dataLoader := &testDataLoader{}
	expiration := 1 * time.Second
	config := cache.BigCacheConfig(expiration)
	name := "testCache"

	c, _ := cache.NewBigCache(config, name, dataLoader)

	key := "new_key"
	value := "new_value"
	err := c.Set(context.TODO(), key, value)
	if err != nil {
		t.Fatalf("failed to set value: %v", err)
	}

	val, err := c.Get(context.Background(), key)
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	if *val != value {
		t.Errorf("expected value %s, got %v", value, val)
	}
}

func TestBigCache_Delete(t *testing.T) {
	dataLoader := &testDataLoader{}
	expiration := 1 * time.Second
	config := cache.BigCacheConfig(expiration)
	name := "testCache"

	c, _ := cache.NewBigCache(config, name, dataLoader)

	key := "hello"
	_ = c.Set(context.Background(), key, "world")

	err := c.Delete(context.TODO(), key)
	if err != nil {
		t.Fatalf("failed to delete value: %v", err)
	}

	_, err = c.Get(context.Background(), key)
	assert.Error(t, err, "expected error")
}

func TestBigCache_Reset(t *testing.T) {
	dataLoader := &testDataLoader{}
	expiration := 1 * time.Second
	config := cache.BigCacheConfig(expiration)
	name := "testCache"

	c, _ := cache.NewBigCache(config, name, dataLoader)

	// 故意set 一個資料沒有TestData的Key
	key := "hello world"
	err := c.Set(context.Background(), key, "value")
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	// 確認有值
	val, _ := c.Get(context.Background(), key)
	if *val != "value" {
		t.Errorf("expected value value, got %v", val)
	}

	err = c.Reset(context.Background())
	if err != nil {
		t.Fatalf("failed to reset cache: %v", err)
	}

	// reset後確認沒有值
	_, err = c.Get(context.Background(), key)
	assert.Error(t, err)
}

func TestBigCache_Preload(t *testing.T) {
	dataLoader := &testDataLoader{}
	expiration := 1 * time.Second
	config := cache.BigCacheConfig(expiration)
	name := "testCache"

	c, _ := cache.NewBigCache(config, name, dataLoader)

	// 測試 Preload
	val, _ := c.Get(context.Background(), "preload_key")

	if *val != "value1" {
		t.Errorf("expected value value1, got %v", val)
	}
}
