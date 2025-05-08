package ginx

import (
	"time"

	"go-tool/pkg/ginx/ginx_middleware"
)

type Config struct {
	Port            int           `mapstructure:"port" yaml:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout"`
	// debug, release, test
	Mode string                     `mapstructure:"mode" yaml:"mode"`
	CORS ginx_middleware.CORSConfig `mapstructure:"cors" yaml:"cors"`
}
