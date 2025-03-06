package cache

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/allegro/bigcache/v3"
	pkgerrors "github.com/pkg/errors"
	"github.com/spf13/cast"
)

// BigCache 是使用BigCache的通用緩存實現
type BigCache[KEY comparable, VALUE any] struct {
	name       string
	dataLoader IDataLoader[KEY, VALUE]
	cache      *bigcache.BigCache
}

// 確保 BigCache 實現了 ICache 接口
var _ ICache[string, string] = (*BigCache[string, string])(nil)

// NewBigCache 創建一個新的BigCache實例
func NewBigCache[KEY comparable, VALUE any](
	config bigcache.Config,
	name string,
	dataLoader IDataLoader[KEY, VALUE],
) (*BigCache[KEY, VALUE], error) {
	if dataLoader == nil {
		return nil, errors.New("NewBigCache dataLoader is required")
	}

	bc, err := bigcache.New(context.TODO(), config)
	if err != nil {
		return nil, pkgerrors.WithMessage(err, "NewBigCache bigcache.New error")
	}

	c := &BigCache[KEY, VALUE]{
		name:       name,
		cache:      bc,
		dataLoader: dataLoader,
	}

	// 開go routine，在後台預加載緩存
	go func() {
		defer func() {
			_ = recover()
		}()
		_ = c.Preload(context.TODO())
	}()

	return c, nil
}

// Name 返回緩存名稱
func (receiver *BigCache[KEY, VALUE]) Name() string {
	return receiver.name
}

// Set 將值存儲到緩存中
func (receiver *BigCache[KEY, VALUE]) Set(_ context.Context, key KEY, val VALUE) error {
	// 將值序列化為JSON
	data, err := json.Marshal(val)
	if err != nil {
		return pkgerrors.WithMessage(err, "BigCache.Set json.Marshal error")
	}

	// 存儲到BigCache
	return receiver.cache.Set(cast.ToString(key), data)
}

func (receiver *BigCache[KEY, VALUE]) Get(ctx context.Context, key KEY) (*VALUE, error) {
	// 嘗試從緩存中獲取值
	data, err := receiver.cache.Get(cast.ToString(key))
	if err == nil {
		// 緩存命中，反序列化數據
		var val VALUE
		if unmarshalErr := json.Unmarshal(data, &val); unmarshalErr != nil {
			return nil, pkgerrors.WithMessagef(unmarshalErr, "BigCache.Get json.Unmarshal error for key %v", key)
		}

		// 返回找到的值
		return &val, nil
	}

	// 緩存未命中，檢查錯誤類型
	if errors.Is(err, bigcache.ErrEntryNotFound) {
		// 從數據加載器獲取數據
		value, loadErr := receiver.dataLoader.Load(ctx, key)
		if loadErr != nil {
			return nil, pkgerrors.WithMessagef(loadErr, "BigCache.Get dataLoader.Load for key %v", key)
		}

		// 將加載的數據存入緩存
		if cacheErr := receiver.Set(ctx, key, *value); cacheErr != nil {
			return nil, pkgerrors.WithMessagef(cacheErr, "BigCache.Get Set error for key %v", key)
		}

		// 返回加載的值
		return value, nil
	}

	// 其他未知錯誤
	return nil, pkgerrors.WithMessagef(err, "BigCache.Get cache.Get error for key %v", key)
}

// Delete 從緩存中刪除項目
func (receiver *BigCache[KEY, VALUE]) Delete(_ context.Context, key KEY) error {
	// 直接刪除指定的鍵
	err := receiver.cache.Delete(cast.ToString(key))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			// 如果條目未找到，就返回 nil（表示無需刪除）
			return nil
		}

		// 如果因其他原因刪除失敗，則記錄錯誤
		return pkgerrors.WithMessagef(err, "delete key failed: %v", key)
	}

	return nil
}

// Reset 清除緩存中的所有項目
func (receiver *BigCache[KEY, VALUE]) Reset(_ context.Context) error {
	return receiver.cache.Reset()
}

// Preload 使用預加載鍵值對預填充緩存
// 返回成功預加載的項目和可能的錯誤
func (receiver *BigCache[KEY, VALUE]) Preload(ctx context.Context) error {
	// 獲取預加載的鍵值對
	preloadItems, err := receiver.dataLoader.PreLoad(ctx)
	if err != nil {
		return pkgerrors.WithMessage(err, "BigCache.Preload dataLoader.PreLoad(ctx) error")
	}

	if len(preloadItems) == 0 {
		return nil
	}

	// 將預加載的項目存入緩存
	for key, value := range preloadItems {
		if err = receiver.Set(ctx, key, value); err != nil {
			return pkgerrors.WithMessage(err, "BigCache.Preload Set error")
		}
	}

	return nil
}
