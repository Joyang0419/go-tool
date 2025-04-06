package ginx

import (
	"time"
)

type Config struct {
	Port            int           `mapstructure:"port" yaml:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout"`
	// debug, release, test
	Mode string `mapstructure:"mode" yaml:"mode"`
}
