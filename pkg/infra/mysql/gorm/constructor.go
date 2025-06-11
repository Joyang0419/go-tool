package gorm

import (
	"fmt"
	"net/url"
	"time"

	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	inframysql "go-tool/pkg/infra/mysql"
)

// New 建立 MySQL 連線，純粹的函數，沒有外部依賴
func New(config inframysql.Config) *gorm.DB {
	if config.Location == "" {
		msg := "gorm.New config.Location is empty"
		panic(msg)
	}
	location, err := time.LoadLocation(config.Location)
	if err != nil {
		msg := fmt.Sprintf("gorm.New time.LoadLocation err: %v", err)
		panic(msg)
	}
	locationStr := url.QueryEscape(location.String())

	// 建立 DSN 格式字串
	dsnFormat := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s&timeout=%s&readTimeout=%s&writeTimeout=%s"

	// 建立主資料庫 DSN
	masterDsn := fmt.Sprintf(
		dsnFormat,
		config.Master.Username,
		config.Master.Password,
		config.Master.Host,
		config.Master.Port,
		config.Database,
		locationStr,
		config.ConnectTimeOut.String(),
		config.ReadTimeOut.String(),
		config.WriteTimeOut.String(),
	)

	// 建立從資料庫 DSN 列表
	replicas := make([]gorm.Dialector, 0, len(config.Slaves))
	for _, slave := range config.Slaves {
		slaveDsn := fmt.Sprintf(
			dsnFormat,
			slave.Username,
			slave.Password,
			slave.Host,
			slave.Port,
			config.Database,
			locationStr,
			config.ConnectTimeOut.String(),
			config.ReadTimeOut.String(),
			config.WriteTimeOut.String(),
		)
		replicas = append(replicas, mysql.Open(slaveDsn))
	}

	// logger 設定; 當config有時，就用config的, 沒有就單純初始化DiscardLogger
	l := lo.Ternary[logger.Interface](lo.IsNil(config.Logger), logger.Discard, config.Logger)

	// GORM 設定
	gormConfig := &gorm.Config{
		SkipDefaultTransaction: false, // 明確設定為 false，保證事務安全
		PrepareStmt:            true,  // 固定為 true，提升效能
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用單數表名
		},
		Logger: l,
	}

	// 開啟主資料庫連線
	db, err := gorm.Open(
		mysql.Open(masterDsn),
		gormConfig,
	)
	if err != nil {
		msg := fmt.Sprintf("gorm.New gorm.Open masterDsn: %s, err: %v", masterDsn, err)
		panic(msg)
	}

	// 統一使用 dbresolver，即使沒有從資料庫也一樣
	// 這樣可以統一設定連線池參數，避免 if-else 的奇怪邏輯
	resolverConfig := dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(masterDsn)},
		Policy:   dbresolver.RandomPolicy{},
		Replicas: replicas,
	}

	err = db.Use(
		dbresolver.Register(resolverConfig).
			SetConnMaxIdleTime(config.MaxIdleTime).
			SetConnMaxLifetime(config.MaxLifeTime).
			SetMaxIdleConns(config.MaxIdleConns).
			SetMaxOpenConns(config.MaxOpenConns))

	if err != nil {
		msg := fmt.Sprintf("gorm.New dbresolver.Register err: %v", err)
		panic(msg)
	}

	return db
}
