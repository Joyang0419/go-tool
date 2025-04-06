package clickhouse

import (
	"testing"
	"time"
)

func TestNewConnection(t *testing.T) {
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
		DialTimeout:     time.Duration(10) * time.Second,
		ReadTimeout:     time.Duration(20) * time.Second,
	}

	NewConnection(config)
}
