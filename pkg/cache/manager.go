package cache

import (
	"context"
	"sync"

	pkgerrors "github.com/pkg/errors"
)

// Manager 是緩存管理器的實現，用於管理多個緩存實例
type Manager struct {
	cachesMap sync.Map // 使用 sync.Map 避免併發問題
}

// 確保 Manager 實現了 IManager 接口
var _ IManager = (*Manager)(nil)

// NewManager 創建一個新的緩存管理器
func NewManager() *Manager {
	return &Manager{
		cachesMap: sync.Map{},
	}
}

// AddCache 將緩存添加到管理器
func (receiver *Manager) AddCache(cache IManageable) {
	if cache == nil {
		return
	}

	// 使用緩存的名稱作為鍵
	receiver.cachesMap.Store(cache.Name(), cache)
}

// ResetCache 重置指定名稱的緩存
func (receiver *Manager) ResetCache(ctx context.Context, name string) error {
	cacheVal, exists := receiver.cachesMap.Load(name)
	if !exists {
		return nil
	}

	cache, ok := cacheVal.(IManageable)
	if !ok {
		return pkgerrors.New("Manager.ResetCache: cacheVal.(IManageable) type assertion failed")
	}

	// 調用緩存的 Reset 方法
	return cache.Reset(ctx)
}

// Caches 返回管理器中的所有緩存
func (receiver *Manager) Caches() []string {
	var result []string
	receiver.cachesMap.Range(func(key, value interface{}) bool {
		cache, ok := value.(IManageable)
		if ok {
			result = append(result, cache.Name())
		}
		return true // 繼續遍歷
	})

	return result
}

// ResetAllCaches 重置所有緩存
func (receiver *Manager) ResetAllCaches(ctx context.Context) error {
	var err error
	receiver.cachesMap.Range(func(key, value interface{}) bool {
		cache, ok := value.(IManageable)
		if ok {
			if errReset := cache.Reset(ctx); errReset != nil {
				err = pkgerrors.WithMessagef(err, "Manager.ResetAllCaches: %s", errReset)
			}
		}
		return true // 繼續遍歷
	})

	return err
}
