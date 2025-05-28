package slog

type Config struct {
	// Path 存放路徑(直接寫: example: /opt/logs/casino-dealer-client/casino-dealer-client.log)
	Path string `mapstructure:"path" yaml:"path"`
	// MaxSize 一個檔案最大的size
	MaxSize int `mapstructure:"max_size" yaml:"max_size"`
	// MaxAge 保留幾天的log
	MaxAge int `mapstructure:"max_age" yaml:"max_age"`
	// EnableWriteFile 是否要寫入檔案
	EnableWriteFile bool `mapstructure:"enable_write_file" yaml:"enable_write_file"`
}
