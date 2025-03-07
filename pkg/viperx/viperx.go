package viperx

import (
	"embed"
	"fmt"
	"log/slog"
	"strings"

	pkgerrors "github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// LoadConfig 載入配置文件
// path: 配置文件路徑
// name: 配置文件名稱（不包含副檔名）
// configType: 配置文件類型（如 yaml, json 等）
// config: 配置結構體的指針
func LoadConfig(path string, name string, configType string, config interface{}) error {
	v := viper.New()

	// 設置配置文件路徑
	v.AddConfigPath(path)

	// 設置配置文件名稱
	v.SetConfigName(name)

	// 設置配置文件類型
	v.SetConfigType(configType)

	// 設置環境變數前綴
	v.SetEnvPrefix(strings.ToUpper(name))

	// 將 . 轉換為 _
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 自動載入環境變數
	v.AutomaticEnv()

	// 讀取配置文件
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("[LoadConfig]v.ReadInConfig error: %w", err)
	}

	// 解析到結構體
	if err := v.Unmarshal(config); err != nil {
		return fmt.Errorf("[LoadConfig]v.Unmarshal: %w", err)
	}

	return nil
}

// LoadEmbeddedConfig 加载嵌入式配置文件到指定的结构体
// config: 配置结构体的指针
func LoadEmbeddedConfig(embed embed.FS, filename string, config interface{}) error {
	// 读取嵌入式配置
	data, err := embed.ReadFile(filename)
	if err != nil {
		return pkgerrors.Errorf("LoadEmbeddedConfig embed.ReadFile error %v", err)
	}

	// 解析YAML到结构体
	if err = yaml.Unmarshal(data, config); err != nil {
		return pkgerrors.Errorf("LoadEmbeddedConfig yaml.Unmarshal error %v", err)
	}

	slog.Info("LoadEmbeddedConfig success", slog.String("filename", filename))

	return nil
}
