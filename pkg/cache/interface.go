package cache

import (
	"context"
)

// IManager 定義緩存管理器的接口; 管理多個緩存; 主要目前是設計來刷新緩存; 好比DB改資料了，線上要即時刷緩存
type IManager interface {
	AddCache(cache IManageable)
	ResetCache(ctx context.Context, name string) error
	ResetAllCaches(ctx context.Context) error
	Caches() []string
}

// IManageable 一個通用的緩存管理接口
type IManageable interface {
	// Name 返回緩存名稱
	Name() string
	// Reset 清除緩存中的所有項目
	Reset(ctx context.Context) error
}

// ICache 定義緩存操作的接口
type ICache[KEY comparable, VALUE any] interface {
	IManageable

	// Set 將值存儲到緩存中
	Set(ctx context.Context, key KEY, val VALUE) error

	// Get 從緩存中獲取項目
	Get(ctx context.Context, key KEY) (*VALUE, error)

	// Delete 從緩存中刪除項目
	Delete(ctx context.Context, key KEY) error

	// Preload 使用預加載鍵值對預填充緩存
	Preload(ctx context.Context) error
}

// IDataLoader 定義數據加載的接口
type IDataLoader[KEY comparable, VALUE any] interface {
	// Load 從數據源加載單個數據項
	Load(ctx context.Context, key KEY) (*VALUE, error)

	// PreLoad 返回需要預加載的鍵值對
	PreLoad(ctx context.Context) (map[KEY]VALUE, error)
}
