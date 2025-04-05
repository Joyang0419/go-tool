package clickhouse

import (
	"time"
)

/*
	使用gorm
	Why? 有支援orm機制。
*/

type Node struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	// Nodes 集群中所有节点的配置列表
	Nodes []Node `mapstructure:"nodes"`
	// Database 要连接的数据库名称
	Database string `mapstructure:"database"`
	// Username 连接用户名
	Username string `mapstructure:"username"`
	// Password 连接密码
	Password string `mapstructure:"password"`
	// Debug 是否启用调试模式（启用GORM详细日志）
	Debug bool `mapstructure:"debug"`
	// ClusterName ClickHouse集群名称，应与ClickHouse服务器配置中的名称一致
	ClusterName string `mapstructure:"cluster_name"`
	// MaxIdleConns 连接池中的最大空闲连接数
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	// MaxOpenConns 连接池中的最大打开连接数
	MaxOpenConns int `mapstructure:"max_open_conns"`
	// ConnMaxLifetime 连接的最大生命周期（格式：1h, 30m等）
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	// DialTimeout 连接超时时间（格式：10s, 1m等）
	DialTimeout time.Duration `mapstructure:"dial_timeout"`
	// ReadTimeout 读取超时时间（格式：20s, 1m等）
	ReadTimeout time.Duration `mapstructure:"read_timeout"`
}
