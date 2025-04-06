package clickhouse

import (
	"fmt"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewConnection 创建一个新的ClickHouse连接
func NewConnection(config Config) *gorm.DB {
	if len(config.Nodes) == 0 {
		panic("clickhouse config nodes is empty")
	}

	// 使用第一个节点作为连接点
	mainNode := config.Nodes[0]
	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%d/%s?dial_timeout=%s&read_timeout=%s",
		config.Username,
		config.Password,
		mainNode.Host,
		mainNode.Port,
		config.Database,
		config.DialTimeout.String(),
		config.ReadTimeout.String(),
	)

	// 配置GORM
	gormConfig := &gorm.Config{}
	if config.Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	// 连接数据库
	db, err := gorm.Open(clickhouse.Open(dsn), gormConfig)
	if err != nil {
		panic(fmt.Sprintf("clickhouse gorm.Open error: %v", err))
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("clickhouse gorm.DB error: %v", err))
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	// 测试连接
	if err = sqlDB.Ping(); err != nil {
		panic(fmt.Sprintf("clickhouse sqlDB.Ping error: %v", err))
	}

	return db
}
