package gorm

import (
	"time"

	"go-tool/pkg/logger"
)

// Config 適配器配置
type Config struct {
	SlowThreshold time.Duration
	Logger        logger.ILogger
	// EnableSuccessSQLLog 成功的SQL, 且不超過 SlowThreshold 要印
	EnableSuccessSQLLog bool
}
