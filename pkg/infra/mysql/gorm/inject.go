package gorm

import (
	"go.uber.org/fx"
	gormLogger "gorm.io/gorm/logger"

	inframysql "go-tool/pkg/infra/mysql"
)

func FxModule(config inframysql.Config, log ...gormLogger.Interface) fx.Option {
	var options []fx.Option
	if len(log) > 0 {
		config.Logger = log[0]
	}

	options = append(options, fx.Supply(config))
	options = append(options, fx.Provide(New))

	return fx.Module("infra_gorm_module", options...)
}
