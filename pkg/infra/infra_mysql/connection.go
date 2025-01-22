package infra_mysql

import (
	"fmt"
	"time"

	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type ConnectionConfig struct {
	Master struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
	} `mapstructure:"master"`

	Slaves []struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
	} `mapstructure:"slaves"`

	/*
		這些連接池參數的建議值：
		MaxIdleConns：通常設置為 MaxOpenConns 的 25%-50%
		MaxOpenConns：根據服務器配置和數據庫負載能力設置
		MaxLifeTime：避免太長，通常 1 小時左右
		MaxIdleTime：通常幾分鐘，可以及時釋放不需要的連接
	*/
	Pool struct {
		// 最大空閒連接數
		MaxIdleConns int `mapstructure:"max_idle_conns"`
		// 打開的最大連接數
		MaxOpenConns int `mapstructure:"max_open_conns"`
		// 連接可重用的最大時間
		MaxLifeTime time.Duration `mapstructure:"max_life_time"`
		// 空閒連接最大存活時間
		MaxIdleTime time.Duration `mapstructure:"max_idle_time"`
	} `mapstructure:"pool"`
}

func NewMYSQLConnection(config ConnectionConfig) *gorm.DB {
	// 連接主庫
	masterDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Master.User,
		config.Master.Password,
		config.Master.Host,
		config.Master.Port,
		config.Master.DBName,
	)

	db, err := gorm.Open(mysql.Open(masterDSN), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("[NewMySQLDB]gorm.Open(infra_mysql.Open(masterDSN) error: %v", err))
	}

	// 準備從庫的配置
	var replicas []gorm.Dialector
	for _, slave := range config.Slaves {
		slaveDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			slave.User,
			slave.Password,
			slave.Host,
			slave.Port,
			slave.DBName,
		)
		replicas = append(replicas, mysql.Open(slaveDSN))
	}

	// 使用 dbresolver 插件配置讀寫分離
	err = db.Use(
		dbresolver.Register(
			dbresolver.Config{
				// sources: 寫操作使用的數據庫
				Sources: []gorm.Dialector{mysql.Open(masterDSN)},
				// replicas: 讀操作使用的數據庫
				Replicas: replicas,
				/*
					Policy: dbresolver.RandomPolicy{} 隨機策略: 隨機選擇一個從庫
					Policy: dbresolver.RoundRobinPolicy{} 輪詢策略: 按順序輪流使用從庫
					Policy: dbresolver.WeightRandomPolicy{ 權重隨機策略: 根據權重隨機選擇從庫
					Policy: dbresolver.WeightRoundRobinPolicy{Weights: weights,} 權重輪詢策略: 根據權重輪詢選擇從庫
				*/
				Policy: dbresolver.RandomPolicy{}, // 隨機策略
				// 定義讀寫分離的操作
			},
		).
			// 設置最大空閒連接數
			SetMaxIdleConns(config.Pool.MaxIdleConns).
			// 設置最大連接數
			SetMaxOpenConns(config.Pool.MaxOpenConns).
			// 設置連接最大存活時間
			SetConnMaxLifetime(config.Pool.MaxLifeTime).
			// 設置空閒連接最大存活時間
			SetConnMaxIdleTime(config.Pool.MaxIdleTime))

	if err != nil {
		panic(fmt.Sprintf("[NewMySQLDB]db.Use(dbresolver.Register error: %v", err))
	}

	return db
}

func InjectMySQLConnection(config ConnectionConfig) fx.Option {
	return fx.Options(
		fx.Supply(config),
		fx.Provide(NewMYSQLConnection),
	)
}
