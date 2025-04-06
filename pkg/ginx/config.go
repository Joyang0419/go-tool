package ginx

import (
	"time"
)

type Config struct {
	Port            int           `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	// debug, release, test
	Mode string `mapstructure:"mode"`
}
