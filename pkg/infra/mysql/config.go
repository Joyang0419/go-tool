package mysql

import (
	"time"

	gormLogger "gorm.io/gorm/logger"
)

type Config struct {
	// 資料庫名稱
	Database string `mapstructure:"database"`

	// 主資料庫連線設定（寫入用）
	Master ConnectConfig `mapstructure:"master"`

	// 從資料庫連線設定列表（讀取用，支援多個從庫）
	Slaves []ConnectConfig `mapstructure:"slaves"`

	// 連線超時時間（建立連線的最大等待時間）
	ConnectTimeOut time.Duration `mapstructure:"connect_timeout"`

	// 讀取操作超時時間
	ReadTimeOut time.Duration `mapstructure:"read_timeout"`

	// 寫入操作超時時間
	WriteTimeOut time.Duration `mapstructure:"write_timeout"`

	// 時區設定（例如：Asia/Taipei), time.Location要會過, 才給用
	Location string `mapstructure:"location"`

	// 連線池中連線的最大閒置時間
	MaxIdleTime time.Duration `mapstructure:"max_idle_time"`

	// 連線的最大生命週期，超過此時間的連線會被關閉並重新建立
	MaxLifeTime time.Duration `mapstructure:"max_life_time"`

	// 連線池中允許的最大閒置連線數
	MaxIdleConns int `mapstructure:"max_idle_conns"`

	// 連線池中允許的最大開啟連線數
	MaxOpenConns int `mapstructure:"max_open_conns"`

	/*
			SkipDefaultTransaction: false → GORM 明確管理事務（多次網路請求);
				begin一次, 實際SQL一次, commit 一次; 效能較慢, 但更安全
			SkipDefaultTransaction: true → MySQL 自動管理事務（單次網路請求);
		    => 不給設定了，註解掉，知道有這件事就好，有需要在開
	*/
	// SkipDefaultTransaction bool `mapstructure:"skip_default_transaction"`

	// 是否啟用 Prepared Statement（true = 預編譯 SQL 語句，提升效能和安全性）
	// PrepareStmt bool `mapstructure:"prepare_stmt"` <= 故意不給設定，直接用True
	Logger gormLogger.Interface
}

type ConnectConfig struct {
	Username string `mapstructure:"user_name"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}
