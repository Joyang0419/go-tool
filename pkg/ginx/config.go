package ginx

import (
	"time"

	"go-tool/pkg/ginx/middleware"
)

type Config struct {
	Port            int                   `mapstructure:"port" yaml:"port"`
	ShutdownTimeout time.Duration         `mapstructure:"shutdown_timeout" yaml:"shutdown_timeout"` // debug, release, test
	Mode            string                `mapstructure:"mode" yaml:"mode"`
	CORS            middleware.CORSConfig `mapstructure:"cors" yaml:"cors"`
	EnableSwagger   bool                  `mapstructure:"enable_swagger" yaml:"enable_swagger"`
}
