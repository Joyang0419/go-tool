package cache

import (
	"time"

	"github.com/allegro/bigcache/v3"
)

type BigCacheOptFn func(*bigcache.Config)

func BigCacheConfig(expiration time.Duration, opts ...BigCacheOptFn) bigcache.Config {
	bcConfig := bigcache.DefaultConfig(expiration)
	for _, opt := range opts {
		opt(&bcConfig)
	}

	return bcConfig
}

// WithShards 設置分片數量
// 值越大，並發性能越好，但內存消耗越大
// 必須是2的冪
// 默認: 1024
func WithShards(shards int) BigCacheOptFn {
	return func(c *bigcache.Config) {
		c.Shards = shards
	}
}

// WithMaxEntriesInWindow 設置在生命週期窗口內的最大條目數
// 影響初始內存分配，較大的值會預分配更多內存，但減少動態擴展
// 默認: 100000
func WithMaxEntriesInWindow(maxEntries int) BigCacheOptFn {
	return func(c *bigcache.Config) {
		c.MaxEntriesInWindow = maxEntries
	}
}

// WithHardMaxCacheSize 設置緩存的最大大小(MB)
// 0表示無限制，當達到此限制時，將開始逐出舊條目
// 默認: 0 (無限制)
func WithHardMaxCacheSize(maxSize int) BigCacheOptFn {
	return func(c *bigcache.Config) {
		c.HardMaxCacheSize = maxSize
	}
}

// WithVerboseLogging 設置是否輸出詳細日誌
// 用於調試，生產環境通常設為false
// 默認: false
func WithVerboseLogging(verbose bool) BigCacheOptFn {
	return func(c *bigcache.Config) {
		c.Verbose = verbose
	}
}

// WithCleanWindow 設置清理過期項目的時間間隔
// 較小的值可以更快地清理過期項目，但會增加CPU使用率
// 默認: 5分鐘
func WithCleanWindow(window time.Duration) BigCacheOptFn {
	return func(c *bigcache.Config) {
		c.CleanWindow = window
	}
}

// WithMaxEntrySize 設置項目在生命週期窗口內的最大大小（字節）
// 較大的值允許存儲更大的對象，但會增加內存使用
// 默認: 500（字節）
func WithMaxEntrySize(size int) BigCacheOptFn {
	return func(c *bigcache.Config) {
		c.MaxEntrySize = size
	}
}
