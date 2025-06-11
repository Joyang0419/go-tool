package redis

import (
	"time"

	"go-tool/pkg/logger"
)

// Config 適配器配置
type Config struct {
	SlowThreshold time.Duration `mapstructure:"slow_threshold" yaml:"slow_threshold"`
	Logger        logger.Interface
	// EnableSuccessSQLLog 成功的SQL, 且不超過 SlowThreshold 要印
	EnableSuccessLog bool `mapstructure:"enable_success_log" yaml:"enable_success_log"`
}
