package clickhouse

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func TestModule(t *testing.T) {
	config := Config{
		Nodes: []Node{
			{Host: "localhost", Port: 9000},
		},
		Database:        "trader",
		Username:        "trader",
		Password:        "",
		Debug:           true,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		DialTimeout:     time.Duration(1) * time.Second,
		ReadTimeout:     time.Duration(1) * time.Second,
	}

	app := fx.New(
		Module(config),
		fx.Invoke(
			func(db *gorm.DB) {
				assert.NotNil(t, db)
			}),
	)

	assert.NoError(t, app.Start(context.TODO()))
}
