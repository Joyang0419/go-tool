package standalone

import (
	"time"

	"go-tool/pkg/infra/redis/logger"
)

type Config struct {
	// Address 服務器地址 (例如: "localhost:6379")
	Address string `mapstructure:"address" yaml:"address"`
	// ClientName 客戶端名稱，用於在 Redis 服務器上識別連接 (可選)
	ClientName string `mapstructure:"client_name" yaml:"client_name"`
	// Username 認證用戶名
	Username string `mapstructure:"username" yaml:"username"`
	// Password 認證密碼
	Password string `mapstructure:"password" yaml:"password"`
	// DB 數據庫編號 (0-15，默認為 0)
	DB int `mapstructure:"db" yaml:"db"`

	// DialTimeout 建立連接的超時時間 (連接到 Redis 服務器的最大等待時間)
	DialTimeout time.Duration `mapstructure:"dial_timeout" yaml:"dial_timeout"`
	// ReadTimeout 讀取操作的超時時間 (從 Redis 讀取數據的最大等待時間)
	ReadTimeout time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
	// WriteTimeout 寫入操作的超時時間 (向 Redis 寫入數據的最大等待時間)
	WriteTimeout time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`

	// PoolSize 連接池的最大連接數 (同時可以打開的最大連接數)
	PoolSize int `mapstructure:"pool_size" yaml:"pool_size"`
	// PoolTimeout 從連接池獲取連接的超時時間 (等待可用連接的最大時間)
	PoolTimeout time.Duration `mapstructure:"pool_timeout" yaml:"pool_timeout"`
	// MinIdleConns 連接池中保持的最小空閒連接數 (提高響應速度)
	MinIdleConns int `mapstructure:"min_idle_conns" yaml:"min_idle_conns"`
	// MaxIdleConns 連接池中允許的最大空閒連接數 (超過此數量的空閒連接會被關閉)
	MaxIdleConns int `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	// ConnMaxIdleTime 連接的最大空閒時間 (空閒超過此時間的連接會被關閉，釋放資源)
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	// 不給使用者設定，目前是直接設定true, 若使用者有傳入Ctx就可以控制timeout
	//ContextTimeoutEnabled bool `mapstructure:"context_timeout_enabled"`

	// Logger, 如果有注入; 就會啟動ProcessPipelineHook, ProcessHook
	Logger logger.Interface
}
