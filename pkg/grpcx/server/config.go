package server

import (
	"time"

	"go-tool/pkg/grpcx/server/logger"
)

type Config struct {
	Port            int           `mapstructure:"port" yaml:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout"`
	Logger          logger.Interface
}
