package slogx

type Config struct {
	// Path 存放路徑(直接寫: example: /opt/logs/casino-dealer-client/casino-dealer-client.log)
	Path string `mapstructure:"path"`
	// MaxSize 一個檔案最大的size
	MaxSize int `mapstructure:"max_size"`
	// MaxAge 保留幾天的log
	MaxAge int `mapstructure:"max_age"`
	// CallerSkip 要跳過幾層(通常是在找caller的時候, 會用到); 目前測下來是CallerSkip: 7, 會剛好是應用層呼叫的位置
	CallerSkip int `mapstructure:"caller_skip"`
	// EnableWriteFile 是否要寫入檔案
	EnableWriteFile bool `mapstructure:"enable_write_file"`
}
